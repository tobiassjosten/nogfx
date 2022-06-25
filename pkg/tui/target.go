package tui

import (
	"fmt"
	"strconv"

	"github.com/tobiassjosten/nogfx/pkg"
)

// RenderTarget renders a vital for the current target.
func (tui *TUI) RenderTarget(width int) Rows {
	if rows, ok := tui.getCache(paneTarget); ok {
		return rows
	}

	if tui.target == nil || tui.target.Health < 0 {
		return Rows{}
	}

	row := RenderVital(
		pkg.CharacterVital{Value: tui.target.Health, Max: 100},
		width, vitalStyles["target"],
	)
	lrow := len(row)

	name := tui.target.Name
	lname := len(name)
	lhealth := len(strconv.Itoa(tui.target.Health))

	if (width-lhealth)/2 > lname+1 {
		if queued := tui.target.Queue() - 1; queued > 0 {
			queue := fmt.Sprintf(" (+%d)", queued)
			lqueue := len(queue)

			if (width-lhealth)/2 > lname+1+lqueue+1 {
				name += queue
				lname += len(queue)
			}
		}

		for i, r := range name {
			row[lrow-1-(lname-i)].Content = r
		}
	}

	rows := Rows{row}

	tui.setCache(paneTarget, rows)

	return rows
}
