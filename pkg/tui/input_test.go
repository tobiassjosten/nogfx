package tui

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
	"github.com/tobiassjosten/nogfx/pkg/mock"
)

func TestRenderInput(t *testing.T) {
	rowToString := func(row Row) (str string) {
		for _, c := range row {
			str += string(c.Content)
		}
		return
	}

	rowsToStrings := func(rows Rows) (strs []string) {
		for _, row := range rows {
			strs = append(strs, rowToString(row))
		}
		return
	}

	tcs := map[string]struct {
		input    *Input
		masked   bool
		unmasked bool
		width    int
		height   int
		rows     []string
		style    tcell.Style
	}{
		"empty": {
			input:  &Input{},
			width:  1,
			height: 1,
			rows:   []string{" "},
		},

		"simple inputting": {
			input: &Input{
				buffer: []rune{'a'},
			},
			width:  1,
			height: 1,
			rows:   []string{"a"},
		},

		"cramped inputting": {
			input: &Input{
				buffer: []rune{'a'},
			},
			width:  0,
			height: 1,
			rows:   nil,
		},

		"masked": {
			input: &Input{
				buffer: []rune{'a'},
			},
			masked: true,
			width:  1,
			height: 1,
			rows:   []string{"*"},
		},

		"unmasked": {
			input: &Input{
				buffer: []rune{'a'},
			},
			masked:   true,
			unmasked: true,
			width:    1,
			height:   1,
			rows:     []string{"a"},
		},

		"padding": {
			input: &Input{
				buffer: []rune{'a'},
			},
			width:  2,
			height: 1,
			rows:   []string{"a "},
		},

		"word wrap": {
			input: &Input{
				buffer: []rune{'a', ' ', 's', 'd', 'f'},
			},
			width:  2,
			height: 1,
			rows:   []string{"a ", "sd", "f "},
		},

		"not inputted": {
			input: &Input{
				buffer: []rune{'a'},
			},
			width:  1,
			height: 1,
			rows:   []string{"a"},
			style: (tcell.Style{}).
				Foreground(tcell.ColorWhite).
				Background(tcell.Color235),
		},

		"inputted": {
			input: &Input{
				buffer:   []rune{'a'},
				inputted: true,
			},
			width:  1,
			height: 1,
			rows:   []string{"a"},
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

			// Dereference so we can repopulate after MaskInput().
			input := *tc.input
			tui.input = &input

			if tc.masked {
				tui.MaskInput()
				tui.input.buffer = tc.input.buffer
				tui.input.cursor = tc.input.cursor
			}

			if tc.unmasked {
				tui.UnmaskInput()
				tui.input.buffer = tc.input.buffer
				tui.input.cursor = tc.input.cursor
			}

			rows := tui.RenderInput(tc.width, tc.height)
			assert.Equal(t, tc.rows, rowsToStrings(rows))

			if tc.style != (tcell.Style{}) {
				assert.Equal(t, tc.style, rows[len(rows)-1][0].Style)
			}
		})
	}
}
