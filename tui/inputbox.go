package tui

import (
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

type InputBox struct {
	screen  *Screen
	written string
}

func (inputBox *InputBox) Add(ch rune) {
	inputBox.written += string(ch)
}

func (inputBox *InputBox) Remove() {
	inputBox.written = inputBox.written[0 : len(inputBox.written)-1]
}

func (inputBox *InputBox) Get() string {
	written := inputBox.written
	inputBox.written = ""
	return written
}

func (inputBox *InputBox) Draw(x int, y int, X int, Y int) {
	for _, c := range []rune(inputBox.written) {
		termbox.SetCell(x, y, c, termbox.ColorWhite, termbox.ColorDefault)
		x += runewidth.RuneWidth(c)
	}

	termbox.SetCursor(x, y)

	for ; x < X; x++ {
		termbox.SetCell(x, y, '_', termbox.ColorWhite, termbox.ColorDefault)
	}
}
