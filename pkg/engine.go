package pkg

import (
	"bufio"
	"fmt"
	"io"
)

// Engine is the orchestrator of the entire application, tying together input
// and output from the different parts.
type Engine struct {
	ui     UI
	client Client
	world  World
}

// UI is the primary user interface for the application.
type UI interface {
	Run(<-chan []byte, chan<- struct{})
	Print([]byte)
	MaskInput()
	UnmaskInput()
}

// Client is the application's main connection to the game server.
type Client interface {
	io.ReadWriter
	Scanner() *bufio.Scanner

	// Telnet utilities.
	Will(byte) error
	Wont(byte) error
	Do(byte) error
	Dont(byte) error
	Subneg(byte, []byte) error
}

// NewEngine creates a new Engine based on the given parts.
func NewEngine(world World, ui UI, client Client) *Engine {
	return &Engine{
		ui:     ui,
		client: client,
		world:  world,
	}
}

// Run is the main loop of the application.
func (engine *Engine) Run(inputs <-chan []byte, commands <-chan []byte) error {
	uiOutput := make(chan []byte)
	uiDone := make(chan struct{})

	go engine.ui.Run(uiOutput, uiDone)

	clientOutput := make(chan []byte)
	clientErr := make(chan error)
	clientDone := make(chan struct{})

	go func() {
		scanner := engine.client.Scanner()

		for scanner.Scan() {
			clientOutput <- scanner.Bytes()
		}

		if err := scanner.Err(); err != nil {
			clientErr <- err
		}

		clientDone <- struct{}{}
	}()

	// @todo Implement proper logging.

	for {
		select {
		case _ = <-uiDone:
			return nil

		case err := <-clientErr:
			return err

		case _ = <-clientDone:
			uiOutput <- []byte("server disconnected")

		case input := <-inputs:
			input = engine.world.Input(input)
			if len(input) > 0 {
				_, err := engine.client.Write(append(input, '\r', '\n'))
				if err != nil {
					return fmt.Errorf(
						"failed sending input: %w", err,
					)
				}
			}

		case output := <-clientOutput:
			output = engine.world.Output(output)
			if len(output) > 0 {
				uiOutput <- output
			}

		case command, ok := <-commands:
			if !ok {
				continue
			}

			err := engine.world.Command(command)
			if err != nil {
				return fmt.Errorf(
					"failed processing command '%s': %w",
					command, err,
				)
			}
		}
	}
}
