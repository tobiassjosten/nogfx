package tui

import (
	"fmt"
	"log"
	"unicode"

	"github.com/gdamore/tcell/v2"
)

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

// NewInputPane creates a new InputPane.
func NewInputPane(inputStyle, inputtedStyle tcell.Style) *InputPane {
	return &InputPane{
		inputStyle:    inputStyle,
		inputtedStyle: inputtedStyle,
		input:         []rune{},
	}
}

// MaskInput replaces input with stars when printed.
func (pane *InputPane) Mask() {
	pane.masked = true
}

// UnmaskInput shows the actual input when printed.
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

	for _, r := range pane.input {
		word = append(word, r)
		wwidth++

		if unicode.IsSpace(r) || wwidth >= pane.width {
			x += wwidth
			word = []rune{}
			wwidth = 0
		}

		if x+wwidth > pane.x+pane.width || r == '\n' {
			x = pane.x
			height++
		}
	}

	return height
}

// HandleEvent reacts on a user event and modifies itself from it.
func (pane *InputPane) HandleEvent(event tcell.Event) (bool, []rune) {
	ev, ok := event.(*tcell.EventKey)
	if !ok {
		return false, nil
	}

	if !pane.inputting {
		if ev.Key() == tcell.KeyRune && ev.Rune() == ' ' {
			pane.inputting = true
			return true, nil
		}

		return pane.HandleBinding(ev)
	}

	if pane.inputted {
		switch ev.Key() {
		case tcell.KeyBackspace, tcell.KeyBackspace2:
			fallthrough
		case tcell.KeyETB: // opt/elt+backspace
			fallthrough
		case tcell.KeyNAK: // cmd/ctrl+backspace
			pane.input = []rune{}
			pane.cursor = 0
			pane.inputted = false
			return true, nil

		case tcell.KeyRune:
			pane.input = []rune{}
			pane.cursor = 0
			pane.inputted = false
		}
	}

	if len(pane.input) == 0 || pane.cursor == 0 {
		switch ev.Key() {
		case tcell.KeyBackspace, tcell.KeyBackspace2:
			fallthrough
		case tcell.KeyETB: // opt/elt+backspace
			fallthrough
		case tcell.KeyNAK: // cmd/ctrl+backspace
			return true, nil
		}
	}

	switch ev.Key() {
	case tcell.KeyCtrlC:
		pane.input = []rune{}
		pane.cursor = 0
		fallthrough

	case tcell.KeyESC:
		pane.inputting = false
		return true, nil

	case tcell.KeyBackspace, tcell.KeyBackspace2:
		pane.input = append(
			pane.input[:pane.cursor-1],
			pane.input[pane.cursor:]...,
		)
		pane.cursor--
		return true, nil

	case tcell.KeyETB: // opt/elt+backspace
		delete := false
		for i := pane.cursor - 1; i > 0; i-- {
			if delete && pane.input[i] == ' ' {
				pane.input = append(
					pane.input[:i+1],
					pane.input[pane.cursor:]...,
				)
				pane.cursor = i + 1
				return true, nil
			}
			if !delete && pane.input[i] != ' ' {
				delete = true
			}
		}
		pane.input = pane.input[pane.cursor:]
		pane.cursor = 0
		return true, nil

	case tcell.KeyNAK: // cmd/ctrl+backspace
		pane.input = pane.input[pane.cursor:]
		pane.cursor = 0
		return true, nil

	case tcell.KeyLeft:
		if pane.cursor > 0 {
			pane.cursor--
		}
		return true, nil

	case tcell.KeyRight:
		if pane.cursor < len(pane.input) {
			pane.cursor++
		}
		return true, nil

	case tcell.KeyUp:
		// Search like fish

	case tcell.KeyDown:
		// search

	case tcell.KeyEnter:
		input := pane.input
		if pane.masked {
			pane.input = []rune{}
			pane.cursor = 0
		} else {
			pane.inputted = true
		}
		return true, append(input, '\n')

	case tcell.KeyRune:
		pane.input = append(pane.input[:pane.cursor], append(
			[]rune{ev.Rune()}, pane.input[pane.cursor:]...,
		)...)
		pane.cursor++

		return true, nil

	default:
		// @todo Remove this when we're done exploring keys and their
		// mappings.
		log.Println(fmt.Sprintf("[Unknown key: '%d']", ev.Key()))
	}

	return false, nil
}

// HandleBinding reacts to keypress events during normal mode.
func (pane *InputPane) HandleBinding(ev *tcell.EventKey) (bool, []rune) {
	if ev.Key() != tcell.KeyRune {
		// This guard here doesn't make sense now but it will when we
		// have non-rune key bindings in the future.
		return false, nil
	}

	switch ev.Rune() {
	case '1':
		return true, []rune{'s', 'w'}

	case '2':
		return true, []rune{'s'}

	case '3':
		return true, []rune{'s', 'e'}

	case '4':
		return true, []rune{'w'}

	case '5':
		return true, []rune{'m', 'a', 'p'}

	case '6':
		return true, []rune{'e'}

	case '7':
		return true, []rune{'n', 'w'}

	case '8':
		return true, []rune{'n'}

	case '9':
		return true, []rune{'n', 'e'}
	}

	return false, nil
}

// Draw prints the contents of the InputPane to the given tcell.Screen.
func (pane *InputPane) Draw(screen tcell.Screen) {
	if !pane.inputting {
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

	for i, r := range pane.input {
		if pane.masked {
			r = '*'
		}

		word = append(word, r)
		wwidth++

		if unicode.IsSpace(r) || wwidth >= pane.width {
			for i, r := range word {
				screen.SetContent(x+i, y, r, nil, style)
			}

			x += wwidth
			word = []rune{}
			wwidth = 0
		}

		if x+wwidth > pane.x+pane.width || r == '\n' {
			for x < pane.x+pane.width {
				screen.SetContent(x, y, ' ', nil, style)
				x++
			}
			x = pane.x
			y++
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
