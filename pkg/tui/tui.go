package tui

import (
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
)

type TUI struct {
	screen tcell.Screen
	input  []rune
	inputs chan []byte
	texts  []Text
}

func NewTUI() (*TUI, <-chan []byte, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, nil, err
	}

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

func (tui *TUI) Run(outputs <-chan []byte) {
	style := tcell.StyleDefault.
		Background(tcell.ColorReset).
		Foreground(tcell.ColorReset)
	tui.screen.SetStyle(style)
	// 16777216 == tui.screen.Colors() // 24 bit

	tui.draw()

	quit := func() {
		tui.screen.Fini()
		os.Exit(0) // @todo Move this to main loop.
	}

	inputs := make(chan []byte)
	go func() {
		for {
			switch ev := tui.screen.PollEvent().(type) {
			case *tcell.EventResize:
				tui.drawSync()

			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyESC, tcell.KeyCtrlC:
					quit()

				case tcell.KeyBackspace, tcell.KeyBackspace2:
					if len(tui.input) > 0 {
						tui.input = tui.input[:len(tui.input)-1]
						tui.draw()
					}

				case tcell.KeyETB: // opt/elt+backspace
					deleted := false
					for i := len(tui.input) - 1; i >= 0; i-- {
						if tui.input[i] == ' ' {
							tui.input = tui.input[0:i]
							deleted = true
						}
					}

					if !deleted {
						tui.input = []rune{}
					}

					tui.draw()

				case tcell.KeyNAK: // cmd/ctrl+backspace
					if len(tui.input) > 0 {
						tui.input = []rune{}
						tui.draw()
					}

				case tcell.KeyEnter:
					inputs <- []byte(string(tui.input))
					tui.input = []rune{}
					tui.draw()

				case tcell.KeyRune:
					tui.input = append(tui.input, ev.Rune())
					tui.draw()

				default:
					log.Println("KEY", ev.Key())
				}
			}
		}
	}()

	for {
		select {
		case input := <-inputs:
			tui.inputs <- input

		case output := <-outputs:
			style = tui.newOutput(output, style)
			tui.draw()
		}
	}
}

func (tui *TUI) newOutput(output []byte, style tcell.Style) tcell.Style {
	text, style := NewText(output, style)
	tui.texts = append([]Text{text}, tui.texts...)

	return style
}

func (tui *TUI) draw() {
	tui.screen.Clear()

	width, height := tui.screen.Size()

	inputHeight := 1
	tui.drawInput(0, height-1, width, inputHeight)

	tui.drawOutput(0, 0, width, height-inputHeight)

	tui.screen.Show()
}

func (tui *TUI) drawSync() {
	tui.draw()
	tui.screen.Sync()
}

func (tui *TUI) drawOutput(x, y, width, height int) {
	line := y + height - 1

	for _, t := range tui.texts {
		b := NewBlock(t, width)
		line = line - b.Height() + 1
		b.draw(tui.screen, x, line)

		line--
		if line < y {
			break
		}
	}
}

func (tui *TUI) drawInput(x, y, width, height int) {
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorGray)

	x2, y2 := tui.screen.Size()
	x1, y1 := 0, y2-1

	for row := y1; row <= y2; row++ {
		for col := x1; col <= x2; col++ {
			r := ' '
			if col < len(tui.input) {
				r = tui.input[col]
			}
			tui.screen.SetContent(col, row, r, nil, style)
		}
	}
}
