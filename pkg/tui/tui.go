package tui

import (
	"context"
	"fmt"
	"math"

	"github.com/gdamore/tcell/v2"
)

type TUI struct {
	screen tcell.Screen

	width  int
	height int

	input *InputPane

	output *OutputPane
}

func NewTUI() (*TUI, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, fmt.Errorf("failed creating screen: %w", err)
	}

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

	tui := &TUI{screen: screen}
	tui.input = NewInputPane(tui, inputStyle, inputtedStyle)
	tui.output = NewOutputPane(tui, outputStyle)

	screen.SetStyle(outputStyle)
	screen.SetCursorStyle(tcell.CursorStyleBlinkingBlock)

	if err := screen.Init(); err != nil {
		return nil, fmt.Errorf("failed initializing screen: %w", err)
	}

	return tui, nil
}

func (tui *TUI) Inputs() <-chan []byte {
	return tui.input.inputs
}

func (tui *TUI) Outputs() chan<- []byte {
	return tui.output.outputs
}

func (tui *TUI) Run(pctx context.Context) {
	ctx, cancel := context.WithCancel(pctx)

	tui.Resize(tui.screen.Size())

	go func() {
		for {
			event := tui.screen.PollEvent()
			if event == nil {
				return
			}

			switch ev := event.(type) {
			case *tcell.EventResize:
				tui.Resize(tui.screen.Size())
				tui.screen.Sync()

			case *tcell.EventKey:
				if ev.Key() == tcell.KeyCtrlD {
					cancel()
					return
				}
			}

			if ok := tui.input.HandleEvent(event); ok {
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
			return
		}
	}
}

func (tui *TUI) Print(output []byte) {
	// @todo Apply default style instead of inheriting whatever's current.
	tui.output.Add(output)
	tui.Draw()
}

func (tui *TUI) MaskInput() {
	tui.input.masked = true
}

func (tui *TUI) UnmaskInput() {
	tui.input.masked = false
}

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

	tui.Draw()
}

func (tui *TUI) Draw() {
	tui.screen.Clear()

	tui.input.Draw(tui.screen)
	tui.output.Draw(tui.screen)

	tui.screen.Show()
}
