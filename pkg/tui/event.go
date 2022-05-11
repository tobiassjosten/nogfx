package tui

import (
	"github.com/gdamore/tcell/v2"
)

const (
	normalMode   = 100_000
	inputMode    = 200_000
	inputtedMode = 300_000
)

func (tui *TUI) eventHandlers() map[int]func(rune) bool {
	return map[int]func(rune) bool{
		int(tcell.KeyRune) + int(' ') + normalMode: tui.handleSpaceNormal,
		int(tcell.KeyRune) + normalMode:            tui.handleRuneNormal,

		int(tcell.KeyRune) + inputtedMode: tui.handleRuneInputted,

		int(tcell.KeyBackspace) + inputtedMode:  tui.handleBackspaceInputted,
		int(tcell.KeyBackspace2) + inputtedMode: tui.handleBackspaceInputted,
		int(tcell.KeyETB) + inputtedMode:        tui.handleBackspaceInputted,
		int(tcell.KeyNAK) + inputtedMode:        tui.handleBackspaceInputted,

		int(tcell.KeyRune) + inputMode:  tui.handleRuneInput,
		int(tcell.KeyEnter) + inputMode: tui.handleEnterInput,
		int(tcell.KeyEsc) + inputMode:   tui.handleEscInput,
		int(tcell.KeyCtrlC) + inputMode: tui.handleCtrlCInput,
		int(tcell.KeyLeft) + inputMode:  tui.handleLeftInput,
		int(tcell.KeyRight) + inputMode: tui.handleRightInput,

		int(tcell.KeyBackspace) + inputMode:  tui.handleBackspaceInput,
		int(tcell.KeyBackspace2) + inputMode: tui.handleBackspaceInput,
		int(tcell.KeyETB) + inputMode:        tui.handleOptBackspaceInput,
		int(tcell.KeyNAK) + inputMode:        tui.handleCmdBackspaceInput,

		int(tcell.KeyUp):                       tui.handleUpInput,
		int(tcell.KeyUp) + int(tcell.ModAlt):   tui.handleAltUpInput,
		int(tcell.KeyDown):                     tui.handleDownInput,
		int(tcell.KeyDown) + int(tcell.ModAlt): tui.handleAltDownInput,
	}
}

// HandleEvent reacts on a user event and modifies itself from it.
func (tui *TUI) HandleEvent(event *tcell.EventKey) bool {
	// List alternatives in order of specificity (descending).
	alts := []int{
		int(event.Rune()) + int(event.Key()),
		int(event.Rune()),
		int(event.Key()),
	}

	// Add modified alternatives based on the three basic ones.
	addAlt := func(n int) {
		for _, alt := range alts[:3] {
			alts = append(alts, alt+n)
		}
	}

	if event.Modifiers()&tcell.ModAlt > 0 {
		addAlt(int(tcell.ModAlt))
	}
	if tui.input.inputting {
		if tui.input.inputted {
			if event.Modifiers()&tcell.ModAlt > 0 {
				addAlt(inputtedMode + int(tcell.ModAlt))
			}
			addAlt(inputtedMode)
		}

		if event.Modifiers()&tcell.ModAlt > 0 {
			addAlt(inputMode + int(tcell.ModAlt))
		}
		addAlt(inputMode)
	} else {
		if event.Modifiers()&tcell.ModAlt > 0 {
			addAlt(normalMode + int(tcell.ModAlt))
		}
		addAlt(normalMode)
	}

	// Move the three basic alternatives last, since they're the least
	// specific ones and only served as templates for addAlt().
	alts = append(alts[2:], alts[:2]...)
	if event.Modifiers()&tcell.ModAlt > 0 {
		alts = append(alts[2:], alts[:2]...)
	}

	for _, alt := range alts {
		if f, ok := tui.eventHandlers()[alt]; ok {
			if handled := f(event.Rune()); handled {
				return true
			}
		}
	}

	return false
}

func (tui *TUI) handleSpaceNormal(_ rune) bool {
	tui.input.inputting = true
	return true
}

func (tui *TUI) handleRuneInputted(r rune) bool {
	_ = tui.handleBackspaceInputted(r)
	_ = tui.handleRuneInput(r)
	return true
}

