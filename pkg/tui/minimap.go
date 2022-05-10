package tui

import (
	"log"

	"github.com/tobiassjosten/nogfx/pkg/navigation"

	"github.com/gdamore/tcell/v2"
)

// Minimap is a map rendition based on the given room.
type Minimap struct {
	room     *navigation.Room
	rows     Rows
	rendered map[int]struct{}
}

// RenderMap renders cascading layers of adjacent rooms, based on the given.
func RenderMap(room *navigation.Room, width, height int) Rows {
	if room == nil || width == 0 || height == 0 {
		return Rows{}
	}

	m := Minimap{room, NewRows(width, height), map[int]struct{}{}}

	return m.render(room, len(m.rows[0])/2, len(m.rows)/2, 5)
}

func (mmap Minimap) render(room *navigation.Room, x, y, depth int) Rows {
	// Make sure we have enough padding to render room borders and exits.
	if x < 2 || y < 2 || y > len(mmap.rows)-3 || x > len(mmap.rows[0])-3 {
		return mmap.rows
	}

	mmap.rows[y][x].Content = ' '
	if room.HasPlayer {
		mmap.rows[y][x].Content = '+'
	}

	mmap.rendered[room.ID] = struct{}{}

	mmap.rows[y][x-1].Content = '['
	mmap.rows[y][x+1].Content = ']'

	// @todo Consider which environments that are actually important for
	// the map to visualise and color the brackets accordingly.

	if !room.Known {
		mmap.rows[y][x-1].Style = mmap.rows[y][x-1].Style.Foreground(tcell.Color237)
		mmap.rows[y][x+1].Style = mmap.rows[y][x+1].Style.Foreground(tcell.Color237)
	}

paths:
	for direction, adjacent := range room.Exits {
		switch direction {
		case "u":
			if room.HasPlayer {
				continue
			}

			if room.HasExit("d") {
				mmap.rows[y][x].Content = '='
				mmap.rows[y][x].Foreground(tcell.Color245)
				continue
			}

			mmap.rows[y][x].Content = '^'
			mmap.rows[y][x].Foreground(tcell.Color245)
			continue

		case "d":
			if room.HasPlayer {
				continue
			}

			if room.HasExit("u") {
				continue
			}

			mmap.rows[y][x].Content = 'v'
			mmap.rows[y][x].Foreground(tcell.Color245)
			continue

		case "out":
			mmap.rows[y][x-1].Content = '{'
			continue

		case "in":
			mmap.rows[y][x+1].Content = '}'
			continue
		}

		diffx, diffy := room.Displacement(direction)
		if diffx == 0 && diffy == 0 {
			continue
		}

		if _, done := mmap.rendered[adjacent.ID]; done {
			continue
		}

		var dirchar rune
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

			if yy < 0 || xx < 0 || yy >= len(mmap.rows) || xx >= len(mmap.rows[yy]) {
				log.Println("overflow", len(mmap.rows), len(mmap.rows[0]), "coords", yy, xx)
				break
			}

			mmap.rows[yy][xx].Content = dirchar
		}
	}

	if depth == 0 {
		return mmap.rows
	}

	for direction, adjacent := range room.Exits {
		if _, done := mmap.rendered[adjacent.ID]; done {
			continue
		}

		if room.Area != nil && adjacent.Area != nil && room.Area.ID != adjacent.Area.ID {
			continue
		}

		diffx, diffy := room.Displacement(direction)
		if diffx == 0 && diffy == 0 {
			continue
		}

		mmap.rows = mmap.render(adjacent, x+4*diffx, y+2*diffy, depth-1)
	}

	return mmap.rows
}
