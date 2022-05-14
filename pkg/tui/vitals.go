package tui

import (
	"log"
	"strconv"

	"github.com/gdamore/tcell/v2"
)

// AddVital adds a new pane to show a vital metric.
func (tui *TUI) AddVital(name string, v interface{}) {
	vital, ok := v.(*Vital)
	if !ok {
		log.Printf("unsupported '%s' vital", name)
		return
	}

	if _, ok := tui.vitals[name]; !ok {
		tui.vitals[name] = vital
		tui.vorder = append(tui.vorder, name)
	}

	tui.Draw()
}

// UpdateVital updates a given Vital with new current and max values.
func (tui *TUI) UpdateVital(name string, value, max int) {
	vital, ok := tui.vitals[name]
	if !ok {
		log.Printf("couldn't update non-existent 'health' vital")
		return
	}

	vital.value = value
	vital.max = max

	tui.Draw()
}

// Default vitals suitable for most games.
var (
	HealthVital = &Vital{
		fullStyle:  tcell.StyleDefault.Background(tcell.ColorGreen).Foreground(tcell.ColorBlack),
		emptyStyle: tcell.StyleDefault.Background(tcell.ColorDarkGreen).Foreground(tcell.ColorBlack),
	}
	ManaVital = &Vital{
		fullStyle:  tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorBlack),
		emptyStyle: tcell.StyleDefault.Background(tcell.ColorDarkBlue).Foreground(tcell.ColorBlack),
	}
	EnduranceVital = &Vital{
		fullStyle:  tcell.StyleDefault.Background(tcell.ColorTeal).Foreground(tcell.ColorBlack),
		emptyStyle: tcell.StyleDefault.Background(tcell.ColorDarkCyan).Foreground(tcell.ColorBlack),
	}
	WillpowerVital = &Vital{
		fullStyle:  tcell.StyleDefault.Background(tcell.ColorFuchsia).Foreground(tcell.ColorBlack),
		emptyStyle: tcell.StyleDefault.Background(tcell.ColorRebeccaPurple).Foreground(tcell.ColorBlack),
	}
	EnergyVital = &Vital{
		fullStyle:  tcell.StyleDefault.Background(tcell.ColorRed).Foreground(tcell.ColorBlack),
		emptyStyle: tcell.StyleDefault.Background(tcell.ColorRosyBrown).Foreground(tcell.ColorBlack),
	}
)

// Vital represents a vital metric (health, mana, etc).
type Vital struct {
	value      int
	max        int
	fullStyle  tcell.Style
	emptyStyle tcell.Style
}

func (tui *TUI) RenderVitals(width int) Rows {
	if len(tui.vorder) == 0 {
		return Rows{}
	}

	row := Row{}
	for i, name := range tui.vorder {
		row = row.append(NewRow(min(1, i), NewCell(' '))...)
		row = row.append(RenderVital(
			tui.vitals[name],
			(width-len(row))/(len(tui.vorder)-i),
		)...)
	}

	return Rows{row}
}

func RenderVital(vital *Vital, width int) Row {
	full := NewRow(
		proc(width, vital.value, vital.max),
		NewCell(' ', vital.fullStyle),
	)
	empty := NewRow(width-len(full), NewCell(' ', vital.emptyStyle))

	row := append(full, empty...)

	value := strconv.Itoa(vital.value)

	if len(value) <= len(row) {
		for i, x := 0, width/2-len(value)/2; i < len(value); i++ {
			row[x+i].Content = rune(value[i])
		}
	}

	return row
}
