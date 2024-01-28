package tui

import (
	"math"
	"strconv"

	"github.com/tobiassjosten/nogfx/pkg"

	"github.com/gdamore/tcell/v2"
)

var (
	// Predefined styles of some common vitals. i=0 is full, i=1 is empty.
	vitalStyles = map[string][]tcell.Style{
		// Fallback style, for when none other matches.
		"": {
			tcell.StyleDefault.
				Background(tcell.Color250).
				Foreground(tcell.ColorBlack),
			tcell.StyleDefault.
				Background(tcell.Color240).
				Foreground(tcell.ColorBlack),
		},
		"health": {
			tcell.StyleDefault.
				Background(tcell.ColorGreen).
				Foreground(tcell.ColorBlack),
			tcell.StyleDefault.
				Background(tcell.ColorDarkGreen).
				Foreground(tcell.ColorBlack),
		},
		"mana": {
			tcell.StyleDefault.
				Background(tcell.ColorBlue).
				Foreground(tcell.ColorBlack),
			tcell.StyleDefault.
				Background(tcell.ColorDarkBlue).
				Foreground(tcell.ColorBlack),
		},
		"endurance": {
			tcell.StyleDefault.
				Background(tcell.ColorTeal).
				Foreground(tcell.ColorBlack),
			tcell.StyleDefault.
				Background(tcell.ColorDarkCyan).
				Foreground(tcell.ColorBlack),
		},
		"willpower": {
			tcell.StyleDefault.
				Background(tcell.ColorFuchsia).
				Foreground(tcell.ColorBlack),
			tcell.StyleDefault.
				Background(tcell.ColorRebeccaPurple).
				Foreground(tcell.ColorBlack),
		},
		"energy": {
			tcell.StyleDefault.
				Background(tcell.ColorYellow).
				Foreground(tcell.ColorBlack),
			tcell.StyleDefault.
				Background(tcell.Color100).
				Foreground(tcell.ColorBlack),
		},
		"target": {
			tcell.StyleDefault.
				Background(tcell.ColorRed).
				Foreground(tcell.ColorBlack),
			tcell.StyleDefault.
				Background(tcell.ColorDarkRed).
				Foreground(tcell.ColorBlack),
		},
	}
)

// RenderVitals renders the current Vitals.
func (tui *TUI) RenderVitals(width int) Rows {
	if rows, ok := tui.getCache(paneVitals); ok {
		return rows
	}

	if len(tui.character.Vitals) == 0 {
		return Rows{}
	}

	vorder := []string{"health", "mana", "endurance", "willpower", "energy"}

	// Remove non-existent vitals.
	for i := 0; i < len(vorder); {
		if _, ok := tui.character.Vitals[vorder[i]]; !ok {
			vorder = append(vorder[:i], vorder[i+1:]...)
			continue
		}
		i++
	}

	// Add missing vitals.
	for name := range tui.character.Vitals {
		exists := false

		for _, nname := range vorder {
			if name == nname {
				exists = true
				break
			}
		}

		if !exists {
			vorder = append(vorder, name)
		}
	}

	gapStyle := (tcell.Style{}).Background(tcell.Color235)

	row := Row{}

	for i, name := range vorder {
		styles, ok := vitalStyles[name]
		if !ok {
			styles = vitalStyles[""]
		}

		row = row.append(NewRow(min(1, i), NewCell(' ', gapStyle))...)
		row = row.append(RenderVital(
			tui.character.Vitals[name],
			(width-len(row))/(len(vorder)-i),
			styles,
		)...)
	}

	rows := Rows{row}

	tui.setCache(paneVitals, rows)

	return rows
}

// RenderVital renders the given Vital.
func RenderVital(vital pkg.CharacterVital, width int, styles []tcell.Style) Row {
	fullWidth := int(math.Round(
		(float64(width) * float64(vital.Value) / float64(vital.Max)) - 0.01,
	))

	row := NewRow(fullWidth, NewCell(' ', styles[0]))
	row = append(row, NewRow(width-len(row), NewCell(' ', styles[1]))...)

	value := strconv.Itoa(vital.Value)

	if len(value) <= len(row) {
		for i, x := 0, (width-len(value))/2; i < len(value); i++ {
			row[x+i].Content = rune(value[i])
		}
	}

	return row
}
