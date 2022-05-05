package tui_test

import (
	"fmt"
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

	tcs := []struct {
		in      []byte
		width   int
		text    tui.Text
		styleIn tcell.Style
	}{
		{
			in:    []byte("xy"),
			width: 2,
			text: tui.Text{
				tui.Cell{'x', baseStyle, 1},
				tui.Cell{'y', baseStyle, 1},
			},
		},
		{
			in:    []byte("yx"),
			width: 2,
			text: tui.Text{
				tui.Cell{'y', baseStyle, 1},
				tui.Cell{'x', baseStyle, 1},
			},
		},
		{
			in:    []byte("yx"),
			width: 2,
			text: tui.Text{
				tui.Cell{'y', redStyle, 1},
				tui.Cell{'x', redStyle, 1},
			},
			styleIn: redStyle,
		},
		{
			in:    []byte("y\r\nx"),
			width: 2,
			text: tui.Text{
				tui.Cell{'y', baseStyle, 1},
				tui.Cell{'\n', baseStyle, 0},
				tui.Cell{'x', baseStyle, 1},
			},
		},
		{
			in:    []byte("y\033[32;44mx"),
			width: 2,
			text: tui.Text{
				tui.Cell{'y', redStyle, 1},
				tui.Cell{'x', greenStyle, 1},
			},
			styleIn: redStyle,
		},
		{
			in:    []byte("y\033[34;46mx"),
			width: 2,
			text: tui.Text{
				tui.Cell{'y', greenStyle, 1},
				tui.Cell{'x', blueStyle, 1},
			},
			styleIn: greenStyle,
		},
		{
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

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			assert := assert.New(t)

			text := tui.NewText(tc.in, tc.styleIn)
			assert.Equal(tc.text, text)
			assert.Equal(tc.width, text.Width())
		})
	}
}
