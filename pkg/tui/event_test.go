package tui

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
	"github.com/tobiassjosten/nogfx/pkg/mock"
)

func TestHandleEvent(t *testing.T) {
	bools := func(c int, b bool) (bs []bool) {
		for i := 0; i < c; i++ {
			bs = append(bs, b)
		}
		return
	}

	tcs := map[string]struct {
		events  []*tcell.EventKey
		setup   func(*TUI)
		inputs  [][]byte
		handled []bool
		test    func(*assert.Assertions, *TUI)
	}{
		"inputting appends buffer and moves cursor": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, 'x', 0),
				tcell.NewEventKey(tcell.KeyRune, 'y', 0),
			},
			test: func(a *assert.Assertions, ui *TUI) {
				a.Equal([]rune("xy"), ui.input.buffer)
				a.Equal(2, ui.input.cursoroff)
			},
		},

		"ctrl+c clears buffer": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, 'x', 0),
				tcell.NewEventKey(tcell.KeyRune, 'y', 0),
				tcell.NewEventKey(tcell.KeyCtrlC, 0, 0),
			},
			test: func(a *assert.Assertions, ui *TUI) {
				a.Equal([]rune(""), ui.input.buffer)
				a.Equal(0, ui.input.cursoroff)
			},
		},

		"left arrow moves cursor left": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, 'x', 0),
				tcell.NewEventKey(tcell.KeyRune, 'y', 0),
				tcell.NewEventKey(tcell.KeyLeft, 0, 0),
			},
			test: func(a *assert.Assertions, ui *TUI) {
				a.Equal(1, ui.input.cursoroff)
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
			test: func(a *assert.Assertions, ui *TUI) {
				a.Equal(0, ui.input.cursoroff)
			},
		},

		"right arrow moves cursor right": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, 'x', 0),
				tcell.NewEventKey(tcell.KeyRune, 'y', 0),
				tcell.NewEventKey(tcell.KeyLeft, 0, 0),
				tcell.NewEventKey(tcell.KeyRight, 0, 0),
			},
			test: func(a *assert.Assertions, ui *TUI) {
				a.Equal(2, ui.input.cursoroff)
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
			test: func(a *assert.Assertions, ui *TUI) {
				a.Equal(2, ui.input.cursoroff)
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
			test: func(a *assert.Assertions, ui *TUI) {
				a.Equal([]rune("asd"), ui.input.buffer)
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
			test: func(a *assert.Assertions, ui *TUI) {
				a.Equal([]rune("asf"), ui.input.buffer)
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
			test: func(a *assert.Assertions, ui *TUI) {
				a.Equal([]rune("as "), ui.input.buffer)
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
			test: func(a *assert.Assertions, ui *TUI) {
				a.Equal([]rune("as f"), ui.input.buffer)
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
			test: func(a *assert.Assertions, ui *TUI) {
				a.Equal([]rune{}, ui.input.buffer)
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
			test: func(a *assert.Assertions, ui *TUI) {
				a.Equal([]rune{}, ui.input.buffer)
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
			test: func(a *assert.Assertions, ui *TUI) {
				a.Equal([]rune("f"), ui.input.buffer)
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
			test: func(a *assert.Assertions, ui *TUI) {
				a.True(ui.input.inputted)
				a.Equal([]rune("asdf"), ui.input.buffer)
			},
		},

		"enter when masked clears the buffer": {
			setup: func(ui *TUI) {
				ui.MaskInput()
			},
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, 'a', 0),
				tcell.NewEventKey(tcell.KeyRune, 's', 0),
				tcell.NewEventKey(tcell.KeyRune, 'd', 0),
				tcell.NewEventKey(tcell.KeyRune, 'f', 0),
				tcell.NewEventKey(tcell.KeyEnter, 0, 0),
			},
			test: func(a *assert.Assertions, ui *TUI) {
				a.Empty(ui.input.buffer)
				a.False(ui.input.inputted)
				a.Equal(0, ui.input.cursoroff)
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
			test: func(a *assert.Assertions, ui *TUI) {
				a.False(ui.input.inputted)
				a.Equal([]rune{}, ui.input.buffer)
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
			test: func(a *assert.Assertions, ui *TUI) {
				a.False(ui.input.inputted)
				a.Equal([]rune("a"), ui.input.buffer)
			},
		},

		"up arrow scrolls back output": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyUp, 0, 0),
			},
			test: func(a *assert.Assertions, ui *TUI) {
				a.Equal(1, ui.output.offset)
			},
		},

		"up up arrow scrolls back output": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyUp, 0, 0),
				tcell.NewEventKey(tcell.KeyUp, 0, 0),
			},
			test: func(a *assert.Assertions, ui *TUI) {
				a.Equal(2, ui.output.offset)
			},
		},

		"alt up arrow scrolls back output more": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModAlt),
			},
			test: func(a *assert.Assertions, ui *TUI) {
				a.Equal(5, ui.output.offset)
			},
		},

		"alt up arrow scrolls back when inputted": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyRune, 'a', 0),
				tcell.NewEventKey(tcell.KeyEnter, 0, 0),
				tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModAlt),
			},
			test: func(a *assert.Assertions, ui *TUI) {
				a.Equal(5, ui.output.offset)
			},
		},

		"down arrow scrolls forward output": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyUp, 0, 0),
				tcell.NewEventKey(tcell.KeyUp, 0, 0),
				tcell.NewEventKey(tcell.KeyDown, 0, 0),
			},
			test: func(a *assert.Assertions, ui *TUI) {
				a.Equal(1, ui.output.offset)
			},
		},

		"down arrow does nothing without scrollback": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyDown, 0, 0),
			},
			test: func(a *assert.Assertions, ui *TUI) {
				a.Equal(0, ui.output.offset)
			},
		},

		"alt down arrow scrolls back output more": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModAlt),
				tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModAlt),
				tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModAlt),
			},
			test: func(a *assert.Assertions, ui *TUI) {
				a.Equal(5, ui.output.offset)
			},
		},

		"alt down arrow does nothing without scrollback": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModAlt),
			},
			test: func(a *assert.Assertions, ui *TUI) {
				a.Equal(0, ui.output.offset)
			},
		},

		"escape resets scrollback": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyUp, 0, 0),
				tcell.NewEventKey(tcell.KeyEsc, 0, 0),
			},
			test: func(a *assert.Assertions, ui *TUI) {
				a.Equal(0, ui.output.offset)
			},
		},

		"unknown keys dont do anything": {
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyCtrlP, 0, 0),
			},
			handled: []bool{false},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			screen := &mock.ScreenMock{
				SetCursorStyleFunc: func(_ tcell.CursorStyle) {},
				SetStyleFunc:       func(_ tcell.Style) {},
			}

			ui := NewTUI(screen)

			if tc.setup != nil {
				tc.setup(ui)
			}

			done := make(chan struct{})

			var inputs [][]byte
			go func() {
				for input := range ui.inputs {
					inputs = append(inputs, input)
				}
				done <- struct{}{}
			}()

			handled := []bool{}
			for _, event := range tc.events {
				handled = append(handled, ui.HandleEvent(event))
			}
			close(ui.inputs)

			<-done

			if tc.inputs != nil {
				assert.Equal(t, tc.inputs, inputs)
			}

			if tc.handled == nil {
				tc.handled = bools(len(tc.events), true)
			}
			assert.Equal(t, tc.handled, handled)

			if tc.test != nil {
				tc.test(assert.New(t), ui)
			}
		})
	}
}
