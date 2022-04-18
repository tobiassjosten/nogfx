package tui_test

import (
	"fmt"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/tui"
)

func TestOutputDraw(t *testing.T) {
	tcs := []struct {
		outputs [][]byte
		pos     []int
		x       int
		y       int
		width   int
		calls   int
		content map[int]map[int]rune
	}{
		{
			outputs: [][]byte{},
			pos:     []int{0, 0, 1, 1},
			calls:   0,
			content: map[int]map[int]rune{},
		},
		{
			outputs: [][]byte{{'x'}},
			pos:     []int{0, 0, 1, 1},
			calls:   1,
			content: map[int]map[int]rune{0: map[int]rune{0: 'x'}},
		},
		{
			outputs: [][]byte{{'x'}},
			pos:     []int{1, 2, 1, 1},
			calls:   1,
			content: map[int]map[int]rune{1: map[int]rune{2: 'x'}},
		},
		{
			outputs: [][]byte{{'a', 's', 'd', 'f'}},
			pos:     []int{5, 3, 2, 2},
			calls:   4,
			content: map[int]map[int]rune{
				5: map[int]rune{
					3: 'a',
					4: 'd',
				},
				6: map[int]rune{
					3: 's',
					4: 'f',
				},
			},
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			assert := assert.New(t)

			pane := tui.NewOutputPane()
			pane.Position(tc.pos[0], tc.pos[1], tc.pos[2], tc.pos[3])

			for _, output := range tc.outputs {
				pane.Add(output)
			}

			content := map[int]map[int]rune{}
			screen := &pkg.ScreenMock{
				SetContentFunc: func(x int, y int, r rune, _ []rune, _ tcell.Style) {
					if _, ok := content[x]; !ok {
						content[x] = map[int]rune{}
					}
					content[x][y] = r
				},
			}

			pane.Draw(screen)

			assert.Equal(tc.calls, len(screen.SetContentCalls()))
			assert.Equal(tc.content, content)
		})
	}
}
