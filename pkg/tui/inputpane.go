package tui

import (
	"fmt"
	"log"
	"unicode"

	"github.com/gdamore/tcell/v2"
)

// InputPane is the pane that takes input and users most often interact with.
type InputPane struct {
	tui *TUI

	inputs chan []byte

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
func NewInputPane(tui *TUI, inputStyle, inputtedStyle tcell.Style) *InputPane {
	return &InputPane{
		tui:           tui,
		inputs:        make(chan []byte),
		inputStyle:    inputStyle,
		inputtedStyle: inputtedStyle,
	}
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

// Resize triggers a new layout to be calculated, if needed.
func (pane *InputPane) Resize() {
	resize := len(pane.input) >= pane.width && pane.height != pane.Height()
	resize = resize || pane.height > 1 && len(pane.input) <= pane.width
	resize = resize || pane.height > 0 && !pane.inputting
	resize = resize || pane.height == 0 && pane.inputting

	if resize {
		pane.tui.Resize(pane.tui.screen.Size())
	}
}

// HandleEvents reacts on a user event and modifies itself from it.
func (pane *InputPane) HandleEvent(event tcell.Event) bool {
	ev, ok := event.(*tcell.EventKey)
	if !ok {
		return false
	}

	if !pane.inputting {
		if ev.Key() == tcell.KeyRune && ev.Rune() == ' ' {
			pane.inputting = true
			pane.Resize()
			return true
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
			return true

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
			return true
		}
	}

	switch ev.Key() {
	case tcell.KeyCtrlC:
		pane.input = []rune{}
		pane.cursor = 0
		fallthrough

	case tcell.KeyESC:
		pane.inputting = false
		pane.Resize()
		return true

	case tcell.KeyBackspace, tcell.KeyBackspace2:
		pane.input = append(
			pane.input[:pane.cursor-1],
			pane.input[pane.cursor:]...,
		)
		pane.cursor--
		pane.Resize()
		return true

	case tcell.KeyETB: // opt/elt+backspace
		delete := false
		for i := pane.cursor - 1; i > 0; i-- {
			if delete && pane.input[i] == ' ' {
				pane.input = append(
					pane.input[:i+1],
					pane.input[pane.cursor:]...,
				)
				pane.cursor = i + 1
				pane.Resize()
				return true
			}
			if !delete && pane.input[i] != ' ' {
				delete = true
			}
		}
		pane.input = pane.input[pane.cursor:]
		pane.cursor = 0
		pane.Resize()
		return true

	case tcell.KeyNAK: // cmd/ctrl+backspace
		pane.input = pane.input[pane.cursor:]
		pane.cursor = 0
		pane.Resize()
		return true

	case tcell.KeyLeft:
		if pane.cursor > 0 {
			pane.cursor--
		}
		return true

	case tcell.KeyRight:
		if pane.cursor < len(pane.input) {
			pane.cursor++
		}
		return true

	case tcell.KeyUp:
		// Search like fish

	case tcell.KeyDown:
		// search

	case tcell.KeyEnter:
		pane.inputs <- []byte(string(pane.input))
		if pane.masked {
			pane.input = []rune{}
			pane.cursor = 0
			pane.Resize()
		} else {
			pane.inputted = true
		}
		return true

	case tcell.KeyRune:
		pane.input = append(pane.input[:pane.cursor], append(
			[]rune{ev.Rune()}, pane.input[pane.cursor:]...,
		)...)
		pane.cursor++
		pane.Resize()

		return true

	default:
		// @todo Remove this when we're done exploring keys and their
		// mappings.
		log.Println(fmt.Sprintf("[Unknown key: '%d']", ev.Key()))
	}

	return false
}

// HandleBinding reacts to keypress events during normal mode.
func (pane *InputPane) HandleBinding(ev *tcell.EventKey) bool {
	if ev.Key() != tcell.KeyRune {
		// This guard here doesn't make sense now but it will when we
		// have non-rune key bindings in the future.
		return false
	}

	switch ev.Rune() {
	case '1':
		pane.inputs <- []byte{'s', 'w'}
		return true

	case '2':
		pane.inputs <- []byte{'s'}
		return true

	case '3':
		pane.inputs <- []byte{'s', 'e'}
		return true

	case '4':
		pane.inputs <- []byte{'w'}
		return true

	case '5':
		pane.inputs <- []byte{'m', 'a', 'p'}
		return true

	case '6':
		pane.inputs <- []byte{'e'}
		return true

	case '7':
		pane.inputs <- []byte{'n', 'w'}
		return true

	case '8':
		pane.inputs <- []byte{'n'}
		return true

	case '9':
		pane.inputs <- []byte{'n', 'e'}
		return true
	}

	return false
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
