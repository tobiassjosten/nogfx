package tui

import (
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
)

func (tui *TUI) VitalsAdd(name string, vital Vital) {
	if _, ok := tui.panes.vitals.vitals[name]; ok {
		return
	}

	tui.panes.vitals.vitals[name] = vital
	tui.panes.vitals.vorder = append(tui.panes.vitals.vorder, name)
}

func (tui *TUI) VitalsUpdate(name string, value, max int) {
	vital, ok := tui.panes.vitals.vitals[name]
	if !ok {
		return
	}

	vital.Value = value
	vital.Max = max
	tui.panes.vitals.vitals[name] = vital

	tui.Draw()
}

type Vital struct {
	Value      int
	Max        int
	FullStyle  tcell.Style
	EmptyStyle tcell.Style
}

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
