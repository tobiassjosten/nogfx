package achaea

import (
	"bytes"
	"fmt"

	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/gmcp"
	"github.com/tobiassjosten/nogfx/pkg/navigation"
	"github.com/tobiassjosten/nogfx/pkg/telnet"
	"github.com/tobiassjosten/nogfx/pkg/tui"
	"github.com/tobiassjosten/nogfx/pkg/world/achaea/agmcp"
	amodule "github.com/tobiassjosten/nogfx/pkg/world/achaea/module"
	"github.com/tobiassjosten/nogfx/pkg/world/module"

	"github.com/icza/gox/gox"
)

// World is an Achaea-specific implementation of the pkg.World interface.
type World struct {
	client pkg.Client

	ui       pkg.UI
	uiVitals map[string]struct{}

	modules []pkg.Module

	Character *Character
	Room      *navigation.Room
	Target    *Target
}

// NewWorld creates a new Achaea-specific pkg.World.
func NewWorld(client pkg.Client, ui pkg.UI) pkg.World {
	var modules []pkg.Module
	for _, constructor := range moduleConstructors {
		modules = append(modules, constructor(client, ui))
	}

	return &World{
		client: client,

		ui:       ui,
		uiVitals: map[string]struct{}{},

		modules: modules,

		Character: &Character{},
		Target:    &Target{client: client},
	}
}

var moduleConstructors = []pkg.ModuleConstructor{
	module.NewRepeatInput,
	amodule.NewLearnMultipleLessons,
}

// ProcessInput processes player input.
func (world *World) ProcessInput(input []byte) [][]byte {
	// @todo Figure out if the `;` character is configurable, so that we
	// have to make this dynamic.
	inputs := bytes.Split(input, []byte(";"))

	inputs = processInputs(inputs, world.modules)
	if inputs == nil {
		return nil
	}

	// @todo Figure out if the `;` character is configurable, so that we
	// have to make this dynamic.
	return [][]byte{bytes.Join(inputs, []byte(";"))}
}

func processInputs(inputs [][]byte, modules []pkg.Module) [][]byte {
	nullinputs := false

	for i := 0; i < len(inputs); i++ {
		input := inputs[i]

		var newinputs [][]byte

		for _, module := range modules {
			newnewinputs := module.ProcessInput(input)
			if newnewinputs == nil {
				nullinputs = true
				continue
			}

			newinputs = append(newinputs, newnewinputs...)
		}

		if nullinputs {
			inputs = append(inputs[:i], inputs[i+1:]...)
		} else if len(newinputs) > 0 {
			inputs = append(inputs[:i], append(
				newinputs,
				inputs[i+1:]...,
			)...)
		}
	}

	if nullinputs && len(inputs) == 0 {
		return nil
	}

	return inputs
}

// ProcessOutput processes game output.
func (world *World) ProcessOutput(output []byte) [][]byte {
	var outputs [][]byte

	for _, module := range world.modules {
		newoutputs := module.ProcessOutput(output)
		if newoutputs == nil {
			return nil
		}
		outputs = append(outputs, newoutputs...)
	}

	if len(outputs) == 0 {
		outputs = [][]byte{output}
	}

	return outputs
}

// ProcessCommand processes telnet commands.
func (world *World) ProcessCommand(command []byte) error {
	if data := gmcp.Unwrap(command); data != nil {
		return world.ProcessGMCP(data)
	}

	switch {
	case bytes.Equal(command, []byte{telnet.IAC, telnet.WILL, telnet.GMCP}):
		err := world.SendGMCP(agmcp.CoreSupportsSet{
			CoreSupports: agmcp.CoreSupports{
				CoreSupports: gmcp.CoreSupports{
					Char:        gox.NewInt(1),
					CharSkills:  gox.NewInt(1),
					CharItems:   gox.NewInt(1),
					CommChannel: gox.NewInt(1),
					Room:        gox.NewInt(1),
				},
				IRERift:   gox.NewInt(1),
				IRETarget: gox.NewInt(1),
			},
		})
		if err != nil {
			return fmt.Errorf("failed GMCP: %w", err)
		}
	}

	return nil
}

