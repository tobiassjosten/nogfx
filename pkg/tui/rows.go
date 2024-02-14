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

// Foreground sets the foreground color of the cell.
func (cell *Cell) Foreground(color tcell.Color) {
	cell.Style = cell.Style.Foreground(color)
}

// Row is a slice of Cell and represents one line to be printed to the screen.
type Row []Cell

// NewRow creates a new row of cells with the given width.
func NewRow(width int, cells ...Cell) Row {
	row := Row{}

	if len(cells) == 0 {
		cells = []Cell{NewCell(' ')}
	}

	for len(row) < width {
		row = row.append(cells[0])
	}

	return row
}

// NewRowFromRunes creates a new Row from a slice of runes.
func NewRowFromRunes(rs []rune, styles ...tcell.Style) Row {
	row := Row{}

	if len(styles) == 0 {
		styles = []tcell.Style{{}}
	}

	for _, r := range rs {
		row = row.append(NewCell(r, styles[0]))
	}

	return row
}

// NewRowFromBytes traveses a raw text with ANSI control sequences and
// transforms that into styled Cells.
func NewRowFromBytes(bs []byte, styles ...tcell.Style) (Row, tcell.Style) {
	row := Row{}

	if len(styles) == 0 {
		styles = []tcell.Style{{}}
	}

	style := styles[0]

	escaped := false
	parsing := false
	ansi := []rune{}

	for _, r := range string(bs) {
		if r == '\033' {
			escaped = true
			continue
		}

		if escaped {
			if r == '[' {
				parsing = true
			} else {
				row = row.append(NewCell('^', style), NewCell(r, style))
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

		row = row.append(NewCell(r, style))
	}

	return row, style
}

// String converts the row to a string.
func (row Row) String() (str string) {
	for _, c := range row {
		str += string(c.Content)
	}

	return
}

// Append adds a new Cell to the end of the Row.
func (row Row) append(cells ...Cell) Row {
	return append(row, cells...)
}

// revIndexSpace finds the first space in the last block of spaces, or -1 if
// there are no spaces.
func (row Row) revIndexSpace() int {
	space := -1

	for i := len(row) - 1; i >= 0; i-- {
		if row[i].Content == ' ' {
			space = i
		} else if space > 0 {
			return space
		}
	}

	return space
}

// indexNospace finds the first non-space, or -1 if there are no non-spaces.
func (row Row) indexNospace() int {
	for i := 0; i < len(row); i++ {
		if row[i].Content != ' ' {
			return i
		}
	}

	return -1
}

// revIndexNospace finds the first non-space in the last block of non-spaces,
// or -1 if there are no non-spaces.
func (row Row) revIndexNospace() int {
	nospace := -1

	for i := len(row) - 1; i >= 0; i-- {
		if row[i].Content != ' ' {
			nospace = i
		} else if nospace > 0 {
			return nospace
		}
	}

	return nospace
}

// Wrap breaks the Row into Rows to fit the given width.
func (row Row) Wrap(width int, padding ...Cell) Rows {
	lrow := len(row)
	if lrow == 0 || lrow <= width {
		if len(padding) > 0 {
			row = row.Pad(width, padding[0])
		}

		return Rows{row}
	}

	rows := Rows{}

wordwrap:
	for i := 0; i < lrow; {
		// If the remains fits the width, we're done.
		if len(row[i:]) <= width {
			rows = append(rows, row[i:])
			break wordwrap
		}

		// Jump to where the width would cut the row and look back
		// from there to know where to wrap.
		ii := min(width, len(row[i:]))
		riSpace := row[i : i+ii].revIndexSpace()
		riNospace := row[i : i+ii].revIndexNospace()

		switch {
		// Either it's all spaces or all non-spaces, or its only space
		// exists in the beginning of the row, which we preseve. So we
		// can't wrap it but take the lot.
		case riSpace < 0 || riNospace < 0 || riSpace == 0:
			rows = append(rows, row[i:i+ii])

		// It ends with spaces. So we wrap at the preceding non-space
		// and skip the following spaces.
		case riSpace > riNospace:
			rows = append(rows, row[i:i+riSpace])

		// It ends with non-spaces and it succeeded by a space. So we
		// can conveniently take the whole row and then skip the spaces.
		case row[i+ii].Content == ' ':
			rows = append(rows, row[i:i+ii])

		case riNospace > riSpace:
			rows = append(rows, row[i:i+riSpace])
			ii = riNospace
		}

		i += ii
		if iii := row[i:].indexNospace(); iii > 0 {
			i += iii
		} else if iii == -1 {
			break
		}
	}

	if len(padding) > 0 {
		for i, row := range rows {
			rows[i] = row.Pad(width, padding[0])
		}
	}

	return rows
}

// Pad adds cells to make the row a certain length.
func (row Row) Pad(width int, padding Cell) Row {
	newrow := append(Row{}, row...)
	for i := len(newrow); i < width; i++ {
		newrow = append(newrow, padding)
	}

	return newrow
}

// Rows is a slice of Row and represents an area to be printed to the screen.
type Rows []Row

// NewRows creates a given number of rows with the given width.
func NewRows(width, height int, cells ...Cell) Rows {
	rows := Rows{}

	cell := NewCell(' ')
	if len(cells) > 0 {
		cell = cells[0]
	}

	for len(rows) < height {
		rows = rows.append(NewRow(width, cell))
	}

	return rows
}

// Strings converts the rows to a slice of strings.
func (rows Rows) Strings() (strs []string) {
	for _, row := range rows {
		strs = append(strs, row.String())
	}

	return
}

// Append adds a new Row to the end of the Rows.
func (rows Rows) append(rowses ...Row) Rows {
	for _, row := range rowses {
		rows = append(rows, row)
	}

	return rows
}

// Prepend adds a new Row to the beginning of the Rows.
func (rows Rows) prepend(rowses ...Row) Rows {
	for _, row := range rowses {
		rows = append(Rows{row}, rows...)
	}

	return rows
}
