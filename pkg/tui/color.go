package tui

import (
	"github.com/gdamore/tcell/v2"
)

var ansiColors = map[int]tcell.Color{
	0: tcell.ColorBlack,
	1: tcell.ColorRed,
	2: tcell.ColorGreen,
	3: tcell.ColorYellow,
	4: tcell.ColorBlue,
	5: tcell.ColorPurple,
	6: tcell.ColorTeal,
	7: tcell.ColorWhite,
	9: tcell.ColorDefault,
}

// ApplyANSI modifies a tcell.Style by the given ANSI code.
func ApplyANSI(style tcell.Style, ansi int) tcell.Style {
	switch ansi {
	case 0:
		return tcell.Style{}

	case 1:
		return style.Bold(true)

	case 2:
		return style.Dim(true)

	case 3:
		return style.Italic(true)

	case 4:
		return style.Underline(true)

	case 5:
		return style.Blink(true)

	case 7:
		return style.Reverse(true)

	case 22:
		return style.Bold(false).Dim(false)

	case 23:
		return style.Italic(false)

	case 24:
		return style.Underline(false)

	case 25:
		return style.Blink(false)

	case 27:
		return style.Reverse(false)

	case 30, 31, 32, 33, 34, 35, 36, 37, 39:
		return style.Foreground(ansiColors[ansi%10])

	case 40, 41, 42, 43, 44, 45, 46, 47, 49:
		return style.Background(ansiColors[ansi%10])

	case 90, 91, 92, 93, 94, 95, 96, 97, 99:
		// There might be a better way to represent "high intensity"
		// colors with tcell but I can't find it.
		return style.
			Foreground(ansiColors[ansi%10]).
			Attributes(tcell.AttrBold)

	case 100, 101, 102, 103, 104, 105, 106, 107, 109:
		// There might be a better way to represent "high intensity"
		// colors with tcell but I can't find it.
		return style.
			Background(ansiColors[ansi%10]).
			Attributes(tcell.AttrBold)
	}

	// @todo Implement 256-color scale: https://en.wikipedia.org/wiki/ANSI_escape_code#8-bit

	return style
}
