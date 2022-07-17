package achaea

import (
	"bytes"
	"fmt"
	"log"

	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/gmcp"
	agmcp "github.com/tobiassjosten/nogfx/pkg/gmcp/achaea"
	igmcp "github.com/tobiassjosten/nogfx/pkg/gmcp/ironrealms"
	"github.com/tobiassjosten/nogfx/pkg/navigation"
	"github.com/tobiassjosten/nogfx/pkg/simpex"
	"github.com/tobiassjosten/nogfx/pkg/telnet"
	amodule "github.com/tobiassjosten/nogfx/pkg/world/achaea/module"
	"github.com/tobiassjosten/nogfx/pkg/world/module"
)

// World is an Achaea-specific implementation of the pkg.World interface.
type World struct {
	client pkg.Client

	ui       pkg.UI
	uiVitals map[string]struct{}

	inputTriggers  []pkg.Trigger[pkg.Input]
	outputTriggers []pkg.Trigger[pkg.Output]

	Character *Character
	Room      *navigation.Room
	Target    *Target
}

// NewWorld creates a new Achaea-specific pkg.World.
func NewWorld(client pkg.Client, ui pkg.UI) pkg.World {
	world := &World{
		client: client,

		ui:       ui,
		uiVitals: map[string]struct{}{},

		inputTriggers:  []pkg.Trigger[pkg.Input]{},
		outputTriggers: []pkg.Trigger[pkg.Output]{},

		Character: &Character{},
		Target:    NewTarget(client),
	}

	// @todo Make sure these are ordered correctly. Potentially by adding a weight
	// property for sorting?
	var moduleConstructors = []pkg.ModuleConstructor{
		module.NewRepeatInput,
		amodule.NewLearnMultipleLessons,
	}

	for _, constructor := range moduleConstructors {
		module := constructor(world)
		world.inputTriggers = append(
			world.inputTriggers, module.InputTriggers()...,
		)
		world.outputTriggers = append(
			world.outputTriggers, module.OutputTriggers()...,
		)
	}

	return world
}

// Print passes data onto the configured UI.
func (world *World) Print(data []byte) {
	world.ui.Print(data)
}

// Send passes data onto the configured Client.
func (world *World) Send(data []byte) {
	world.client.Send(data)
}

// ProcessInput processes player input.
func (world *World) ProcessInput(input pkg.Input) pkg.Input {
	// @todo Read the CommandSeparator configuration option and use that.
	sep := []byte{';'}

	input = input.Split(sep)

	for _, trigger := range world.inputTriggers {
		for i, command := range input {
			match := simpex.Match(trigger.Pattern, command.Text)
			if match == nil {
				continue
			}

			input = trigger.Callback(pkg.TriggerMatch[pkg.Input]{
				Captures: match,
				Content:  input,
				Index:    i,
			})
		}
	}

	return input.Join(sep)
}

// ProcessOutput processes game output.
func (world *World) ProcessOutput(output pkg.Output) pkg.Output {
	// If the first line is empty (save for ANSI colors) we remove it, as
	// it's been added to compensate for echoed player input.
	if len(output) >= 3 && len(output[0].Text) == 0 {
		if len(output[0].Raw) > 0 {
			output[1].Raw = append(output[0].Raw, output[1].Raw...)
		}
		output = output.Remove(0)
	}

	for _, trigger := range world.outputTriggers {
		for i, line := range output {
			match := simpex.Match(trigger.Pattern, line.Text)
			if match == nil {
				continue
			}

			output = trigger.Callback(pkg.TriggerMatch[pkg.Output]{
				Captures: match,
				Content:  output,
				Index:    i,
			})
		}
	}

	// If only the prompt remains, we omit the whole paragraph.
	if len(output) == 1 {
		output = pkg.Output{}
	}

	return output
}

// ProcessCommand processes telnet commands.
func (world *World) ProcessCommand(command []byte) {
	if data := gmcp.Unwrap(command); data != nil {
		err := world.processGMCP(data)
		if err != nil {
			log.Printf("failed processing gmcp: %s", err)
		}

		return
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
			log.Printf("failed sending gmcp: %s", err)
		}
	}
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
