package navigation

import (
	"strings"

	"github.com/icza/gox/gox"
	"github.com/tobiassjosten/nogfx/pkg/gmcp"
)

var (
	areas = map[int]*Area{}
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
func RoomFromGMCP(msg *gmcp.RoomInfo) *Room {
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

	area, ok := areas[msg.AreaNumber]
	if !ok {
		area = &Area{
			ID:   msg.AreaNumber,
			Name: msg.AreaName,
		}
		areas[msg.AreaNumber] = area
	}
	room.Area = area

	return room
}

// HasExit determines whether the room has a specific exit or not. It supports
// a sequence of exits as well, like "s se e", to determine whether a chain of
// adjacent rooms have the specific exits.
func (room *Room) HasExit(directions string) (exists bool) {
	for _, direction := range strings.Split(directions, " ") {
		if room, exists = room.Exits[direction]; !exists {
			return false
		}
	}

	return true
}

// HasAnyExits determine whether the room has ANY of the specific exit sequences.
func (room *Room) HasAnyExits(directionses ...string) bool {
	for _, directions := range directionses {
		if room.HasExit(directions) {
			return true
		}
	}

	return false
}

// Displacement calculates the coordinate offset of the room at the given exit
// or 0, 0 if an offset couldn't be calculated.
func (room *Room) Displacement(direction string) (int, int) {
	if !room.HasExit(direction) {
		return 0, 0
	}

	adjacent := room.Exits[direction]

	// Coordinates doesn't translate between areas, so we only use them for
	// room in the same area.
	if room.Area == nil || adjacent.Area == nil || room.Area.ID == adjacent.Area.ID {
		if room.X != nil && room.Y != nil && adjacent.X != nil && adjacent.Y != nil {
			// Achaea's coordinate system has north increasing Y
			// and south decreasing it, so we reverse that.
			return *adjacent.X - *room.X, *room.Y - *adjacent.Y
		}
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

	case "out":
		if !room.HasAnyExits("w", "n sw", "s nw") {
			return -1, 0
		}
		return 0, 0

	case "in":
		if !room.HasAnyExits("e", "n se", "s ne") {
			return 1, 0
		}
		return 0, 0
	}

	return 0, 0
}
