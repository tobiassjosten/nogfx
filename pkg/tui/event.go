package tui

import (
	"github.com/gdamore/tcell/v2"
)

const (
	// @todo With only this mode left, it's a little obsolete. Remove and
	// add the distinction in each event handler instead.
	inputtedMode = 100_000
)

func (tui *TUI) eventHandlers() map[int]func(rune) bool {
	alt := int(tcell.ModAlt)

	return map[int]func(rune) bool{
		int(tcell.KeyRune) + inputtedMode: tui.handleRuneInputted,

		int(tcell.KeyBackspace) + inputtedMode:  tui.handleBackspaceInputted,
		int(tcell.KeyBackspace2) + inputtedMode: tui.handleBackspaceInputted,
		int(tcell.KeyETB) + inputtedMode:        tui.handleBackspaceInputted,
		int(tcell.KeyNAK) + inputtedMode:        tui.handleBackspaceInputted,

		int(tcell.KeyRune):  tui.handleRuneInput,
		int(tcell.KeyEnter): tui.handleEnterInput,
		int(tcell.KeyEsc):   tui.handleEscInput,
		int(tcell.KeyCtrlC): tui.handleCtrlCInput,
		int(tcell.KeyLeft):  tui.handleLeftInput,
		int(tcell.KeyRight): tui.handleRightInput,

		int(tcell.KeyBackspace):  tui.handleBackspaceInput,
		int(tcell.KeyBackspace2): tui.handleBackspaceInput,
		int(tcell.KeyETB):        tui.handleOptBackspaceInput,
		int(tcell.KeyNAK):        tui.handleCmdBackspaceInput,

		int(tcell.KeyUp):         tui.handleUpInput,
		int(tcell.KeyUp) + alt:   tui.handleAltUpInput,
		int(tcell.KeyDown):       tui.handleDownInput,
		int(tcell.KeyDown) + alt: tui.handleAltDownInput,
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
	if tui.input.inputted {
		if event.Modifiers()&tcell.ModAlt > 0 {
			addAlt(inputtedMode + int(tcell.ModAlt))
		}
		addAlt(inputtedMode)
	}

	if event.Modifiers()&tcell.ModAlt > 0 {
		addAlt(int(tcell.ModAlt))
	}

	// Move the three basic alternatives last, since they're the least
	// specific ones and only served as templates for addAlt().
	alts = append(alts[3:], alts[:3]...)
	if event.Modifiers()&tcell.ModAlt > 0 {
		alts = append(alts[3:], alts[:3]...)
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

func (tui *TUI) handleRuneInputted(r rune) bool {
	_ = tui.handleBackspaceInputted(r)
	_ = tui.handleRuneInput(r)
	return true
}

func (tui *TUI) handleBackspaceInputted(_ rune) bool {
	tui.setCache(paneInput, nil)
	tui.input.buffer = []rune{}
	tui.input.cursoroff = 0
	tui.input.inputted = false
	return true
}

func (tui *TUI) handleBackspaceInput(_ rune) bool {
	tui.setCache(paneInput, nil)
	cursor := int(max(0, tui.input.cursoroff-1))
	tui.input.buffer = append(
		tui.input.buffer[:cursor],
		tui.input.buffer[tui.input.cursoroff:]...,
	)
	tui.input.cursoroff = cursor
	return true
}

func (tui *TUI) handleOptBackspaceInput(_ rune) bool {
	tui.setCache(paneInput, nil)
	delete := false
	for i := tui.input.cursoroff - 1; i > 0; i-- {
		if delete && tui.input.buffer[i] == ' ' {
			tui.input.buffer = append(
				tui.input.buffer[:i+1],
				tui.input.buffer[tui.input.cursoroff:]...,
			)
			tui.input.cursoroff = i + 1
			return true
		}
		if !delete && tui.input.buffer[i] != ' ' {
			delete = true
		}
	}
	tui.input.buffer = tui.input.buffer[tui.input.cursoroff:]
	tui.input.cursoroff = 0
	return true
}

func (tui *TUI) handleCmdBackspaceInput(_ rune) bool {
	tui.setCache(paneInput, nil)
	tui.input.buffer = tui.input.buffer[tui.input.cursoroff:]
	tui.input.cursoroff = 0
	return true
}

func (tui *TUI) handleRuneInput(r rune) bool {
	tui.setCache(paneInput, nil)
	tui.input.buffer = append(tui.input.buffer[:tui.input.cursoroff], append(
		[]rune{r}, tui.input.buffer[tui.input.cursoroff:]...,
	)...)
	tui.input.cursoroff++
	return true
}

func (tui *TUI) handleEnterInput(_ rune) bool {
	tui.setCache(paneInput, nil)
	input := tui.input.buffer
	tui.input.inputted = true
	if tui.input.masked {
		tui.input.inputted = false
		tui.input.buffer = []rune{}
		tui.input.cursoroff = 0
	}

	tui.inputs <- []byte(string(input))

	return true
}

func (tui *TUI) handleEscInput(_ rune) bool {
	tui.setCache(paneInput, nil)
	tui.setCache(paneOutput, nil)
	tui.output.offset = 0
	return true
}

func (tui *TUI) handleCtrlCInput(_ rune) bool {
	tui.setCache(paneInput, nil)
	tui.input.buffer = []rune{}
	tui.input.cursoroff = 0
	return true
}

func (tui *TUI) handleLeftInput(_ rune) bool {
	tui.setCache(paneInput, nil)
	tui.input.inputted = false
	tui.input.cursoroff = int(max(0, tui.input.cursoroff-1))
	return true
}

func (tui *TUI) handleRightInput(_ rune) bool {
	tui.setCache(paneInput, nil)
	tui.input.inputted = false
	tui.input.cursoroff = int(min(len(tui.input.buffer), tui.input.cursoroff+1))
	return true
}

func (tui *TUI) handleUpInput(_ rune) bool {
	tui.setCache(paneOutput, nil)
	tui.output.offset++
	return true
}

func (tui *TUI) handleAltUpInput(_ rune) bool {
	tui.setCache(paneOutput, nil)
	tui.output.offset += 5
	return true
}

func (tui *TUI) handleDownInput(_ rune) bool {
	tui.setCache(paneOutput, nil)
	tui.output.offset--
	if tui.output.offset < 0 {
		tui.output.offset = 0
	}
	return true
}

func (tui *TUI) handleAltDownInput(_ rune) bool {
	tui.setCache(paneOutput, nil)
	tui.output.offset -= 5
	if tui.output.offset < 0 {
		tui.output.offset = 0
	}
	return true
}
