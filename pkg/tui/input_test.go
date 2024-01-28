package tui

import (
	"testing"

	"github.com/tobiassjosten/nogfx/pkg/mock"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
)

func TestRenderInput(t *testing.T) {
	tcs := map[string]struct {
		input     *Input
		width     int
		height    int
		rows      []string
		cursorpos []int
		style     tcell.Style
	}{
		"empty": {
			input:     &Input{},
			width:     1,
			height:    1,
			rows:      []string{" "},
			cursorpos: []int{0, 0},
		},

		"simple inputting": {
			input: &Input{
				buffer:    []rune{'a'},
				cursoroff: 1,
			},
			width:     2,
			height:    1,
			rows:      []string{"a "},
			cursorpos: []int{1, 0},
		},

		"cursor back one": {
			input: &Input{
				buffer:    []rune{'a'},
				cursoroff: 0,
			},
			width:     2,
			height:    1,
			rows:      []string{"a "},
			cursorpos: []int{0, 0},
		},

		"cursor back two": {
			input: &Input{
				buffer:    []rune("asdf"),
				cursoroff: 2,
			},
			width:     5,
			height:    1,
			rows:      []string{"asdf "},
			cursorpos: []int{2, 0},
		},

		"no width": {
			input: &Input{
				buffer:    []rune{'a'},
				cursoroff: 1,
			},
			width:     0,
			height:    1,
			rows:      nil,
			cursorpos: []int{0, 0},
		},

		"no height": {
			input: &Input{
				buffer:    []rune{'a'},
				cursoroff: 1,
			},
			width:     1,
			height:    0,
			rows:      nil,
			cursorpos: []int{0, 0},
		},

		"padding": {
			input: &Input{
				buffer:    []rune{'a'},
				cursoroff: 1,
			},
			width:     2,
			height:    1,
			rows:      []string{"a "},
			cursorpos: []int{1, 0},
		},

		"word wrap": {
			input: &Input{
				buffer:    []rune{'a', ' ', 's', 'd', 'f'},
				cursoroff: 5,
			},
			width:     2,
			height:    3,
			rows:      []string{"a ", "sd", "f "},
			cursorpos: []int{1, 2},
		},

		"line wrap": {
			input: &Input{
				buffer:    []rune{'a'},
				cursoroff: 1,
			},
			width:     1,
			height:    2,
			rows:      []string{"a", " "},
			cursorpos: []int{0, 1},
		},

		"not inputted": {
			input: &Input{
				buffer:    []rune{'a'},
				cursoroff: 1,
			},
			width:     2,
			height:    1,
			rows:      []string{"a "},
			cursorpos: []int{1, 0},
			style: (tcell.Style{}).
				Foreground(tcell.ColorWhite).
				Background(tcell.Color235),
		},

		"inputted": {
			input: &Input{
				buffer:    []rune{'a'},
				cursoroff: 1,
				inputted:  true,
			},
			width:     2,
			height:    1,
			rows:      []string{"a "},
			cursorpos: []int{1, 0},
			style: (tcell.Style{}).
				Foreground(tcell.ColorWhite).
				Background(tcell.Color235).
				Attributes(tcell.AttrDim),
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			tui := NewTUI(&mock.ScreenMock{
				HideCursorFunc:     func() {},
				SetCursorStyleFunc: func(_ tcell.CursorStyle) {},
				SetStyleFunc:       func(_ tcell.Style) {},
			})

			tui.input = tc.input

			rows, cx, cy := tui.RenderInput(tc.width, tc.height)

			assert.Equal(t, tc.rows, rows.Strings())
			if tc.cursorpos != nil {
				assert.Equal(t, tc.cursorpos, []int{cx, cy})
			}

			if tc.style != (tcell.Style{}) {
				assert.Equal(t, tc.style, rows[len(rows)-1][0].Style)
			}
		})
	}
}

func TestMasking(t *testing.T) {
	width, height := 5, 1
	tui := NewTUI(&mock.ScreenMock{
		HideCursorFunc:     func() {},
		SetCursorStyleFunc: func(_ tcell.CursorStyle) {},
		SetStyleFunc:       func(_ tcell.Style) {},
	})

	tui.input = &Input{
		buffer: []rune("asdf"),
	}

	tui.MaskInput()

	rows, _, _ := tui.RenderInput(width, height)

	// Masking input deletes the buffer.
	assert.Equal(t, []string{"     "}, rows.Strings())

	tui.setCache(paneInput, nil)
	tui.input.buffer = []rune("asdf")

	rows, _, _ = tui.RenderInput(width, height)

	// Masked rendering transforms the buffer.
	assert.Equal(t, []string{"**** "}, rows.Strings())

	tui.UnmaskInput()

	rows, _, _ = tui.RenderInput(width, height)

	// Unmasking input deletes the buffer.
	assert.Equal(t, []string{"     "}, rows.Strings())
}
