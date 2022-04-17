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
func NewTUI(screen tcell.Screen) *TUI {
	var (
		outputStyle = tcell.Style{}
		inputStyle  = (tcell.Style{}).
				Foreground(tcell.ColorWhite).
				Background(tcell.ColorGray)
		inputtedStyle = (tcell.Style{}).
				Foreground(tcell.ColorWhite).
				Background(tcell.ColorGray).
				Attributes(tcell.AttrDim)
	)

	tui := &TUI{
		screen: screen,
		inputs: make(chan []byte),
		input:  NewInputPane(inputStyle, inputtedStyle),
	}
	tui.output = NewOutputPane(tui, outputStyle)

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

	// tui.Resize(tui.screen.Size())

	go func() {
		for {
			event := tui.screen.PollEvent()
			if event == nil {
				return
			}

			switch ev := event.(type) {
			case *tcell.EventResize:
				// tui.Resize(tui.screen.Size())
				tui.Draw()
				tui.screen.Sync()

			case *tcell.EventKey:
				if ev.Key() == tcell.KeyCtrlD {
					cancel()
					return
				}
			}

			if ok, input := tui.input.HandleEvent(event); ok {
				if input != nil {
					tui.inputs <- []byte(string(input))
				}
				tui.Draw()
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

	tui.input.width = width
	tui.input.height = int(math.Min(
		float64(height),
		float64(tui.input.Height()),
	))
	tui.input.x, tui.input.y = 0, height-tui.input.height

	tui.output.x, tui.output.y = 0, 0
	tui.output.width, tui.output.height = width, height-tui.input.height

	// tui.Draw()
}

// Draw updates the terminal and prints the contents of the panes.
func (tui *TUI) Draw() {
	tui.screen.Clear()

	tui.Resize(tui.screen.Size())

	tui.input.Draw(tui.screen)
	tui.output.Draw(tui.screen)

	tui.screen.Show()
}
