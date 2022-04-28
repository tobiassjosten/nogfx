package tui_test

import (
	"fmt"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
	"github.com/tobiassjosten/nogfx/pkg/mock"
	"github.com/tobiassjosten/nogfx/pkg/tui"
)

func TestBlockProcessing(t *testing.T) {
	newText := func(output string) tui.Text {
		text, _ := tui.NewText([]byte(output), tcell.Style{})
		return text
	}

	tcs := []struct {
		text   tui.Text
		width  int
		height int
	}{
		{
			text:   tui.Text{},
			width:  0,
			height: 0,
		},
		{
			text:   newText("x"),
			width:  0,
			height: 0,
		},
		{
			text:   newText("x"),
			width:  1,
			height: 1,
		},
		{
			text:   newText("asdf"),
			width:  1,
			height: 4,
		},
		{
			text:   newText("asdf\n"),
			width:  1,
			height: 5,
		},
		{
			text:   newText("a\ns\nd\nf"),
			width:  1,
			height: 4,
		},
		{
			text:   newText("a\ns\nd\nf\n"),
			width:  2,
			height: 5,
		},
		{
			text:   newText("a s d f "),
			width:  2,
			height: 4,
		},
		{
			text:   newText("a s d f"),
			width:  2,
			height: 4,
		},
		{
			text:   newText("a s d f "),
			width:  1,
			height: 8,
		},
		{
			text:   newText("a \ns \nd \nf "),
			width:  2,
			height: 4,
		},
		{
			text:   newText("a \ns \nd \nf \n"),
			width:  2,
			height: 5,
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			assert := assert.New(t)

			block := tui.NewBlock(tc.text, tc.width)
			assert.Equal(tc.width, block.Width())
			assert.Equal(tc.height, block.Height())
		})
	}
}

func TestBlockDraw(t *testing.T) {
	newText := func(output string) tui.Text {
		text, _ := tui.NewText([]byte(output), tcell.Style{})
		return text
	}

	tcs := []struct {
		text    tui.Text
		x       int
		y       int
		width   int
		calls   int
		content map[int]map[int]rune
	}{
		{
			text:    tui.Text{},
			x:       0,
			y:       0,
			width:   1,
			calls:   0,
			content: map[int]map[int]rune{},
		},
		{
			text:    newText("x"),
			x:       0,
			y:       0,
			width:   1,
			calls:   1,
			content: map[int]map[int]rune{0: map[int]rune{0: 'x'}},
		},
		{
			text:    newText("x"),
			x:       1,
			y:       2,
			width:   1,
			calls:   1,
			content: map[int]map[int]rune{1: map[int]rune{2: 'x'}},
		},
		{
			text:  newText("asdf"),
			x:     5,
			y:     3,
			width: 2,
			calls: 4,
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

			block := tui.NewBlock(tc.text, tc.width)

			content := map[int]map[int]rune{}
			screen := &mock.ScreenMock{
				SetContentFunc: func(x int, y int, r rune, _ []rune, _ tcell.Style) {
					if _, ok := content[x]; !ok {
						content[x] = map[int]rune{}
					}
					content[x][y] = r
				},
			}

			block.Draw(screen, tc.x, tc.y)

			assert.Equal(tc.calls, len(screen.SetContentCalls()))
			assert.Equal(tc.content, content)
		})
	}
}
