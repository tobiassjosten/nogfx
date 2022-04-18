package tui

import "github.com/gdamore/tcell/v2"

// OutputPane is the pane where primary game output is shown.
type OutputPane struct {
	outputs chan []byte

	x      int
	y      int
	width  int
	height int

	outputStyle tcell.Style

	texts []Text
}

// NewOutputPane creates a new OutputPane.
func NewOutputPane() *OutputPane {
	var (
		outputStyle = tcell.Style{}
	)

	return &OutputPane{
		outputs:     make(chan []byte),
		outputStyle: outputStyle,
	}
}

// Position sets the x.y coordinates for and resizes the pane.
func (pane *OutputPane) Position(x, y, width, height int) {
	pane.x, pane.y = x, y
	pane.width, pane.height = width, height
}

// Add appends new paragraphs of text to be show to the user.
func (pane *OutputPane) Add(output []byte) {
	text, style := NewText(output, pane.outputStyle)
	pane.outputStyle = style
	pane.texts = append([]Text{text}, pane.texts...)

	// @todo Cap tui.texts so it doesn't grow indefinitely.
}

// Draw prints the contents of the OutputPane to the given tcell.Screen.
func (pane *OutputPane) Draw(screen tcell.Screen) {
	x, y := pane.x, pane.y+pane.height-1

	for _, t := range pane.texts {
		b := NewBlock(t, pane.width)
		y = y - b.Height() + 1
		b.Draw(screen, x, y)

		y--
		if y < pane.y {
			break
		}
	}

	// @todo Fixa stöd för att kunna scrolla upp.
}
