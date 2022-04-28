package tui

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"
)

// Panes is a collection of various panes used throughout the user interface.
type Panes struct {
	Input  *InputPane
	Output *OutputPane
	Vitals *VitalsPane
}

// NewPanes creates a new Panes with the standard collection of panes.
func NewPanes() Panes {
	return Panes{
		Input:  NewInputPane(),
		Output: NewOutputPane(),
		Vitals: NewVitalsPane(),
	}
}

// TUI orchestrates different panes to make up the primary user interface.
type TUI struct {
	screen  tcell.Screen
	panes   Panes
	inputs  chan []byte
	running bool
}

// NewTUI creates a new TUI.
func NewTUI(screen tcell.Screen, panes Panes) *TUI {
	var (
		outputStyle = tcell.Style{}
	)

	tui := &TUI{
		screen: screen,
		inputs: make(chan []byte),
		panes:  panes,
	}

	screen.SetStyle(outputStyle)
	screen.SetCursorStyle(tcell.CursorStyleBlinkingBlock)

	return tui
}

// Inputs exposes the outgoing channel for player input.
func (tui *TUI) Inputs() <-chan []byte {
	return tui.inputs
}

// Run is the main loop of the user interface, where everything is orchestrated.
func (tui *TUI) Run(pctx context.Context) error {
	ctx, cancel := context.WithCancel(pctx)

	tui.running = true
	defer func() { tui.running = false }()

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
				if ok, input := tui.panes.Input.HandleEvent(ev); ok {
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
		case output := <-tui.panes.Output.outputs:
			tui.panes.Output.Add(output)
			tui.Draw()

		case <-ctx.Done():
			tui.screen.Fini()
			return nil
		}
	}
}

// Resize calculates the layout of the various panes.
func (tui *TUI) Resize(width, height int) {
	outputWidth := int(min(120, width))

	inputWidth := outputWidth
	inputHeight := int(min(height, tui.panes.Input.Height()))
	inputX, inputY := 0, height-inputHeight
	tui.panes.Input.Position(inputX, inputY, inputWidth, inputHeight)

	// Draw VitalsPane if OutputPane has at least one row.
	vitalsHeight := min(tui.panes.Vitals.Height(), max(0, height-inputHeight-1))
	tui.panes.Vitals.Position(inputX, inputY-1, outputWidth, vitalsHeight)

	tui.panes.Output.Position(0, 0, outputWidth, height-inputHeight-vitalsHeight)
}

// Draw updates the terminal and prints the contents of the panes.
func (tui *TUI) Draw() {
	if !tui.running {
		return
	}

	tui.screen.Clear()

	tui.Resize(tui.screen.Size())

	tui.panes.Input.Draw(tui.screen)
	tui.panes.Vitals.Draw(tui.screen)
	tui.panes.Output.Draw(tui.screen)

	tui.screen.Show()
}
