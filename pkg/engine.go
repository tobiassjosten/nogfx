package pkg

import (
	"bufio"
	"bytes"
	"context"
	"log"
	"net"

	"github.com/tobiassjosten/nogfx/pkg/gmcp"
	"github.com/tobiassjosten/nogfx/pkg/process"
	"github.com/tobiassjosten/nogfx/pkg/telnet"
)

// Conn extends net.Conn with a built-in bufio.SplitFunc.
type Conn interface {
	net.Conn
	SplitFunc([]byte, bool) (int, []byte, error)
}

// Engine is the orchestrator of all the cogs of this machinery.
type Engine struct {
	conn      Conn
	ui        UI
	Processor process.Processor
}

// Run is the main loop of the application, where everything is orchestrated.
func (engine *Engine) Run(pctx context.Context) error {
	ctx, cancel := context.WithCancel(pctx)

	connOutput := make(chan []byte)
	connErrs := make(chan error)
	connDone := make(chan struct{})
	// @todo Make sure this can gracefully shut down as well.
	go engine.readConn(connOutput, connErrs, connDone)

	uiErrs := make(chan error)
	go engine.runUI(ctx, uiErrs, cancel)

	var outs [][]byte

	for {
		select {
		case <-ctx.Done():
			return nil

		case err := <-connErrs:
			return err

		case err := <-uiErrs:
			return err

		case <-connDone:
			engine.ui.Outputs() <- []byte("server disconnected")

		case in := <-engine.ui.Inputs():
			engine.process([][]byte{in})

		case out := <-connOutput:
			// @todo Move this processing out of here. We tell the
			// scanner to look for either \r\n or GA and then break
			// up what we get into [][]byte by splitting on \r\n.
			// Will be fixed with our refactoring of telnet!

			out = bytes.TrimRight(out, "\r\n")

			// Add output to a buffer until we see a GA. Then strip
			// that and procees with processing.
			// @todo Make this dynamic, based on Telnet negotiation.
			if len(out) == 0 || out[len(out)-1] != telnet.GA {
				out = out.Append(out)
				continue
			}
			out = out.Append(out[:len(out)-1])

			inout := Inoutput{Output: out}
			engine.process(inout)

			out = Exput{}

		case data, ok := <-engine.Client.Commands():
			// @todo Figure out if we really need this here. Could
			// it be solved by adding a default case?
			if !ok {
				continue
			}

			inout := Inoutput{Input: NewExput(data)}
			engine.process(inout)
		}
	}
}

// readConn reads data from the conn and reports back output and potential
// errors to the given channels, before marking its completion.
func (engine *Engine) readConn(outputs chan []byte, errs chan error, done chan struct{}) {
	scanner := bufio.NewScanner(engine.conn)
	scanner.Split(engine.conn.SplitFunc)

	for scanner.Scan() {
		// Scanner.Bytes() returns a byte slice, a reference type. We
		// dereference here to allow for later modification.
		var output = make([]byte, len(scanner.Bytes()))
		// @todo Verify that this dereference is still needed.
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
func (engine *Engine) runUI(ctx context.Context, errs chan error, cancel func()) {
	if err := engine.ui.Run(ctx); err != nil {
		errs <- err
	}

	cancel()
}

func (engine *Engine) process(inout Inoutput) {
	for _, cmd := range inout.Commands() {
		switch {
		case bytes.Equal(cmd, telnet.IAC_WILL_ECHO):
			engine.UI.MaskInput()

		case bytes.Equal(cmd, telnet.IAC_WONT_ECHO):
			engine.UI.UnmaskInput()

		case bytes.Equal(cmd, telnet.IAC_WILL_GMCP):
			inout.Input = inout.Input.Append([]byte((&gmcp.CoreHello{
				Client:  "nogfx",
				Version: Version,
			}).Marshal()))
		}
	}

	// @todo Remove Telnet commands from `outs`. We put them there for easy
	// processing but they need to be deleted before printed by the UI.

	if engine.Processor != nil {
		inout = engine.Processor.Process(inout)
	}

	engine.dispatch(inout)
}

func (engine *Engine) dispatch(inout Inoutput) {
	for _, data := range inout.Input.Bytes() {
		if _, err := engine.Client.Write(data); err != nil {
			log.Printf("failed sending command: %s", err)
		}
	}

	for _, data := range inout.Output.Bytes() {
		engine.UI.Outputs() <- data
	}
}
