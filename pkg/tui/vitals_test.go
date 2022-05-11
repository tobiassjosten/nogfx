package tui_test

import (
	"fmt"
	"testing"

	"github.com/tobiassjosten/nogfx/pkg/mock"
	"github.com/tobiassjosten/nogfx/pkg/tui"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
)

func TestVitalsDraw(t *testing.T) {
	tcs := []struct {
		pos     []int
		vitals  []map[string]tui.Vital
		updates map[string][]int
		content map[int]map[int]rune
	}{
		{
			pos:     []int{0, 0, 1},
			content: map[int]map[int]rune{},
		},
		{
			pos: []int{0, 0, 1},
			vitals: []map[string]tui.Vital{
				{"health": tui.HealthVital},
			},
			updates: map[string][]int{
				"health": []int{4, 4},
			},
			content: map[int]map[int]rune{
				0: map[int]rune{
					0: '4',
				},
			},
		},
		{
			pos: []int{0, 0, 2},
			vitals: []map[string]tui.Vital{
				{"health": tui.HealthVital},
			},
			updates: map[string][]int{
				"health": []int{4, 4},
			},
			content: map[int]map[int]rune{
				0: map[int]rune{
					0: '4',
				},
				1: map[int]rune{
					0: ' ',
				},
			},
		},
		{
			pos: []int{0, 0, 2},
			vitals: []map[string]tui.Vital{
				{"health": tui.HealthVital},
			},
			updates: map[string][]int{
				"health": []int{1, 4},
			},
			content: map[int]map[int]rune{
				0: map[int]rune{
					0: '1',
				},
				1: map[int]rune{
					0: ' ',
				},
			},
		},
		{
			pos: []int{0, 0, 1},
			vitals: []map[string]tui.Vital{
				{"health": tui.HealthVital},
			},
			updates: map[string][]int{
				"health": []int{4, 4},
				"asdf":   []int{5, 5},
			},
			content: map[int]map[int]rune{
				0: map[int]rune{
					0: '4',
				},
			},
		},
		{
			pos: []int{0, 0, 3},
			vitals: []map[string]tui.Vital{
				{"health": tui.HealthVital},
				{"mana": tui.ManaVital},
			},
			updates: map[string][]int{
				"health": []int{4, 4},
				"mana":   []int{5, 5},
			},
			content: map[int]map[int]rune{
				0: map[int]rune{
					0: '4',
				},
				1: map[int]rune{
					0: ' ',
				},
				2: map[int]rune{
					0: '5',
				},
			},
		},
		{
			pos: []int{0, 0, 3},
			vitals: []map[string]tui.Vital{
				{"mana": tui.ManaVital},
				{"health": tui.HealthVital},
			},
			updates: map[string][]int{
				"health": []int{4, 4},
				"mana":   []int{5, 5},
			},
			content: map[int]map[int]rune{
				0: map[int]rune{
					0: '5',
				},
				1: map[int]rune{
					0: ' ',
				},
				2: map[int]rune{
					0: '4',
				},
			},
		},
		{
			pos: []int{0, 0, 4},
			vitals: []map[string]tui.Vital{
				{"mana": tui.ManaVital},
				{"health": tui.HealthVital},
			},
			updates: map[string][]int{
				"health": []int{4, 4},
				"mana":   []int{5, 5},
			},
			content: map[int]map[int]rune{
				0: map[int]rune{
					0: '5',
				},
				1: map[int]rune{
					0: ' ',
				},
				2: map[int]rune{
					0: ' ',
				},
				3: map[int]rune{
					0: '4',
				},
			},
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			assert := assert.New(t)

			pane := tui.NewVitalsPane()

			for _, pairs := range tc.vitals {
				for name, vital := range pairs {
					pane.AddVital(name, vital)
					// Prove idempotency.
					pane.AddVital(name, vital)
				}
			}

			for name, update := range tc.updates {
				pane.UpdateVital(name, update[0], update[1])
			}

			content := map[int]map[int]rune{}

			screen := &mock.ScreenMock{
				SetContentFunc: func(x int, y int, r rune, _ []rune, _ tcell.Style) {
					if _, ok := content[x]; !ok {
						content[x] = map[int]rune{}
					}
					content[x][y] = r
				},
			}

			pane.Position(tc.pos[0], tc.pos[1], tc.pos[2], pane.Height())
			pane.Draw(screen)

			assert.Equal(tc.content, content)
		})
	}
}
