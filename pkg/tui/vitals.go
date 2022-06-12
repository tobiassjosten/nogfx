package tui

import (
	"math"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/tobiassjosten/nogfx/pkg"
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
	}
)

// RenderVitals renders the current Vitals.
func (tui *TUI) RenderVitals(width int) Rows {
	if len(tui.character.Vitals) == 0 {
		return Rows{}
	}

	vorder := []string{"health", "mana", "endurance", "willpower", "energy"}

	// Remove non-existant vitals.
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

	row := Row{}
	for i, name := range vorder {
		styles, ok := vitalStyles[name]
		if !ok {
			styles = vitalStyles[""]
		}

		row = row.append(NewRow(min(1, i), NewCell(' '))...)
		row = row.append(RenderVital(
			tui.character.Vitals[name],
			(width-len(row))/(len(vorder)-i),
			styles,
		)...)
	}

	return Rows{row}
}

// RenderVital renders the given Vital.
func RenderVital(vital pkg.CharacterVital, width int, styles []tcell.Style) Row {
	fullWidth := int(math.Round(
		(float64(width) * float64(vital.Value) / float64(vital.Max)) - 0.01,
	))

	full := NewRow(fullWidth, NewCell(' ', styles[0]))
	empty := NewRow(width-len(full), NewCell(' ', styles[1]))

	row := append(full, empty...)

	value := strconv.Itoa(vital.Value)

	if len(value) <= len(row) {
		for i, x := 0, (width-len(value))/2; i < len(value); i++ {
			row[x+i].Content = rune(value[i])
		}
	}

	return row
}
