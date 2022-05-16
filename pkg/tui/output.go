package tui

import (
	"github.com/gdamore/tcell/v2"
)

const (
	outputMinWidth = 80
	outputMaxWidth = 120
)

// Print shows a message to the user.
func (tui *TUI) Print(output []byte) {
	tui.output.Append(output)
	tui.Draw()
}

// Output is the widget where game output is shown.
type Output struct {
	buffer Rows
	offset int
	pwidth int
}

// Append adds a new paragraph to the Output.
func (output *Output) Append(data []byte) {
	row := NewRowFromBytes(data, output.lastStyle())
	output.buffer = output.buffer.prepend(row)

	if output.offset > 0 && output.pwidth > 0 {
		output.offset += len(row.Wrap(output.pwidth))
	}

	// @todo Completely arbitrary. Evaluate.
	if len(output.buffer) > 5000 {
		output.buffer = output.buffer[0:5000]
	}
}

func (output *Output) lastStyle() tcell.Style {
	for _, row := range output.buffer {
		if len(row) == 0 {
			continue
		}
		return row[len(row)-1].Style
	}

	return tcell.Style{}
}

// RenderOutput renders the current Output.
func (tui *TUI) RenderOutput(width, height int) Rows {
	return RenderOutput(tui.output, width, height)
}

// RenderOutput renders the given Output.
func RenderOutput(output *Output, width, height int) Rows {
	rows := Rows{}

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
		paragraph := row.Wrap(width)

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
		return rows[max(0, length-height):]
	}

	// Cap offset to the last row in the buffer.
	output.offset = min(length-height, output.offset)

	hheight := length - height - output.offset
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

	rows = rows[length-(height-height/2)+1:]

	return append(history, append(Rows{divider}, rows...)...)
}
