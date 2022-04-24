package tui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
)

// AddVital adds a new Vital to be displayed in the VitalsPane.
func (tui *TUI) AddVital(name string, v interface{}) error {
	vital, ok := v.(Vital)
	if !ok {
		return fmt.Errorf("only tui.Vital vitals are supported")
	}

	return tui.panes.vitals.AddVital(name, vital)
}

// UpdateVital updates a given Vital with new current and max values.
func (tui *TUI) UpdateVital(name string, value, max int) error {
	defer tui.Draw()
	return tui.panes.vitals.UpdateVital(name, value, max)
}

// Vital is a metric (health, mana, etc) visualised in a VitalsPane.
type Vital struct {
	Value      int
	Max        int
	FullStyle  tcell.Style
	EmptyStyle tcell.Style
}

// Default vitals suitable for most games.
var (
	HealthVital = Vital{
		FullStyle:  tcell.StyleDefault.Background(tcell.ColorGreen).Foreground(tcell.ColorBlack),
		EmptyStyle: tcell.StyleDefault.Background(tcell.ColorDarkGreen).Foreground(tcell.ColorBlack),
	}
	ManaVital = Vital{
		FullStyle:  tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorBlack),
		EmptyStyle: tcell.StyleDefault.Background(tcell.ColorDarkBlue).Foreground(tcell.ColorBlack),
	}
	EnduranceVital = Vital{
		FullStyle:  tcell.StyleDefault.Background(tcell.ColorTeal).Foreground(tcell.ColorBlack),
		EmptyStyle: tcell.StyleDefault.Background(tcell.ColorDarkCyan).Foreground(tcell.ColorBlack),
	}
	WillpowerVital = Vital{
		FullStyle:  tcell.StyleDefault.Background(tcell.ColorFuchsia).Foreground(tcell.ColorBlack),
		EmptyStyle: tcell.StyleDefault.Background(tcell.ColorRebeccaPurple).Foreground(tcell.ColorBlack),
	}
	EnergyVital = Vital{
		FullStyle:  tcell.StyleDefault.Background(tcell.ColorRed).Foreground(tcell.ColorBlack),
		EmptyStyle: tcell.StyleDefault.Background(tcell.ColorRosyBrown).Foreground(tcell.ColorBlack),
	}
)

// VitalsPane shows vital metrics like health and energies.
type VitalsPane struct {
	x      int
	y      int
	width  int
	height int

	vitals map[string]Vital
	vorder []string
}

// NewVitalsPane creates a new VitalsPane.
func NewVitalsPane() *VitalsPane {
	return &VitalsPane{
		vitals: map[string]Vital{},
	}
}

// AddVital adds a new Vital to be displayed in the VitalsPane.
func (pane *VitalsPane) AddVital(name string, vital Vital) error {
	if _, ok := pane.vitals[name]; ok {
		return fmt.Errorf("vital already added '%s'", name)
	}

	pane.vitals[name] = vital
	pane.vorder = append(pane.vorder, name)

	return nil
}

// UpdateVital updates a given Vital with new current and max values.
func (pane *VitalsPane) UpdateVital(name string, value, max int) error {
	vital, ok := pane.vitals[name]
	if !ok {
		return fmt.Errorf("non-existent vital '%s'", name)
	}

	vital.Value = value
	vital.Max = max
	pane.vitals[name] = vital

	return nil
}

// Position sets the x.y coordinates for and resizes the pane.
func (pane *VitalsPane) Position(x, y, width, height int) {
	pane.x, pane.y = x, y
	pane.width, pane.height = width, height
}

// Height is the actual height that a full rendition of VitalsPane would need,
// as opposed to its `height` property, which is what it's afforded.
func (pane *VitalsPane) Height() int {
	for _, vital := range pane.vitals {
		if vital.Value > 0 && vital.Max > 0 {
			return 1
		}
	}

	return 0
}

// Draw prints the contents of the VitalsPane to the given tcell.Screen.
func (pane *VitalsPane) Draw(screen tcell.Screen) {
	if pane.height == 0 {
		return
	}

	lvitals := len(pane.vitals)
	x, y, pwidth := pane.x, pane.y, pane.width

	vwidth := (pwidth - lvitals + 1) / lvitals
	remains := (pwidth - lvitals + 1) % lvitals

	for ii, name := range pane.vorder {
		if ii > 0 {
			screen.SetContent(x, y, ' ', nil, tcell.StyleDefault)
			x++
		}

		vital := pane.vitals[name]

		vwidth := vwidth
		if remains > 0 {
			vwidth++
			remains--
		}

		fullBreak := int(float64(vital.Value) / float64(vital.Max) * float64(vwidth))

		text := []rune(strconv.Itoa(vital.Value))
		text = append(text, []rune(strings.Repeat(
			" ", max(0, vwidth-len(text)),
		))...)

		for i := 0; i < vwidth; i++ {
			style := vital.FullStyle
			if vital.Max == 0 || fullBreak < i {
				style = vital.EmptyStyle
			}

			screen.SetContent(x+i, y, text[i], nil, style)
		}

		x += vwidth
	}
}
