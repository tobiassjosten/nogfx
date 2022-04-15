package tui

import (
	"bytes"
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

// Text is a slice of Cells (e.g. a line of characters).
type Text []Cell

// NewText parses a byte slice and creates a Text, with ANSI color codes
// abstracted into Cell styles.
func NewText(output []byte, style tcell.Style) (Text, tcell.Style) {
	var text Text

	escaped := false
	parsing := false
	ansi := []rune{}

	// We replace GA with a newline so we can throw away these otherwise
	// useful newlines.
	output = bytes.TrimLeft(output, "\r\n")

	for _, r := range []rune(string(output)) {
		if r == '\r' {
			continue
		}

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

	return text, style
}

// Width calculates the sum of all the containing Cells' width.
func (text *Text) Width() int {
	var width int

	for _, c := range *text {
		width += c.Width
	}

	return width
}

func (text *Text) Bytes() []byte {
	var runes []rune
	for _, c := range *text {
		runes = append(runes, c.Content)
	}

	return []byte(string(runes))
}
