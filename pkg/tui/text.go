package tui

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

type Cell struct {
	// @todo Figure out how this main+combined work in tcell so that we
	// model it better on our side.
	Content rune
	Style   tcell.Style
	Width   int
}

func NewCell(r rune, style tcell.Style) Cell {
	return Cell{
		Content: r,
		Style:   style,
		Width:   runewidth.RuneWidth(r),
	}
}

type Text []Cell

func NewText(output []byte, style tcell.Style) (Text, tcell.Style) {
	var text Text

	escaped := false
	parsing := false
	ansi := []rune{}

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
					style = applyANSI(style, ansii)
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

func (text *Text) Width() int {
	var width int

	for _, c := range *text {
		width += c.Width
	}

	return width
}
