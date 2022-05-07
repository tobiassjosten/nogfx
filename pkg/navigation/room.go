package navigation

import (
	"github.com/icza/gox/gox"
	"github.com/tobiassjosten/nogfx/pkg/gmcp"
)

var (
	rooms = map[int]*Room{}
)

// Area is a region covering a number of rooms.
type Area struct {
	ID   int
	Name string
}

// Room represents a location within a game.
type Room struct {
	ID   int
	Name string

	X *int
	Y *int

	HasPlayer bool

	Known bool

	Area *Area

	Exits map[string]*Room
}

// RoomFromGMCP creates a Room from a GMCP Room.Info message.
func RoomFromGMCP(msg gmcp.RoomInfo) *Room {
	room, ok := rooms[msg.Number]
	if ok && room.Known {
		return room
	} else if !ok {
		room = &Room{ID: msg.Number}
		rooms[msg.Number] = room
	}

	room.Name = msg.Name
	room.X = gox.NewInt(msg.X)
	room.Y = gox.NewInt(msg.Y)
	room.Known = true

	if msg.Exits != nil {
		room.Exits = map[string]*Room{}

		for direction, number := range msg.Exits {
			adjacent, ok := rooms[number]
			if !ok {
				adjacent = &Room{ID: number}
				rooms[number] = adjacent
			}

			room.Exits[direction] = adjacent
		}
	}

	return room
}

// HasExit determines whether the room has a specific exit or not.
func (room *Room) HasExit(direction string) bool {
	_, exists := room.Exits[direction]
	return exists
}

// Displacement calculates the coordinate offset of the room at the given exit
// or 0, 0 if an offset couldn't be calculated.
func (room *Room) Displacement(direction string) (int, int) {
	if !room.HasExit(direction) {
		return 0, 0
	}

	adjacent := room.Exits[direction]

	if room.X != nil && room.Y != nil && adjacent.X != nil && adjacent.Y != nil {
		// Achaea's coordinate system, through GMCP, puts the origin
		// 0,0 in the middle of the map, instead of the top left (nw)
		// corner like when we draw the map. So north increases Y and
		// south decreases it, reverse of our needs.
		return *adjacent.X - *room.X, *room.Y - *adjacent.Y
	}

	switch direction {
	case "n":
		return 0, -1

	case "ne":
		return 1, -1

	case "e":
		return 1, 0

	case "se":
		return 1, 1

	case "s":
		return 0, 1

	case "sw":
		return -1, 1

	case "w":
		return -1, 0

	case "nw":
		return -1, -1

	case "u":
		if !room.HasExit("n") {
			return 0, -1
		}
		return 0, 0

	case "d":
		if !room.HasExit("s") {
			return 0, 1
		}
		return 0, 0

	case "out":
		if !room.HasExit("w") {
			return -1, 0
		}
		return 0, 0

	case "in":
		if !room.HasExit("e") {
			return 1, 0
		}
		return 0, 0
	}

	return 0, 0
}
