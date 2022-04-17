package tui_test

import (
	"fmt"
	"testing"

	"github.com/tobiassjosten/nogfx/pkg/tui"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
)

func TestApplyANSI(t *testing.T) {
	tcs := []struct {
		in   tcell.Style
		ansi int
		out  tcell.Style
	}{
		{
			in:   tcell.Style{},
			ansi: 123,
			out:  tcell.Style{},
		},
		{
			in: (tcell.Style{}).
				Foreground(tcell.ColorRed).
				Background(tcell.ColorBlue).
				Attributes(tcell.AttrBold),
			ansi: 0,
			out:  tcell.Style{},
		},
		{
			in:   tcell.Style{},
			ansi: 1,
			out:  (tcell.Style{}).Attributes(tcell.AttrBold),
		},
		{
			in:   tcell.Style{},
			ansi: 2,
			out:  (tcell.Style{}).Attributes(tcell.AttrDim),
		},
		{
			in:   tcell.Style{},
			ansi: 3,
			out:  (tcell.Style{}).Attributes(tcell.AttrItalic),
		},
		{
			in:   tcell.Style{},
			ansi: 4,
			out:  (tcell.Style{}).Attributes(tcell.AttrUnderline),
		},
		{
			in:   tcell.Style{},
			ansi: 5,
			out:  (tcell.Style{}).Attributes(tcell.AttrBlink),
		},
		{
			in:   tcell.Style{},
			ansi: 7,
			out:  (tcell.Style{}).Attributes(tcell.AttrReverse),
		},
		{
			in: (tcell.Style{}).Attributes(
				tcell.AttrBold | tcell.AttrDim | tcell.AttrBlink,
			),
			ansi: 22,
			out:  (tcell.Style{}).Attributes(tcell.AttrBlink),
		},
		{
			in: (tcell.Style{}).
				Attributes(tcell.AttrItalic | tcell.AttrBlink),
			ansi: 23,
			out:  (tcell.Style{}).Attributes(tcell.AttrBlink),
		},
		{
			in: (tcell.Style{}).
				Attributes(tcell.AttrUnderline | tcell.AttrBlink),
			ansi: 24,
			out:  (tcell.Style{}).Attributes(tcell.AttrBlink),
		},
		{
			in: (tcell.Style{}).
				Attributes(tcell.AttrBlink | tcell.AttrBold),
			ansi: 25,
			out:  (tcell.Style{}).Attributes(tcell.AttrBold),
		},
		{
			in: (tcell.Style{}).
				Attributes(tcell.AttrReverse | tcell.AttrBlink),
			ansi: 27,
			out:  (tcell.Style{}).Attributes(tcell.AttrBlink),
		},
		{
			in:   tcell.Style{},
			ansi: 91,
			out: (tcell.Style{}).
				Foreground(tcell.ColorRed).
				Attributes(tcell.AttrBold),
		},
		{
			in:   tcell.Style{},
			ansi: 102,
			out: (tcell.Style{}).
				Background(tcell.ColorGreen).
				Attributes(tcell.AttrBold),
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			assert := assert.New(t)
			assert.Equal(tc.out, tui.ApplyANSI(tc.in, tc.ansi))
		})
	}
}
