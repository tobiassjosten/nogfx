package tui_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tobiassjosten/nogfx/pkg/tui"

	"github.com/gdamore/tcell/v2"
)

func rowToString(row tui.Row) (str string) {
	for _, c := range row {
		str += string(c.Content)
	}
	return
}

func rowsToStrings(rows tui.Rows) (strs []string) {
	for _, row := range rows {
		strs = append(strs, rowToString(row))
	}
	return
}

func TestNewRow(t *testing.T) {
	baseStyle := tcell.Style{}
	redStyle := baseStyle.
		Foreground(tcell.ColorRed).
		Background(tcell.ColorYellow)
	greenStyle := baseStyle.
		Foreground(tcell.ColorGreen).
		Background(tcell.ColorBlue)

	tcs := map[string]struct {
		width    int
		cell     tui.Cell
		rs       []rune
		bs       []byte
		stylein  *tcell.Style
		row      tui.Row
		styleout *tcell.Style
	}{
		"simple row no cell": {
			width: 1,
			row:   tui.Row{tui.NewCell(' ', baseStyle)},
		},

		"simple row 'a' cell": {
			width: 1,
			cell:  tui.NewCell('a', redStyle),
			row:   tui.Row{tui.NewCell('a', redStyle)},
		},

		"simple runes default style": {
			rs:  []rune{'a'},
			row: tui.Row{tui.NewCell('a', baseStyle)},
		},

		"simple runes red style": {
			rs:      []rune{'a'},
			stylein: &redStyle,
			row:     tui.Row{tui.NewCell('a', redStyle)},
		},

		"simple runes green style": {
			rs:      []rune{'a'},
			stylein: &greenStyle,
			row:     tui.Row{tui.NewCell('a', greenStyle)},
		},

		"simple bytes default style": {
			bs:  []byte("a"),
			row: tui.Row{tui.NewCell('a', baseStyle)},
		},

		"simple bytes red style": {
			bs:      []byte("\033[31;43ma"),
			stylein: &redStyle,
			row:     tui.Row{tui.NewCell('a', redStyle)},
		},

		"simple bytes green style": {
			bs:      []byte("\033[32;44ma"),
			stylein: &greenStyle,
			row:     tui.Row{tui.NewCell('a', greenStyle)},
		},

		"invalid ansi color": {
			bs:      []byte("\033{32ma"),
			stylein: &greenStyle,
			row: tui.Row{
				tui.NewCell('^', greenStyle),
				tui.NewCell('{', greenStyle),
				tui.NewCell('3', greenStyle),
				tui.NewCell('2', greenStyle),
				tui.NewCell('m', greenStyle),
				tui.NewCell('a', greenStyle),
			},
		},

		// @todo Test different style inputs/outputs.
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			var row tui.Row
			var style tcell.Style

			if tc.width > 0 {
				if tc.cell != (tui.Cell{}) {
					row = tui.NewRow(tc.width, tc.cell)
				} else {
					row = tui.NewRow(tc.width)
				}
			}

			if len(tc.rs) > 0 {
				if tc.stylein != nil {
					row = tui.NewRowFromRunes(tc.rs, *tc.stylein)
				} else {
					row = tui.NewRowFromRunes(tc.rs)
				}
			}

			if len(tc.bs) > 0 {
				if tc.stylein != nil {
					row, style = tui.NewRowFromBytes(tc.bs, *tc.stylein)
				} else {
					row, style = tui.NewRowFromBytes(tc.bs)
				}
			}

			assert.Equal(t, tc.row, row)
			if tc.styleout != nil {
				assert.Equal(t, *tc.styleout, style)
			}
		})
	}
}

func TestWrap(t *testing.T) {
	tcs := map[string]struct {
		line    string
		row     tui.Row
		width   int
		padding rune
		rows    []string
	}{
		"asdf in one column": {
			line:  "asdf",
			width: 1,
			rows:  []string{"a", "s", "d", "f"},
		},

		"asdf in two columns": {
			line:  "asdf",
			width: 2,
			rows:  []string{"as", "df"},
		},

		"asdf in five columns": {
			line:  "asdf",
			width: 5,
			rows:  []string{"asdf"},
		},

		"simple wordwrap exact single first": {
			line:  "as df",
			width: 3,
			rows:  []string{"as", "df"},
		},

		"simple wordwrap exact single second": {
			line:  "a s d f",
			width: 3,
			rows:  []string{"a s", "d f"},
		},

		"simple wordwrap exact single third": {
			line:  "a sd f",
			width: 3,
			rows:  []string{"a", "sd", "f"},
		},

		"simple wordwrap exact multiple": {
			line:  "as  df",
			width: 3,
			rows:  []string{"as", "df"},
		},

		"simple wordwrap trailing": {
			line:  "as df",
			width: 2,
			rows:  []string{"as", "df"},
		},

		"multi-space wordwrap": {
			line:  "as   df",
			width: 3,
			rows:  []string{"as", "df"},
		},

		"non-destructive": {
			line:  "  as",
			width: 4,
			rows:  []string{"  as"},
		},

		"strip only wrapping word": {
			line:  "as  ",
			width: 4,
			rows:  []string{"as  "},
		},

		"space prefix wordwrap": {
			line:  "  as   df",
			width: 3,
			rows:  []string{"  a", "s", "df"},
		},

		"space suffix wordwrap one": {
			line:  "as   df  ",
			width: 3,
			rows:  []string{"as", "df"},
		},

		"space suffix wordwrap two": {
			line:  "asd   ",
			width: 3,
			rows:  []string{"asd"},
		},

		"asdf": {
			line:    "a",
			width:   3,
			padding: ' ',
			rows:    []string{"a  "},
		},

		"qwer": {
			line:    "asd fgh ij",
			width:   5,
			padding: ' ',
			rows:    []string{"asd  ", "fgh  ", "ij   "},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			row, _ := tui.NewRowFromBytes([]byte(tc.line))

			var rows tui.Rows
			if tc.padding != 0 {
				rows = row.Wrap(tc.width, tui.NewCell(tc.padding))
			} else {
				rows = row.Wrap(tc.width)
			}
			assert.Equal(t, tc.rows, rowsToStrings(rows))
		})
	}
}

func TestNewRows(t *testing.T) {
	baseStyle := tcell.Style{}
	redStyle := baseStyle.
		Foreground(tcell.ColorRed).
		Background(tcell.ColorYellow)

	tcs := map[string]struct {
		width  int
		height int
		cell   tui.Cell
		rows   tui.Rows
	}{
		"one column, no row, no cell": {
			width: 1,
			rows:  tui.Rows{},
		},

		"no column, one row, no cell": {
			height: 1,
			rows:   tui.Rows{tui.Row{}},
		},

		"no column, no row, one cell": {
			cell: tui.NewCell('a', redStyle),
			rows: tui.Rows{},
		},

		"one column, one row, no cell": {
			width:  1,
			height: 1,
			rows:   tui.Rows{tui.Row{tui.NewCell(' ')}},
		},

		"one column, one row, one cell": {
			width:  1,
			height: 1,
			cell:   tui.NewCell('a', redStyle),
			rows:   tui.Rows{tui.Row{tui.NewCell('a', redStyle)}},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			var rows tui.Rows
			if tc.cell != (tui.Cell{}) {
				rows = tui.NewRows(tc.width, tc.height, tc.cell)
			} else {
				rows = tui.NewRows(tc.width, tc.height)
			}

			assert.Equal(t, tc.rows, rows)
		})
	}
}
