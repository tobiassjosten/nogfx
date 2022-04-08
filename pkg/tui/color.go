package tui

import "github.com/gdamore/tcell/v2"

var ansiColors = map[int]tcell.Color{
	0: tcell.ColorBlack,
	1: tcell.ColorRed,
	2: tcell.ColorGreen,
	3: tcell.ColorYellow,
	4: tcell.ColorBlue,
	5: tcell.ColorMaroon, // magenta
	6: tcell.ColorTeal,   // cyan
	7: tcell.ColorWhite,
	9: tcell.ColorDefault,
}

func applyANSI(style tcell.Style, ansi int) tcell.Style {
	switch {
	case ansi >= 30 && ansi < 50:
		change := style.Foreground
		if ansi >= 40 {
			change = style.Background
		}
		return change(ansiColors[ansi%10])
	}

	// @todo Implement support for light/normal and the other ANSI
	// attributes.

	return style
}
