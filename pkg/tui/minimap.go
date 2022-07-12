package tui

import (
	"github.com/tobiassjosten/nogfx/pkg/navigation"

	"github.com/gdamore/tcell/v2"
)

// Minimap is a map rendition based on the given room.
type Minimap struct {
	room     *navigation.Room
	rows     Rows
	rendered map[int]struct{}
}

type maproom struct {
	room *navigation.Room
	x    int
	y    int
}

// RenderMap renders a map from the current room.
func (tui *TUI) RenderMap(width, height int) Rows {
	if rows, ok := tui.getCache(paneMap); ok {
		return rows
	}

	rows := RenderMap(tui.room, width, height)

	tui.setCache(paneMap, rows)

	return rows
}

// RenderMap renders cascading layers of adjacent rooms, based on the given.
func RenderMap(room *navigation.Room, width, height int) Rows {
	if room == nil || width == 0 || height == 0 {
		return Rows{}
	}

	mmap := Minimap{room, NewRows(width, height), map[int]struct{}{}}

	rooms := []maproom{{
		room: room,
		x:    len(mmap.rows[0]) / 2,
		y:    len(mmap.rows) / 2,
	}}

	for len(rooms) > 0 {
		q := rooms[0]
		rooms = append(rooms[1:], mmap.render(q.room, q.x, q.y)...)
	}

	return mmap.rows
}

func (mmap Minimap) render(room *navigation.Room, x, y int) []maproom {
	// Make sure we have enough padding to render room and exits.
	if x < 2 || y < 1 || y > len(mmap.rows)-2 || x > len(mmap.rows[0])-3 {
		return nil
	}

	if room.ID != 0 {
		if _, done := mmap.rendered[room.ID]; done {
			return nil
		}
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

	var adjacents []maproom

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

		case "in":
			mmap.rows[y][x+1].Content = '}'

		case "out":
			mmap.rows[y][x-1].Content = '{'
		}

		if adjacent.ID != 0 {
			if _, done := mmap.rendered[adjacent.ID]; done {
				continue
			}
		}

		diffx, diffy := room.Displacement(direction)
		if diffx == 0 && diffy == 0 {
			continue
		}

		if room.Area == nil || adjacent.Area == nil || room.Area.ID == adjacent.Area.ID {
			adjacents = append(adjacents, maproom{
				room: adjacent,
				x:    x + 4*diffx,
				y:    y + 2*diffy,
			})
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

		case "in", "out":
			continue
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
				break
			}

			if (dirchar == '/' || dirchar == '\\') && mmap.rows[yy][xx].Content != ' ' {
				dirchar = 'X'
			}

			mmap.rows[yy][xx].Content = dirchar
		}
	}

	return adjacents
}
