package tui

import (
	"testing"

	"github.com/tobiassjosten/nogfx/pkg/mock"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
)

func TestRenderOutput(t *testing.T) {
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
		buffer string
		datas  [][]byte
		prints [][]byte
		width  int
		height int
		pwidth int
		offset int
		rows   []string
	}{
		"empty": {
			width:  1,
			height: 1,
			rows:   []string{" "},
		},

		"simple output": {
			buffer: "a",
			width:  1,
			height: 1,
			rows:   []string{"a"},
		},

		"cramped height output": {
			buffer: "a",
			width:  1,
			height: 0,
			rows:   nil,
		},

		"cramped width output": {
			buffer: "a",
			width:  0,
			height: 1,
			rows:   nil,
		},

		"output without padding": {
			buffer: "a",
			width:  2,
			height: 1,
			rows:   []string{"a "},
		},

		"word wrap simple": {
			buffer: "a sdf",
			width:  2,
			height: 3,
			rows:   []string{"a ", "sd", "f "},
		},

		"word wrap overflow": {
			buffer: "a sdf",
			width:  2,
			height: 2,
			rows:   []string{"sd", "f "},
		},

		"history scrollback even": {
			buffer: "a sdfgh",
			width:  2,
			height: 3,
			offset: 1,
			rows:   []string{"a ", "──", "h "},
		},

		"history scrollback odd": {
			buffer: "a sdfghjk",
			width:  2,
			height: 4,
			offset: 1,
			rows:   []string{"a ", "sd", "──", "k "},
		},

		"resize kill history scrollback": {
			buffer: "a sdfgh",
			width:  2,
			height: 3,
			pwidth: 1,
			offset: 1,
			rows:   []string{"sd", "fg", "h "},
		},

		"maintain history scrollback": {
			buffer: "a sdfghjk",
			datas:  [][]byte{[]byte("xy")},
			width:  2,
			height: 4,
			offset: 1,
			rows:   []string{"a ", "sd", "──", "xy"},
		},

		"print message": {
			buffer: "asdf",
			prints: [][]byte{[]byte("xy")},
			width:  2,
			height: 3,
			rows:   []string{"as", "df", "xy"},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			tui := NewTUI(&mock.ScreenMock{
				HideCursorFunc:     func() {},
				SetCursorStyleFunc: func(_ tcell.CursorStyle) {},
				SetStyleFunc:       func(_ tcell.Style) {},
			})
			tui.output.pwidth = tc.pwidth

			if tc.buffer != "" {
				tui.output.buffer = Rows{NewRowFromRunes([]rune(tc.buffer))}
			}

			if tc.offset > 0 {
				tui.output.offset = tc.offset
			}

			if len(tc.datas) > 0 {
				tui.output.pwidth = tc.width
				for _, data := range tc.datas {
					tui.output.Append(data)
				}
			}

			for _, print := range tc.prints {
				tui.Print(print)
			}

			rows := tui.RenderOutput(tc.width, tc.height)
			assert.Equal(t, tc.rows, rowsToStrings(rows))
		})
	}
}

func TestOutputAppend(t *testing.T) {
	redStyle := (tcell.Style{}).
		Foreground(tcell.ColorGreen).
		Background(tcell.ColorBlue)

	tcs := map[string]struct {
		datas  [][]byte
		buffer Rows
	}{
		"plain text xy": {
			datas:  [][]byte{[]byte("xy")},
			buffer: Rows{Row{NewCell('x'), NewCell('y')}},
		},

		"plain text yz": {
			datas:  [][]byte{[]byte("yx")},
			buffer: Rows{Row{NewCell('y'), NewCell('x')}},
		},

		"reverse order": {
			datas:  [][]byte{[]byte("x"), []byte("y")},
			buffer: Rows{Row{NewCell('y')}, Row{NewCell('x')}},
		},

		"maintain color": {
			datas: [][]byte{
				[]byte(""),
				[]byte("\033[32;44my"),
				[]byte("z"),
			},
			buffer: Rows{
				Row{NewCell('z', redStyle)},
				Row{NewCell('y', redStyle)},
				Row{}, // sometimes, somehow, we get empties
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			output := &Output{}
			for _, data := range tc.datas {
				output.Append(data)
			}

			assert.Equal(t, tc.buffer, output.buffer)
		})
	}
}
