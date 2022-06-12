package tui_test

import (
	"context"
	"testing"

	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/mock"
	"github.com/tobiassjosten/nogfx/pkg/tui"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	healthFullStyle = tcell.StyleDefault.
			Background(tcell.ColorGreen).
			Foreground(tcell.ColorBlack)

	healthEmptyStyle = tcell.StyleDefault.
				Background(tcell.ColorDarkGreen).
				Foreground(tcell.ColorBlack)

	manaFullStyle = tcell.StyleDefault.
			Background(tcell.ColorBlue).
			Foreground(tcell.ColorBlack)

	manaEmptyStyle = tcell.StyleDefault.
			Background(tcell.ColorDarkBlue).
			Foreground(tcell.ColorBlack)

	enduranceFullStyle = tcell.StyleDefault.
				Background(tcell.ColorTeal).
				Foreground(tcell.ColorBlack)

	enduranceEmptyStyle = tcell.StyleDefault.
				Background(tcell.ColorDarkCyan).
				Foreground(tcell.ColorBlack)

	willpowerFullStyle = tcell.StyleDefault.
				Background(tcell.ColorFuchsia).
				Foreground(tcell.ColorBlack)

	willpowerEmptyStyle = tcell.StyleDefault.
				Background(tcell.ColorRebeccaPurple).
				Foreground(tcell.ColorBlack)

	energyFullStyle = tcell.StyleDefault.
			Background(tcell.ColorRed).
			Foreground(tcell.ColorBlack)

	energyEmptyStyle = tcell.StyleDefault.
				Background(tcell.ColorRosyBrown).
				Foreground(tcell.ColorBlack)

	unknownFullStyle = tcell.StyleDefault.
				Background(tcell.ColorYellow).
				Foreground(tcell.ColorBlack)

	unknownEmptyStyle = tcell.StyleDefault.
				Background(tcell.Color100).
				Foreground(tcell.ColorBlack)
)

func TestRenderVital(t *testing.T) {
	fullStyle := healthFullStyle
	emptyStyle := healthEmptyStyle

	tcs := map[string]struct {
		vital  pkg.CharacterVital
		width  int
		styles []tcell.Style
		row    tui.Row
	}{
		"1/1 one cell": {
			vital:  pkg.CharacterVital{Value: 1, Max: 1},
			width:  1,
			styles: []tcell.Style{fullStyle, emptyStyle},
			row: tui.Row{
				tui.NewCell('1', fullStyle),
			},
		},

		"0/1 one cell": {
			vital:  pkg.CharacterVital{Value: 0, Max: 1},
			width:  1,
			styles: []tcell.Style{fullStyle, emptyStyle},
			row: tui.Row{
				tui.NewCell('0', emptyStyle),
			},
		},

		"1/2 two cells": {
			vital:  pkg.CharacterVital{Value: 1, Max: 2},
			width:  2,
			styles: []tcell.Style{fullStyle, emptyStyle},
			row: tui.Row{
				tui.NewCell('1', fullStyle),
				tui.NewCell(' ', emptyStyle),
			},
		},

		"1/2 three cells": {
			vital:  pkg.CharacterVital{Value: 1, Max: 2},
			width:  3,
			styles: []tcell.Style{fullStyle, emptyStyle},
			row: tui.Row{
				tui.NewCell(' ', fullStyle),
				tui.NewCell('1', emptyStyle),
				tui.NewCell(' ', emptyStyle),
			},
		},

		"3/4 four cells": {
			vital:  pkg.CharacterVital{Value: 3, Max: 4},
			width:  4,
			styles: []tcell.Style{fullStyle, emptyStyle},
			row: tui.Row{
				tui.NewCell(' ', fullStyle),
				tui.NewCell('3', fullStyle),
				tui.NewCell(' ', fullStyle),
				tui.NewCell(' ', emptyStyle),
			},
		},

		"10/40 four cells": {
			vital:  pkg.CharacterVital{Value: 10, Max: 40},
			width:  4,
			styles: []tcell.Style{fullStyle, emptyStyle},
			row: tui.Row{
				tui.NewCell(' ', fullStyle),
				tui.NewCell('1', emptyStyle),
				tui.NewCell('0', emptyStyle),
				tui.NewCell(' ', emptyStyle),
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			row := tui.RenderVital(tc.vital, tc.width, tc.styles)
			assert.Equal(t, tc.row, row)
		})
	}
}

