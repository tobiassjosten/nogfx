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
	// @todo Figure out why we have empty rows and revert this to a normal
	// slice access for the last cell in the most recent row.
	for _, row := range pane.rows {
		if len(row) == 0 {
			continue
		}
		return row[len(row)-1].Style
	}

	return tcell.Style{}
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

// Render processes the output buffer and distributes its rows within the given
// width and height confines, potentially with a history scrollback split.
func (pane *OutputPane) Render(width, height int) Rows {
	rows := Rows{}

	if width == 0 || height == 0 {
		return rows
	}

	// @todo Make resizing maintain history scrollback. Resetting it is a
	// temporary workaround because calculating and maintaining scrollback
	// state through resizing is a pain in the ass.
	if pane.pwidth > 0 && pane.pwidth != width {
		pane.offset = 0
	}
	pane.pwidth = width

	// Make sure to render enough for a history scrollback split.
	height += pane.offset

	for _, row := range pane.rows {
		paragraph := row.Wrap(width)

		// Rows are ordered with the most recent one first, so we
		// prepend older paragraphs to the rows.
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
	if height <= 2 || length <= height || pane.offset == 0 {
		return rows[max(0, length-height):]
	}

	// Cap offset to the last row in the buffer.
	pane.offset = min(length-height, pane.offset)

	hheight := length - height - pane.offset
	history := rows[hheight : hheight+height/2]

	// Color history scrollback split background.
	hstyle := (tcell.Style{}).Background(tcell.Color236)
	for y, row := range history {
		for x := range row {
			history[y][x].Background(tcell.Color236)
		}
		for i := len(row); i < width; i++ {
			history[y] = append(history[y], NewCell(' ', hstyle))
		}
	}

	divider := NewRow(width, NewCell(tcell.RuneHLine))

	output := rows[length-(height-height/2)+1:]

	return append(history, append(Rows{divider}, output...)...)
}
