package tui_test

import (
	"testing"

	"github.com/tobiassjosten/nogfx/pkg/tui"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewText(t *testing.T) {
	baseStyle := tcell.StyleDefault
	redStyle := baseStyle.
		Foreground(tcell.ColorRed).
		Background(tcell.ColorYellow)
	greenStyle := baseStyle.
		Foreground(tcell.ColorGreen).
		Background(tcell.ColorBlue)
	blueStyle := baseStyle.
		Foreground(tcell.ColorBlue).
		Background(tcell.ColorTeal)

	tcs := map[string]struct {
		in      []byte
		width   int
		text    tui.Text
		styleIn tcell.Style
	}{
		"plain xy": {
			in:    []byte("xy"),
			width: 2,
			text: tui.Text{
				tui.Cell{'x', baseStyle, 1},
				tui.Cell{'y', baseStyle, 1},
			},
		},
		"plain yx": {
			in:    []byte("yx"),
			width: 2,
			text: tui.Text{
				tui.Cell{'y', baseStyle, 1},
				tui.Cell{'x', baseStyle, 1},
			},
		},
		"red style": {
			in:    []byte("yx"),
			width: 2,
			text: tui.Text{
				tui.Cell{'y', redStyle, 1},
				tui.Cell{'x', redStyle, 1},
			},
			styleIn: redStyle,
		},
		"change to green text": {
			in:    []byte("y\033[32;44mx"),
			width: 2,
			text: tui.Text{
				tui.Cell{'y', redStyle, 1},
				tui.Cell{'x', greenStyle, 1},
			},
			styleIn: redStyle,
		},
		"change to blue text": {
			in:    []byte("y\033[34;46mx"),
			width: 2,
			text: tui.Text{
				tui.Cell{'y', greenStyle, 1},
				tui.Cell{'x', blueStyle, 1},
			},
			styleIn: greenStyle,
		},
		"invalid ascii color": {
			in:    []byte("y\033{x"),
			width: 4,
			text: tui.Text{
				tui.Cell{'y', baseStyle, 1},
				tui.Cell{'^', baseStyle, 1},
				tui.Cell{'{', baseStyle, 1},
				tui.Cell{'x', baseStyle, 1},
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			text := tui.NewText(tc.in, tc.styleIn)
			assert.Equal(t, tc.text, text)
			assert.Equal(t, tc.width, text.Width())
		})
	}
}
