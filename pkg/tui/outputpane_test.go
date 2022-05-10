package tui_test

import (
	"fmt"
	"testing"

	"github.com/tobiassjosten/nogfx/pkg/mock"
	"github.com/tobiassjosten/nogfx/pkg/tui"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
)

func TestOutputRows(t *testing.T) {
	rowToRunes := func(row tui.Row) (rs []rune) {
		for _, cell := range row {
			rs = append(rs, cell.Content)
		}
		return
	}

	rowsToRunes := func(rows tui.Rows) (rss [][]rune) {
		for _, row := range rows {
			rss = append(rss, rowToRunes(row))
		}
		return
	}

	tcs := map[string]struct {
		outputs [][]byte
		events  []*tcell.EventKey
		size    []int
		output  [][]rune
		history [][]rune
	}{
		"no output": {
			outputs: [][]byte{},
			size:    []int{1, 1},
		},

		"single cell output": {
			outputs: [][]byte{{'x'}},
			size:    []int{1, 1},
			output:  [][]rune{{'x'}},
		},

		"correct order plain": {
			outputs: [][]byte{{'x'}, {'y'}, {'z'}},
			size:    []int{1, 3},
			output:  [][]rune{{'x'}, {'y'}, {'z'}},
		},

		"correct order scroll back one": {
			outputs: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}, {'e'}},
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyUp, 0, 0),
			},
			size:    []int{1, 3},
			output:  [][]rune{{'e'}},
			history: [][]rune{{'b'}, {tcell.RuneHLine}},
		},

		"correct order scroll back two": {
			outputs: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}, {'e'}},
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyUp, 0, 0),
				tcell.NewEventKey(tcell.KeyUp, 0, 0),
			},
			size:    []int{1, 3},
			output:  [][]rune{{'e'}},
			history: [][]rune{{'a'}, {tcell.RuneHLine}},
		},

		"no width": {
			outputs: [][]byte{{'x'}},
			size:    []int{0, 1},
		},

		"no height": {
			outputs: [][]byte{{'x'}},
			size:    []int{1, 0},
		},

		"linebreak": {
			outputs: [][]byte{{'a', 's', 'd', 'f'}},
			size:    []int{2, 2},
			output:  [][]rune{{'a', 's'}, {'d', 'f'}},
		},

		"wordwrap one": {
			outputs: [][]byte{{'a', 's', ' ', 'd', 'f'}},
			size:    []int{2, 2},
			output:  [][]rune{{'a', 's'}, {'d', 'f'}},
		},

		"wordwrap two": {
			outputs: [][]byte{{'a', 's', ' ', 'd', 'f'}},
			size:    []int{3, 2},
			output:  [][]rune{{'a', 's'}, {'d', 'f'}},
		},

		"wordwrap long word": {
			outputs: [][]byte{{'a', 's', ' ', 'd', 'f', 'g'}},
			size:    []int{2, 3},
			output:  [][]rune{{'a', 's'}, {'d', 'f'}, {'g'}},
		},

		"wordwrap long word scroll back one": {
			outputs: [][]byte{{'a', 's', ' ', 'd', 'f', 'g'}},
			events: []*tcell.EventKey{
				tcell.NewEventKey(tcell.KeyUp, 0, 0),
			},
			size:    []int{2, 2},
			output:  [][]rune{{'g'}},
			history: [][]rune{{'a', 's'}},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			pane := tui.NewOutputPane()

			for _, output := range tc.outputs {
				pane.Add(output)
			}

			for _, event := range tc.events {
				pane.HandleEvent(event)
			}

			outputRows, historyRows := pane.Rows(tc.size[0], tc.size[1])
			output := rowsToRunes(outputRows)
			history := rowsToRunes(historyRows)

			assert.Equal(t, tc.output, output)
			assert.Equal(t, tc.history, history)
		})
	}
}

func TestOutputsChannel(t *testing.T) {
	tcs := []struct {
		output []byte
	}{
		{
			output: []byte("asdf"),
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			screen := &mock.ScreenMock{
				SetStyleFunc:       func(_ tcell.Style) {},
				SetCursorStyleFunc: func(_ tcell.CursorStyle) {},
			}

			pane := tui.NewOutputPane()

			ui := tui.NewTUI(screen, tui.Panes{Output: pane})

			output := []byte{}

			done := make(chan struct{})
			go func() {
				output = <-pane.Outputs()
				done <- struct{}{}
			}()

			ui.Outputs() <- tc.output

			assert.Equal(t, tc.output, output)
		})
	}
}
