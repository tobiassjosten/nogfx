package tui

import (
	"github.com/gdamore/tcell/v2"
)

const nbsp = '\u00A0' // Non-breaking space.

// MaskInput hides the content of the InputPane.
func (tui *TUI) MaskInput() {
	tui.input.buffer = []rune{}
	tui.input.cursoroff = 0
	tui.input.masked = true
	tui.setCache(paneInput, nil)
}

// UnmaskInput shows the content of the InputPane.
func (tui *TUI) UnmaskInput() {
	tui.input.buffer = []rune{}
	tui.input.cursoroff = 0
	tui.input.masked = false
	tui.setCache(paneInput, nil)
}

// Input is the widget where the player types what's sent to the game.
type Input struct {
	buffer    []rune
	inputted  bool
	masked    bool
	cursoroff int
	cursorpos []int
}

// RenderInput renders the current Input.
func (tui *TUI) RenderInput(width, height int) (Rows, int, int) {
	if rows, ok := tui.getCache(paneInput); ok {
		return rows, tui.input.cursorpos[0], tui.input.cursorpos[1]
	}

	rows, cx, cy := RenderInput(tui.input, width, height)

	tui.setCache(paneInput, rows)
	tui.input.cursorpos = []int{cx, cy}

	return rows, cx, cy
}

// RenderInput renders the given Input.
func RenderInput(input *Input, width, height int) (Rows, int, int) {
	if width == 0 {
		return nil, 0, 0
	}

	style := (tcell.Style{}).
		Foreground(tcell.ColorWhite).
		Background(tcell.Color235)

	if input.inputted {
		style = style.Attributes(tcell.AttrDim)
	}

	padding := NewCell(nbsp, style)

	buffer := make([]rune, len(input.buffer))
	copy(buffer, input.buffer)

	if input.masked {
		for i := range buffer {
			buffer[i] = '*'
		}
	}

	row := NewRowFromRunes(buffer, style)
	rows := row.Wrap(width, padding)

	// Add a new, empty line if the last one is full, to show ahead where
	// new input will show up.
	if last := rows[len(rows)-1]; last[len(last)-1].Content != nbsp {
		rows = append(rows, NewRow(width, padding))
	}

	cursorpos := cursorPosition(rows, input.cursoroff)

	// Adhere to the max height, adjusting rows output and cursor position.
	if lrows := len(rows); lrows > height {
		start := min(lrows-height, cursorpos[1])
		end := start + height

		rows = rows[start:end]
		cursorpos[1] -= start
	}

	return rows, cursorpos[0], cursorpos[1]
}

func cursorPosition(rows Rows, offset int) []int {
	if offset == 0 {
		return []int{0, 0}
	}

	cursorpos := []int{0, -1} // -1 simplifies the algorithm below.

outer:
	for _, row := range rows {
		cursorpos[0] = 0
		cursorpos[1]++

		for _, cell := range row {
			if cell.Content == nbsp {
				break
			}

			offset--
			cursorpos[0]++

			if offset == 0 {
				break outer
			}
		}
	}

	if len(rows) > 0 && cursorpos[0] == len(rows[0]) && cursorpos[1] < len(rows) {
		cursorpos = []int{0, cursorpos[1] + 1}
	}

	return cursorpos
}
