package tui

import (
	"github.com/gdamore/tcell/v2"
)

func (pane *OutputPane) handlers() map[int]func(rune) bool {
	return map[int]func(rune) bool{
		int(tcell.KeyUp):                       pane.handleUpInput,
		int(tcell.KeyUp) + int(tcell.ModAlt):   pane.handleAltUpInput,
		int(tcell.KeyDown):                     pane.handleDownInput,
		int(tcell.KeyDown) + int(tcell.ModAlt): pane.handleAltDownInput,
	}
}

// HandleEvent reacts on a user event and modifies itself from it.
func (pane *OutputPane) HandleEvent(event *tcell.EventKey) bool {
	alts := []int{}

	if event.Modifiers()&tcell.ModAlt > 0 {
		alts = append(alts, int(event.Key())+int(event.Rune())+int(tcell.ModAlt))
		alts = append(alts, int(event.Key())+int(tcell.ModAlt))
	}
	alts = append(alts, int(event.Key())+int(event.Rune()))
	alts = append(alts, int(event.Key()))

	for _, alt := range alts {
		if f, ok := pane.handlers()[alt]; ok {
			if handled := f(event.Rune()); handled {
				return true
			}
		}
	}

	return false
}

func (pane *OutputPane) handleUpInput(_ rune) bool {
	pane.offset++
	return true
}

func (pane *OutputPane) handleAltUpInput(_ rune) bool {
	pane.offset += 5
	return true
}

func (pane *OutputPane) handleDownInput(_ rune) bool {
	pane.offset--
	if pane.offset < 0 {
		pane.offset = 0
	}
	return true
}

func (pane *OutputPane) handleAltDownInput(_ rune) bool {
	pane.offset -= 5
	if pane.offset < 0 {
		pane.offset = 0
	}
	return true
}
