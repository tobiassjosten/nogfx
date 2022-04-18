package tui

import (
	"context"
	"fmt"
	"math"

	"github.com/gdamore/tcell/v2"
)

// TUI orchestrates different panes to make up the primary user interface.
type TUI struct {
	screen tcell.Screen

	inputs chan []byte

	width  int
	height int

	input *InputPane

	output *OutputPane
}

// NewTUI creates a new TUI.
func NewTUI(screen tcell.Screen, input *InputPane, output *OutputPane) *TUI {
	var (
		outputStyle = tcell.Style{}
	)

	tui := &TUI{
		screen: screen,
		inputs: make(chan []byte),
		input:  input,
		output: output,
	}

	screen.SetStyle(outputStyle)
	screen.SetCursorStyle(tcell.CursorStyleBlinkingBlock)

	return tui
}

// Inputs exposes the outgoing channel for player input.
func (tui *TUI) Inputs() <-chan []byte {
	return tui.inputs
}

// Outputs exposes the incoming channel for server output.
func (tui *TUI) Outputs() chan<- []byte {
	return tui.output.outputs
}

// Run is the main loop of the user interface, where everything is orchestrated.
func (tui *TUI) Run(pctx context.Context) error {
	ctx, cancel := context.WithCancel(pctx)

	if err := tui.screen.Init(); err != nil {
		cancel()
		return fmt.Errorf("failed initializing screen: %w", err)
	}

	go func() {
		for {
			event := tui.screen.PollEvent()
			if event == nil {
				return
			}

			switch ev := event.(type) {
			case *tcell.EventResize:
				tui.Draw()
				tui.screen.Sync()

			case *tcell.EventKey:
				if ev.Key() == tcell.KeyCtrlD {
					cancel()
					return
				}
			}

			if ev, ok := event.(*tcell.EventKey); ok {
				if ok, input := tui.input.HandleEvent(ev); ok {
					if input != nil {
						tui.inputs <- []byte(string(input))
					}
					tui.Draw()
				}
			}
		}
	}()

	for {
		select {
		case output := <-tui.output.outputs:
			tui.output.Add(output)
			tui.Draw()

		case <-ctx.Done():
			tui.screen.Fini()
			return nil
		}
	}
}

// Print shows a text to the user.
func (tui *TUI) Print(output []byte) {
	// @todo Apply default style instead of inheriting whatever's current.
	tui.output.Add(output)
	tui.Draw()
}

// MaskInput hides the content of the InputPane.
func (tui *TUI) MaskInput() {
	tui.input.Mask()
}

// UnmaskInput shows the content of the InputPane.
func (tui *TUI) UnmaskInput() {
	tui.input.Unmask()
}

// Resize calculates the layout of the various panes.
func (tui *TUI) Resize(width, height int) {
	tui.width, tui.height = width, height

	inputWidth := width
	inputHeight := int(math.Min(
		float64(height),
		float64(tui.input.Height()),
	))
	inputX, inputY := 0, height-inputHeight
	tui.input.Position(inputX, inputY, inputWidth, inputHeight)

	tui.output.Position(0, 0, width, height-inputHeight)
}

// Draw updates the terminal and prints the contents of the panes.
func (tui *TUI) Draw() {
	tui.screen.Clear()

	tui.Resize(tui.screen.Size())

	tui.input.Draw(tui.screen)
	tui.output.Draw(tui.screen)

	tui.screen.Show()
}
