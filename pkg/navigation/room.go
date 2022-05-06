package navigation

import (
	"github.com/icza/gox/gox"
	"github.com/tobiassjosten/nogfx/pkg/gmcp"
)

var (
	areas = map[int]*Area{}
	rooms = map[int]*Room{}
)

type Area struct {
	ID   int
	Name string
}

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

func (room *Room) HasExit(direction string) bool {
	_, exists := room.Exits[direction]
	return exists
}

func (room *Room) Displacement(room2 *Room, direction string) (int, int) {
	if room.X != nil && room.Y != nil && room2.X != nil && room2.Y != nil {
		// Achaea's coordinate system, through GMCP, puts the origin
		// 0,0 in the middle of the map, instead of the top left (nw)
		// corner like when we draw the map. So north increases Y and
		// south decreases it, reverse of our needs.
		return *room2.X - *room.X, *room.Y - *room2.Y
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
