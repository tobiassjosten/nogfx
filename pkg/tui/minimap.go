package tui

import (
	"github.com/gdamore/tcell"
)

func (t *TUI) drawMinimap(x int, y int, X int, Y int) {
	character := 'X'
	if t.world == nil {
		character = '?'
	}

	for yy := y; yy <= Y; yy += 1 {
		for xx := x; xx <= X; xx += 1 {
			t.screen.SetContent(xx, yy, character, nil, tcell.StyleDefault)
		}
	}
}
