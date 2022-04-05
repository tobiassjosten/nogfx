package tui

import (
	"os"

	"github.com/gdamore/tcell/v2"
)

type TUI struct {
	screen tcell.Screen
	input  []rune
	inputs chan []byte
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

	// 16777216 == tui.screen.Colors() // 24 bit

	return tui, inputs, nil
}

func (tui *TUI) Run(outputs <-chan []byte) {
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

	tui.screen.SetStyle(defStyle)
	tui.screen.Clear()

	tui.drawInputBox()

	quit := func() {
		tui.screen.Fini()
		os.Exit(0) // @todo Move this to main loop.
	}

	inputs := make(chan []byte)
	go func() {
		for {
			tui.screen.Show()

			switch ev := tui.screen.PollEvent().(type) {
			case *tcell.EventResize:
				tui.screen.Sync()

			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyESC, tcell.KeyCtrlC:
					quit()

				case tcell.KeyCtrlL:
					tui.screen.Sync()

				case tcell.KeyEnter:
					inputs <- []byte(string(tui.input))
					tui.input = []rune{}
					tui.drawInputBox()

				case tcell.KeyRune:
					tui.input = append(tui.input, ev.Rune())
					tui.drawInputBox()
				}
			}
		}
	}()

	for {
		select {
		case input := <-inputs:
			tui.inputs <- input

		case output := <-outputs:
			drawText(tui.screen, 0, 0, 100, 1, defStyle, string(output))
		}
	}
}

func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}

func (tui *TUI) drawInputBox() {
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
