package tui

import (
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

type Block struct {
	width int
	Rows  []Text
}

func (block *Block) addCell(cell Cell) {
	block.Rows[len(block.Rows)-1] = append(
		block.Rows[len(block.Rows)-1],
		cell,
	)
}

func (block *Block) addRow() {
	block.Rows = append(block.Rows, Text{})
}

func (block *Block) rowWidth() int {
	rwidth := 0
	for _, r := range block.Rows[len(block.Rows)-1] {
		rwidth += runewidth.RuneWidth(r.Content)
	}
	return rwidth
}

func (block *Block) Width() int {
	return block.width
}

func (block *Block) Height() int {
	return len(block.Rows)
}

func (block *Block) Size() (int, int) {
	return block.Width(), block.Height()
}

func NewBlock(text Text, width int) Block {
	block := Block{width: width, Rows: []Text{{}}}

	word := Text{}

	for _, cell := range text {
		word = append(word, cell)
		wwidth := word.Width()

		rwidth := block.rowWidth()

		if unicode.IsSpace(cell.Content) {
			word = word[:len(word)-1]
			for i := 0; i < len(word); i++ {
				block.addCell(word[i])
			}

			word = Text{}

			if cell.Width > 0 {
				block.addCell(cell)
			}
		}

		if (rwidth+wwidth > width && rwidth > wwidth) || cell.Content == '\n' {
			block.addRow()
		}

		if rwidth+wwidth > width && rwidth <= wwidth {
			word = word[:len(word)-1]
			for i := 0; i < len(word); i++ {
				block.addCell(word[i])
			}

			word = Text{}

			block.addRow()
		}
	}

	for _, c := range word {
		block.addCell(c)
	}

	return block
}

func (block *Block) draw(screen tcell.Screen, x, y int) {
	for iy, row := range block.Rows {
		for ix, cell := range row {
			screen.SetContent(x+ix, y+iy, cell.Content, nil, cell.Style)
		}
	}
}