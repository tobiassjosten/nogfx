package gmcp

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// ClientMessage is a GMCP message sent from the client.
type ClientMessage interface {
	String() string
}

// ServerMessage is a GMCP message sent from the server.
type ServerMessage interface {
	Hydrate([]byte) (ServerMessage, error)
}

// @todo Consider turning all messages into structs, for consistency.
var serverMessages = map[string]ServerMessage{
	"Comm.Channel.End":     CommChannelEnd(""),
	"Comm.Channel.List":    CommChannelList{},
	"Comm.Channel.Players": CommChannelPlayers{},
	"Comm.Channel.Start":   CommChannelStart(""),
	"Comm.Channel.Text":    CommChannelText{},

	"Char.Afflictions.Add":    CharAfflictionsAdd{},
	"Char.Afflictions.List":   CharAfflictionsList{},
	"Char.Afflictions.Remove": CharAfflictionsRemove{},

	"Char.Defences.Add":    CharDefencesAdd{},
	"Char.Defences.List":   CharDefencesList{},
	"Char.Defences.Remove": CharDefencesRemove{},

	"Char.Items.Add":    CharItemsAdd{},
	"Char.Items.List":   CharItemsList{},
	"Char.Items.Remove": CharItemsRemove{},
	"Char.Items.Update": CharItemsUpdate{},

	"Char.Name": CharName{},

	"Char.Skills.Groups": CharSkillsGroups{},
	"Char.Skills.Info":   CharSkillsInfo{},
	"Char.Skills.List":   CharSkillsList{},

	"Char.Status":     CharStatus{},
	"Char.StatusVars": CharStatusVars{},

	"Char.Vitals": CharVitals{},

	"Core.Goodbye": CoreGoodbye{},
	"Core.Ping":    CorePing{},

	"IRE.Rift.Change": IRERiftChange{},
	"IRE.Rift.List":   IRERiftList{},

	"IRE.Target.Set":  IRETargetSet(""),
	"IRE.Target.Info": &IRETargetInfo{},

	"Room.Info":         RoomInfo{},
	"Room.Players":      RoomPlayers{},
	"Room.AddPlayer":    RoomAddPlayer{},
	"Room.RemovePlayer": RoomRemovePlayer{},
	"Room.WrongDir":     RoomWrongDir(""),
}

// Parse converts a byte slice into a GMCP message.
func Parse(command []byte) (ServerMessage, error) {
	parts := bytes.SplitN(command, []byte{' '}, 2)

	message, ok := serverMessages[string(parts[0])]
	if !ok {
		return nil, fmt.Errorf("unknown message '%s'", parts[0])
	}

	if len(parts) == 1 {
		parts = append(parts, []byte{})
	}

	msg, err := message.Hydrate(parts[1])
	if err != nil {
		return nil, fmt.Errorf(
			"failed hydrating %T (%s): %w",
			message, parts[1], err,
		)
	}

	return msg, nil
}

func splitRank(str string) (string, *int) {
	parts := strings.SplitN(str, "(", 2)
	name := strings.Trim(parts[0], " ")

	var rank *int
	if len(parts) > 1 {
		r, err := strconv.Atoi(strings.Trim(parts[1], "%)"))
		if err == nil {
			rank = &r
		}
	}

	return name, rank
}

func splitLevelRank(str string) (int, *int) {
	name, rank := splitRank(str)
	level, _ := strconv.Atoi(name)

	return level, rank
}
