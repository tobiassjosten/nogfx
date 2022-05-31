package gmcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// RoomInfo is a server-sent GMCP message containing information about the
// room that the player is in.
type RoomInfo struct {
	Number      int            `json:"num"`
	Name        string         `json:"name"`
	AreaName    string         `json:"area"`
	AreaNumber  int            `json:"-"`
	Environment string         `json:"environment"`
	X           int            `json:"-"`
	Y           int            `json:"-"`
	Building    int            `json:"-"`
	Map         string         `json:"map"`
	Exits       map[string]int `json:"exits"`
	Details     []string       `json:"details"`
}

func (msg *RoomInfo) HasDetail(wanted string) bool {
	for _, detail := range msg.Details {
		if detail == wanted {
			return true
		}
	}

	return false
}

func (msg *RoomInfo) IsBank() bool {
	return msg.HasDetail("bank")
}

func (msg *RoomInfo) IsIndoors() bool {
	return msg.HasDetail("indoors")
}

func (msg *RoomInfo) IsOutdoors() bool {
	return msg.HasDetail("outdoors")
}

func (msg *RoomInfo) IsSewer() bool {
	return msg.HasDetail("sewer")
}

func (msg *RoomInfo) IsShop() bool {
	return msg.HasDetail("shop")
}

func (msg *RoomInfo) IsSubdivision() bool {
	return msg.HasDetail("subdivision")
}

func (msg *RoomInfo) IsWilderness() bool {
	return msg.HasDetail("wilderness")
}

// ID is the prefix before the message's data.
func (msg *RoomInfo) ID() string {
	return "Room.Info"
}

// Marshal converts the message to a string.
func (msg *RoomInfo) Marshal() string {
	proxy := struct {
		*RoomInfo
		PCoords  string         `json:"coords"`
		PExits   map[string]int `json:"exits"`
		PDetails []string       `json:"details"`
	}{
		RoomInfo: msg,
		PExits:   msg.Exits,
		PDetails: msg.Details,
	}

	if msg.AreaNumber != 0 && msg.X != 0 && msg.Y != 0 {
		proxy.PCoords = fmt.Sprintf("%d,%d,%d", msg.AreaNumber, msg.X, msg.Y)
		if msg.Building != 0 {
			proxy.PCoords += fmt.Sprintf(",%d", msg.Building)
		}
	}

	if msg.Exits == nil {
		proxy.PExits = map[string]int{}
	}

	if msg.Details == nil {
		proxy.PDetails = []string{}
	}

	data, _ := json.Marshal(proxy)
	return fmt.Sprintf("%s %s", msg.ID(), string(data))
}

// Unmarshal populates the message with data.
func (msg *RoomInfo) Unmarshal(data []byte) error {
	data = bytes.TrimSpace(bytes.TrimPrefix(data, []byte(msg.ID())))

	if msg == nil {
		*msg = RoomInfo{}
	}

	proxy := struct {
		*RoomInfo
		PCoords string `json:"coords"`
	}{
		RoomInfo: msg,
	}

	err := json.Unmarshal(data, &proxy)
	if err != nil {
		return err
	}

	*msg = (RoomInfo)(*proxy.RoomInfo)

	coords := strings.Split(proxy.PCoords, ",")
	switch {
	case proxy.PCoords == "":
		break

	case len(coords) == 4:
		building, err := strconv.Atoi(coords[3])
		if err != nil {
			return fmt.Errorf("failed parsing building from coords: %w", err)
		}
		msg.Building = building

		fallthrough

	case len(coords) == 3:
		areaNumber, err := strconv.Atoi(coords[0])
		if err != nil {
			return fmt.Errorf("failed parsing area number from coords: %w", err)
		}
		msg.AreaNumber = areaNumber

		x, err := strconv.Atoi(coords[1])
		if err != nil {
			return fmt.Errorf("failed parsing x from coords: %w", err)
		}
		msg.X = x

		y, err := strconv.Atoi(coords[2])
		if err != nil {
			return fmt.Errorf("failed parsing y from coords: %w", err)
		}
		msg.Y = y

	default:
		return fmt.Errorf("failed parsing coords '%s'", coords)
	}

	return nil
}

// RoomPlayer is a player entering, being in, or leaving a room.
type RoomPlayer struct {
	Name     string `json:"name"`
	Fullname string `json:"fullname"`
}

// RoomPlayers is a server-sent GMCP message containing basic information about
// players in the room.
type RoomPlayers []RoomPlayer

// ID is the prefix before the message's data.
func (msg *RoomPlayers) ID() string {
	return "Room.Players"
}

// Marshal converts the message to a string.
func (msg *RoomPlayers) Marshal() string {
	return Marshal(msg)
}

// Unmarshal populates the message with data.
func (msg *RoomPlayers) Unmarshal(data []byte) error {
	return Unmarshal(data, msg)
}

// RoomAddPlayer is a server-sent GMCP message containing basic information about
// players in the room.
type RoomAddPlayer RoomPlayer

// ID is the prefix before the message's data.
func (msg *RoomAddPlayer) ID() string {
	return "Room.AddPlayer"
}

// Marshal converts the message to a string.
func (msg *RoomAddPlayer) Marshal() string {
	return Marshal(msg)
}

// Unmarshal populates the message with data.
func (msg *RoomAddPlayer) Unmarshal(data []byte) error {
	return Unmarshal(data, msg)
}

// RoomRemovePlayer is a server-sent GMCP message containing basic information about
// players in the room.
type RoomRemovePlayer RoomPlayer

// ID is the prefix before the message's data.
func (msg *RoomRemovePlayer) ID() string {
	return "Room.RemovePlayer"
}

// Marshal converts the message to a string.
func (msg *RoomRemovePlayer) Marshal() string {
	return Marshal(msg)
}

// Unmarshal populates the message with data.
func (msg *RoomRemovePlayer) Unmarshal(data []byte) error {
	return Unmarshal(data, msg)
}

// RoomWrongDir is a server-sent GMCP message giving feedback when the player
// has tried a currently non-functional exit.
type RoomWrongDir struct {
	Direction string
}

// ID is the prefix before the message's data.
func (msg *RoomWrongDir) ID() string {
	return "Room.WrongDir"
}

// Marshal converts the message to a string.
func (msg *RoomWrongDir) Marshal() string {
	return fmt.Sprintf(`%s "%s"`, msg.ID(), msg.Direction)
}

// Unmarshal populates the message with data.
func (msg *RoomWrongDir) Unmarshal(data []byte) error {
	data = bytes.TrimSpace(bytes.TrimPrefix(data, []byte(msg.ID())))

	err := json.Unmarshal(data, &msg.Direction)
	if err != nil {
		return err
	}

	return nil
}
