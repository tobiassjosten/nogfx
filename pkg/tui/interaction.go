package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

func (tui *TUI) InteractKey(ev *tcell.EventKey) []byte {
	switch ev.Key() {
	case tcell.KeyESC:
		tui.inputting = false
		tui.screen.HideCursor()
		tui.draw()

	case tcell.KeyBackspace, tcell.KeyBackspace2:
		if !tui.inputting {
			return []byte{}
		}

		if len(tui.input) > 0 {
			tui.input = tui.input[:len(tui.input)-1]
			tui.draw()
		}

	case tcell.KeyETB: // opt/elt+backspace
		if !tui.inputting {
			return []byte{}
		}

		deleted := false
		for i := len(tui.input) - 1; i >= 0; i-- {
			if tui.input[i] == ' ' {
				tui.input = tui.input[0:i]
				deleted = true
				break
			}
		}

		if !deleted {
			tui.input = []rune{}
		}

		tui.draw()

	case tcell.KeyNAK: // cmd/ctrl+backspace
		if !tui.inputting {
			return []byte{}
		}

		if len(tui.input) > 0 {
			tui.input = []rune{}
			tui.draw()
		}

	case tcell.KeyEnter:
		if !tui.inputting {
			return []byte{}
		}

		tui.input = []rune{}
		tui.draw()

		return []byte(string(tui.input))

	case tcell.KeyRune:
		if !tui.inputting {
			switch ev.Rune() {
			case ' ':
				tui.inputting = true
				tui.draw()

			case '1':
				return []byte{'s', 'w'}

			case '2':
				return []byte{'s'}

			case '3':
				return []byte{'s', 'e'}

			case '4':
				return []byte{'w'}

			case '5':
				return []byte{'m', 'a', 'p'}

			case '6':
				return []byte{'e'}

			case '7':
				return []byte{'n', 'w'}

			case '8':
				return []byte{'n'}

			case '9':
				return []byte{'n', 'e'}
			}

			return []byte{}
		}

		tui.input = append(tui.input, ev.Rune())
		tui.draw()

	default:
		// @todo Remove this when we're done exploring keys and their
		// mappings.
		tui.Print([]byte(fmt.Sprintf(
			"[Unknown key pressed: '%d']", ev.Key(),
		)))
	}

	return []byte{}
}
