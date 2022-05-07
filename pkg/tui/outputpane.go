package tui

import (
	"github.com/gdamore/tcell/v2"
)

// Outputs exposes the incoming channel for server output.
func (tui *TUI) Outputs() chan<- []byte {
	return tui.panes.Output.Outputs()
}

// Print shows a text to the user.
func (tui *TUI) Print(output []byte) {
	tui.panes.Output.Add(output)
	tui.Draw()
}

// OutputPane is the pane where primary game output is shown.
type OutputPane struct {
	outputs chan []byte

	texts  []Text
	offset int

	pwidth int
}

// NewOutputPane creates a new OutputPane.
func NewOutputPane() *OutputPane {
	return &OutputPane{outputs: make(chan []byte)}
}

func (pane *OutputPane) lastStyle() tcell.Style {
	// @todo Figure out why we have empty texts and revert this to a normal
	// slice access for the last cell in the first text.
	for _, text := range pane.texts {
		if len(text) == 0 {
			continue
		}
		return text[len(text)-1].Style
	}

	return tcell.Style{}
}

// Outputs exposes the incoming channel for server output.
func (pane *OutputPane) Outputs() chan []byte {
	return pane.outputs
}

// Add appends new paragraphs of text to be show to the user.
func (pane *OutputPane) Add(output []byte) {
	text := NewText(output, pane.lastStyle())
	pane.texts = append([]Text{text}, pane.texts...)

	if pane.offset > 0 && pane.pwidth > 0 {
		pane.offset += len(text.Wrap(pane.pwidth))
	}

	// @todo Completely arbitrary. Evaluate.
	if len(pane.texts) > 5000 {
		pane.texts = pane.texts[0:5000]
	}
}

// Texts distributes Cells from the Text buffer to be printed to the screen, in
// the form of an output area and an optional history scrollback area.
func (pane *OutputPane) Texts(width, height int) ([]Text, []Text) {
	rows := []Text{}

	if width == 0 || height == 0 {
		return rows, rows
	}

	// Resizing the window resets history scrollback, simply because it's a
	// pain in the ass to calculate and maintain that state. For now.
	// @todo Make resizing maintain history scrollback.
	if pane.pwidth > 0 && pane.pwidth != width {
		pane.offset = 0
	}
	pane.pwidth = width

	// Make sure to calculate enough for a history subpane.
	height += pane.offset

	for _, text := range pane.texts {
		lines := text.Wrap(width)

		// Texts are ordered with the most recent one first, so we
		// prepend older paragraphs to the area.
		for i := len(lines) - 1; i >= 0; i-- {
			rows = append([]Text{lines[i]}, rows...)
		}

		if len(rows) >= height {
			break
		}
	}

	// Reset back to actual height, for finalization.
	height -= pane.offset
	length := len(rows)

	// For simpler cases we just return the full buffer.
	if height == 1 || length <= height || pane.offset == 0 {
		return rows[max(0, length-height):], []Text{}
	}

	// Cap offset to the last row in the text buffer.
	pane.offset = min(length-height, pane.offset)

	history := length - height - pane.offset
	historyPane := rows[history : history+(height-height/2)]

	output := rows[length-height/2:]

	// @todo Replace with NewRow() after merging that PR.
	var divider Text
	for i := 0; i < width; i++ {
		divider = append(divider, Cell{Content: tcell.RuneHLine})
	}

	if height > 2 {
		hlength := len(historyPane)
		if hlength > len(output) {
			historyPane[hlength-1] = divider
		} else {
			output[0] = divider
		}
	}

	return output, historyPane
}
