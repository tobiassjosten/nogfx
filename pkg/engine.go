package pkg

import (
	"bufio"
	"context"
	"fmt"
	"io"
)

// UI is the primary user interface for the application.
type UI interface {
	Inputs() <-chan []byte
	Outputs() chan<- []byte
	Run(context.Context) error
	Print([]byte)
	MaskInput()
	UnmaskInput()
}

// Client is the application's main connection to the game server.
type Client interface {
	io.ReadWriter
	Commands() <-chan []byte
	Scanner() *bufio.Scanner

	// Telnet utilities.
	Will(byte) error
	Wont(byte) error
	Do(byte) error
	Dont(byte) error
	Subneg(byte, []byte) error
}

// Run is the main loop of the application, where everything is orchestrated.
func Run(pctx context.Context, client Client, ui UI, world World) error {
	ctx, cancel := context.WithCancel(pctx)

	go func() {
		// @todo Feed potential errors back so they can be returned.
		_ = ui.Run(ctx)
		cancel()
	}()

	clientOutput := make(chan []byte)
	clientErr := make(chan error)
	clientDone := make(chan struct{})

	go func() {
		scanner := client.Scanner()

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
		case _ = <-ctx.Done():
			return nil

		case err := <-clientErr:
			return err

		case _ = <-clientDone:
			ui.Outputs() <- []byte("server disconnected")

		case input := <-ui.Inputs():
			input = world.Input(input)
			if input == nil {
				continue
			}

			if _, err := client.Write(input); err != nil {
				return fmt.Errorf(
					"failed sending input: %w", err,
				)
			}

		case output := <-clientOutput:
			output = world.Output(output)
			if output == nil {
				continue
			}

			ui.Outputs() <- output

		case command, ok := <-client.Commands():
			if !ok {
				continue
			}

			if err := world.Command(command); err != nil {
				return fmt.Errorf(
					"failed processing command '%s': %w",
					command, err,
				)
			}
		}
	}

}
