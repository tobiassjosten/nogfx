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
			return nil
		}

		if len(tui.input) > 0 {
			tui.input = tui.input[:len(tui.input)-1]
			tui.draw()
		}

	case tcell.KeyETB: // opt/elt+backspace
		if !tui.inputting {
			return nil
		}

		deleted := false
		for i := len(tui.input) - 1; i >= 0; i-- {
			if tui.input[i].Content == ' ' {
				tui.input = tui.input[0:i]
				deleted = true
				break
			}
		}

		if !deleted {
			tui.input = Text{}
		}

		tui.draw()

	case tcell.KeyNAK: // cmd/ctrl+backspace
		if !tui.inputting {
			return nil
		}

		if len(tui.input) > 0 {
			tui.input = Text{}
			tui.draw()
		}

	case tcell.KeyEnter:
		if !tui.inputting {
			// @todo Keep the text to enable quick repetition.
			return nil
		}

		input := tui.input.Bytes()

		tui.input = Text{}
		tui.draw()

		return input

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

			return nil
		}

		tui.input = append(tui.input, NewCell(ev.Rune(), inputStyle))
		tui.draw()

	default:
		// @todo Use arrow keys (wht opt/cmd mods) to jump the cursor.

		// @todo Remove this when we're done exploring keys and their
		// mappings.
		tui.Print([]byte(fmt.Sprintf(
			"[Unknown key pressed: '%d']", ev.Key(),
		)))
	}

	return nil
}
