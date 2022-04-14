package achaea

import (
	"bytes"
	"fmt"

	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/telnet"
	"github.com/tobiassjosten/nogfx/pkg/world/achaea/gmcp"
)

type World struct {
	ui     pkg.UI
	client pkg.Client

	character Character
}

func NewWorld(ui pkg.UI, client pkg.Client) *World {
	return &World{
		ui:     ui,
		client: client,

		character: Character{},
	}
}

func (world *World) Input(input []byte) []byte {
	return input
}

func (world *World) Output(output []byte) []byte {
	return output
}

// func(command []byte) (responses [][]byte, error)

func (world *World) Command(command []byte) error {
	willEcho := []byte{telnet.IAC, telnet.WILL, telnet.ECHO}
	wontEcho := []byte{telnet.IAC, telnet.WONT, telnet.ECHO}

	willGMCP := []byte{telnet.IAC, telnet.WILL, telnet.GMCP}
	prefixGMCP := []byte{telnet.IAC, telnet.SB, telnet.GMCP}
	suffixGMCP := []byte{telnet.IAC, telnet.SE}

	if !bytes.HasPrefix(command, prefixGMCP) {
		world.ui.Print([]byte(fmt.Sprintf(
			"[Telnet command: %s]",
			telnet.CommandToString(command),
		)))
	}

	switch {
	case bytes.Equal(command, willEcho):
		world.ui.MaskInput()

	case bytes.Equal(command, wontEcho):
		world.ui.UnmaskInput()

	case bytes.Equal(command, willGMCP):
		// @todo Use the actual version number when we have one.
		err := world.gmcp(gmcp.CoreHello{
			Client:  "NoGFX",
			Version: "0.0.1",
		})
		if err != nil {
			return fmt.Errorf("failed GMCP: %w", err)
		}

		err = world.gmcp(gmcp.CoreSupportsSet{
			Char:        true,
			CharSkills:  true,
			CharItems:   true,
			CommChannel: true,
			Room:        true,
			IRERift:     true,
		})
		if err != nil {
			return fmt.Errorf("failed GMCP: %w", err)
		}

	case bytes.HasPrefix(command, prefixGMCP):
		data := command[len(prefixGMCP) : len(command)-len(suffixGMCP)]
		message, err := gmcp.Parse(data)
		if err != nil {
			world.ui.Print([]byte(fmt.Sprintf("[Invalid GMCP: %s]", err)))
			return nil
		}

		switch msg := message.(type) {
		case gmcp.CharName:
			world.character.fromCharName(msg)

			// We have just logged in, so let's do an inventory.
			err := world.gmcp(gmcp.IRERiftRequest{})
			if err != nil {
				return fmt.Errorf("failed GMCP: %w", err)
			}

			err = world.gmcp(gmcp.CommChannelPlayers{})
			if err != nil {
				return fmt.Errorf("failed GMCP: %w", err)
			}

			err = world.gmcp(gmcp.CharItemsInv{})
			if err != nil {
				return fmt.Errorf("failed GMCP: %w", err)
			}

		case gmcp.CharVitals:
			world.character.fromCharVitals(msg)
		}
	}

	return nil
}

func (world *World) gmcp(value gmcp.ClientMessage) error {
	_, err := world.client.Write(append(append(
		[]byte{telnet.IAC, telnet.SB, telnet.GMCP},
		[]byte(value.String())...,
	), telnet.IAC, telnet.SE))
	return err
}
