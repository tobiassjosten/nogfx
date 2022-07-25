package tui

import (
	"github.com/gdamore/tcell/v2"
)

// Print shows a message to the user.
func (tui *TUI) Print(output []byte) {
	// @todo Make it set its own color (or ^[37m) before resetting back to
	// the previous last seen style.
	tui.output.Append(output)
	tui.setCache(paneOutput, nil)
	tui.Draw()
}

// Output is the widget where game output is shown.
type Output struct {
	buffer Rows
	offset int
	pwidth int
	style  tcell.Style
}

// Append adds a new paragraph to the Output.
func (output *Output) Append(data []byte) {
	row, style := NewRowFromBytes(data, output.style)
	output.style = style

	output.buffer = output.buffer.prepend(row)

	if output.offset > 0 && output.pwidth > 0 {
		output.offset += len(row.Wrap(output.pwidth))
	}

	// @todo Completely arbitrary. Evaluate.
	if len(output.buffer) > 5000 {
		output.buffer = output.buffer[0:5000]
	}
}

// RenderOutput renders the current Output.
func (tui *TUI) RenderOutput(width, height int) Rows {
	if rows, ok := tui.getCache(paneOutput); ok {
		return rows
	}

	rows := RenderOutput(tui.output, width, height)

	tui.setCache(paneOutput, rows)

	return rows
}

// RenderOutput renders the given Output.
func RenderOutput(output *Output, width, height int) Rows {
	rows := Rows{}

	padding := NewCell(' ')

	if width == 0 || height == 0 {
		return rows
	}

	// @todo Make resizing maintain history scrollback. Resetting it is a
	// temporary workaround because calculating and maintaining scrollback
	// state through resizing is a pain in the ass.
	if output.pwidth > 0 && output.pwidth != width {
		output.offset = 0
	}
	output.pwidth = width

	// Make sure to render enough for a history scrollback split.
	height += output.offset

	for _, row := range output.buffer {
		paragraph := row.Wrap(width, padding)

		// Rows are ordered with the most recent one first, so we
		// prepend older paragraphs to the rows.
		for i := len(paragraph) - 1; i >= 0; i-- {
			rows = rows.prepend(paragraph[i])
		}

		if len(rows) >= height {
			break
		}
	}

	// Reset back to actual height, for finalization.
	height -= output.offset
	length := len(rows)

	// For simpler cases we just return the full buffer.
	if height <= 2 || length <= height || output.offset == 0 {
		rows = rows[max(0, length-height):]
		rows = append(NewRows(width, height-length, padding), rows...)
		return rows
	}

	// Cap offset to the last row in the buffer.
	output.offset = min(length-height, output.offset)

	hheight := length - height - output.offset
	history := rows[hheight : hheight+height/2]

	// @todo Mark this divider better, with colors and flair.
	divider := NewRow(width, NewCell(tcell.RuneHLine))

	rows = rows[length-(height-height/2)+1:]

	return append(history, append(Rows{divider}, rows...)...)
}
