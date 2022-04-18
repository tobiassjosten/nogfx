package tui

import (
	"unicode"

	"github.com/gdamore/tcell/v2"
)

// Block is a Text word-wrapped at a specific width.
type Block struct {
	width int
	rows  []Text
}

func (block *Block) addCell(cell Cell) {
	block.rows[len(block.rows)-1] = append(
		block.rows[len(block.rows)-1],
		cell,
	)
}

func (block *Block) addRow() {
	block.rows = append(block.rows, Text{})
}

func (block *Block) rowWidth() int {
	rwidth := 0
	for _, r := range block.rows[len(block.rows)-1] {
		rwidth += r.Width
	}
	return rwidth
}

// Width returns the max width of the Block.
func (block *Block) Width() int {
	return block.width
}

// Height returns the actual height of the Block.
func (block *Block) Height() int {
	return len(block.rows)
}

// NewBlock parses a Text and performs word wrapping at the given width.
func NewBlock(text Text, width int) Block {
	block := Block{width: width, rows: []Text{}}

	if len(text) == 0 || width == 0 {
		return block
	}
	block.addRow()

	word := Text{}
	wwidth := 0

	for _, cell := range text {
		rwidth := block.rowWidth()

		if rwidth >= width {
			block.addRow()
			rwidth = 0

			if cell.Content == '\n' {
				continue
			}
		}

		word = append(word, cell)
		wwidth += cell.Width

		if unicode.IsSpace(cell.Content) || wwidth >= width {
			for _, c := range word {
				block.addCell(c)
			}

			word = Text{}
			wwidth = 0
		}

		if rwidth+wwidth > width || cell.Content == '\n' {
			block.addRow()
		}
	}

	for _, c := range word {
		block.addCell(c)
	}

	return block
}

// Draw prints the contents of the Block to the given tcell.Screen.
func (block *Block) Draw(screen tcell.Screen, x, y int) {
	for yy, row := range block.rows {
		for xx, cell := range row {
			content := cell.Content
			screen.SetContent(x+xx, y+yy, content, nil, cell.Style)
		}

	}
}
