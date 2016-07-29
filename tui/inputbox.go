package tui

import (
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

type InputBox struct {
	screen *Screen
	buffer string
	kept   bool
}

func (inputBox *InputBox) Add(ch rune) {
	if inputBox.kept {
		inputBox.Clear()
	}
	inputBox.buffer += string(ch)
}

func (inputBox *InputBox) Remove() {
	if inputBox.kept {
		inputBox.Clear()
	} else if 0 < len(inputBox.buffer) {
		inputBox.buffer = inputBox.buffer[0 : len(inputBox.buffer)-1]
	}
}

func (inputBox *InputBox) Clear() {
	inputBox.buffer = ""
	inputBox.kept = false
}

func (inputBox *InputBox) Get() string {
	inputBox.kept = true
	return inputBox.buffer
}

func (inputBox *InputBox) Draw(x int, y int, X int, Y int) {
	fg := termbox.ColorWhite
	if inputBox.kept {
		fg = termbox.ColorCyan
	}

	for _, c := range []rune(inputBox.buffer) {
		termbox.SetCell(x, y, c, fg, termbox.ColorDefault)
		x += runewidth.RuneWidth(c)
	}

	termbox.SetCursor(x, y)

	for ; x < X; x++ {
		termbox.SetCell(x, y, '_', termbox.ColorWhite, termbox.ColorDefault)
	}
}
