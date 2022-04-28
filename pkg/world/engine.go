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
	"github.com/tobiassjosten/nogfx/pkg/world/imperian"
)

var worlds = map[string]func(pkg.Client, pkg.UI) pkg.World{
	// Dummy world just for testing.
	"example.com:23": NewExampleWorld,

	"achaea.com:23":  achaea.NewWorld,
	"50.31.100.8:23": achaea.NewWorld,

	"imperian.com:23":  imperian.NewWorld,
	"67.202.121.44:23": imperian.NewWorld,

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

	clientOutput := make(chan []byte)
	clientErrs := make(chan error)
	clientDone := make(chan struct{})
	go engine.RunClient(clientOutput, clientErrs, clientDone)

	uiErrs := make(chan error)
	go engine.RunUI(ctx, uiErrs, cancel)

	go func() {
	}()

	for {
		select {
		case _ = <-ctx.Done():
			return nil

		case err := <-clientErrs:
			return err

		case err := <-uiErrs:
			return err

		case _ = <-clientDone:
			engine.ui.Outputs() <- []byte("server disconnected")

		case output := <-clientOutput:
			output = engine.world.ProcessOutput(output)
			if output == nil {
				continue
			}

			engine.ui.Outputs() <- output

		case input := <-engine.ui.Inputs():
			input = engine.world.ProcessInput(input)
			if input == nil {
				continue
			}

			if _, err := engine.client.Write(input); err != nil {
				return fmt.Errorf("failed sending: %w", err)
			}

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

			if err := engine.world.ProcessCommand(command); err != nil {
				log.Printf(
					"Failed processing command '%s': %s",
					command, err.Error(),
				)
			}
		}
	}
}

// RunClient reads data from the client and reports back output and potential
// errors to the given channels, before marking its completion.
func (engine *Engine) RunClient(output chan []byte, errs chan error, done chan struct{}) {
	scanner := engine.client.Scanner()

	for scanner.Scan() {
		output <- scanner.Bytes()
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
		err := engine.SendGMCP(gmcp.CoreHello{
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
func (engine *Engine) SendGMCP(message gmcp.ClientMessage) error {
	data := []byte(message.String())
	if _, err := engine.client.Write(gmcp.Wrap(data)); err != nil {
		return err
	}

	return nil
}
