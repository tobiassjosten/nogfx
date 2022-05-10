package tui

import (
	"github.com/gdamore/tcell/v2"
)

// Outputs exposes the incoming channel for server output.
func (tui *TUI) Outputs() chan<- []byte {
	return tui.panes.Output.Outputs()
}

// Print shows a message to the user.
func (tui *TUI) Print(output []byte) {
	tui.panes.Output.Add(output)
	tui.Draw()
}

// OutputPane is the pane where primary game output is shown.
type OutputPane struct {
	outputs chan []byte

	rows   Rows
	offset int

	pwidth int
}

// NewOutputPane creates a new OutputPane.
func NewOutputPane() *OutputPane {
	return &OutputPane{outputs: make(chan []byte)}
}

func (pane *OutputPane) style() tcell.Style {
	if len(pane.rows) == 0 {
		return tcell.Style{}
	}
	return pane.rows[0][0].Style
}

// Outputs exposes the incoming channel for server output.
func (pane *OutputPane) Outputs() chan []byte {
	return pane.outputs
}

// Add appends new paragraphs of text to be show to the user.
func (pane *OutputPane) Add(output []byte) {
	row := (Row{}).Parse(output, pane.style())
	pane.rows = pane.rows.Prepend(row)

	if pane.offset > 0 && pane.pwidth > 0 {
		pane.offset += len(row.Wrap(pane.pwidth))
	}

	// @todo Completely arbitrary. Evaluate.
	if len(pane.rows) > 5000 {
		pane.rows = pane.rows[0:5000]
	}
}

// Rows distributes Cells from the buffer to be printed to the screen, in the
// form of an output area and an optional history scrollback area.
func (pane *OutputPane) Rows(width, height int) (Rows, Rows) {
	rows := Rows{}

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

	for _, row := range pane.rows {
		paragraph := row.Wrap(width)

		// Rows are ordered with the most recent one first, so we
		// prepend older paragraphs to the area.
		for i := len(paragraph) - 1; i >= 0; i-- {
			rows = rows.Prepend(paragraph[i])
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
		return rows[max(0, length-height):], Rows{}
	}

	// Cap offset to the last row in the buffer.
	pane.offset = min(length-height, pane.offset)

	history := length - height - pane.offset
	historyPane := rows[history : history+(height-height/2)]

	output := rows[length-height/2:]

	divider := NewRow(width, NewCell(tcell.RuneHLine))

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
