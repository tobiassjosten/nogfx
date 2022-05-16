package tui

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
	"github.com/tobiassjosten/nogfx/pkg/mock"
)

func TestHandleEvent(t *testing.T) {
	tcs := map[string]struct {
		events []*tcell.EventKey
		inputs [][]byte
		f      func(*testing.T, *TUI)
	}{
		"space enables inputting": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, ' ', 0),
			},
			f: func(t *testing.T, ui *TUI) {
				assert.True(t, ui.input.inputting)
			},
		},

		"escape disables inputting": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, ' ', 0),
				tcell.NewEventKey(tcell.KeyEsc, 0, 0),
			},
			f: func(t *testing.T, ui *TUI) {
				assert.False(t, ui.input.inputting)
			},
		},

		"ctrl+c disables inputting": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, ' ', 0),
				tcell.NewEventKey(tcell.KeyCtrlC, 0, 0),
			},
			f: func(t *testing.T, ui *TUI) {
				assert.False(t, ui.input.inputting)
			},
		},

		"inputting moves cursor right": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, ' ', 0),
				tcell.NewEventKey(tcell.KeyRune, 'x', 0),
				tcell.NewEventKey(tcell.KeyRune, 'y', 0),
			},
			f: func(t *testing.T, ui *TUI) {
				assert.Equal(t, 2, ui.input.cursor)
			},
		},

		"left arrow moves cursor left": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, ' ', 0),
				tcell.NewEventKey(tcell.KeyRune, 'x', 0),
				tcell.NewEventKey(tcell.KeyRune, 'y', 0),
				tcell.NewEventKey(tcell.KeyLeft, 0, 0),
			},
			f: func(t *testing.T, ui *TUI) {
				assert.Equal(t, 1, ui.input.cursor)
			},
		},

		"left arrow stops at buffer end": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, ' ', 0),
				tcell.NewEventKey(tcell.KeyRune, 'x', 0),
				tcell.NewEventKey(tcell.KeyRune, 'y', 0),
				tcell.NewEventKey(tcell.KeyLeft, 0, 0),
				tcell.NewEventKey(tcell.KeyLeft, 0, 0),
				tcell.NewEventKey(tcell.KeyLeft, 0, 0),
			},
			f: func(t *testing.T, ui *TUI) {
				assert.Equal(t, 0, ui.input.cursor)
			},
		},

		"right arrow moves cursor right": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, ' ', 0),
				tcell.NewEventKey(tcell.KeyRune, 'x', 0),
				tcell.NewEventKey(tcell.KeyRune, 'y', 0),
				tcell.NewEventKey(tcell.KeyLeft, 0, 0),
				tcell.NewEventKey(tcell.KeyRight, 0, 0),
			},
			f: func(t *testing.T, ui *TUI) {
				assert.Equal(t, 2, ui.input.cursor)
			},
		},

		"right arrow stops at buffer end": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, ' ', 0),
				tcell.NewEventKey(tcell.KeyRune, 'x', 0),
				tcell.NewEventKey(tcell.KeyRune, 'y', 0),
				tcell.NewEventKey(tcell.KeyLeft, 0, 0),
				tcell.NewEventKey(tcell.KeyRight, 0, 0),
				tcell.NewEventKey(tcell.KeyRight, 0, 0),
			},
			f: func(t *testing.T, ui *TUI) {
				assert.Equal(t, 2, ui.input.cursor)
			},
		},

		"backspace deletes character from buffer": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, ' ', 0),
				tcell.NewEventKey(tcell.KeyRune, 'a', 0),
				tcell.NewEventKey(tcell.KeyRune, 's', 0),
				tcell.NewEventKey(tcell.KeyRune, 'd', 0),
				tcell.NewEventKey(tcell.KeyRune, 'f', 0),
				tcell.NewEventKey(tcell.KeyBackspace, 0, 0),
			},
			f: func(t *testing.T, ui *TUI) {
				assert.True(t, ui.input.inputting)
				assert.Equal(t, []rune("asd"), ui.input.buffer)
			},
		},

		"backspace deletes based on cursor": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, ' ', 0),
				tcell.NewEventKey(tcell.KeyRune, 'a', 0),
				tcell.NewEventKey(tcell.KeyRune, 's', 0),
				tcell.NewEventKey(tcell.KeyRune, 'd', 0),
				tcell.NewEventKey(tcell.KeyRune, 'f', 0),
				tcell.NewEventKey(tcell.KeyLeft, 0, 0),
				tcell.NewEventKey(tcell.KeyBackspace, 0, 0),
			},
			f: func(t *testing.T, ui *TUI) {
				assert.True(t, ui.input.inputting)
				assert.Equal(t, []rune("asf"), ui.input.buffer)
			},
		},

		"opt + backspace deletes word from buffer": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, ' ', 0),
				tcell.NewEventKey(tcell.KeyRune, 'a', 0),
				tcell.NewEventKey(tcell.KeyRune, 's', 0),
				tcell.NewEventKey(tcell.KeyRune, ' ', 0),
				tcell.NewEventKey(tcell.KeyRune, 'd', 0),
				tcell.NewEventKey(tcell.KeyRune, 'f', 0),
				tcell.NewEventKey(tcell.KeyETB, 0, 0),
			},
			f: func(t *testing.T, ui *TUI) {
				assert.True(t, ui.input.inputting)
				assert.Equal(t, []rune("as "), ui.input.buffer)
			},
		},

		"opt + backspace deletes based on cursor": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, ' ', 0),
				tcell.NewEventKey(tcell.KeyRune, 'a', 0),
				tcell.NewEventKey(tcell.KeyRune, 's', 0),
				tcell.NewEventKey(tcell.KeyRune, ' ', 0),
				tcell.NewEventKey(tcell.KeyRune, 'd', 0),
				tcell.NewEventKey(tcell.KeyRune, 'f', 0),
				tcell.NewEventKey(tcell.KeyLeft, 0, 0),
				tcell.NewEventKey(tcell.KeyETB, 0, 0),
			},
			f: func(t *testing.T, ui *TUI) {
				assert.True(t, ui.input.inputting)
				assert.Equal(t, []rune("as f"), ui.input.buffer)
			},
		},

		"opt + backspace deletes last word from buffer": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, ' ', 0),
				tcell.NewEventKey(tcell.KeyRune, 'a', 0),
				tcell.NewEventKey(tcell.KeyRune, 's', 0),
				tcell.NewEventKey(tcell.KeyRune, 'd', 0),
				tcell.NewEventKey(tcell.KeyRune, 'f', 0),
				tcell.NewEventKey(tcell.KeyETB, 0, 0),
			},
			f: func(t *testing.T, ui *TUI) {
				assert.True(t, ui.input.inputting)
				assert.Equal(t, []rune{}, ui.input.buffer)
			},
		},

		"cmd + backspace deletes everything from buffer": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, ' ', 0),
				tcell.NewEventKey(tcell.KeyRune, 'a', 0),
				tcell.NewEventKey(tcell.KeyRune, 's', 0),
				tcell.NewEventKey(tcell.KeyRune, ' ', 0),
				tcell.NewEventKey(tcell.KeyRune, 'd', 0),
				tcell.NewEventKey(tcell.KeyRune, 'f', 0),
				tcell.NewEventKey(tcell.KeyNAK, 0, 0),
			},
			f: func(t *testing.T, ui *TUI) {
				assert.True(t, ui.input.inputting)
				assert.Equal(t, []rune{}, ui.input.buffer)
			},
		},

		"cmd + backspace deletes based on cursor": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, ' ', 0),
				tcell.NewEventKey(tcell.KeyRune, 'a', 0),
				tcell.NewEventKey(tcell.KeyRune, 's', 0),
				tcell.NewEventKey(tcell.KeyRune, ' ', 0),
				tcell.NewEventKey(tcell.KeyRune, 'd', 0),
				tcell.NewEventKey(tcell.KeyRune, 'f', 0),
				tcell.NewEventKey(tcell.KeyLeft, 0, 0),
				tcell.NewEventKey(tcell.KeyNAK, 0, 0),
			},
			f: func(t *testing.T, ui *TUI) {
				assert.True(t, ui.input.inputting)
				assert.Equal(t, []rune("f"), ui.input.buffer)
			},
		},

		"enter sends the buffer": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, ' ', 0),
				tcell.NewEventKey(tcell.KeyRune, 'a', 0),
				tcell.NewEventKey(tcell.KeyRune, 's', 0),
				tcell.NewEventKey(tcell.KeyRune, 'd', 0),
				tcell.NewEventKey(tcell.KeyRune, 'f', 0),
				tcell.NewEventKey(tcell.KeyEnter, 0, 0),
			},
			inputs: [][]byte{
				[]byte("asdf"),
			},
			f: func(t *testing.T, ui *TUI) {
				assert.True(t, ui.input.inputted)
				assert.Equal(t, []rune("asdf"), ui.input.buffer)
			},
		},

		"backspace resets inputted": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, ' ', 0),
				tcell.NewEventKey(tcell.KeyRune, 'a', 0),
				tcell.NewEventKey(tcell.KeyRune, 's', 0),
				tcell.NewEventKey(tcell.KeyRune, 'd', 0),
				tcell.NewEventKey(tcell.KeyRune, 'f', 0),
				tcell.NewEventKey(tcell.KeyEnter, 0, 0),
				tcell.NewEventKey(tcell.KeyBackspace, 0, 0),
			},
			f: func(t *testing.T, ui *TUI) {
				assert.False(t, ui.input.inputted)
				assert.Equal(t, []rune{}, ui.input.buffer)
			},
		},

		"new input resets inputted": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, ' ', 0),
				tcell.NewEventKey(tcell.KeyRune, 'a', 0),
				tcell.NewEventKey(tcell.KeyRune, 's', 0),
				tcell.NewEventKey(tcell.KeyRune, 'd', 0),
				tcell.NewEventKey(tcell.KeyRune, 'f', 0),
				tcell.NewEventKey(tcell.KeyEnter, 0, 0),
				tcell.NewEventKey(tcell.KeyRune, 'a', 0),
			},
			f: func(t *testing.T, ui *TUI) {
				assert.False(t, ui.input.inputted)
				assert.Equal(t, []rune("a"), ui.input.buffer)
			},
		},

		"keypad 1 sends sw": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, '1', 0),
			},
			inputs: [][]byte{
				[]byte("sw"),
			},
		},

		"keypad 2 sends s": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, '2', 0),
			},
			inputs: [][]byte{
				[]byte("s"),
			},
		},

		"keypad 3 sends se": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, '3', 0),
			},
			inputs: [][]byte{
				[]byte("se"),
			},
		},

		"keypad 4 sends w": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, '4', 0),
			},
			inputs: [][]byte{
				[]byte("w"),
			},
		},

		"keypad 5 sends map": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, '5', 0),
			},
			inputs: [][]byte{
				[]byte("map"),
			},
		},

		"keypad 6 sends e": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, '6', 0),
			},
			inputs: [][]byte{
				[]byte("e"),
			},
		},

		"keypad 7 sends nw": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, '7', 0),
			},
			inputs: [][]byte{
				[]byte("nw"),
			},
		},

		"keypad 8 sends n": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, '8', 0),
			},
			inputs: [][]byte{
				[]byte("n"),
			},
		},

		"keypad 9 sends ne": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, '9', 0),
			},
			inputs: [][]byte{
				[]byte("ne"),
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			screen := &mock.ScreenMock{
				SetCursorStyleFunc: func(_ tcell.CursorStyle) {},
				SetStyleFunc:       func(_ tcell.Style) {},
			}

			ui := NewTUI(screen)

			done := make(chan struct{})

			var inputs [][]byte
			go func() {
				for input := range ui.inputs {
					inputs = append(inputs, input)
				}
				done <- struct{}{}
			}()

			for _, event := range tc.events {
				_ = ui.HandleEvent(event)
			}
			close(ui.inputs)

			<-done

			if len(tc.inputs) > 0 {
				assert.Equal(t, tc.inputs, inputs)
			}

			if tc.f != nil {
				tc.f(t, ui)
			}
		})
	}
}
