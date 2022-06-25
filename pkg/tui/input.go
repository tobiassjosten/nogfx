package tui

import (
	"github.com/gdamore/tcell/v2"
)

const nbsp = '\u00A0' // Non-breaking space.

// MaskInput hides the content of the InputPane.
func (tui *TUI) MaskInput() {
	tui.input.buffer = []rune{}
	tui.input.cursor = 0
	tui.input.masked = true
	tui.setCache(paneInput, nil)
}

// UnmaskInput shows the content of the InputPane.
func (tui *TUI) UnmaskInput() {
	tui.input.buffer = []rune{}
	tui.input.cursor = 0
	tui.input.masked = false
	tui.setCache(paneInput, nil)
}

// Input is the widget where the player types what's sent to the game.
type Input struct {
	buffer   []rune
	inputted bool
	masked   bool
	cursor   int
}

// RenderInput renders the current Input.
func (tui *TUI) RenderInput(width, height int) Rows {
	if rows, ok := tui.getCache(paneInput); ok {
		return rows
	}

	rows := RenderInput(tui.input, width, height)

	tui.setCache(paneInput, rows)

	return rows
}

// RenderInput renders the given Input.
func RenderInput(input *Input, width, height int) Rows {
	if width == 0 {
		return nil
	}

	style := (tcell.Style{}).
		Foreground(tcell.ColorWhite).
		Background(tcell.Color235)

	if input.inputted {
		style = style.Attributes(tcell.AttrDim)
	}

	buffer := make([]rune, len(input.buffer))
	copy(buffer, input.buffer)

	if input.masked {
		for i := range buffer {
			buffer[i] = '*'
		}
	}

	row := NewRowFromRunes(buffer, style)

	// Pad with non-breaking spaces to simplify cursor positioning later.
	rows := row.Wrap(width, NewCell(nbsp, style))

	return rows
}

func cursorPosition(rows Rows, offset, x, y int) []int {
	cursor := []int{x, y - 1} // y-1 simplifies the algorithm below.

	lastrow := rows[len(rows)-1]
	if lastrow[len(lastrow)-1].Content != nbsp {
		return nil
	}

outer:
	for _, row := range rows {
		cursor[0] = 0
		cursor[1]++

		for _, cell := range row {
			if cell.Content == nbsp {
				break
			}

			offset--
			cursor[0]++

			if offset == 0 {
				break outer
			}
		}
	}

	return cursor
}
