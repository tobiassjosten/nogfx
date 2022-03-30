package tui

// import (
// 	"github.com/gdamore/tcell"
// )

// func (t *TUI) drawOutputbox(events [][]rune, x int, y int, X int, Y int) {
// 	width, height := X-x+1, Y-y+1
// 	lines := wrapLines(events, width)

// 	offset := len(lines) - height
// 	if offset < 0 {
// 		y -= offset
// 		offset = 0
// 	}

// 	for i := offset; i < len(lines); i++ {
// 		for ii := 0; ii < len(lines[i]); ii++ {
// 			t.screen.SetContent(x+ii, y+i-offset, lines[i][ii], nil, tcell.StyleDefault)
// 		}
// 	}
// }

func wrapLines(events [][]rune, width int) [][]rune {
	lines := [][]rune{}

	for i := 0; i < len(events); i++ {
		for ii := 0; ii < len(events[i]); ii++ {
			if ii%width == 0 {
				lines = append(lines, []rune{})
			}
			lines[len(lines)-1] = append(lines[len(lines)-1], events[i][ii])
		}
	}

	return lines
}
