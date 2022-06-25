package tui_test

import (
	"context"
	"testing"

	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/mock"
	"github.com/tobiassjosten/nogfx/pkg/tui"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
)

var (
	gapStyle = tcell.StyleDefault.Background(tcell.Color235)

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
			Background(tcell.ColorYellow).
			Foreground(tcell.ColorBlack)

	energyEmptyStyle = tcell.StyleDefault.
				Background(tcell.Color100).
				Foreground(tcell.ColorBlack)

	unknownFullStyle = tcell.StyleDefault.
				Background(tcell.Color250).
				Foreground(tcell.ColorBlack)

	unknownEmptyStyle = tcell.StyleDefault.
				Background(tcell.Color240).
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
		row    tui.Row
		err    string
	}{
		"too short": {
			width:  1,
			height: 4,
			row:    tui.Row{tui.NewCell(' ')},
		},

		"width 1 health 1/1": {
			vitals: map[string]pkg.CharacterVital{
				"health": {Value: 1, Max: 1},
			},
			width:  1,
			height: 5,
			row:    tui.Row{tui.NewCell('1', healthFullStyle)},
		},

		"width 2 mana 1/2": {
			vitals: map[string]pkg.CharacterVital{
				"mana": {Value: 1, Max: 2},
			},
			width:  2,
			height: 5,
			row: tui.Row{
				tui.NewCell('1', manaFullStyle),
				tui.NewCell(' ', manaEmptyStyle),
			},
		},

		"width 3 endurance 1/3": {
			vitals: map[string]pkg.CharacterVital{
				"endurance": {Value: 1, Max: 3},
			},
			width:  3,
			height: 5,
			row: tui.Row{
				tui.NewCell(' ', enduranceFullStyle),
				tui.NewCell('1', enduranceEmptyStyle),
				tui.NewCell(' ', enduranceEmptyStyle),
			},
		},

		"width 4 willpower 3/4": {
			vitals: map[string]pkg.CharacterVital{
				"willpower": {Value: 3, Max: 4},
			},
			width:  4,
			height: 5,
			row: tui.Row{
				tui.NewCell(' ', willpowerFullStyle),
				tui.NewCell('3', willpowerFullStyle),
				tui.NewCell(' ', willpowerFullStyle),
				tui.NewCell(' ', willpowerEmptyStyle),
			},
		},

		"width 2 energy 100/200": {
			vitals: map[string]pkg.CharacterVital{
				"energy": {Value: 100, Max: 200},
			},
			width:  2,
			height: 5,
			row: tui.Row{
				tui.NewCell(' ', energyFullStyle),
				tui.NewCell(' ', energyEmptyStyle),
			},
		},

		"width 2 unknown 1/2": {
			vitals: map[string]pkg.CharacterVital{
				"unknown": {Value: 1, Max: 2},
			},
			width:  2,
			height: 5,
			row: tui.Row{
				tui.NewCell('1', unknownFullStyle),
				tui.NewCell(' ', unknownEmptyStyle),
			},
		},

		"width 9 health 1/2 mana 3/4": {
			vitals: map[string]pkg.CharacterVital{
				"health": {Value: 1, Max: 2},
				"mana":   {Value: 3, Max: 4},
			},
			width:  9,
			height: 5,
			row: tui.Row{
				tui.NewCell(' ', healthFullStyle),
				tui.NewCell('1', healthFullStyle),
				tui.NewCell(' ', healthEmptyStyle),
				tui.NewCell(' ', healthEmptyStyle),
				tui.NewCell(' ', gapStyle),
				tui.NewCell(' ', manaFullStyle),
				tui.NewCell('3', manaFullStyle),
				tui.NewCell(' ', manaFullStyle),
				tui.NewCell(' ', manaEmptyStyle),
			},
		},

		"width 7 endurance 10/20 willpower 10/40": {
			vitals: map[string]pkg.CharacterVital{
				"endurance": {Value: 10, Max: 20},
				"willpower": {Value: 10, Max: 40},
			},
			width:  7,
			height: 5,
			row: tui.Row{
				tui.NewCell('1', enduranceFullStyle),
				tui.NewCell('0', enduranceEmptyStyle),
				tui.NewCell(' ', enduranceEmptyStyle),
				tui.NewCell(' ', gapStyle),
				tui.NewCell('1', willpowerFullStyle),
				tui.NewCell('0', willpowerEmptyStyle),
				tui.NewCell(' ', willpowerEmptyStyle),
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
				ShowCursorFunc:     func(_, _ int) {},
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

			vitalsi := len(rows) - 2
			if vitalsi < 0 {
				vitalsi = 0
			}
			assert.Equal(t, tc.row, rows[vitalsi])
		})
	}
}
