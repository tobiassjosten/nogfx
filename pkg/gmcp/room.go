package gmcp

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
)

var (
	_ ServerMessage = &RoomInfo{}
	_ ServerMessage = &RoomPlayers{}
	_ ServerMessage = &RoomAddPlayer{}
	_ ServerMessage = &RoomRemovePlayer{}
	_ ServerMessage = RoomWrongDir("")
)

// RoomDetails is a set of flags denoting how to interact with a room.
type RoomDetails struct {
	Bank        bool
	Indoors     bool
	Outdoors    bool
	Sewer       bool
	Shop        bool
	Subdivision bool
}

// UnmarshalJSON hydrates RoomDetails from a list of unstructured strings.
func (details *RoomDetails) UnmarshalJSON(data []byte) error {
	var list []string

	// This should only be invoked from RoomInfo.UnmarshalJSON(), so any
	// formatting errors will be caught there.
	_ = json.Unmarshal(data, &list)

	for _, item := range list {
		switch item {
		case "bank":
			details.Bank = true

		case "indoors":
			details.Outdoors = true

		case "outdoors":
			details.Outdoors = true

		case "sewer":
			details.Sewer = true

		case "shop":
			details.Shop = true

		case "subdivision":
			details.Subdivision = true

		default:
			log.Printf("unknown Room.Info detail '%s'", item)
		}
	}

	return nil
}

// RoomInfo is a server-sent GMCP message containing information about the
// room that the player is in.
type RoomInfo struct {
	Number      int    `json:"num"`
	Name        string `json:"name"`
	Description string `json:"desc"`
	AreaName    string `json:"area"`
	AreaNumber  int
	Environment string `json:"environment"`
	X           int
	Y           int
	Building    int
	Map         string         `json:"map"`
	Exits       map[string]int `json:"exits"`
	Details     RoomDetails    `json:"details"`
}

// Hydrate populates the message with data.
func (msg RoomInfo) Hydrate(data []byte) (ServerMessage, error) {
	type RoomInfoAlias RoomInfo
	var child struct {
		RoomInfoAlias
		Coords string `json:"coords"`
	}

	err := json.Unmarshal(data, &child)
	if err != nil {
		return nil, err
	}

	msg = (RoomInfo)(child.RoomInfoAlias)

	coords := strings.Split(child.Coords, ",")
	switch {
	case len(coords) == 4:
		building, err := strconv.Atoi(coords[3])
		if err != nil {
			return nil, fmt.Errorf(
				"failed parsing building from coords '%s': %w",
				coords, err,
			)
		}
		msg.Building = building

		fallthrough

	case len(coords) == 3:
		areaNumber, err := strconv.Atoi(coords[0])
		if err != nil {
			return nil, fmt.Errorf(
				"failed parsing area number from coords '%s': %w",
				coords, err,
			)
		}
		msg.AreaNumber = areaNumber

		x, err := strconv.Atoi(coords[1])
		if err != nil {
			return nil, fmt.Errorf(
				"failed parsing x from coords '%s': %w",
				coords, err,
			)
		}
		msg.X = x

		y, err := strconv.Atoi(coords[2])
		if err != nil {
			return nil, fmt.Errorf(
				"failed parsing y from coords '%s': %w",
				coords, err,
			)
		}
		msg.Y = y

	default:
		return nil, fmt.Errorf("failed parsing coords '%s'", coords)
	}

	return msg, nil
}

// RoomPlayer is a player joining, exiting in, or leaving a room.
type RoomPlayer struct {
	Name     string `json:"name"`
	Fullname string `json:"fullname"`
}

// RoomPlayers is a server-sent GMCP message containing basic information about
// players in the room.
type RoomPlayers []RoomPlayer

// Hydrate populates the message with data.
func (msg RoomPlayers) Hydrate(data []byte) (ServerMessage, error) {
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// RoomAddPlayer is a server-sent GMCP message containing basic information about
// players in the room.
type RoomAddPlayer RoomPlayer

// Hydrate populates the message with data.
func (msg RoomAddPlayer) Hydrate(data []byte) (ServerMessage, error) {
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// RoomRemovePlayer is a server-sent GMCP message containing basic information about
// players in the room.
type RoomRemovePlayer RoomPlayer

// Hydrate populates the message with data.
func (msg RoomRemovePlayer) Hydrate(data []byte) (ServerMessage, error) {
	var name string

	err := json.Unmarshal(data, &name)
	if err != nil {
		return nil, err
	}

	msg.Name = name

	return msg, nil
}

// RoomWrongDir is a server-sent GMCP message giving feedback when the player
// has tried a currently non-functional exit.
type RoomWrongDir string

// Hydrate populates the message with data.
func (msg RoomWrongDir) Hydrate(data []byte) (ServerMessage, error) {
	var exit string

	err := json.Unmarshal(data, &exit)
	if err != nil {
		return nil, err
	}

	return RoomWrongDir(exit), nil
}
