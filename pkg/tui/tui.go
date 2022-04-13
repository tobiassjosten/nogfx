package tui

import (
	"log"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
)

type TUI struct {
	screen tcell.Screen

	inputting bool
	input     []rune
	inputs    chan []byte
	inputMask bool

	texts []Text
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
				case tcell.KeyCtrlC:
					quit()

				case tcell.KeyESC:
					tui.inputting = false
					tui.screen.HideCursor()
					tui.draw()

				case tcell.KeyBackspace, tcell.KeyBackspace2:
					if !tui.inputting {
						continue
					}

					if len(tui.input) > 0 {
						tui.input = tui.input[:len(tui.input)-1]
						tui.draw()
					}

				case tcell.KeyETB: // opt/elt+backspace
					if !tui.inputting {
						continue
					}

					deleted := false
					for i := len(tui.input) - 1; i >= 0; i-- {
						if tui.input[i] == ' ' {
							tui.input = tui.input[0:i]
							deleted = true
							break
						}
					}

					if !deleted {
						tui.input = []rune{}
					}

					tui.draw()

				case tcell.KeyNAK: // cmd/ctrl+backspace
					if !tui.inputting {
						continue
					}

					if len(tui.input) > 0 {
						tui.input = []rune{}
						tui.draw()
					}

				case tcell.KeyEnter:
					if !tui.inputting {
						continue
					}

					inputs <- []byte(string(tui.input))
					tui.input = []rune{}
					tui.draw()

				case tcell.KeyRune:
					if !tui.inputting {
						if ev.Rune() == ' ' {
							tui.inputting = true
							tui.draw()
						}

						if ev.Rune() == '1' {
							inputs <- []byte{'s', 'w'}
						}
						if ev.Rune() == '2' {
							inputs <- []byte{'s'}
						}
						if ev.Rune() == '3' {
							inputs <- []byte{'s', 'e'}
						}
						if ev.Rune() == '4' {
							inputs <- []byte{'w'}
						}
						if ev.Rune() == '5' {
							inputs <- []byte{'m', 'a', 'p'}
						}
						if ev.Rune() == '6' {
							inputs <- []byte{'e'}
						}
						if ev.Rune() == '7' {
							inputs <- []byte{'n', 'w'}
						}
						if ev.Rune() == '8' {
							inputs <- []byte{'n'}
						}
						if ev.Rune() == '9' {
							inputs <- []byte{'n', 'e'}
						}

						continue
					}

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

func (tui *TUI) MaskInput() {
	tui.inputMask = true
}

func (tui *TUI) UnmaskInput() {
	tui.inputMask = false
}

func (tui *TUI) newOutput(output []byte, style tcell.Style) tcell.Style {
	text, style := NewText(output, style)
	tui.texts = append([]Text{text}, tui.texts...)

	return style
}

func (tui *TUI) draw() {
	tui.screen.Clear()

	width, height := tui.screen.Size()

	inputHeight := 0
	if tui.inputting {
		inputHeight = 1
		tui.drawInput(0, height-1, width, inputHeight)
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
		b.draw(tui.screen, x, line)

		line--
		if line < y {
			break
		}
	}
}

func (tui *TUI) drawInput(x, y, width, height int) {
	style := tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorGray)

	// @todo Fixa stöd för flera rader.

	// @todo Fixa stöd för att hoppa med cursorn.

	// @todo Behåll texten för att lätt kunna repetera.

	input := append(tui.input, []rune(strings.Repeat(" ", width-len(tui.input)))...)

	for i, r := range input {
		if tui.inputMask && i < len(tui.input) {
			r = '*'
		}
		tui.screen.SetContent(x+i, y, r, nil, style)
	}

	tui.screen.ShowCursor(x+len(tui.input), y)
}