func TestRenderVitals(t *testing.T) {
	tcs := map[string]struct {
		vitals map[string]pkg.CharacterVital
		width  int
		height int
		rows   tui.Rows
		err    string
	}{
		"no vitals": {
			width:  1,
			height: 2,
			rows: tui.Rows{
				tui.Row{tui.NewCell(' ')},
				tui.Row{tui.NewCell(' ')},
			},
		},

		"1x2 health 1/1": {
			vitals: map[string]pkg.CharacterVital{
				"health": {Value: 1, Max: 1},
			},
			width:  1,
			height: 2,
			rows: tui.Rows{
				tui.Row{tui.NewCell(' ')},
				tui.Row{tui.NewCell('1', healthFullStyle)},
			},
		},

		"2x2 mana 1/2": {
			vitals: map[string]pkg.CharacterVital{
				"mana": {Value: 1, Max: 2},
			},
			width:  2,
			height: 2,
			rows: tui.Rows{
				tui.NewRow(2, tui.NewCell(' ')),
				tui.Row{
					tui.NewCell('1', manaFullStyle),
					tui.NewCell(' ', manaEmptyStyle),
				},
			},
		},

		"3x2 endurance 1/3": {
			vitals: map[string]pkg.CharacterVital{
				"endurance": {Value: 1, Max: 3},
			},
			width:  3,
			height: 2,
			rows: tui.Rows{
				tui.NewRow(3, tui.NewCell(' ')),
				tui.Row{
					tui.NewCell(' ', enduranceFullStyle),
					tui.NewCell('1', enduranceEmptyStyle),
					tui.NewCell(' ', enduranceEmptyStyle),
				},
			},
		},

		"4x2 willpower 3/4": {
			vitals: map[string]pkg.CharacterVital{
				"willpower": {Value: 3, Max: 4},
			},
			width:  4,
			height: 2,
			rows: tui.Rows{
				tui.NewRow(4, tui.NewCell(' ')),
				tui.Row{
					tui.NewCell(' ', willpowerFullStyle),
					tui.NewCell('3', willpowerFullStyle),
					tui.NewCell(' ', willpowerFullStyle),
					tui.NewCell(' ', willpowerEmptyStyle),
				},
			},
		},

		"2x2 energy 100/200": {
			vitals: map[string]pkg.CharacterVital{
				"energy": {Value: 100, Max: 200},
			},
			width:  2,
			height: 2,
			rows: tui.Rows{
				tui.NewRow(2, tui.NewCell(' ')),
				tui.Row{
					tui.NewCell(' ', energyFullStyle),
					tui.NewCell(' ', energyEmptyStyle),
				},
			},
		},

		"1x2 unknown 1/2": {
			vitals: map[string]pkg.CharacterVital{
				"unknown": {Value: 1, Max: 2},
			},
			width:  2,
			height: 2,
			rows: tui.Rows{
				tui.NewRow(2, tui.NewCell(' ')),
				tui.Row{
					tui.NewCell('1', unknownFullStyle),
					tui.NewCell(' ', unknownEmptyStyle),
				},
			},
		},

		"9x2 health 1/2 mana 3/4": {
			vitals: map[string]pkg.CharacterVital{
				"health": {Value: 1, Max: 2},
				"mana":   {Value: 3, Max: 4},
			},
			width:  9,
			height: 2,
			rows: tui.Rows{
				tui.NewRow(9, tui.NewCell(' ')),
				tui.Row{
					tui.NewCell(' ', healthFullStyle),
					tui.NewCell('1', healthFullStyle),
					tui.NewCell(' ', healthEmptyStyle),
					tui.NewCell(' ', healthEmptyStyle),
					tui.NewCell(' '),
					tui.NewCell(' ', manaFullStyle),
					tui.NewCell('3', manaFullStyle),
					tui.NewCell(' ', manaFullStyle),
					tui.NewCell(' ', manaEmptyStyle),
				},
			},
		},

		"9x2 endurance 10/20 willpower 10/40": {
			vitals: map[string]pkg.CharacterVital{
				"endurance": {Value: 10, Max: 20},
				"willpower": {Value: 10, Max: 40},
			},
			width:  7,
			height: 2,
			rows: tui.Rows{
				tui.NewRow(7, tui.NewCell(' ')),
				tui.Row{
					tui.NewCell('1', enduranceFullStyle),
					tui.NewCell('0', enduranceEmptyStyle),
					tui.NewCell(' ', enduranceEmptyStyle),
					tui.NewCell(' '),
					tui.NewCell('1', willpowerFullStyle),
					tui.NewCell('0', willpowerEmptyStyle),
					tui.NewCell(' ', willpowerEmptyStyle),
				},
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			rows := tui.NewRows(tc.width, tc.height)

			screen := &mock.ScreenMock{
				ClearFunc:      func() {},
				FiniFunc:       func() {},
				HideCursorFunc: func() {},
				InitFunc: func() error {
					return nil
				},
				PollEventFunc: func() tcell.Event {
					return nil
				},
				SetContentFunc: func(x, y int, r rune, rs []rune, style tcell.Style) {
					rows[y][x] = tui.NewCell(r, style)
				},
				SetCursorStyleFunc: func(_ tcell.CursorStyle) {},
				SetStyleFunc:       func(_ tcell.Style) {},
				ShowFunc:           func() {},
				SizeFunc: func() (int, int) {
					return tc.width, tc.height
				},
			}

			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			ui := tui.NewTUI(screen)

			ui.SetCharacter(pkg.Character{Vitals: tc.vitals})

			_ = ui.Run(ctx)

			require.Equal(t, len(tc.rows), len(rows))
			for i, row := range rows {
				require.Equalf(t,
					len(tc.rows[i]), len(row),
					"row %d", i,
				)

				for ii, cell := range row {
					tccell := tc.rows[i][ii]
					assert.Equalf(t,
						tccell, cell,
						"row %d, cell %d", i, ii,
					)
				}
			}
		})
	}
}

