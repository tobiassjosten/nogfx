package pkg

import (
	"bufio"
	"bytes"
	"io"
	"log"

	"github.com/tobiassjosten/nogfx/pkg/gmcp"
	"github.com/tobiassjosten/nogfx/pkg/telnet"
)

type World struct {
	ui     UI
	client Client
}

type UI interface {
	Run(<-chan []byte)
}

type Client interface {
	io.ReadWriter
	Will(byte) error
	Wont(byte) error
	Do(byte) error
	Dont(byte) error
	Subneg(byte, []byte) error
}

func NewWorld(ui UI, client Client) *World {
	return &World{
		ui:     ui,
		client: client,
	}
}

func (world *World) Run(inputs <-chan []byte, commands <-chan []byte) {
	uiOutput := make(chan []byte)
	go world.ui.Run(uiOutput)

	serverOutput := make(chan []byte)
	serverDone := make(chan struct{})
	go func() {
		scanner := bufio.NewScanner(world.client)
		scanner.Split(telnet.ScanGA)

		for scanner.Scan() {
			serverOutput <- scanner.Bytes()
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		serverDone <- struct{}{}
	}()

	for {
		select {
		case _ = <-serverDone:
			uiOutput <- []byte("server disconnected")

		case input := <-inputs:
			// @todo Process input.
			world.client.Write(append(input, '\n'))

		case output := <-serverOutput:
			// @todo Process output.
			uiOutput <- output

		case command, ok := <-commands:
			if !ok {
				continue
			}

			world.processCommand(command)
		}
	}
}

func (world *World) processCommand(command []byte) {
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
			log.Printf(
				"failed parsing GMCP command '%s': %s",
				string(data), err,
			)
			return
		}

		switch message.(type) {
		case gmcp.CharName:
			world.gmcp(gmcp.IRERiftRequest{})
			world.gmcp(gmcp.CommChannelPlayers{})
			world.gmcp(gmcp.CharItemsInv{})
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
