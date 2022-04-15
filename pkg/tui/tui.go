package tui

import (
	"github.com/gdamore/tcell/v2"
)

type TUI struct {
	screen tcell.Screen

	input     Text
	inputting bool
	inputted  bool
	inputs    chan []byte
	inputMask bool

	outputs chan []byte

	style tcell.Style

	texts []Text
}

func NewTUI() (*TUI, <-chan []byte, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, nil, err
	}
	screen.SetStyle(outputStyle)

	inputs := make(chan []byte)

	tui := &TUI{
		screen: screen,
		inputs: inputs,
	}

	if err := tui.screen.Init(); err != nil {
		return nil, nil, err
	}

	return tui, inputs, nil
}

func (tui *TUI) Run(outputs <-chan []byte, done chan<- struct{}) {
	tui.draw()

	quit := func() {
		tui.screen.Fini()
		done <- struct{}{}
	}

	inputs := make(chan []byte)
	go func() {
		for {
			switch ev := tui.screen.PollEvent().(type) {
			case *tcell.EventResize:
				tui.drawSync()

			case *tcell.EventKey:
				if ev.Key() == tcell.KeyCtrlC {
					quit()
				}

				input := tui.InteractKey(ev)
				if len(input) > 0 {
					inputs <- input
				}
			}
		}
	}()

	for {
		select {
		case input := <-inputs:
			tui.inputs <- input

		case output := <-outputs:
			text, style := NewText(output, tui.style)

			// @todo Cap tui.texts so it doesn't grow indefinitely.

			tui.texts = append([]Text{text}, tui.texts...)
			tui.style = style

			tui.draw()
		}
	}
}

func (tui *TUI) Print(output []byte) {
	text, _ := NewText(output, tcell.Style{})
	tui.texts = append([]Text{text}, tui.texts...)
	tui.draw()
}

func (tui *TUI) MaskInput() {
	tui.inputMask = true
}

func (tui *TUI) UnmaskInput() {
	tui.inputMask = false
}

func (tui *TUI) draw() {
	tui.screen.Clear()

	// @todo Workaround for https://github.com/gdamore/tcell/issues/522.
	tui.screen.Sync()

	width, height := tui.screen.Size()

	inputHeight := 0
	if tui.inputting {
		b := NewBlock(tui.input, width)
		b.draw(tui.screen, DrawOptions{
			X:       0,
			Y:       height - b.Height(),
			Masked:  tui.inputMask,
			Filling: NewCell(' ', inputStyle),
		})

		tui.screen.ShowCursor(len(tui.input), height-1)

		inputHeight = b.Height()
	}

	tui.drawOutput(0, 0, width, height-inputHeight)

	tui.screen.Show()
}

func (tui *TUI) drawSync() {
	tui.draw()
	tui.screen.Sync()
}

func (tui *TUI) drawOutput(x, y, width, height int) {
	line := y + height - 1

	// @todo Fixa stöd för att kunna scrolla upp.

	for _, t := range tui.texts {
		b := NewBlock(t, width)
		line = line - b.Height() + 1
		b.draw(tui.screen, DrawOptions{X: x, Y: line})

		line--
		if line < y {
			break
		}
	}
}
