package gmcp

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/tobiassjosten/nogfx/pkg/telnet"
)

// ClientMessage is a GMCP message sent from the client.
type ClientMessage interface {
	String() string
}

// ServerMessage is a GMCP message sent from the server.
type ServerMessage interface {
	Hydrate([]byte) (ServerMessage, error)
}

// ServerMessages maps GMCP messages to associated structs.
// @todo Consider turning all messages into structs, for consistency.
var ServerMessages = map[string]ServerMessage{
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
func Parse(command []byte, messages map[string]ServerMessage) (ServerMessage, error) {
	parts := bytes.SplitN(command, []byte{' '}, 2)

	message, ok := messages[string(parts[0])]
	if !ok {
		if _, ok := ServerMessages[string(parts[0])]; ok {
			return Parse(command, ServerMessages)
		}
		return nil, fmt.Errorf("unknown message '%s'", parts[0])
	}

	// Some messages don't have a message body but we want each message to
	// be responsible for its own hydration and validation. So we mock
	// missing bodies and proceed with hydration as normal.
	if len(parts) == 1 {
		parts = append(parts, []byte{})
	}

	msg, err := message.Hydrate(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed hydrating %T: %w", message, err)
	}

	return msg, nil
}

var (
	gmcpPrefix = []byte{telnet.IAC, telnet.SB, telnet.GMCP}
	gmcpSuffix = []byte{telnet.IAC, telnet.SE}
)

// Wrap embeds a GMCP message in a telnet negotiation sequence.
func Wrap(gmcp []byte) []byte {
	return append(append(gmcpPrefix, gmcp...), gmcpSuffix...)
}

// Unwrap removes telnet control codes from a GMCP message. Returns nil if the
// command actually isn't a GMCP message.
func Unwrap(command []byte) []byte {
	if !bytes.HasPrefix(command, gmcpPrefix) {
		return nil
	}
	if !bytes.HasSuffix(command, gmcpSuffix) {
		return nil
	}

	return command[len(gmcpPrefix) : len(command)-len(gmcpSuffix)]
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
