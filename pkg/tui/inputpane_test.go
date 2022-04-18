package tui_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/tui"
)

func TestInput(t *testing.T) {
	newEventKey := func(key tcell.Key, r rune) tcell.Event {
		return tcell.NewEventKey(key, r, tcell.ModNone)
	}

	tcss := map[string][]struct {
		masked bool
		events []tcell.Event
		inputs [][]rune
	}{
		"invalid": {
			{
				events: []tcell.Event{
					tcell.NewEventResize(1, 1),
				},
				inputs: [][]rune{},
			},
			{
				events: []tcell.Event{
					newEventKey(tcell.KeyRune, '}'),
				},
				inputs: [][]rune{},
			},
			{
				events: []tcell.Event{
					newEventKey(tcell.KeyRune, ' '),
					newEventKey(tcell.KeyPrint, 0),
				},
				inputs: [][]rune{},
			},
		},
		"normal/input mode": {
			{
				events: []tcell.Event{
					newEventKey(tcell.KeyRune, 'a'),
					newEventKey(tcell.KeyEnter, 0),
				},
				inputs: [][]rune{},
			},
			{
				events: []tcell.Event{
					newEventKey(tcell.KeyRune, ' '),
					newEventKey(tcell.KeyRune, 'a'),
					newEventKey(tcell.KeyEnter, 0),
				},
				inputs: [][]rune{{'a'}},
			},
			{
				events: []tcell.Event{
					newEventKey(tcell.KeyRune, ' '),
					newEventKey(tcell.KeyEnter, 0),
				},
				inputs: [][]rune{{}},
			},
			{
				events: []tcell.Event{
					newEventKey(tcell.KeyRune, ' '),
					newEventKey(tcell.KeyRune, 'a'),
					newEventKey(tcell.KeyEnter, 0),
					newEventKey(tcell.KeyRune, 's'),
					newEventKey(tcell.KeyEnter, 0),
				},
				inputs: [][]rune{{'a'}, {'s'}},
			},
			{
				events: []tcell.Event{
					newEventKey(tcell.KeyRune, ' '),
					newEventKey(tcell.KeyRune, 'a'),
					newEventKey(tcell.KeyEsc, 0),
					newEventKey(tcell.KeyEnter, 0),
				},
				inputs: [][]rune{},
			},
			{
				events: []tcell.Event{
					newEventKey(tcell.KeyRune, ' '),
					newEventKey(tcell.KeyRune, 'a'),
					newEventKey(tcell.KeyEsc, 0),
					newEventKey(tcell.KeyRune, ' '),
					newEventKey(tcell.KeyEnter, 0),
				},
				inputs: [][]rune{{'a'}},
			},
			{
				events: []tcell.Event{
					newEventKey(tcell.KeyRune, ' '),
					newEventKey(tcell.KeyRune, 'a'),
					newEventKey(tcell.KeyCtrlC, 0),
					newEventKey(tcell.KeyRune, ' '),
					newEventKey(tcell.KeyEnter, 0),
				},
				inputs: [][]rune{{}},
			},
			{
				events: []tcell.Event{
					newEventKey(tcell.KeyRune, ' '),
					newEventKey(tcell.KeyRune, 'a'),
					newEventKey(tcell.KeyRune, 's'),
					newEventKey(tcell.KeyRune, 'd'),
					newEventKey(tcell.KeyRune, 'f'),
					newEventKey(tcell.KeyEnter, 0),
				},
				inputs: [][]rune{{'a', 's', 'd', 'f'}},
			},
			{
				events: []tcell.Event{
					newEventKey(tcell.KeyRune, ' '),
					newEventKey(tcell.KeyRune, 'a'),
					newEventKey(tcell.KeyRune, 's'),
					newEventKey(tcell.KeyEsc, 0),
					newEventKey(tcell.KeyRune, 'd'),
					newEventKey(tcell.KeyRune, 'f'),
					newEventKey(tcell.KeyRune, ' '),
					newEventKey(tcell.KeyEnter, 0),
				},
				inputs: [][]rune{{'a', 's'}},
			},
		},

		"inputted": {
			{
				events: []tcell.Event{
					newEventKey(tcell.KeyRune, ' '),
					newEventKey(tcell.KeyRune, 'a'),
					newEventKey(tcell.KeyEnter, 0),
					newEventKey(tcell.KeyEnter, 0),
				},
				inputs: [][]rune{{'a'}, {'a'}},
			},
			{
				masked: true,
				events: []tcell.Event{
					newEventKey(tcell.KeyRune, ' '),
					newEventKey(tcell.KeyRune, 'a'),
					newEventKey(tcell.KeyEnter, 0),
					newEventKey(tcell.KeyEnter, 0),
				},
				inputs: [][]rune{{'a'}, {}},
			},
		},

		"backspaces": {
			{
				events: []tcell.Event{
					newEventKey(tcell.KeyRune, ' '),
					newEventKey(tcell.KeyRune, 'a'),
					newEventKey(tcell.KeyRune, 's'),
					newEventKey(tcell.KeyRune, 'd'),
					newEventKey(tcell.KeyRune, 'f'),
					newEventKey(tcell.KeyBackspace, 0),
					newEventKey(tcell.KeyEnter, 0),
				},
				inputs: [][]rune{{'a', 's', 'd'}},
			},
			{
				events: []tcell.Event{
					newEventKey(tcell.KeyRune, ' '),
					newEventKey(tcell.KeyRune, 'a'),
					newEventKey(tcell.KeyRune, 's'),
					newEventKey(tcell.KeyRune, ' '),
					newEventKey(tcell.KeyRune, 'd'),
					newEventKey(tcell.KeyRune, 'f'),
					newEventKey(tcell.KeyETB, 0),
					newEventKey(tcell.KeyEnter, 0),
				},
				inputs: [][]rune{{'a', 's', ' '}},
			},
			{
				events: []tcell.Event{
					newEventKey(tcell.KeyRune, ' '),
					newEventKey(tcell.KeyRune, 'a'),
					newEventKey(tcell.KeyRune, 's'),
					newEventKey(tcell.KeyRune, 'd'),
					newEventKey(tcell.KeyRune, 'f'),
					newEventKey(tcell.KeyETB, 0),
					newEventKey(tcell.KeyEnter, 0),
				},
				inputs: [][]rune{{}},
			},
			{
				events: []tcell.Event{
					newEventKey(tcell.KeyRune, ' '),
					newEventKey(tcell.KeyRune, 'a'),
					newEventKey(tcell.KeyRune, 's'),
					newEventKey(tcell.KeyRune, ' '),
					newEventKey(tcell.KeyRune, 'd'),
					newEventKey(tcell.KeyRune, 'f'),
					newEventKey(tcell.KeyNAK, 0),
					newEventKey(tcell.KeyEnter, 0),
				},
				inputs: [][]rune{{}},
			},
		},

		"backspaces from start": {
			{
				events: []tcell.Event{
					newEventKey(tcell.KeyRune, ' '),
					newEventKey(tcell.KeyRune, 'a'),
					newEventKey(tcell.KeyLeft, 0),
					newEventKey(tcell.KeyBackspace, 0),
					newEventKey(tcell.KeyEnter, 0),
				},
				inputs: [][]rune{{'a'}},
			},
			{
				events: []tcell.Event{
					newEventKey(tcell.KeyRune, ' '),
					newEventKey(tcell.KeyRune, 'a'),
					newEventKey(tcell.KeyLeft, 0),
					newEventKey(tcell.KeyRight, 0),
					newEventKey(tcell.KeyBackspace, 0),
					newEventKey(tcell.KeyEnter, 0),
				},
				inputs: [][]rune{{}},
			},
			{
				events: []tcell.Event{
					newEventKey(tcell.KeyRune, ' '),
					newEventKey(tcell.KeyRune, 'a'),
					newEventKey(tcell.KeyLeft, 0),
					newEventKey(tcell.KeyETB, 0),
					newEventKey(tcell.KeyEnter, 0),
				},
				inputs: [][]rune{{'a'}},
			},
			{
				events: []tcell.Event{
					newEventKey(tcell.KeyRune, ' '),
					newEventKey(tcell.KeyRune, 'a'),
					newEventKey(tcell.KeyLeft, 0),
					newEventKey(tcell.KeyRight, 0),
					newEventKey(tcell.KeyETB, 0),
					newEventKey(tcell.KeyEnter, 0),
				},
				inputs: [][]rune{{}},
			},
			{
				events: []tcell.Event{
					newEventKey(tcell.KeyRune, ' '),
					newEventKey(tcell.KeyRune, 'a'),
					newEventKey(tcell.KeyLeft, 0),
					newEventKey(tcell.KeyNAK, 0),
					newEventKey(tcell.KeyEnter, 0),
				},
				inputs: [][]rune{{'a'}},
			},
			{
				events: []tcell.Event{
					newEventKey(tcell.KeyRune, ' '),
					newEventKey(tcell.KeyRune, 'a'),
					newEventKey(tcell.KeyLeft, 0),
					newEventKey(tcell.KeyRight, 0),
					newEventKey(tcell.KeyNAK, 0),
					newEventKey(tcell.KeyEnter, 0),
				},
				inputs: [][]rune{{}},
			},
		},

		"backspaces with inputted": {
			{
				events: []tcell.Event{
					newEventKey(tcell.KeyRune, ' '),
					newEventKey(tcell.KeyRune, 'a'),
					newEventKey(tcell.KeyRune, 's'),
					newEventKey(tcell.KeyEnter, 0),
					newEventKey(tcell.KeyBackspace, 0),
					newEventKey(tcell.KeyEnter, 0),
				},
				inputs: [][]rune{{'a', 's'}, {}},
			},
			{
				events: []tcell.Event{
					newEventKey(tcell.KeyRune, ' '),
					newEventKey(tcell.KeyRune, 'a'),
					newEventKey(tcell.KeyRune, 's'),
					newEventKey(tcell.KeyEnter, 0),
					newEventKey(tcell.KeyETB, 0),
					newEventKey(tcell.KeyEnter, 0),
				},
				inputs: [][]rune{{'a', 's'}, {}},
			},
			{
				events: []tcell.Event{
					newEventKey(tcell.KeyRune, ' '),
					newEventKey(tcell.KeyRune, 'a'),
					newEventKey(tcell.KeyRune, 's'),
					newEventKey(tcell.KeyEnter, 0),
					newEventKey(tcell.KeyNAK, 0),
					newEventKey(tcell.KeyEnter, 0),
				},
				inputs: [][]rune{{'a', 's'}, {}},
			},
		},

		"bindings": {
			{
				events: []tcell.Event{
					newEventKey(tcell.KeyRune, '1'),
				},
				inputs: [][]rune{{'s', 'w'}},
			},
			{
				events: []tcell.Event{
					newEventKey(tcell.KeyRune, '2'),
				},
				inputs: [][]rune{{'s'}},
			},
			{
				events: []tcell.Event{
					newEventKey(tcell.KeyRune, '3'),
				},
				inputs: [][]rune{{'s', 'e'}},
			},
			{
				events: []tcell.Event{
					newEventKey(tcell.KeyRune, '4'),
				},
				inputs: [][]rune{{'w'}},
			},
			{
				events: []tcell.Event{
					newEventKey(tcell.KeyRune, '5'),
				},
				inputs: [][]rune{{'m', 'a', 'p'}},
			},
			{
				events: []tcell.Event{
					newEventKey(tcell.KeyRune, '6'),
				},
				inputs: [][]rune{{'e'}},
			},
			{
				events: []tcell.Event{
					newEventKey(tcell.KeyRune, '7'),
				},
				inputs: [][]rune{{'n', 'w'}},
			},
			{
				events: []tcell.Event{
					newEventKey(tcell.KeyRune, '8'),
				},
				inputs: [][]rune{{'n'}},
			},
			{
				events: []tcell.Event{
					newEventKey(tcell.KeyRune, '9'),
				},
				inputs: [][]rune{{'n', 'e'}},
			},
		},
	}

	for group, tcs := range tcss {
		for i, tc := range tcs {
			t.Run(fmt.Sprintf("%s/case %d", group, i), func(t *testing.T) {
				assert := assert.New(t)

				pane := tui.NewInputPane()

				if tc.masked {
					pane.Mask()
				}

				inputs := [][]rune{}
				for _, event := range tc.events {
					handled, input := pane.HandleEvent(event)
					if handled && len(input) > 0 {
						inputs = append(inputs, []rune(strings.TrimRight(string(input), "\n")))
					}
				}

				assert.Equal(tc.inputs, inputs)
			})
		}
	}
}

