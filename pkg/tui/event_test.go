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
		"inputting appends buffer and moves cursor": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, 'x', 0),
				tcell.NewEventKey(tcell.KeyRune, 'y', 0),
			},
			f: func(t *testing.T, ui *TUI) {
				assert.Equal(t, []rune("xy"), ui.input.buffer)
				assert.Equal(t, 2, ui.input.cursoroff)
			},
		},

		"ctrl+c clears buffer": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, 'x', 0),
				tcell.NewEventKey(tcell.KeyRune, 'y', 0),
				tcell.NewEventKey(tcell.KeyCtrlC, 0, 0),
			},
			f: func(t *testing.T, ui *TUI) {
				assert.Equal(t, []rune(""), ui.input.buffer)
				assert.Equal(t, 0, ui.input.cursoroff)
			},
		},

		"left arrow moves cursor left": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, 'x', 0),
				tcell.NewEventKey(tcell.KeyRune, 'y', 0),
				tcell.NewEventKey(tcell.KeyLeft, 0, 0),
			},
			f: func(t *testing.T, ui *TUI) {
				assert.Equal(t, 1, ui.input.cursoroff)
			},
		},

		"left arrow stops at buffer end": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, 'x', 0),
				tcell.NewEventKey(tcell.KeyRune, 'y', 0),
				tcell.NewEventKey(tcell.KeyLeft, 0, 0),
				tcell.NewEventKey(tcell.KeyLeft, 0, 0),
				tcell.NewEventKey(tcell.KeyLeft, 0, 0),
			},
			f: func(t *testing.T, ui *TUI) {
				assert.Equal(t, 0, ui.input.cursoroff)
			},
		},

		"right arrow moves cursor right": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, 'x', 0),
				tcell.NewEventKey(tcell.KeyRune, 'y', 0),
				tcell.NewEventKey(tcell.KeyLeft, 0, 0),
				tcell.NewEventKey(tcell.KeyRight, 0, 0),
			},
			f: func(t *testing.T, ui *TUI) {
				assert.Equal(t, 2, ui.input.cursoroff)
			},
		},

		"right arrow stops at buffer end": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, 'x', 0),
				tcell.NewEventKey(tcell.KeyRune, 'y', 0),
				tcell.NewEventKey(tcell.KeyLeft, 0, 0),
				tcell.NewEventKey(tcell.KeyRight, 0, 0),
				tcell.NewEventKey(tcell.KeyRight, 0, 0),
			},
			f: func(t *testing.T, ui *TUI) {
				assert.Equal(t, 2, ui.input.cursoroff)
			},
		},

		"backspace deletes character from buffer": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, 'a', 0),
				tcell.NewEventKey(tcell.KeyRune, 's', 0),
				tcell.NewEventKey(tcell.KeyRune, 'd', 0),
				tcell.NewEventKey(tcell.KeyRune, 'f', 0),
				tcell.NewEventKey(tcell.KeyBackspace, 0, 0),
			},
			f: func(t *testing.T, ui *TUI) {
				assert.Equal(t, []rune("asd"), ui.input.buffer)
			},
		},

		"backspace deletes based on cursor": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, 'a', 0),
				tcell.NewEventKey(tcell.KeyRune, 's', 0),
				tcell.NewEventKey(tcell.KeyRune, 'd', 0),
				tcell.NewEventKey(tcell.KeyRune, 'f', 0),
				tcell.NewEventKey(tcell.KeyLeft, 0, 0),
				tcell.NewEventKey(tcell.KeyBackspace, 0, 0),
			},
			f: func(t *testing.T, ui *TUI) {
				assert.Equal(t, []rune("asf"), ui.input.buffer)
			},
		},

		"opt + backspace deletes word from buffer": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, 'a', 0),
				tcell.NewEventKey(tcell.KeyRune, 's', 0),
				tcell.NewEventKey(tcell.KeyRune, ' ', 0),
				tcell.NewEventKey(tcell.KeyRune, 'd', 0),
				tcell.NewEventKey(tcell.KeyRune, 'f', 0),
				tcell.NewEventKey(tcell.KeyETB, 0, 0),
			},
			f: func(t *testing.T, ui *TUI) {
				assert.Equal(t, []rune("as "), ui.input.buffer)
			},
		},

		"opt + backspace deletes based on cursor": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, 'a', 0),
				tcell.NewEventKey(tcell.KeyRune, 's', 0),
				tcell.NewEventKey(tcell.KeyRune, ' ', 0),
				tcell.NewEventKey(tcell.KeyRune, 'd', 0),
				tcell.NewEventKey(tcell.KeyRune, 'f', 0),
				tcell.NewEventKey(tcell.KeyLeft, 0, 0),
				tcell.NewEventKey(tcell.KeyETB, 0, 0),
			},
			f: func(t *testing.T, ui *TUI) {
				assert.Equal(t, []rune("as f"), ui.input.buffer)
			},
		},

		"opt + backspace deletes last word from buffer": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, 'a', 0),
				tcell.NewEventKey(tcell.KeyRune, 's', 0),
				tcell.NewEventKey(tcell.KeyRune, 'd', 0),
				tcell.NewEventKey(tcell.KeyRune, 'f', 0),
				tcell.NewEventKey(tcell.KeyETB, 0, 0),
			},
			f: func(t *testing.T, ui *TUI) {
				assert.Equal(t, []rune{}, ui.input.buffer)
			},
		},

		"cmd + backspace deletes everything from buffer": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, 'a', 0),
				tcell.NewEventKey(tcell.KeyRune, 's', 0),
				tcell.NewEventKey(tcell.KeyRune, ' ', 0),
				tcell.NewEventKey(tcell.KeyRune, 'd', 0),
				tcell.NewEventKey(tcell.KeyRune, 'f', 0),
				tcell.NewEventKey(tcell.KeyNAK, 0, 0),
			},
			f: func(t *testing.T, ui *TUI) {
				assert.Equal(t, []rune{}, ui.input.buffer)
			},
		},

		"cmd + backspace deletes based on cursor": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, 'a', 0),
				tcell.NewEventKey(tcell.KeyRune, 's', 0),
				tcell.NewEventKey(tcell.KeyRune, ' ', 0),
				tcell.NewEventKey(tcell.KeyRune, 'd', 0),
				tcell.NewEventKey(tcell.KeyRune, 'f', 0),
				tcell.NewEventKey(tcell.KeyLeft, 0, 0),
				tcell.NewEventKey(tcell.KeyNAK, 0, 0),
			},
			f: func(t *testing.T, ui *TUI) {
				assert.Equal(t, []rune("f"), ui.input.buffer)
			},
		},

		"enter sends the buffer": {
			events: []*tcell.EventKey{
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

		// @todo Test scrollback (up/down/escape).
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
