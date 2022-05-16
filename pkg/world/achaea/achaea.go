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

	"github.com/icza/gox/gox"
)

// World is an Achaea-specific implementation of the pkg.World interface.
type World struct {
	client pkg.Client

	ui       pkg.UI
	uiVitals map[string]struct{}

	modules []pkg.Module

	character *Character
	room      *navigation.Room
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

		character: &Character{},
	}
}

// ProcessInput processes player input.
func (world *World) ProcessInput(input []byte) []byte {
	for _, module := range world.modules {
		if input = module.ProcessInput(input); input == nil {
			break
		}
	}

	return input
}

// ProcessOutput processes game output.
func (world *World) ProcessOutput(output []byte) []byte {
	for _, module := range world.modules {
		if output = module.ProcessOutput(output); output == nil {
			break
		}
	}

	return output
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
				IRERift: gox.NewInt(1),
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
	case gmcp.CharName:
		world.character.FromCharName(msg)

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
		world.character.FromCharStatus(msg)

	case agmcp.CharVitals:
		world.character.FromCharVitals(msg)
		if err := world.UpdateVitals(); err != nil {
			return err
		}

	case gmcp.RoomInfo:
		if world.room != nil {
			world.room.HasPlayer = false
		}
		world.room = navigation.RoomFromGMCP(msg)
		world.room.HasPlayer = true

		world.ui.SetRoom(world.room)

		// @todo Implement this to download the official map.
		// case gmcp.ClientMap:
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
		"health":    {world.character.Health, world.character.MaxHealth},
		"mana":      {world.character.Mana, world.character.MaxMana},
		"endurance": {world.character.Endurance, world.character.MaxEndurance},
		"willpower": {world.character.Willpower, world.character.MaxWillpower},
	}

	for name, value := range values {
		err := world.ui.UpdateVital(name, value[0], value[1])
		if err != nil {
			return fmt.Errorf("failed updating vital: %w", err)
		}
	}

	return nil
}
