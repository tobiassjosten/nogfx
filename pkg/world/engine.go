package world

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/gmcp"
	"github.com/tobiassjosten/nogfx/pkg/telnet"
	"github.com/tobiassjosten/nogfx/pkg/world/achaea"
)

var worlds = map[string]func(pkg.Client, pkg.UI) pkg.World{
	"achaea.com:23":  achaea.NewWorld,
	"50.31.100.8:23": achaea.NewWorld,
}

// Engine is the orchestrator of all the cogs of this machinery.
type Engine struct {
	client  pkg.Client
	ui      pkg.UI
	world   pkg.World
	address string
}

// NewEngine creates a new Engine.
func NewEngine(client pkg.Client, ui pkg.UI, address string) *Engine {
	engine := &Engine{
		client:  client,
		ui:      ui,
		address: address,
	}

	if constructor, ok := worlds[address]; ok {
		engine.world = constructor(client, ui)
	}

	return engine
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

	gamelog := engine.openGamelog(ctx)
	if gamelog != nil {
		defer gamelog.Close()
	}

	out := pkg.Exput{}

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
			in := (pkg.Exput{}).Add(data)
			inout := in.Inoutput(pkg.Input)

			if engine.world != nil {
				inout = engine.world.OnInoutput(inout)
			}

			engine.OnInoutput(inout)

		case data := <-serverOutput:
			if gamelog != nil {
				if _, err := gamelog.Write(data); err != nil {
					log.Printf("failed writing game log: %s", err)
				}
			}

			data = bytes.TrimRight(data, "\r\n")

			// Consider it a full capture and proceed only after a
			// Go Ahead termination.
			// @todo Make this dynamic, based on Telnet negotiation.
			if len(data) == 0 || data[len(data)-1] != telnet.GA {
				out = out.Add(data)
				continue
			}

			// Strip the GA and proceed with processing.
			out = out.Add(data[:len(data)-1])
			inout := out.Inoutput(pkg.Output)

			if engine.world != nil {
				inout = engine.world.OnInoutput(inout)
			}

			engine.OnInoutput(inout)

			out = pkg.Exput{}

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

			inout := engine.world.OnCommand(command)

			engine.OnInoutput(inout)
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

// OnInoutput dispatches input and output to the client and UI respectively.
func (engine *Engine) OnInoutput(inout pkg.Inoutput) {
	for _, data := range inout.Input.Bytes() {
		if _, err := engine.client.Write(data); err != nil {
			log.Printf("failed sending command: %s", err)
		}
	}

	for _, data := range inout.Output.Bytes() {
		engine.ui.Outputs() <- data
	}
}

// SendGMCP writes a GMCP message to the client.
func (engine *Engine) SendGMCP(msg gmcp.Message) error {
	data := []byte(msg.Marshal())
	if _, err := engine.client.Write(gmcp.Wrap(data)); err != nil {
		return err
	}

	return nil
}

func (engine *Engine) openGamelog(ctx context.Context) *os.File {
	ctxLogdir := ctx.Value(pkg.CtxLogdir)
	logdir, ok := ctxLogdir.(string)
	if !ok || logdir == "" {
		log.Printf("missing logdir context: '%s'", logdir)
		return nil
	}

	game := strings.Split(engine.address, ":")[0]

	start := time.Now().Format("20060102-150405")
	path := fmt.Sprintf("%s/%s-%s.log", logdir, game, start)

	gamelog, err := os.Create(path)
	if err != nil {
		log.Printf("failed creating gamelog file: %s", err)
		return nil
	}

	return gamelog
}
