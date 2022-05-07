package tui

import (
	"unicode"

	"github.com/gdamore/tcell/v2"
)

// MaskInput hides the content of the InputPane.
func (tui *TUI) MaskInput() {
	tui.panes.Input.Mask()
}

// UnmaskInput shows the content of the InputPane.
func (tui *TUI) UnmaskInput() {
	tui.panes.Input.Unmask()
}

// InputPane is the pane that takes input and users most often interact with.
type InputPane struct {
	x      int
	y      int
	width  int
	height int

	inputStyle    tcell.Style
	inputtedStyle tcell.Style

	input     []rune
	inputting bool
	inputted  bool
	masked    bool
	cursor    int
}

const (
	normalMode   = 100_000
	inputMode    = 200_000
	inputtedMode = 300_000
)

// NewInputPane creates a new InputPane.
func NewInputPane() *InputPane {
	var (
		inputStyle = (tcell.Style{}).
				Foreground(tcell.ColorWhite).
				Background(tcell.ColorGray)
		inputtedStyle = (tcell.Style{}).
				Foreground(tcell.ColorWhite).
				Background(tcell.ColorGray).
				Attributes(tcell.AttrDim)
	)

	pane := &InputPane{
		inputStyle:    inputStyle,
		inputtedStyle: inputtedStyle,
		input:         []rune{},
	}

	return pane
}

// Position sets the x.y coordinates for and resizes the pane.
func (pane *InputPane) Position(x, y, width, height int) {
	pane.x, pane.y = x, y
	pane.width, pane.height = width, height
}

// Mask replaces input with stars when printed.
func (pane *InputPane) Mask() {
	pane.input = []rune{}
	pane.cursor = 0
	pane.masked = true
}

// Unmask shows the actual input when printed.
func (pane *InputPane) Unmask() {
	pane.masked = false
}

// Height is the actual height that a full rendition of InputPane would need,
// as opposed to its `height` property, which is what it's afforded.
func (pane *InputPane) Height() int {
	if !pane.inputting {
		return 0
	}

	x, height := 0, 1

	word := []rune{}
	wwidth := 0
	isword := false

	for _, r := range pane.input {
		if !isword && r != ' ' {
			isword = true
		}

		word = append(word, r)
		wwidth++

		if (isword && unicode.IsSpace(r)) || wwidth >= pane.width {
			x += wwidth
			word = []rune{}
			wwidth = 0
			isword = false
		}

		if x+wwidth > pane.x+pane.width {
			x = pane.x
			height++

			// Swallow a space if it comes right after wrapping.
			if len(word) == 1 && word[0] == ' ' {
				word = []rune{}
				wwidth = 0
				isword = false
			}
		}
	}

	return height

	// @todo Replace this and Draw with one uniform way of painting.
}

// Draw prints the contents of the InputPane to the given tcell.Screen.
func (pane *InputPane) Draw(screen tcell.Screen) {
	if pane.height == 0 {
		screen.HideCursor()
		return
	}

	cursorX, cursorY := pane.x, pane.y

	style := pane.inputStyle
	if pane.inputted {
		style = pane.inputtedStyle
	}

	x, y := pane.x, pane.y

	word := []rune{}
	wwidth := 0
	isword := false

	for i, r := range pane.input {
		if pane.masked {
			r = '*'
		}
		if !isword && r != ' ' {
			isword = true
		}

		word = append(word, r)
		wwidth++

		if (isword && unicode.IsSpace(r)) || wwidth >= pane.width {
			for i, r := range word {
				screen.SetContent(x+i, y, r, nil, style)
			}

			x += wwidth
			word = []rune{}
			wwidth = 0
			isword = false
		}

		if x+wwidth > pane.x+pane.width {
			for x < pane.x+pane.width {
				screen.SetContent(x, y, ' ', nil, style)
				x++
			}
			x = pane.x
			y++

			// Swallow a space if it comes right after wrapping.
			if len(word) == 1 && word[0] == ' ' {
				word = []rune{}
				wwidth = 0
				isword = false
			}
		}

		if i+1 == pane.cursor {
			cursorX, cursorY = x+wwidth, y
		}
	}

	for i, r := range word {
		screen.SetContent(x+i, y, r, nil, style)
	}
	for x+wwidth < pane.x+pane.width {
		screen.SetContent(x+wwidth, y, ' ', nil, style)
		x++
	}

	screen.ShowCursor(cursorX, cursorY)
}
