package tui_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
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

				pane := tui.NewInputPane(tcell.Style{}, tcell.Style{})

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

// @todo TestDraw
