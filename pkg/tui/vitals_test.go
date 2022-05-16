package tui_test

import (
	"context"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/tobiassjosten/nogfx/pkg/mock"
	"github.com/tobiassjosten/nogfx/pkg/tui"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenderVital(t *testing.T) {
	healthFullStyle := tcell.StyleDefault.
		Background(tcell.ColorGreen).
		Foreground(tcell.ColorBlack)

	manaFullStyle := tcell.StyleDefault.
		Background(tcell.ColorBlue).
		Foreground(tcell.ColorBlack)

	manaEmptyStyle := tcell.StyleDefault.
		Background(tcell.ColorDarkBlue).
		Foreground(tcell.ColorBlack)

	enduranceFullStyle := tcell.StyleDefault.
		Background(tcell.ColorTeal).
		Foreground(tcell.ColorBlack)

	enduranceEmptyStyle := tcell.StyleDefault.
		Background(tcell.ColorDarkCyan).
		Foreground(tcell.ColorBlack)

	willpowerFullStyle := tcell.StyleDefault.
		Background(tcell.ColorFuchsia).
		Foreground(tcell.ColorBlack)

	willpowerEmptyStyle := tcell.StyleDefault.
		Background(tcell.ColorRebeccaPurple).
		Foreground(tcell.ColorBlack)

	energyEmptyStyle := tcell.StyleDefault.
		Background(tcell.ColorRosyBrown).
		Foreground(tcell.ColorBlack)

	tcs := map[string]struct {
		vital *tui.Vital
		width int
		row   tui.Row
	}{
		"100% health, even": {
			vital: tui.NewHealthVital(20, 20),
			width: 4,
			row: tui.Row{
				tui.NewCell(' ', healthFullStyle),
				tui.NewCell('2', healthFullStyle),
				tui.NewCell('0', healthFullStyle),
				tui.NewCell(' ', healthFullStyle),
			},
		},

		"100% health, odd": {
			vital: tui.NewHealthVital(20, 20),
			width: 5,
			row: tui.Row{
				tui.NewCell(' ', healthFullStyle),
				tui.NewCell('2', healthFullStyle),
				tui.NewCell('0', healthFullStyle),
				tui.NewCell(' ', healthFullStyle),
				tui.NewCell(' ', healthFullStyle),
			},
		},

		"100% health, cramped odd": {
			vital: tui.NewHealthVital(20, 20),
			width: 3,
			row: tui.Row{
				tui.NewCell('2', healthFullStyle),
				tui.NewCell('0', healthFullStyle),
				tui.NewCell(' ', healthFullStyle),
			},
		},

		"100% health, cramped even": {
			vital: tui.NewHealthVital(20, 20),
			width: 2,
			row: tui.Row{
				tui.NewCell('2', healthFullStyle),
				tui.NewCell('0', healthFullStyle),
			},
		},

		"100% health, too cramped": {
			vital: tui.NewHealthVital(200, 200),
			width: 2,
			row: tui.Row{
				tui.NewCell(' ', healthFullStyle),
				tui.NewCell(' ', healthFullStyle),
			},
		},

		"75% mana": {
			vital: tui.NewManaVital(15, 20),
			width: 4,
			row: tui.Row{
				tui.NewCell(' ', manaFullStyle),
				tui.NewCell('1', manaFullStyle),
				tui.NewCell('5', manaFullStyle),
				tui.NewCell(' ', manaEmptyStyle),
			},
		},

		"50% endurance": {
			vital: tui.NewEnduranceVital(10, 20),
			width: 4,
			row: tui.Row{
				tui.NewCell(' ', enduranceFullStyle),
				tui.NewCell('1', enduranceFullStyle),
				tui.NewCell('0', enduranceEmptyStyle),
				tui.NewCell(' ', enduranceEmptyStyle),
			},
		},

		"25% willpower": {
			vital: tui.NewWillpowerVital(5, 20),
			width: 4,
			row: tui.Row{
				tui.NewCell(' ', willpowerFullStyle),
				tui.NewCell('5', willpowerEmptyStyle),
				tui.NewCell(' ', willpowerEmptyStyle),
				tui.NewCell(' ', willpowerEmptyStyle),
			},
		},

		"0% energy": {
			vital: tui.NewEnergyVital(0, 20),
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
			row := tui.RenderVital(tc.vital, tc.width)
			assert.Equal(t, tc.row, row)
		})
	}
}

type MockVital struct {
}

func TestRenderVitals(t *testing.T) {
	healthFullStyle := tcell.StyleDefault.
		Background(tcell.ColorGreen).
		Foreground(tcell.ColorBlack)

	healthEmptyStyle := tcell.StyleDefault.
		Background(tcell.ColorDarkGreen).
		Foreground(tcell.ColorBlack)

	manaFullStyle := tcell.StyleDefault.
		Background(tcell.ColorBlue).
		Foreground(tcell.ColorBlack)

	manaEmptyStyle := tcell.StyleDefault.
		Background(tcell.ColorDarkBlue).
		Foreground(tcell.ColorBlack)

	enduranceFullStyle := tcell.StyleDefault.
		Background(tcell.ColorTeal).
		Foreground(tcell.ColorBlack)

	enduranceEmptyStyle := tcell.StyleDefault.
		Background(tcell.ColorDarkCyan).
		Foreground(tcell.ColorBlack)

	willpowerFullStyle := tcell.StyleDefault.
		Background(tcell.ColorFuchsia).
		Foreground(tcell.ColorBlack)

	willpowerEmptyStyle := tcell.StyleDefault.
		Background(tcell.ColorRebeccaPurple).
		Foreground(tcell.ColorBlack)

	tcs := map[string]struct {
		vorder  []string
		vitals  map[string]interface{}
		updates map[string][]int
		width   int
		height  int
		rows    tui.Rows
		err     string
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
				"health": tui.NewHealthVital(1, 2),
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
				"health": tui.NewHealthVital(100, 200),
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
				"health": tui.NewHealthVital(1, 2),
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
				"health": tui.NewHealthVital(2, 2),
				"mana":   tui.NewManaVital(1, 2),
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
				"endurance": tui.NewEnduranceVital(1, 2),
				"willpower": tui.NewWillpowerVital(2, 2),
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
				"endurance": tui.NewEnduranceVital(1, 2),
				"willpower": tui.NewWillpowerVital(2, 2),
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

			var err error
			for _, name := range tc.vorder {
				if vital, ok := tc.vitals[name]; ok {
					if e := ui.AddVital(name, vital); e != nil {
						err = e
					}
				}
			}
			for _, name := range tc.vorder {
				if update, ok := tc.updates[name]; ok {
					if e := ui.UpdateVital(name, update[0], update[1]); e != nil {
						err = e
					}
				}
			}

			if tc.err != "" {
				assert.Equal(t, tc.err, err.Error())
				return
			}
			require.Nil(t, err)

			_ = ui.Run(ctx)

			assert.Equal(t, tc.rows, rows)
		})
	}
}