/*
var (
	healthFullStyle = tcell.StyleDefault.
			Background(tcell.ColorGreen).
			Foreground(tcell.ColorBlack)

	healthEmptyStyle = tcell.StyleDefault.
				Background(tcell.ColorDarkGreen).
				Foreground(tcell.ColorBlack)

	manaFullStyle = tcell.StyleDefault.
			Background(tcell.ColorBlue).
			Foreground(tcell.ColorBlack)

	manaEmptyStyle = tcell.StyleDefault.
			Background(tcell.ColorDarkBlue).
			Foreground(tcell.ColorBlack)

	enduranceFullStyle = tcell.StyleDefault.
				Background(tcell.ColorTeal).
				Foreground(tcell.ColorBlack)

	enduranceEmptyStyle = tcell.StyleDefault.
				Background(tcell.ColorDarkCyan).
				Foreground(tcell.ColorBlack)

	willpowerFullStyle = tcell.StyleDefault.
				Background(tcell.ColorFuchsia).
				Foreground(tcell.ColorBlack)

	willpowerEmptyStyle = tcell.StyleDefault.
				Background(tcell.ColorRebeccaPurple).
				Foreground(tcell.ColorBlack)

	energyEmptyStyle = tcell.StyleDefault.
				Background(tcell.ColorRosyBrown).
				Foreground(tcell.ColorBlack)
)

func TestRenderVital(t *testing.T) {
	fullStyle := healthFullStyle
	emptyStyle := healthEmptyStyle

	tcs := map[string]struct {
		vital  pkg.CharacterVital
		width  int
		styles []tcell.Style
		row    tui.Row
	}{
		"100% health, even": {
			vital:  pkg.CharacterVital{Value: 20, Max: 20},
			width:  4,
			styles: []tcell.Style{fullStyle, emptyStyle},
			row: tui.Row{
				tui.NewCell(' ', fullStyle),
				tui.NewCell('2', fullStyle),
				tui.NewCell('0', fullStyle),
				tui.NewCell(' ', fullStyle),
			},
		},

		"100% health, odd": {
			vital:  pkg.CharacterVital{Value: 20, Max: 20},
			width:  5,
			styles: []tcell.Style{fullStyle, emptyStyle},
			row: tui.Row{
				tui.NewCell(' ', fullStyle),
				tui.NewCell('2', fullStyle),
				tui.NewCell('0', fullStyle),
				tui.NewCell(' ', fullStyle),
				tui.NewCell(' ', fullStyle),
			},
		},

		"100% health, cramped odd": {
			vital: pkg.CharacterVital{Value: 20, Max: 20},
			width: 3,
			row: tui.Row{
				tui.NewCell('2', healthFullStyle),
				tui.NewCell('0', healthFullStyle),
				tui.NewCell(' ', healthFullStyle),
			},
		},

		"100% health, cramped even": {
			vital: pkg.CharacterVital{Value: 20, Max: 20},
			width: 2,
			row: tui.Row{
				tui.NewCell('2', healthFullStyle),
				tui.NewCell('0', healthFullStyle),
			},
		},

		"100% health, too cramped": {
			vital: pkg.CharacterVital{Value: 200, Max: 200},
			width: 2,
			row: tui.Row{
				tui.NewCell(' ', healthFullStyle),
				tui.NewCell(' ', healthFullStyle),
			},
		},

		"75% mana": {
			vital: pkg.CharacterVital{Value: 15, Max: 20},
			width: 4,
			row: tui.Row{
				tui.NewCell(' ', manaFullStyle),
				tui.NewCell('1', manaFullStyle),
				tui.NewCell('5', manaFullStyle),
				tui.NewCell(' ', manaEmptyStyle),
			},
		},

		"50% endurance": {
			vital: pkg.CharacterVital{Value: 10, Max: 20},
			width: 4,
			row: tui.Row{
				tui.NewCell(' ', enduranceFullStyle),
				tui.NewCell('1', enduranceFullStyle),
				tui.NewCell('0', enduranceEmptyStyle),
				tui.NewCell(' ', enduranceEmptyStyle),
			},
		},

		"25% willpower": {
			vital: pkg.CharacterVital{Value: 5, Max: 20},
			width: 4,
			row: tui.Row{
				tui.NewCell(' ', willpowerFullStyle),
				tui.NewCell('5', willpowerEmptyStyle),
				tui.NewCell(' ', willpowerEmptyStyle),
				tui.NewCell(' ', willpowerEmptyStyle),
			},
		},

		"0% energy": {
			vital: pkg.CharacterVital{Value: 0, Max: 20},
			width: 4,
			row: tui.Row{
				tui.NewCell(' ', energyEmptyStyle),
				tui.NewCell('0', energyEmptyStyle),
				tui.NewCell(' ', energyEmptyStyle),
				tui.NewCell(' ', energyEmptyStyle),
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			row := tui.RenderVital(tc.vital, tc.width, tc.styles)
			assert.Equal(t, tc.row, row)
		})
	}
}

type MockVital struct {
}

func TestRenderVitals(t *testing.T) {
	tcs := map[string]struct {
		character pkg.Character
		vorder    []string
		vitals    map[string]interface{}
		updates   map[string][]int
		width     int
		height    int
		rows      tui.Rows
		err       string
	}{
		"no vitals": {
			width:  2,
			height: 2,
			rows: tui.Rows{
				tui.NewRow(2, tui.NewCell(' ')),
				tui.NewRow(2, tui.NewCell(' ')),
			},
		},

		"too short for vitals": {
			vorder: []string{"health"},
			vitals: map[string]interface{}{
				"health": pkg.CharacterVital{Value: 1, Max: 2},
			},
			width:  2,
			height: 1,
			rows: tui.Rows{
				tui.NewRow(2, tui.NewCell(' ')),
			},
		},

		"too narrow for numbers": {
			vorder: []string{"health"},
			vitals: map[string]interface{}{
				"health": pkg.CharacterVital{Value: 100, Max: 200},
			},
			width:  2,
			height: 2,
			rows: tui.Rows{
				tui.NewRow(2, tui.NewCell(' ')),
				tui.Row{
					tui.NewCell(' ', healthFullStyle),
					tui.NewCell(' ', healthEmptyStyle),
				},
			},
		},

		"half health": {
			vorder: []string{"health"},
			vitals: map[string]interface{}{
				"health": pkg.CharacterVital{Value: 1, Max: 2},
			},
			width:  2,
			height: 2,
			rows: tui.Rows{
				tui.NewRow(2, tui.NewCell(' ')),
				tui.Row{
					tui.NewCell('1', healthFullStyle),
					tui.NewCell(' ', healthEmptyStyle),
				},
			},
		},

		"full health, half mana": {
			vorder: []string{"health", "mana"},
			vitals: map[string]interface{}{
				"health": pkg.CharacterVital{Value: 2, Max: 2},
				"mana":   pkg.CharacterVital{Value: 1, Max: 2},
			},
			width:  5,
			height: 2,
			rows: tui.Rows{
				tui.NewRow(5, tui.NewCell(' ')),
				tui.Row{
					tui.NewCell('2', healthFullStyle),
					tui.NewCell(' ', healthFullStyle),
					tui.NewCell(' '),
					tui.NewCell('1', manaFullStyle),
					tui.NewCell(' ', manaEmptyStyle),
				},
			},
		},

		"half endurance, full willpower": {
			vorder: []string{"endurance", "willpower"},
			vitals: map[string]interface{}{
				"endurance": pkg.CharacterVital{Value: 1, Max: 2},
				"willpower": pkg.CharacterVital{Value: 2, Max: 2},
			},
			width:  5,
			height: 2,
			rows: tui.Rows{
				tui.NewRow(5, tui.NewCell(' ')),
				tui.Row{
					tui.NewCell('1', enduranceFullStyle),
					tui.NewCell(' ', enduranceEmptyStyle),
					tui.NewCell(' '),
					tui.NewCell('2', willpowerFullStyle),
					tui.NewCell(' ', willpowerFullStyle),
				},
			},
		},

		"full endurance, half willpower": {
			vorder: []string{"endurance", "willpower"},
			vitals: map[string]interface{}{
				"endurance": pkg.CharacterVital{Value: 1, Max: 2},
				"willpower": pkg.CharacterVital{Value: 2, Max: 2},
			},
			updates: map[string][]int{
				"endurance": []int{3, 3},
				"willpower": []int{2, 4},
			},
			width:  5,
			height: 2,
			rows: tui.Rows{
				tui.NewRow(5, tui.NewCell(' ')),
				tui.Row{
					tui.NewCell('3', enduranceFullStyle),
					tui.NewCell(' ', enduranceFullStyle),
					tui.NewCell(' '),
					tui.NewCell('2', willpowerFullStyle),
					tui.NewCell(' ', willpowerEmptyStyle),
				},
			},
		},

		"add wrong type": {
			vorder: []string{"wrong-type"},
			vitals: map[string]interface{}{
				"wrong-type": MockVital{},
			},
			err: "unsupported vital 'wrong-type'",
		},

		"update non-existant": {
			vorder: []string{"non-existant"},
			updates: map[string][]int{
				"non-existant": []int{1, 1},
			},
			err: "couldn't update non-existent 'non-existant' vital",
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			rows := tui.NewRows(tc.width, tc.height)

			screen := &mock.ScreenMock{
				ClearFunc:      func() {},
				FiniFunc:       func() {},
				HideCursorFunc: func() {},
				InitFunc: func() error {
					return nil
				},
				PollEventFunc: func() tcell.Event {
					return nil
				},
				SetContentFunc: func(x, y int, r rune, rs []rune, style tcell.Style) {
					rows[y][x] = tui.NewCell(r, style)
				},
				SetCursorStyleFunc: func(_ tcell.CursorStyle) {},
				SetStyleFunc:       func(_ tcell.Style) {},
				ShowFunc:           func() {},
				SizeFunc: func() (int, int) {
					return tc.width, tc.height
				},
			}

			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			ui := tui.NewTUI(screen)

			ui.SetCharacter(tc.character)

			_ = ui.Run(ctx)

			assert.Equal(t, tc.rows, rows)
		})
	}
}
*/
