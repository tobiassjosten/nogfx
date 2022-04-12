package achaea

import (
	"bytes"
	"log"

	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/telnet"
	"github.com/tobiassjosten/nogfx/pkg/world/achaea/gmcp"
)

type World struct {
	client pkg.Client
}

func NewWorld(client pkg.Client) *World {
	return &World{
		client: client,
	}
}

func (world *World) Input(input []byte) []byte {
	return input
}

func (world *World) Output(output []byte) []byte {
	return output
}

func (world *World) Command(command []byte) {
	willGMCP := []byte{telnet.IAC, telnet.WILL, telnet.GMCP}
	prefixGMCP := []byte{telnet.IAC, telnet.SB, telnet.GMCP}
	suffixGMCP := []byte{telnet.IAC, telnet.SE}

	switch {
	case bytes.Equal(command, willGMCP):
		// @todo Use the actual version number when we have one.
		world.gmcp(gmcp.CoreHello{Client: "NoGFX", Version: "0.0.1"})
		world.gmcp(gmcp.CoreSupportsSet{
			Char:        true,
			CharSkills:  true,
			CharItems:   true,
			CommChannel: true,
			Room:        true,
			IRERift:     true,
		})

	case bytes.HasPrefix(command, prefixGMCP):
		data := command[len(prefixGMCP) : len(command)-len(suffixGMCP)]
		message, err := gmcp.Parse(data)
		if err != nil {
			log.Printf("failed parsing GMCP command: %s", err)
			return
		}

		switch msg := message.(type) {
		case gmcp.CharName:
			world.gmcp(gmcp.IRERiftRequest{})
			world.gmcp(gmcp.CommChannelPlayers{})
			world.gmcp(gmcp.CharItemsInv{})
			// @todo Update `world` with `msg`.

		default: // Noop.
		}

	default: // Noop.
	}
}

func (world *World) gmcp(value gmcp.Message) error {
	_, err := world.client.Write(append(append(
		[]byte{telnet.IAC, telnet.SB, telnet.GMCP},
		[]byte(value.String())...,
	), telnet.IAC, telnet.SE))
	return err
}
