package achaea

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/gmcp"
	agmcp "github.com/tobiassjosten/nogfx/pkg/gmcp/achaea"
	igmcp "github.com/tobiassjosten/nogfx/pkg/gmcp/ironrealms"
	"github.com/tobiassjosten/nogfx/pkg/navigation"
	"github.com/tobiassjosten/nogfx/pkg/process"
	"github.com/tobiassjosten/nogfx/pkg/telnet"
)

type world struct {
	client pkg.Client
	ui     pkg.UI

	Character *Character
	Room      *navigation.Room
	Target    *Target
}

// Processor contains all the game logic for Achaea.
func Processor(client pkg.Client, ui pkg.UI) (process.Processor, error) {
	world := &world{
		client: client,
		ui:     ui,

		Character: &Character{},
		Target:    NewTarget(client),
	}

	now := time.Now().Format("20060102-150405")

	rawLogProcessor, err := process.LogProcessor(
		filepath.Join(pkg.Directory, "logs"),
		fmt.Sprintf("achaea.com-%s.raw.log", now),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create log processor: %w", err)
	}

	logProcessor, err := process.LogProcessor(
		filepath.Join(pkg.Directory, "logs"),
		fmt.Sprintf("achaea.com-%s.log", now),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create log processor: %w", err)
	}

	// @todo Read the CommandSeparator configuration and use that instead.
	sep := []byte{';'}

	tv := TunnelVision{character: world.Character}

	return process.ChainProcessor(
		rawLogProcessor,
		process.SplitInputProcessor(sep),
		process.RepeatInputProcessor(),
		RewriteOutput,

		process.ProcessorFunc(world.cmdprocess),

		(&Learning{}).Processor(),
		// NewBashingMod(world),

		tv.rewriteProcessor(),
		logProcessor,
	), nil
}

func RewriteOutput(outs [][]byte) ([][]byte, [][]byte, error) {
	// Empty first lines (save for ANSI codes) are added to compensate for
	// echoed player input. With our output control, we can omit that.
	if len(outs) >= 3 && len(stripANSI(outs[0])) == 0 {
		// Relay ANSI codes before removing the line.
		if len(outs[0]) > 0 {
			outs[1] = append(outs[0], outs[1]...)
		}
		outs = outs[1:]
	}

	return nil, outs, nil
}

func (world *world) cmdprocess(ioin pkg.Inoutput) (io pkg.Inoutput) {
	io = ioin

	for _, cmd := range ioin.Commands() {
		switch {
		case bytes.Equal(cmd, telnet.IAC_WILL_GMCP):
			io.Input = io.Input.Append([]byte((&gmcp.CoreSupportsSet{
				"Char":         1,
				"Char.Items":   1,
				"Char.Skills":  1,
				"Comm.Channel": 1,
				"Room":         1,
				"IRE.Rift":     1,
				"IRE.Target":   1,
			}).Marshal()))
		}
	}

	for _, data := range ioin.GMCPs() {
		message, err := agmcp.Parse(data)
		if err != nil {
			log.Printf("failed parsing GMCP: %s", err)
			return
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

			inputs := [][]byte{
				[]byte((&gmcp.CharItemsInv{}).Marshal()),
				[]byte((&gmcp.CommChannelPlayers{}).Marshal()),
				[]byte((&igmcp.IRERiftRequest{}).Marshal()),
			}
			for _, input := range inputs {
				io.Input = io.Input.Append(input)
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

		case *igmcp.IRETargetSet:
			world.Target.FromIRETargetSet(msg)
			world.ui.SetTarget(world.Target.PkgTarget())

		case *igmcp.IRETargetInfo:
			world.Target.FromIRETargetInfo(msg)
			world.ui.SetTarget(world.Target.PkgTarget())
		}
		// @todo Implement gmcp.ClientMap to trigger downloading
		// official map data.
	}

	return
}

// @todo Move this somewhere else. It's too general to belong with Achaea
// specifics but we currently don't have a good other place to put it.
func stripANSI(text []byte) (clean []byte) {
	var sequence bool
	for _, c := range text {
		if c == 0x1b {
			sequence = true
			continue
		}

		if sequence {
			if c == 'm' {
				sequence = false
			}
			continue
		}

		clean = append(clean, c)
	}

	return clean
}