// ServerMessages maps GMCP messages to associated structs.
var ServerMessages = map[string]gmcp.ServerMessage{
	"Char.Status": agmcp.CharStatus{},
	"Char.Vitals": agmcp.CharVitals{},
}

// ProcessGMCP processes GMCP messages.
func (world *World) ProcessGMCP(data []byte) error {
	message, err := gmcp.Parse(data, ServerMessages)
	if err != nil {
		return fmt.Errorf("failed parsing GMCP: %w", err)
	}

	switch msg := message.(type) {
	case gmcp.CharItemsList:
		world.Target.FromCharItemsList(msg)

	case gmcp.CharItemsAdd:
		world.Target.FromCharItemsAdd(msg)

	case gmcp.CharItemsRemove:
		world.Target.FromCharItemsRemove(msg)

	case gmcp.CharName:
		world.Character.FromCharName(msg)

		if err := world.SendGMCP(gmcp.CharItemsInv{}); err != nil {
			return fmt.Errorf("failed GMCP: %w", err)
		}

		if err := world.SendGMCP(gmcp.CommChannelPlayers{}); err != nil {
			return fmt.Errorf("failed GMCP: %w", err)
		}

		if err := world.SendGMCP(gmcp.IRERiftRequest{}); err != nil {
			return fmt.Errorf("failed GMCP: %w", err)
		}

	case agmcp.CharStatus:
		world.Character.FromCharStatus(msg)
		world.Target.FromCharStatus(msg)

	case agmcp.CharVitals:
		world.Character.FromCharVitals(msg)
		if err := world.UpdateVitals(); err != nil {
			return err
		}

	case gmcp.RoomInfo:
		world.Target.FromRoomInfo(msg)

		if world.Room != nil {
			world.Room.HasPlayer = false
		}
		world.Room = navigation.RoomFromGMCP(msg)
		world.Room.HasPlayer = true

		world.ui.SetRoom(world.Room)

		// @todo Implement this to download the official map.
		// case gmcp.ClientMap:

	case gmcp.IRETargetSet:
		world.Target.FromIRETargetSet(msg)

	case gmcp.IRETargetInfo:
		world.Target.FromIRETargetInfo(msg)
	}

	return nil
}

// SendGMCP writes a GMCP message to the client.
func (world *World) SendGMCP(message gmcp.ClientMessage) error {
	data := []byte(message.String())
	if _, err := world.client.Write(gmcp.Wrap(data)); err != nil {
		return err
	}

	return nil
}

var (
	vorder = []string{"health", "mana", "endurance", "willpower"}
	vitals = map[string]*tui.Vital{
		"health":    tui.NewHealthVital(),
		"mana":      tui.NewManaVital(),
		"endurance": tui.NewEnduranceVital(),
		"willpower": tui.NewWillpowerVital(),
	}
)

// UpdateVitals creates sends new current and max values to UI's VitalPanes.
func (world *World) UpdateVitals() error {
	for len(vorder) > 0 {
		err := world.ui.AddVital(vorder[0], vitals[vorder[0]])
		if err != nil {
			return fmt.Errorf("failed adding vital: %w", err)
		}
		vorder = vorder[1:]
	}

	values := map[string][]int{
		"health":    {world.Character.Health, world.Character.MaxHealth},
		"mana":      {world.Character.Mana, world.Character.MaxMana},
		"endurance": {world.Character.Endurance, world.Character.MaxEndurance},
		"willpower": {world.Character.Willpower, world.Character.MaxWillpower},
	}

	for name, value := range values {
		err := world.ui.UpdateVital(name, value[0], value[1])
		if err != nil {
			return fmt.Errorf("failed updating vital: %w", err)
		}
	}

	return nil
}
