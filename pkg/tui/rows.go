package tui

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

// Cell represents one character worth of output in a terminal output.
type Cell struct {
	Content rune
	Style   tcell.Style
	Width   int
}

// NewCell wraps a rune and creates a Cell with an optional style.
func NewCell(r rune, styles ...tcell.Style) Cell {
	if len(styles) == 0 {
		styles = []tcell.Style{{}}
	}

	return Cell{
		Content: r,
		Style:   styles[0],
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

// Row is a slice of Cell and represents one line to be printed to the screen.
type Row []Cell

// NewRow creates a new row of cells with the given width.
func NewRow(width int, cells ...Cell) Row {
	row := Row{}

	cell := Cell{Content: ' '}
	if len(cells) > 0 {
		cell = cells[0]
	}

	for len(row) < width {
		row = row.Append(cell)
	}

	return row
}

func NewRowFromRunes(rs []rune, styles ...tcell.Style) Row {
	row := Row{}

	style := tcell.Style{}
	if len(styles) > 0 {
		style = styles[0]
	}

	for _, r := range rs {
		row = row.Append(NewCell(r, style))
	}

	return row
}

// Append adds a new Cell to the end of the Row.
func (row Row) Append(cells ...Cell) Row {
	for _, cell := range cells {
		row = append(row, cell)
	}

	return row
}

// Prepend adds a new Cell to the beginning of the Row.
func (row Row) Prepend(cells ...Cell) Row {
	for _, cell := range cells {
		row = append(Row{cell}, row...)
	}

	return row
}

// NewRowFromBytes traveses a raw text with ANSI control sequences and
// transforms that into styled Cells.
func NewRowFromBytes(bs []byte, style tcell.Style) Row {
	row := Row{}

	escaped := false
	parsing := false
	ansi := []rune{}

	for _, r := range []rune(string(bs)) {
		if r == '\033' {
			escaped = true
			continue
		}

		if escaped {
			if r == '[' {
				parsing = true
			} else {
				row = row.Append(NewCell('^', style), NewCell(r, style))
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

		row = row.Append(NewCell(r, style))
	}

	return row
}

func (row Row) Copy() Row {
	newrow := Row{}

	for _, cell := range row {
		newrow = append(newrow, cell)
	}

	return newrow
}

// Wrap breaks a Row down into lines to fit a specified width.
func (row Row) Wrap(width int) Rows {
	lines := Rows{}

wordwrap:
	for i := 0; i < len(row); {
		// Avoid multiple consecutive spaces bleeding over
		// after being wrapped, causing indendation.
		for i > 0 && row[i].Content == ' ' {
			i++
		}

		// If the remains fits the width, we're done.
		if len(row[i:]) <= width {
			lines = append(lines, row[i:].Copy())
			break wordwrap
		}

		// Jump to where the width would cut the row and look
		// for the first preceding space, to wrap it there.
		rowwidth := min(width, len(row[i:]))
		for ii := rowwidth; ii >= 0; ii-- {
			if row[i+ii].Content == ' ' {
				lines = append(lines, row[i:i+ii].Copy())
				i += ii + 1
				continue wordwrap
			}
		}

		// No space found, so we cut it as is and move on.
		lines = append(lines, row[i:i+rowwidth].Copy())
		i += rowwidth
	}

	return lines
}

// Rows is a slice of Row and represents an area to be printed to the screen.
type Rows []Row

// NewRows creates a given number of rows with the given width.
func NewRows(width, height int, cells ...Cell) Rows {
	rows := Rows{}

	cell := Cell{Content: ' '}
	if len(cells) > 0 {
		cell = cells[0]
	}

	for len(rows) < height {
		rows = rows.Append(NewRow(width, cell))
	}

	return rows
}

// Append adds a new Row to the end of the Rows.
func (rows Rows) Append(rowses ...Row) Rows {
	for _, row := range rowses {
		rows = append(rows, row)
	}

	return rows
}

// Prepend adds a new Row to the beginning of the Rows.
func (rows Rows) Prepend(rowses ...Row) Rows {
	for _, row := range rowses {
		rows = append(Rows{row}, rows...)
	}

	return rows
}
