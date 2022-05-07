package tui

import (
	"github.com/gdamore/tcell/v2"
)

func (pane *InputPane) handlers() map[int]func(rune) (bool, []rune) {
	return map[int]func(rune) (bool, []rune){
		int(tcell.KeyRune) + int(' ') + normalMode: pane.handleSpaceNormal,
		int(tcell.KeyRune) + normalMode:            pane.handleRuneNormal,

		int(tcell.KeyRune) + inputtedMode: pane.handleRuneInputted,

		int(tcell.KeyBackspace) + inputtedMode:  pane.handleBackspaceInputted,
		int(tcell.KeyBackspace2) + inputtedMode: pane.handleBackspaceInputted,
		int(tcell.KeyETB) + inputtedMode:        pane.handleBackspaceInputted,
		int(tcell.KeyNAK) + inputtedMode:        pane.handleBackspaceInputted,

		int(tcell.KeyRune) + inputMode:  pane.handleRuneInput,
		int(tcell.KeyEnter) + inputMode: pane.handleEnterInput,
		int(tcell.KeyEsc) + inputMode:   pane.handleEscInput,
		int(tcell.KeyCtrlC) + inputMode: pane.handleCtrlCInput,
		int(tcell.KeyLeft) + inputMode:  pane.handleLeftInput,
		int(tcell.KeyRight) + inputMode: pane.handleRightInput,

		int(tcell.KeyBackspace) + inputMode:  pane.handleBackspaceInput,
		int(tcell.KeyBackspace2) + inputMode: pane.handleBackspaceInput,
		int(tcell.KeyETB) + inputMode:        pane.handleOptBackspaceInput,
		int(tcell.KeyNAK) + inputMode:        pane.handleCmdBackspaceInput,
	}
}

// HandleEvent reacts on a user event and modifies itself from it.
func (pane *InputPane) HandleEvent(event *tcell.EventKey) (bool, []rune) {
	alts := []int{}

	if pane.inputting {
		if pane.inputted {
			alts = append(alts, int(event.Key())+int(event.Rune())+inputtedMode)
			alts = append(alts, int(event.Key())+inputtedMode)
		}
		alts = append(alts, int(event.Key())+int(event.Rune())+inputMode)
		alts = append(alts, int(event.Key())+inputMode)
	} else {
		alts = append(alts, int(event.Key())+int(event.Rune())+normalMode)
		alts = append(alts, int(event.Key())+normalMode)
	}
	alts = append(alts, int(event.Key())+int(event.Rune()))
	alts = append(alts, int(event.Key()))

	for _, alt := range alts {
		if f, ok := pane.handlers()[alt]; ok {
			if handled, rs := f(event.Rune()); handled {
				return true, rs
			}
		}
	}

	return false, nil
}

func (pane *InputPane) handleSpaceNormal(_ rune) (bool, []rune) {
	pane.inputting = true
	return true, nil
}

func (pane *InputPane) handleRuneInputted(r rune) (bool, []rune) {
	_, _ = pane.handleBackspaceInputted(r)
	return false, nil
}

func (pane *InputPane) handleBackspaceInputted(_ rune) (bool, []rune) {
	pane.input = []rune{}
	pane.cursor = 0
	pane.inputted = false
	return true, nil
}

func (pane *InputPane) handleBackspaceInput(_ rune) (bool, []rune) {
	cursor := int(max(0, pane.cursor-1))
	pane.input = append(
		pane.input[:cursor],
		pane.input[pane.cursor:]...,
	)
	pane.cursor = cursor
	return true, nil
}

func (pane *InputPane) handleOptBackspaceInput(_ rune) (bool, []rune) {
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
}

func (pane *InputPane) handleCmdBackspaceInput(_ rune) (bool, []rune) {
	pane.input = pane.input[pane.cursor:]
	pane.cursor = 0
	return true, nil
}

func (pane *InputPane) handleRuneInput(r rune) (bool, []rune) {
	pane.input = append(pane.input[:pane.cursor], append(
		[]rune{r}, pane.input[pane.cursor:]...,
	)...)
	pane.cursor++
	return true, nil
}

func (pane *InputPane) handleEnterInput(_ rune) (bool, []rune) {
	input := pane.input
	pane.inputted = true
	if pane.masked {
		pane.inputted = false
		pane.input = []rune{}
		pane.cursor = 0
	}
	return true, input
}

func (pane *InputPane) handleEscInput(_ rune) (bool, []rune) {
	pane.inputting = false
	return true, nil
}

func (pane *InputPane) handleCtrlCInput(_ rune) (bool, []rune) {
	pane.input = []rune{}
	pane.cursor = 0
	return pane.handleEscInput(0)
}

func (pane *InputPane) handleLeftInput(_ rune) (bool, []rune) {
	pane.inputted = false
	pane.cursor = int(max(0, pane.cursor-1))
	return true, nil
}

func (pane *InputPane) handleRightInput(_ rune) (bool, []rune) {
	pane.inputted = false
	pane.cursor = int(min(len(pane.input), pane.cursor+1))
	return true, nil
}

func (pane *InputPane) handleRuneNormal(r rune) (bool, []rune) {
	switch r {
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
