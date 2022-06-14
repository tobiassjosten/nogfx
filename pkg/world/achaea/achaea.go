package achaea

import (
	"bytes"
	"fmt"

	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/gmcp"
	agmcp "github.com/tobiassjosten/nogfx/pkg/gmcp/achaea"
	igmcp "github.com/tobiassjosten/nogfx/pkg/gmcp/ironrealms"
	"github.com/tobiassjosten/nogfx/pkg/navigation"
	"github.com/tobiassjosten/nogfx/pkg/telnet"
	amodule "github.com/tobiassjosten/nogfx/pkg/world/achaea/module"
	"github.com/tobiassjosten/nogfx/pkg/world/module"
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
		Target:    NewTarget(client),
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
		return world.processGMCP(data)
	}

	switch {
	case bytes.Equal(command, []byte{telnet.IAC, telnet.WILL, telnet.GMCP}):
		err := world.SendGMCP(&gmcp.CoreSupportsSet{
			"Char":         1,
			"Char.Items":   1,
			"Char.Skills":  1,
			"Comm.Channel": 1,
			"Room":         1,
			"IRE.Rift":     1,
			"IRE.Target":   1,
		})
		if err != nil {
			return fmt.Errorf("failed GMCP: %w", err)
		}
	}

	return nil
}

func (world *World) processGMCP(data []byte) error {
	message, err := agmcp.Parse(data)
	if err != nil {
		return fmt.Errorf("failed parsing GMCP: %w", err)
	}

	switch msg := message.(type) {
	case *gmcp.CharItemsList:
		world.Target.FromCharItemsList(msg)
		world.ui.SetTarget(world.Target.PkgTarget())

	case *gmcp.CharItemsAdd:
		world.Target.FromCharItemsAdd(msg)
		world.ui.SetTarget(world.Target.PkgTarget())

	case *gmcp.CharItemsRemove:
		world.Target.FromCharItemsRemove(msg)
		world.ui.SetTarget(world.Target.PkgTarget())

	case *gmcp.CharName:
		world.Character.FromCharName(msg)

		msgs := []gmcp.Message{
			&gmcp.CharItemsInv{},
			&gmcp.CommChannelPlayers{},
			&igmcp.IRERiftRequest{},
		}
		for _, msg := range msgs {
			data := gmcp.Wrap([]byte(msg.ID()))
			if _, err := world.client.Write(data); err != nil {
				return fmt.Errorf("failed GMCP: %w", err)
			}
		}

	case *agmcp.CharStatus:
		world.Character.FromCharStatus(msg)
		world.ui.SetCharacter(world.Character.PkgCharacter())

		world.Target.FromCharStatus(msg)
		world.ui.SetTarget(world.Target.PkgTarget())

	case *agmcp.CharVitals:
		world.Character.FromCharVitals(msg)
		world.ui.SetCharacter(world.Character.PkgCharacter())

	case *gmcp.RoomInfo:
		world.Target.FromRoomInfo(msg)
		world.ui.SetTarget(world.Target.PkgTarget())

		if world.Room != nil {
			world.Room.HasPlayer = false
		}
		world.Room = navigation.RoomFromGMCP(msg)
		world.Room.HasPlayer = true

		world.ui.SetRoom(world.Room)

		// @todo Implement this to download the official map.
		// case gmcp.ClientMap:

	case *igmcp.IRETargetSet:
		world.Target.FromIRETargetSet(msg)
		world.ui.SetTarget(world.Target.PkgTarget())

	case *igmcp.IRETargetInfo:
		world.Target.FromIRETargetInfo(msg)
		world.ui.SetTarget(world.Target.PkgTarget())
	}

	return nil
}

// SendGMCP writes a GMCP message to the client.
func (world *World) SendGMCP(msg gmcp.Message) error {
	data := gmcp.Wrap([]byte(msg.Marshal()))
	if _, err := world.client.Write(data); err != nil {
		return err
	}

	return nil
}
