package world

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/gmcp"
	"github.com/tobiassjosten/nogfx/pkg/telnet"
	"github.com/tobiassjosten/nogfx/pkg/world/achaea"
)

var worlds = map[string]func(pkg.Client, pkg.UI) pkg.World{
	// Dummy world just for testing.
	"example.com:23": NewExampleWorld,

	"achaea.com:23":  achaea.NewWorld,
	"50.31.100.8:23": achaea.NewWorld,

	// @todo Extend this when we support more games. For now, we list these
	// two so as to force more general, shared functionality.
}

// Engine is the orchestrator of all the cogs of this machinery.
type Engine struct {
	client pkg.Client
	ui     pkg.UI
	world  pkg.World
}

// NewEngine creates a new Engine.
func NewEngine(client pkg.Client, ui pkg.UI, address string) *Engine {
	newWorld := NewGenericWorld
	if constructor, ok := worlds[address]; ok {
		newWorld = constructor
	}

	return &Engine{
		client: client,
		ui:     ui,
		world:  newWorld(client, ui),
	}
}

// Run is the main loop of the application, where everything is orchestrated.
func (engine *Engine) Run(pctx context.Context) error {
	ctx, cancel := context.WithCancel(pctx)

	serverOutput := make(chan []byte)
	serverErrs := make(chan error)
	serverDone := make(chan struct{})
	go engine.RunClient(serverOutput, serverErrs, serverDone)

	uiErrs := make(chan error)
	go engine.RunUI(ctx, uiErrs, cancel)

	output := pkg.Output{}

	for {
		select {
		case <-ctx.Done():
			return nil

		case err := <-serverErrs:
			return err

		case err := <-uiErrs:
			return err

		case <-serverDone:
			engine.ui.Outputs() <- []byte("server disconnected")

		case data := <-engine.ui.Inputs():
			input := (pkg.Input{}).Add(data)

			input = engine.world.ProcessInput(input)

			for _, command := range input {
				if _, err := engine.client.Write(command.Text); err != nil {
					return fmt.Errorf("failed sending: %w", err)
				}
			}

		case data := <-serverOutput:
			data = bytes.TrimRight(data, "\r\n")
			output = output.Add(data)

			// Consider it a full capture and proceed only after a
			// Go Ahead termination.
			// @todo Make this dynamic, based on Telnet negotiation.
			if len(data) == 0 || data[len(data)-1] != telnet.GA {
				continue
			}

			output = engine.world.ProcessOutput(output)
			for _, line := range output.Lines() {
				engine.ui.Outputs() <- line.Raw
			}

			output = pkg.Output{}

		case command, ok := <-engine.client.Commands():
			if !ok {
				continue
			}

			err := engine.ProcessCommand(command)
			if err != nil {
				log.Printf(
					"Failed processing command '%s': %s",
					command, err.Error(),
				)
			}

			engine.world.ProcessCommand(command)
		}
	}
}

// RunClient reads data from the client and reports back output and potential
// errors to the given channels, before marking its completion.
func (engine *Engine) RunClient(outputs chan []byte, errs chan error, done chan struct{}) {
	scanner := engine.client.Scanner()

	for scanner.Scan() {
		// Scanner.Bytes() returns a byte slice, a reference type. We
		// dereference here to allow for later modification.
		var output = make([]byte, len(scanner.Bytes()))
		copy(output, scanner.Bytes())
		outputs <- output
	}

	if err := scanner.Err(); err != nil {
		errs <- err
	}

	done <- struct{}{}
}

// RunUI starts the main loop of the user interface and reports back potential
// errors, before marking its completion.
func (engine *Engine) RunUI(ctx context.Context, errs chan error, cancel func()) {
	if err := engine.ui.Run(ctx); err != nil {
		errs <- err
	}

	cancel()
}

// ProcessCommand processes telnet commands.
func (engine *Engine) ProcessCommand(command []byte) error {
	switch {
	case bytes.Equal(command, []byte{telnet.IAC, telnet.WILL, telnet.ECHO}):
		engine.ui.MaskInput()

	case bytes.Equal(command, []byte{telnet.IAC, telnet.WONT, telnet.ECHO}):
		engine.ui.UnmaskInput()

	case bytes.Equal(command, []byte{telnet.IAC, telnet.WILL, telnet.GMCP}):
		err := engine.SendGMCP(&gmcp.CoreHello{
			Client:  "nogfx",
			Version: pkg.Version,
		})
		if err != nil {
			return fmt.Errorf("failed GMCP: %w", err)
		}
	}

	return nil
}

// SendGMCP writes a GMCP message to the client.
func (engine *Engine) SendGMCP(msg gmcp.Message) error {
	data := []byte(msg.Marshal())
	if _, err := engine.client.Write(gmcp.Wrap(data)); err != nil {
		return err
	}

	return nil
}
