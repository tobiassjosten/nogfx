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
		width int
		cell  tui.Cell
		rs    []rune
		bs    []byte
		style tcell.Style
		row   tui.Row
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
			rs:    []rune{'a'},
			style: redStyle,
			row:   tui.Row{tui.NewCell('a', redStyle)},
		},

		"simple runes green style": {
			rs:    []rune{'a'},
			style: greenStyle,
			row:   tui.Row{tui.NewCell('a', greenStyle)},
		},

		"simple bytes default style": {
			bs:  []byte("a"),
			row: tui.Row{tui.NewCell('a', baseStyle)},
		},

		"simple bytes red style": {
			bs:    []byte("\033[31;43ma"),
			style: redStyle,
			row:   tui.Row{tui.NewCell('a', redStyle)},
		},

		"simple bytes green style": {
			bs:    []byte("\033[32;44ma"),
			style: greenStyle,
			row:   tui.Row{tui.NewCell('a', greenStyle)},
		},

		"invalid ansi color": {
			bs:    []byte("\033{32ma"),
			style: greenStyle,
			row: tui.Row{
				tui.NewCell('^', greenStyle),
				tui.NewCell('{', greenStyle),
				tui.NewCell('3', greenStyle),
				tui.NewCell('2', greenStyle),
				tui.NewCell('m', greenStyle),
				tui.NewCell('a', greenStyle),
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			var row tui.Row
			if tc.width > 0 {
				if tc.cell != (tui.Cell{}) {
					row = tui.NewRow(tc.width, tc.cell)
				} else {
					row = tui.NewRow(tc.width)
				}
			}
			if len(tc.rs) > 0 {
				if tc.style != (tcell.Style{}) {
					row = tui.NewRowFromRunes(tc.rs, tc.style)
				} else {
					row = tui.NewRowFromRunes(tc.rs)
				}
			}
			if len(tc.bs) > 0 {
				if tc.style != (tcell.Style{}) {
					row = tui.NewRowFromBytes(tc.bs, tc.style)
				} else {
					row = tui.NewRowFromBytes(tc.bs)
				}
			}

			assert.Equal(t, tc.row, row)
		})
	}
}

func TestWrap(t *testing.T) {
	tcs := map[string]struct {
		row   tui.Row
		width int
		rows  []string
	}{
		"asdf in one column": {
			row:   tui.NewRowFromBytes([]byte("asdf")),
			width: 1,
			rows:  []string{"a", "s", "d", "f"},
		},

		"asdf in two columns": {
			row:   tui.NewRowFromBytes([]byte("asdf")),
			width: 2,
			rows:  []string{"as", "df"},
		},

		"asdf in five columns": {
			row:   tui.NewRowFromBytes([]byte("asdf")),
			width: 5,
			rows:  []string{"asdf"},
		},

		"simple wordwrap exact single first": {
			row:   tui.NewRowFromBytes([]byte("as df")),
			width: 3,
			rows:  []string{"as", "df"},
		},

		"simple wordwrap exact single second": {
			row:   tui.NewRowFromBytes([]byte("a s d f")),
			width: 3,
			rows:  []string{"a s", "d f"},
		},

		"simple wordwrap exact single third": {
			row:   tui.NewRowFromBytes([]byte("a sd f")),
			width: 3,
			rows:  []string{"a", "sd", "f"},
		},

		"simple wordwrap exact multiple": {
			row:   tui.NewRowFromBytes([]byte("as  df")),
			width: 3,
			rows:  []string{"as", "df"},
		},

		"simple wordwrap trailing": {
			row:   tui.NewRowFromBytes([]byte("as df")),
			width: 2,
			rows:  []string{"as", "df"},
		},

		"multi-space wordwrap": {
			row:   tui.NewRowFromBytes([]byte("as   df")),
			width: 3,
			rows:  []string{"as", "df"},
		},

		"non-destructive": {
			row:   tui.NewRowFromBytes([]byte("  as")),
			width: 4,
			rows:  []string{"  as"},
		},

		"strip only wrapping word": {
			row:   tui.NewRowFromBytes([]byte("as  ")),
			width: 4,
			rows:  []string{"as  "},
		},

		"space prefix wordwrap": {
			row:   tui.NewRowFromBytes([]byte("  as   df")),
			width: 3,
			rows:  []string{"  a", "s", "df"},
		},

		"space suffix wordwrap one": {
			row:   tui.NewRowFromBytes([]byte("as   df  ")),
			width: 3,
			rows:  []string{"as", "df"},
		},

		"space suffix wordwrap two": {
			row:   tui.NewRowFromBytes([]byte("asd   ")),
			width: 3,
			rows:  []string{"asd"},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			rows := tc.row.Wrap(tc.width)
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
