package tui

import (
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
)

type TUI struct {
	screen  tcell.Screen
	input   []rune
	inputs  chan []byte
	outputs [][]rune
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

	tui.drawInput()
	tui.drawOutput()

	quit := func() {
		tui.screen.Fini()
		os.Exit(0) // @todo Move this to main loop.
	}

	inputs := make(chan []byte)
	go func() {
		for {
			switch ev := tui.screen.PollEvent().(type) {
			case *tcell.EventResize:
				tui.screen.Sync()

			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyESC, tcell.KeyCtrlC:
					quit()

				case tcell.KeyEnter:
					inputs <- []byte(string(tui.input))
					tui.input = []rune{}
					tui.drawInput()

				case tcell.KeyRune:
					tui.input = append(tui.input, ev.Rune())
					tui.drawInput()
				}
			}

			tui.screen.Show()
		}
	}()

	for {
		select {
		case input := <-inputs:
			tui.inputs <- input

		case output := <-outputs:
			// Reverse the list, with the most recent on top, so
			// that we don't have to do that at every draw.
			tui.outputs = append([][]rune{[]rune(string(output))}, tui.outputs...)
			tui.drawOutput()
		}

		tui.screen.Show()
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

func (tui *TUI) drawOutput() {
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorReset)

	width, height := tui.screen.Size()
	y := height - 2

	log.Println("-- redraw --")

	for _, output := range tui.outputs {
		lines := countLines(output, width)

		x := 0
		y = y - lines + 1

		log.Println(string(output))
		log.Println("\t", lines, "lines, ", y, "y")

		for _, r := range output {
			// @todo Kolla istället efter "printable rune". Har
			// kanske tcell det konceptet?
			if r != '\n' && r != '\r' {
				tui.screen.SetContent(x, y, r, nil, style)
			}

			x++
			if r == '\n' || r == '\r' || x == width {
				log.Println("\t", width-x, "newline fills")
				for ; x < width; x++ {
					tui.screen.SetContent(x, y, ' ', nil, style)
				}

				x = 0
				y++
			}
		}

		log.Println("\t", width-x, "done fills")
		for ; x < width; x++ {
			tui.screen.SetContent(x, y, ' ', nil, style)
		}

		y = y - lines
	}
}

func (tui *TUI) drawInput() {
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

func countLines(text []rune, width int) int {
	lines := 1

	col := 0
	for _, r := range text {
		if r == '\n' || r == '\r' {
			col = 0
			lines++
		} else if col == width {
			col = 1
			lines++
		} else {
			col++
		}
	}

	return lines
}
