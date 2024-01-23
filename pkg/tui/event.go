package tui

import (
	"github.com/gdamore/tcell/v2"
)

// @todo Map more keys, like 271 (delete), 268 (home), 269 (end), 267 (pgdn),
// 266 (pgup), and alt/cmd+left/right.

func (tui *TUI) eventHandlers() map[int]func(rune) bool {
	alt := int(tcell.ModAlt)

	return map[int]func(rune) bool{
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

		int(keyNum1): tui.handleNum1,
		int(keyNum2): tui.handleNum2,
		int(keyNum3): tui.handleNum3,
		int(keyNum4): tui.handleNum4,
		int(keyNum6): tui.handleNum6,
		int(keyNum7): tui.handleNum7,
		int(keyNum8): tui.handleNum8,
		int(keyNum9): tui.handleNum9,
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

func (tui *TUI) handleBackspaceInput(_ rune) bool {
	if tui.input.inputted {
		return tui.handleCtrlCInput(0)
	}

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
	if tui.input.inputted {
		return tui.handleCtrlCInput(0)
	}

	tui.setCache(paneInput, nil)
	del := false
	for i := tui.input.cursoroff - 1; i > 0; i-- {
		if del && tui.input.buffer[i] == ' ' {
			tui.input.buffer = append(
				tui.input.buffer[:i+1],
				tui.input.buffer[tui.input.cursoroff:]...,
			)
			tui.input.cursoroff = i + 1
			return true
		}
		if !del && tui.input.buffer[i] != ' ' {
			del = true
		}
	}
	tui.input.buffer = tui.input.buffer[tui.input.cursoroff:]
	tui.input.cursoroff = 0
	return true
}

func (tui *TUI) handleCmdBackspaceInput(_ rune) bool {
	if tui.input.inputted {
		return tui.handleCtrlCInput(0)
	}

	tui.setCache(paneInput, nil)
	tui.input.buffer = tui.input.buffer[tui.input.cursoroff:]
	tui.input.cursoroff = 0
	return true
}

func (tui *TUI) handleRuneInput(r rune) bool {
	if tui.input.inputted {
		_ = tui.handleCtrlCInput(0)
	}

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
		tui.input.buffer = []rune{}
		tui.input.cursoroff = 0
		tui.input.inputted = false
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
	tui.input.inputted = false
	return true
}

func (tui *TUI) handleLeftInput(_ rune) bool {
	tui.setCache(paneInput, nil)
	tui.input.cursoroff = int(max(0, tui.input.cursoroff-1))
	tui.input.inputted = false
	return true
}

func (tui *TUI) handleRightInput(_ rune) bool {
	tui.setCache(paneInput, nil)
	tui.input.cursoroff = int(min(len(tui.input.buffer), tui.input.cursoroff+1))
	tui.input.inputted = false
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

func (tui *TUI) handleNum1(r rune) bool {
	tui.inputs <- []byte{'s', 'w'}
	return true
}

func (tui *TUI) handleNum2(r rune) bool {
	tui.inputs <- []byte{'s'}
	return true
}

func (tui *TUI) handleNum3(r rune) bool {
	tui.inputs <- []byte{'s', 'e'}
	return true
}

func (tui *TUI) handleNum4(r rune) bool {
	tui.inputs <- []byte{'w'}
	return true
}

func (tui *TUI) handleNum6(r rune) bool {
	tui.inputs <- []byte{'e'}
	return true
}

func (tui *TUI) handleNum7(r rune) bool {
	tui.inputs <- []byte{'n', 'w'}
	return true
}

func (tui *TUI) handleNum8(r rune) bool {
	tui.inputs <- []byte{'n'}
	return true
}

func (tui *TUI) handleNum9(r rune) bool {
	tui.inputs <- []byte{'n', 'e'}
	return true
}
