package pkg

import (
	"bufio"
	"io"
	"log"

	"github.com/tobiassjosten/nogfx/pkg/telnet"
)

type Engine struct {
	ui     UI
	client Client
	world  World
}

type UI interface {
	Run(<-chan []byte)
}

type Client interface {
	io.ReadWriter
	TelnetClient
}

type TelnetClient interface {
	Will(byte) error
	Wont(byte) error
	Do(byte) error
	Dont(byte) error
	Subneg(byte, []byte) error
}

func NewEngine(world World, ui UI, client Client) *Engine {
	return &Engine{
		ui:     ui,
		client: client,
		world:  world,
	}
}

func (engine *Engine) Run(inputs <-chan []byte, commands <-chan []byte) {
	uiOutput := make(chan []byte)
	go engine.ui.Run(uiOutput)

	serverOutput := make(chan []byte)
	serverDone := make(chan struct{})
	go func() {
		scanner := bufio.NewScanner(engine.client)
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
			input = engine.world.Input(input)
			if len(input) > 0 {
				engine.client.Write(append(input, '\n'))
			}

		case output := <-serverOutput:
			output = engine.world.Output(output)
			if len(output) > 0 {
				uiOutput <- output
			}

		case command, ok := <-commands:
			if !ok {
				continue
			}

			engine.world.Command(command)
		}
	}
}
