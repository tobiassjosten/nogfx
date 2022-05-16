package tui

import (
	"github.com/gdamore/tcell/v2"
)

// MaskInput hides the content of the InputPane.
func (tui *TUI) MaskInput() {
	tui.input.buffer = []rune{}
	tui.input.cursor = 0
	tui.input.masked = true
}

// UnmaskInput shows the content of the InputPane.
func (tui *TUI) UnmaskInput() {
	tui.input.buffer = []rune{}
	tui.input.cursor = 0
	tui.input.masked = false
}

// Input is the widget where the player types what's sent to the game.
type Input struct {
	buffer    []rune
	inputting bool
	inputted  bool
	masked    bool
	cursor    int
}

// RenderInput renders the current Input.
func (tui *TUI) RenderInput(width int) Rows {
	rows := RenderInput(tui.input, width)

	if len(rows) == 0 {
		tui.screen.HideCursor()
	}

	return rows
}

// RenderInput renders the given Input.
func RenderInput(input *Input, width int) Rows {
	if width == 0 {
		return nil
	}

	if !input.inputting {
		return nil
	}

	style := (tcell.Style{}).
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorGray)

	if input.inputted {
		style = (tcell.Style{}).
			Foreground(tcell.ColorWhite).
			Background(tcell.ColorGray).
			Attributes(tcell.AttrDim)
	}

	buffer := make([]rune, len(input.buffer))
	copy(buffer, input.buffer)

	if input.masked {
		for i := range buffer {
			buffer[i] = '*'
		}
	}

	row := NewRowFromRunes(buffer, style)
	rows := row.Wrap(width)

	for y, row := range rows {
		for i := len(row); i < width; i++ {
			rows[y] = append(rows[y], NewCell(' ', style))
		}
	}

	return rows
}
