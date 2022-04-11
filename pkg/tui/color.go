package tui

import "github.com/gdamore/tcell/v2"

var ansiColors = map[int]tcell.Color{
	0: tcell.ColorBlack,
	1: tcell.ColorRed,
	2: tcell.ColorGreen,
	3: tcell.ColorYellow, // väldigt gul - mer brunaktigt i telnet
	4: tcell.ColorBlue,   // som bakgrund ser det mer ut som lavender än blå
	5: tcell.ColorPurple,
	6: tcell.ColorTeal, // cyan
	7: tcell.ColorWhite,
	9: tcell.ColorDefault,
}

func applyANSI(style tcell.Style, ansi int) tcell.Style {
	// fg, bg, attrs := style.Decompose()

	switch {
	case ansi >= 30 && ansi < 50:
		change := style.Foreground
		if ansi >= 40 {
			change = style.Background
		}
		return change(ansiColors[ansi%10])
	}

	// @todo Implement reset (0).
	// @todo Implement support for bold (1) and the other ANSI attributes.
	// Underline (4), high intensity (90), bold high intensity (1;90),
	// high intensity background (100).

	return style
}
