package tui

import (
	"github.com/gdamore/tcell"
)

var inputboxSentStyle tcell.Style = tcell.StyleDefault.Background(tcell.ColorMaroon)

func (t *TUI) drawInputbox(input []rune, x int, y int, X int, Y int) {
	// ska också kunna hantera input som sträcker sig längre än X - x.
	style := tcell.StyleDefault
	if t.input.sent {
		style = inputboxSentStyle
	}

	for k, v := range input {
		t.screen.SetContent(x+k, y, v, nil, style)
	}
}
