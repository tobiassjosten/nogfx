package tui

import (
	"log"

	"github.com/tobiassjosten/nogfx/pkg/navigation"

	"github.com/gdamore/tcell/v2"
)

// Outputs exposes the incoming channel for server output.
func (tui *TUI) SetRoom(room *navigation.Room) {
	tui.panes.Minimap.room = room
	tui.Draw()
}

type MinimapPane struct {
	room *navigation.Room
}

func NewMinimapPane() *MinimapPane {
	return &MinimapPane{}
}

func (pane *MinimapPane) Texts(width, height int) []Text {
	if pane.room == nil || width == 0 || height == 0 {
		return []Text{}
	}

	rows := NewRows(width, height)

	x, y := len(rows[0])/2, len(rows)/2

	depth := 3 // @todo bestäm baserat på width/height och hur mycket som får plast

	rendered := map[int]struct{}{}
	rows = renderRoom(rows, pane.room, x, y, depth, rendered)

	return rows
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

	// @todo Use characters and colors to represent the room's type.

	rows[y][x-1].Content = '['
	rows[y][x+1].Content = ']'

	if !room.Known {
		rows[y][x-1].Style = rows[y][x-1].Style.Foreground(tcell.Color237)
		rows[y][x+1].Style = rows[y][x+1].Style.Foreground(tcell.Color237)
	}

	for direction, adjacent := range room.Exits {
		// These special directions paints the room, not paths between
		// rooms, and so we exclude them from the outwards rule below.
		switch direction {
		case "u":
			if room.HasPlayer {
				continue
			}

			if room.HasExit("d") {
				rows[y][x].Content = '=' // dimma färgen
				continue
			}

			rows[y][x].Content = '^' // dimma färgen
			continue

		case "d":
			if room.HasPlayer {
				continue
			}

			if room.HasExit("u") {
				continue
			}

			rows[y][x].Content = 'v' // dimma färgen
			continue

		case "out":
			rows[y][x-1].Content = '{' // dimma färgen
			continue

		case "in":
			rows[y][x+1].Content = '}' // dimma färgen
			continue
		}

		// We're only interested in paths relative from where we are,
		// so we only paint exits outwards.
		if _, done := rendered[adjacent.ID]; done {
			continue
		}

		dirchar := ' '

		diffx, diffy := room.Displacement(direction)
		if diffx == 0 && diffy == 0 {
			continue
		}

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
			log.Println("unknown exit:", direction)
		}

		// One point offset, when rendered, equals four or two cells.
		if diffx > 0 || diffx < 0 {
			diffx = diffx*4 - rel(3, diffx)
		}
		if diffy > 0 || diffy < 0 {
			diffy = diffy*2 - rel(1, diffy)
		}

		steps := max(1, max(abs(diffx), abs(diffy)))
		for i := 1; i <= steps; i++ {
			// X is offset by one because rooms are three cells
			// wide and paths need to start "one out".
			xx := x + rel(1, diffx) + diffx/steps*i
			yy := y + diffy/steps*i
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

		diffx, diffy := room.Displacement(direction)
		if diffx == 0 && diffy == 0 {
			continue
		}

		rows = renderRoom(rows, adjacent, x+4*diffx, y+2*diffy, depth-1, rendered)
	}

	return rows
}
