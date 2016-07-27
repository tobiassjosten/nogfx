package tui

import (
	"github.com/nsf/termbox-go"
)

type OutputBox struct {
	screen *Screen
	buffer []string
}

func (outputBox *OutputBox) SetScreen(screen *Screen) {
	outputBox.screen = screen
}

func (outputBox *OutputBox) Add(line string) {
	outputBox.buffer = append(outputBox.buffer, line)
}

func (outputBox *OutputBox) Draw() {
	for i := len(outputBox.buffer) - 1; i >= 0; i-- {
		for ii, c := range outputBox.buffer[i] {
			termbox.SetCell(ii, i, c, termbox.ColorWhite, termbox.ColorDefault)
		}
	}
}