func (tui *TUI) handleBackspaceInputted(_ rune) bool {
	tui.input.buffer = []rune{}
	tui.input.cursor = 0
	tui.input.inputted = false
	return true
}

func (tui *TUI) handleBackspaceInput(_ rune) bool {
	cursor := int(max(0, tui.input.cursor-1))
	tui.input.buffer = append(
		tui.input.buffer[:cursor],
		tui.input.buffer[tui.input.cursor:]...,
	)
	tui.input.cursor = cursor
	return true
}

func (tui *TUI) handleOptBackspaceInput(_ rune) bool {
	delete := false
	for i := tui.input.cursor - 1; i > 0; i-- {
		if delete && tui.input.buffer[i] == ' ' {
			tui.input.buffer = append(
				tui.input.buffer[:i+1],
				tui.input.buffer[tui.input.cursor:]...,
			)
			tui.input.cursor = i + 1
			return true
		}
		if !delete && tui.input.buffer[i] != ' ' {
			delete = true
		}
	}
	tui.input.buffer = tui.input.buffer[tui.input.cursor:]
	tui.input.cursor = 0
	return true
}

func (tui *TUI) handleCmdBackspaceInput(_ rune) bool {
	tui.input.buffer = tui.input.buffer[tui.input.cursor:]
	tui.input.cursor = 0
	return true
}

func (tui *TUI) handleRuneInput(r rune) bool {
	tui.input.buffer = append(tui.input.buffer[:tui.input.cursor], append(
		[]rune{r}, tui.input.buffer[tui.input.cursor:]...,
	)...)
	tui.input.cursor++
	return true
}

func (tui *TUI) handleEnterInput(_ rune) bool {
	input := tui.input.buffer
	tui.input.inputted = true
	if tui.input.masked {
		tui.input.inputted = false
		tui.input.buffer = []rune{}
		tui.input.cursor = 0
	}

	tui.inputs <- []byte(string(input))

	return true
}

func (tui *TUI) handleEscInput(_ rune) bool {
	tui.input.inputting = false
	return true
}

func (tui *TUI) handleCtrlCInput(_ rune) bool {
	tui.input.buffer = []rune{}
	tui.input.cursor = 0
	tui.input.inputting = false
	return true
}

func (tui *TUI) handleLeftInput(_ rune) bool {
	tui.input.inputted = false
	tui.input.cursor = int(max(0, tui.input.cursor-1))
	return true
}

func (tui *TUI) handleRightInput(_ rune) bool {
	tui.input.inputted = false
	tui.input.cursor = int(min(len(tui.input.buffer), tui.input.cursor+1))
	return true
}

func (tui *TUI) handleRuneNormal(r rune) bool {
	switch r {
	case '1':
		tui.inputs <- []byte{'s'}
		tui.inputs <- []byte{'w'}
		return true

	case '2':
		tui.inputs <- []byte{'s'}
		return true

	case '3':
		tui.inputs <- []byte{'s'}
		tui.inputs <- []byte{'e'}
		return true

	case '4':
		tui.inputs <- []byte{'w'}
		return true

	case '5':
		tui.inputs <- []byte{'m', 'a', 'p'}
		return true

	case '6':
		tui.inputs <- []byte{'e'}
		return true

	case '7':
		tui.inputs <- []byte{'n'}
		tui.inputs <- []byte{'w'}
		return true

	case '8':
		tui.inputs <- []byte{'n'}
		return true

	case '9':
		tui.inputs <- []byte{'n'}
		tui.inputs <- []byte{'e'}
		return true
	}

	return false
}

func (tui *TUI) handleUpInput(_ rune) bool {
	tui.output.offset++
	return true
}

func (tui *TUI) handleAltUpInput(_ rune) bool {
	tui.output.offset += 5
	return true
}

func (tui *TUI) handleDownInput(_ rune) bool {
	tui.output.offset--
	if tui.output.offset < 0 {
		tui.output.offset = 0
	}
	return true
}

func (tui *TUI) handleAltDownInput(_ rune) bool {
	tui.output.offset -= 5
	if tui.output.offset < 0 {
		tui.output.offset = 0
	}
	return true
}
