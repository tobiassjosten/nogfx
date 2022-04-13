package tui_test

import (
	"fmt"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
	"github.com/tobiassjosten/nogfx/pkg/tui"
)

func TestNewText(t *testing.T) {
	assert := assert.New(t)

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

	tcs := []struct {
		output   []byte
		text     tui.Text
		styleIn  tcell.Style
		styleOut tcell.Style
	}{
		{
			output: []byte("xy"),
			text: tui.Text{
				tui.Cell{'x', baseStyle, 1},
				tui.Cell{'y', baseStyle, 1},
			},
		},
		{
			output: []byte("yx"),
			text: tui.Text{
				tui.Cell{'y', baseStyle, 1},
				tui.Cell{'x', baseStyle, 1},
			},
		},
		{
			output: []byte("yx"),
			text: tui.Text{
				tui.Cell{'y', redStyle, 1},
				tui.Cell{'x', redStyle, 1},
			},
			styleIn:  redStyle,
			styleOut: redStyle,
		},
		{
			output: []byte("y\033[32;44mx"),
			text: tui.Text{
				tui.Cell{'y', redStyle, 1},
				tui.Cell{'x', greenStyle, 1},
			},
			styleIn:  redStyle,
			styleOut: greenStyle,
		},
		{
			output: []byte("y\033[34;46mx"),
			text: tui.Text{
				tui.Cell{'y', greenStyle, 1},
				tui.Cell{'x', blueStyle, 1},
			},
			styleIn:  greenStyle,
			styleOut: blueStyle,
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			text, style := tui.NewText(tc.output, tc.styleIn)
			assert.Equal(tc.text, text)
			assert.Equal(tc.styleOut, style)
		})
	}
}
