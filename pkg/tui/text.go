package tui

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

// Cell represents one character worth of output in a terminal output.
type Cell struct {
	// @todo Figure out how this main+combined work in tcell so that we
	// model it better on our side.
	Content rune
	Style   tcell.Style
	Width   int
}

// NewCell wraps a rune and creates a Cell.
func NewCell(r rune, style tcell.Style) Cell {
	return Cell{
		Content: r,
		Style:   style,
		Width:   runewidth.RuneWidth(r),
	}
}

// Background sets the background color of the cell.
func (cell *Cell) Background(color tcell.Color) {
	cell.Style = cell.Style.Background(color)
}

// Foreground sets the foreground color of the cell.
func (cell *Cell) Foreground(color tcell.Color) {
	cell.Style = cell.Style.Foreground(color)
}

// Text is a slice of Cells (e.g. a line of characters).
// @todo Rename this `Cells`, which is more apt for non-textual output, like
// the minimap and other visuals.
type Text []Cell

// NewRow creates a new row of cells with the given width.
func NewRow(width int, cells ...Cell) Text {
	row := Text{}

	cell := Cell{Content: ' '}
	if len(cells) > 0 {
		cell = cells[0]
	}

	for len(row) < width {
		row = append(row, cell)
	}

	return row
}

// NewRows creates a given number of rows with the given width.
func NewRows(width, height int, cells ...Cell) []Text {
	rows := []Text{}

	cell := Cell{Content: ' '}
	if len(cells) > 0 {
		cell = cells[0]
	}

	for len(rows) < height {
		rows = append(rows, NewRow(width, cell))
	}

	return rows
}

// NewText parses a byte slice and creates a Text, with ANSI color codes
// abstracted into Cell styles.
func NewText(output []byte, style tcell.Style) Text {
	var text Text

	escaped := false
	parsing := false
	ansi := []rune{}

	for _, r := range []rune(string(output)) {
		if r == '\033' {
			escaped = true
			continue
		}

		if escaped {
			if r == '[' {
				parsing = true
			} else {
				text = append(text, NewCell('^', style), NewCell(r, style))
			}
			escaped = false
			continue
		}

		if parsing {
			if r == ';' || r == 'm' {
				ansii, err := strconv.Atoi(string(ansi))
				if err == nil {
					style = ApplyANSI(style, ansii)
				}

				ansi = []rune{}

				if r == 'm' {
					parsing = false
				}
			} else {
				ansi = append(ansi, r)
			}
			continue
		}

		text = append(text, NewCell(r, style))
	}

	return text
}

// Width calculates the sum of all the containing Cells' width.
func (text *Text) Width() int {
	var width int

	for _, c := range *text {
		width += c.Width
	}

	return width
}

// Wrap breaks a Text down into lines to fit a specified width.
func (text Text) Wrap(width int) []Text {
	lines := []Text{}

wordwrap:
	for i := 0; i < len(text); {
		// Avoid multiple consecutive spaces bleeding over
		// after being wrapped, causing indendation.
		for i > 0 && text[i].Content == ' ' {
			i++
		}

		// If the remains fits the width, we're done.
		if len(text[i:]) <= width {
			lines = append(lines, text[i:])
			break wordwrap
		}

		// Jump to where the width would cut the text and look
		// for the first preceding space, to wrap it there.
		rowwidth := min(width, len(text[i:]))
		for ii := rowwidth; ii >= 0; ii-- {
			if text[i+ii].Content == ' ' {
				lines = append(lines, text[i:i+ii])
				i += ii + 1
				continue wordwrap
			}
		}

		// No space found, so we cut it as is and move on.
		lines = append(lines, text[i:i+rowwidth])
		i += rowwidth
	}

	return lines
}
