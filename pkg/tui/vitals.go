package tui

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell/v2"
)

// AddVital adds a new pane to show a vital metric.
func (tui *TUI) AddVital(name string, v interface{}) error {
	vital, ok := v.(*Vital)
	if !ok {
		return fmt.Errorf("unsupported vital '%s'", name)
	}

	if _, ok := tui.vitals[name]; !ok {
		tui.vitals[name] = vital
		tui.vorder = append(tui.vorder, name)
	}

	tui.Draw()

	return nil
}

// UpdateVital updates a given Vital with new current and max values.
func (tui *TUI) UpdateVital(name string, value, max int) error {
	vital, ok := tui.vitals[name]
	if !ok {
		return fmt.Errorf("couldn't update non-existent '%s' vital", name)
	}

	vital.value = value
	vital.max = max

	tui.Draw()

	return nil
}

// Vital represents a vital metric (health, mana, etc).
type Vital struct {
	value      int
	max        int
	fullStyle  tcell.Style
	emptyStyle tcell.Style
}

func setVitalNumbers(vital *Vital, numbers ...int) *Vital {
	if len(numbers) >= 1 {
		vital.value = numbers[0]
	}
	if len(numbers) >= 2 {
		vital.max = numbers[1]
	}
	return vital
}

// NewHealthVital creates a new Vital for a 'health' type.
func NewHealthVital(numbers ...int) *Vital {
	return setVitalNumbers(&Vital{
		fullStyle: tcell.StyleDefault.
			Background(tcell.ColorGreen).
			Foreground(tcell.ColorBlack),
		emptyStyle: tcell.StyleDefault.
			Background(tcell.ColorDarkGreen).
			Foreground(tcell.ColorBlack),
	}, numbers...)
}

// NewManaVital creates a new Vital for a 'mana' type.
func NewManaVital(numbers ...int) *Vital {
	return setVitalNumbers(&Vital{
		fullStyle: tcell.StyleDefault.
			Background(tcell.ColorBlue).
			Foreground(tcell.ColorBlack),
		emptyStyle: tcell.StyleDefault.
			Background(tcell.ColorDarkBlue).
			Foreground(tcell.ColorBlack),
	}, numbers...)
}

// NewEnduranceVital creates a new Vital for a 'endurance' type.
func NewEnduranceVital(numbers ...int) *Vital {
	return setVitalNumbers(&Vital{
		fullStyle: tcell.StyleDefault.
			Background(tcell.ColorTeal).
			Foreground(tcell.ColorBlack),
		emptyStyle: tcell.StyleDefault.
			Background(tcell.ColorDarkCyan).
			Foreground(tcell.ColorBlack),
	}, numbers...)
}

// NewWillpowerVital creates a new Vital for a 'willpower' type.
func NewWillpowerVital(numbers ...int) *Vital {
	return setVitalNumbers(&Vital{
		fullStyle: tcell.StyleDefault.
			Background(tcell.ColorFuchsia).
			Foreground(tcell.ColorBlack),
		emptyStyle: tcell.StyleDefault.
			Background(tcell.ColorRebeccaPurple).
			Foreground(tcell.ColorBlack),
	}, numbers...)
}

// NewEnergyVital creates a new Vital for an 'energy' type.
func NewEnergyVital(numbers ...int) *Vital {
	return setVitalNumbers(&Vital{
		fullStyle: tcell.StyleDefault.
			Background(tcell.ColorRed).
			Foreground(tcell.ColorBlack),
		emptyStyle: tcell.StyleDefault.
			Background(tcell.ColorRosyBrown).
			Foreground(tcell.ColorBlack),
	}, numbers...)
}

// RenderVitals renders the current Vitals.
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

// RenderVital renders the given Vital.
func RenderVital(vital *Vital, width int) Row {
	full := NewRow(
		proc(width, vital.value, vital.max),
		NewCell(' ', vital.fullStyle),
	)
	empty := NewRow(width-len(full), NewCell(' ', vital.emptyStyle))

	row := append(full, empty...)

	value := strconv.Itoa(vital.value)

	if len(value) <= len(row) {
		for i, x := 0, (width-len(value))/2; i < len(value); i++ {
			row[x+i].Content = rune(value[i])
		}
	}

	return row
}
