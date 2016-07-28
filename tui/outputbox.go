package tui

import (
	"github.com/nsf/termbox-go"
	"math"
)

type OutputBox struct {
	screen *Screen
	buffer []string
}

func (outputBox *OutputBox) Add(line string) {
	outputBox.buffer = append(outputBox.buffer, line)
}

func (outputBox *OutputBox) Draw(x int, y int, X int, Y int) {
	offset := int(math.Dim(float64(len(outputBox.buffer)), float64(Y-1)))
	for i := 0; i <= Y && i+offset < len(outputBox.buffer); i++ {
		for ii, c := range outputBox.buffer[i+offset] {
			termbox.SetCell(ii, i, c, termbox.ColorWhite, termbox.ColorDefault)
		}
	}
}
