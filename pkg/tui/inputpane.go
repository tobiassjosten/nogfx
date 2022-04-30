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

var (
	inputEvents = map[int]func(rune) (bool, []rune){}
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

	inputEvents = map[int]func(rune) (bool, []rune){
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
		if f, ok := inputEvents[alt]; ok {
			if handled, rs := f(event.Rune()); handled {
				return true, rs
			}
		}
	}

	return false, nil
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
