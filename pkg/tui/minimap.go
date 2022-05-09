package tui

import (
	"log"

	"github.com/tobiassjosten/nogfx/pkg/navigation"

	"github.com/gdamore/tcell/v2"
)

// SetRoom updates the current room and causes a repaint.
func (tui *TUI) SetRoom(room *navigation.Room) {
	tui.panes.Minimap.room = room
	tui.Draw()
}

// MinimapPane is a map rendition based on the current room.
type MinimapPane struct {
	room *navigation.Room
}

// NewMinimapPane creates a new MinimapPane.
func NewMinimapPane() *MinimapPane {
	return &MinimapPane{}
}

// Texts renders cascading layers of adjacent rooms, based on the current one.
func (pane *MinimapPane) Texts(width, height int) []Text {
	if pane.room == nil || width == 0 || height == 0 {
		return []Text{}
	}

	rows := NewRows(width, height)

	return renderRoom(
		rows, pane.room,
		len(rows[0])/2, len(rows)/2,
		5, map[int]struct{}{},
	)
}

func renderRoom(rows []Text, room *navigation.Room, x, y, depth int, rendered map[int]struct{}) []Text {
	// Make sure we have enough padding to render room borders and exits.
	if x < 2 || y < 2 || y > len(rows)-3 || x > len(rows[0])-3 {
		return rows
	}

	rows[y][x].Content = ' '
	if room.HasPlayer {
		rows[y][x].Content = '+'
	}

	rendered[room.ID] = struct{}{}

	rows[y][x-1].Content = '['
	rows[y][x+1].Content = ']'

	// @todo Consider which environments that are actually important for
	// the map to visualise and color the brackets accordingly.

	if !room.Known {
		rows[y][x-1].Style = rows[y][x-1].Style.Foreground(tcell.Color237)
		rows[y][x+1].Style = rows[y][x+1].Style.Foreground(tcell.Color237)
	}

paths:
	for direction, adjacent := range room.Exits {
		switch direction {
		case "u":
			if room.HasPlayer {
				continue
			}

			if room.HasExit("d") {
				rows[y][x].Content = '='
				rows[y][x].Foreground(tcell.Color245)
				continue
			}

			rows[y][x].Content = '^'
			rows[y][x].Foreground(tcell.Color245)
			continue

		case "d":
			if room.HasPlayer {
				continue
			}

			if room.HasExit("u") {
				continue
			}

			rows[y][x].Content = 'v'
			rows[y][x].Foreground(tcell.Color245)
			continue

		case "out":
			rows[y][x-1].Content = '{'
			continue

		case "in":
			rows[y][x+1].Content = '}'
			continue
		}

		diffx, diffy := room.Displacement(direction)
		if diffx == 0 && diffy == 0 {
			continue
		}

		if _, done := rendered[adjacent.ID]; done {
			continue
		}

		dirchar := ' '
		switch direction {
		case "n":
			dirchar = '|'

		case "ne":
			dirchar = '/'

		case "e":
			dirchar = '-'

		case "se":
			dirchar = '\\'

		case "s":
			dirchar = '|'

		case "sw":
			dirchar = '/'

		case "w":
			dirchar = '-'

		case "nw":
			dirchar = '\\'

		default:
			continue paths
		}

		// Calculate the number of steps, taking into account the fact
		// that every offset beyond the first equals one path and one
		// room. So 1 offset = 1 step, 2 offsets = 3 steps, etc.
		steps := max(1, max(abs(diffx), abs(diffy)))
		steps = (steps-1)*2 + 1

		// One point offset, when rendered, equals four or two cells.
		if diffx > 0 || diffx < 0 {
			diffx = diffx*4 - rel(3, diffx)
		}
		if diffy > 0 || diffy < 0 {
			diffy = diffy*2 - rel(1, diffy)
		}

		for i := 1; i <= steps; i++ {
			// X is offset by one because rooms are three cells
			// wide and paths need to start "one out".
			xx := x + rel(1, diffx) + proc(diffx, i, steps)
			yy := y + proc(diffy, i, steps)

			if yy < 0 || xx < 0 || yy >= len(rows) || xx >= len(rows[yy]) {
				log.Println("overflow", len(rows), len(rows[0]), "coords", yy, xx)
				break
			}

			rows[yy][xx].Content = dirchar
		}
	}

	if depth == 0 {
		return rows
	}

	for direction, adjacent := range room.Exits {
		if _, done := rendered[adjacent.ID]; done {
			continue
		}

		if room.Area != nil && adjacent.Area != nil && room.Area.ID != adjacent.Area.ID {
			continue
		}

		diffx, diffy := room.Displacement(direction)
		if diffx == 0 && diffy == 0 {
			continue
		}

		rows = renderRoom(rows, adjacent, x+4*diffx, y+2*diffy, depth-1, rendered)
	}

	return rows
}
