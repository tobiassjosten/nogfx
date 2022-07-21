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
	"github.com/tobiassjosten/nogfx/pkg/telnet"
	amodule "github.com/tobiassjosten/nogfx/pkg/world/achaea/module"
	"github.com/tobiassjosten/nogfx/pkg/world/module"
)

// World is an Achaea-specific implementation of the pkg.World interface.
type World struct {
	client pkg.Client

	ui       pkg.UI
	uiVitals map[string]struct{}

	triggers []pkg.Trigger

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

		Character: &Character{},
		Target:    NewTarget(client),
	}

	// @todo Make sure these are ordered correctly. Potentially by adding a weight
	// property for sorting?
	var modules = []pkg.Module{
		module.NewRepeatInput(),
		amodule.NewLearnMultipleLessons(),
	}

	for _, module := range modules {
		world.triggers = append(world.triggers, module.Triggers()...)
	}

	return world
}

// OnInoutput reacts to player input and server output.
func (world *World) OnInoutput(inout pkg.Inoutput) pkg.Inoutput {
	// @todo Read the CommandSeparator configuration option and use that.
	sep := []byte{';'}
	inout.Input = inout.Input.Split(sep)

	// @todo Append a prompt if there isn't one (e.g. detect that first).

	// If the first line is empty (save for ANSI colors) we remove it, as
	// it's been added to compensate for echoed player input.
	if len(inout.Output) >= 3 && len(inout.Output[0].Text.Clean()) == 0 {
		// Move potential control codes to the next line.
		if len(inout.Output[0].Text) > 0 {
			inout.Output[1].Text = append(
				inout.Output[0].Text,
				inout.Output[1].Text...,
			)
		}
		inout.Output = inout.Output.Omit(0)
	}

	for _, trigger := range world.triggers {
		if trigger.Kind == pkg.Input && len(inout.Input) > 0 {
			inout = trigger.Match(inout.Input.Bytes(), inout)
		}
		if trigger.Kind == pkg.Output && len(inout.Output) > 0 {
			inout = trigger.Match(inout.Output.Bytes(), inout)
		}
	}

	// If only the prompt remains, we omit the whole paragraph.
	// @todo Use above prompt detection instead of relying on lone lines
	// being prompts.
	if len(inout.Output.Bytes()) == 1 {
		inout.Output = pkg.Exput{}
	}

	return inout
}

// OnCommand reacts to telnet commands.
// @todo Consider merging this with OnOutput() or making it a callback for GMCP
// only. Telnet commands are cool and all but YAGNI, evidently.
func (world *World) OnCommand(cmd []byte) (inout pkg.Inoutput) {
	if data := gmcp.Unwrap(cmd); data != nil {
		err := world.onGMCP(data)
		if err != nil {
			log.Printf("failed processing gmcp: %s", err)
		}

		return
	}

	switch {
	case bytes.Equal(cmd, []byte{telnet.IAC, telnet.WILL, telnet.GMCP}):
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

	return
}

func (world *World) onGMCP(data []byte) error {
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