func TestInputDraw(t *testing.T) {
	newEventKey := func(key tcell.Key, r rune) tcell.Event {
		return tcell.NewEventKey(key, r, tcell.ModNone)
	}

	tcs := []struct {
		pos     []int
		events  []tcell.Event
		content map[int]map[int]rune
		cursor  []int
		masked  bool
		height  int
	}{
		{ // Normal mode gives no output.
			pos:     []int{1, 2, 2, 2},
			content: map[int]map[int]rune{},
			cursor:  []int{-1, -1},
			height:  0,
		},
		{ // Returning to normal mode gives no output.
			pos: []int{1, 2, 2, 2},
			events: []tcell.Event{
				newEventKey(tcell.KeyRune, ' '),
				newEventKey(tcell.KeyEsc, 0),
			},
			content: map[int]map[int]rune{},
			cursor:  []int{-1, -1},
			height:  0,
		},
		{ // Pane is padded with spaces.
			pos: []int{1, 2, 2, 2},
			events: []tcell.Event{
				newEventKey(tcell.KeyRune, ' '),
				newEventKey(tcell.KeyRune, 'a'),
				newEventKey(tcell.KeyRune, 'b'),
			},
			content: map[int]map[int]rune{
				1: map[int]rune{
					2: 'a',
				},
				2: map[int]rune{
					2: 'b',
				},
			},
			cursor: []int{3, 2},
			height: 1,
		},
		{ // Hitting enter doesn't change output.
			pos: []int{1, 2, 2, 2},
			events: []tcell.Event{
				newEventKey(tcell.KeyRune, ' '),
				newEventKey(tcell.KeyRune, 'a'),
				newEventKey(tcell.KeyRune, 'b'),
				newEventKey(tcell.KeyEnter, 0),
			},
			content: map[int]map[int]rune{
				1: map[int]rune{
					2: 'a',
				},
				2: map[int]rune{
					2: 'b',
				},
			},
			cursor: []int{3, 2},
			height: 1,
		},
		{ // Pane position controls output coordinates.
			pos: []int{2, 1, 2, 2},
			events: []tcell.Event{
				newEventKey(tcell.KeyRune, ' '),
				newEventKey(tcell.KeyRune, 'a'),
				newEventKey(tcell.KeyRune, 'b'),
			},
			content: map[int]map[int]rune{
				2: map[int]rune{
					1: 'a',
				},
				3: map[int]rune{
					1: 'b',
				},
			},
			cursor: []int{4, 1},
			height: 1,
		},
		{ // Masked mode turns input (bot not padding) into stars.
			pos: []int{1, 2, 2, 2},
			events: []tcell.Event{
				newEventKey(tcell.KeyRune, ' '),
				newEventKey(tcell.KeyRune, 'a'),
			},
			content: map[int]map[int]rune{
				1: map[int]rune{
					2: '*',
				},
				2: map[int]rune{
					2: ' ',
				},
			},
			cursor: []int{2, 2},
			masked: true,
			height: 1,
		},
		{ // Words are wrapped to new lines.
			pos: []int{0, 0, 3, 3},
			events: []tcell.Event{
				newEventKey(tcell.KeyRune, ' '),
				newEventKey(tcell.KeyRune, 'a'),
				newEventKey(tcell.KeyRune, ' '),
				newEventKey(tcell.KeyRune, 's'),
				newEventKey(tcell.KeyRune, 'd'),
			},
			content: map[int]map[int]rune{
				0: map[int]rune{
					0: 'a',
					1: 's',
				},
				1: map[int]rune{
					0: ' ',
					1: 'd',
				},
				2: map[int]rune{
					0: ' ',
					1: ' ',
				},
			},
			cursor: []int{2, 1},
			height: 2,
		},
		{ // Line-length words also wrap to new lines.
			pos: []int{0, 0, 3, 3},
			events: []tcell.Event{
				newEventKey(tcell.KeyRune, ' '),
				newEventKey(tcell.KeyRune, 'a'),
				newEventKey(tcell.KeyRune, 's'),
				newEventKey(tcell.KeyRune, 'd'),
				newEventKey(tcell.KeyRune, ' '),
				newEventKey(tcell.KeyRune, 'f'),
			},
			content: map[int]map[int]rune{
				0: map[int]rune{
					0: 'a',
					1: 'f',
				},
				1: map[int]rune{
					0: 's',
					1: ' ',
				},
				2: map[int]rune{
					0: 'd',
					1: ' ',
				},
			},
			cursor: []int{1, 1},
			height: 2,
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			assert := assert.New(t)

			pane := tui.NewInputPane()
			pane.Position(tc.pos[0], tc.pos[1], tc.pos[2], tc.pos[3])

			if tc.masked {
				pane.Mask()
			}

			content := map[int]map[int]rune{}
			cursor := []int{}

			screen := &pkg.ScreenMock{
				HideCursorFunc: func() {
					cursor = []int{-1, -1}
				},
				SetContentFunc: func(x int, y int, r rune, _ []rune, _ tcell.Style) {
					if _, ok := content[x]; !ok {
						content[x] = map[int]rune{}
					}
					content[x][y] = r
				},
				ShowCursorFunc: func(x, y int) {
					cursor = []int{x, y}
				},
			}

			for _, event := range tc.events {
				_, _ = pane.HandleEvent(event)
			}

			pane.Draw(screen)

			assert.Equal(tc.content, content)
			assert.Equal(tc.height, pane.Height())
			assert.Equal(tc.cursor, cursor)
		})
	}
}
