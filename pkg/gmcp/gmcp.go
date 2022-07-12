package gmcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/tobiassjosten/nogfx/pkg/telnet"
)

// Message is a GMCP data object.
type Message interface {
	ID() string
	Marshal() string
	Unmarshal([]byte) error
}

var messages = map[string]func() Message{
	(&CharLogin{}).ID():      func() Message { return &CharLogin{} },
	(&CharName{}).ID():       func() Message { return &CharName{} },
	(&CharStatusVars{}).ID(): func() Message { return &CharStatusVars{} },

	(&CharAfflictionsList{}).ID():   func() Message { return &CharAfflictionsList{} },
	(&CharAfflictionsAdd{}).ID():    func() Message { return &CharAfflictionsAdd{} },
	(&CharAfflictionsRemove{}).ID(): func() Message { return &CharAfflictionsRemove{} },

	(&CharDefencesList{}).ID():   func() Message { return &CharDefencesList{} },
	(&CharDefencesAdd{}).ID():    func() Message { return &CharDefencesAdd{} },
	(&CharDefencesRemove{}).ID(): func() Message { return &CharDefencesRemove{} },

	(&CharItemsContents{}).ID(): func() Message { return &CharItemsContents{} },
	(&CharItemsInv{}).ID():      func() Message { return &CharItemsInv{} },
	(&CharItemsRoom{}).ID():     func() Message { return &CharItemsRoom{} },
	(&CharItemsList{}).ID():     func() Message { return &CharItemsList{} },
	(&CharItemsAdd{}).ID():      func() Message { return &CharItemsAdd{} },
	(&CharItemsRemove{}).ID():   func() Message { return &CharItemsRemove{} },
	(&CharItemsUpdate{}).ID():   func() Message { return &CharItemsUpdate{} },

	(&CharSkillsGet{}).ID():    func() Message { return &CharSkillsGet{} },
	(&CharSkillsGroups{}).ID(): func() Message { return &CharSkillsGroups{} },
	(&CharSkillsInfo{}).ID():   func() Message { return &CharSkillsInfo{} },
	(&CharSkillsList{}).ID():   func() Message { return &CharSkillsList{} },

	(&CommChannelEnable{}).ID():  func() Message { return &CommChannelEnable{} },
	(&CommChannelList{}).ID():    func() Message { return &CommChannelList{} },
	(&CommChannelPlayers{}).ID(): func() Message { return &CommChannelPlayers{} },
	(&CommChannelText{}).ID():    func() Message { return &CommChannelText{} },

	(&CoreGoodbye{}).ID():        func() Message { return &CoreGoodbye{} },
	(&CoreHello{}).ID():          func() Message { return &CoreHello{} },
	(&CoreKeepAlive{}).ID():      func() Message { return &CoreKeepAlive{} },
	(&CorePing{}).ID():           func() Message { return &CorePing{} },
	(&CoreSupportsSet{}).ID():    func() Message { return &CoreSupportsSet{} },
	(&CoreSupportsAdd{}).ID():    func() Message { return &CoreSupportsAdd{} },
	(&CoreSupportsRemove{}).ID(): func() Message { return &CoreSupportsRemove{} },

	(&RoomInfo{}).ID():         func() Message { return &RoomInfo{} },
	(&RoomPlayers{}).ID():      func() Message { return &RoomPlayers{} },
	(&RoomAddPlayer{}).ID():    func() Message { return &RoomAddPlayer{} },
	(&RoomRemovePlayer{}).ID(): func() Message { return &RoomRemovePlayer{} },
}

// Parse converts a byte slice into a GMCP message.
func Parse(data []byte) (Message, error) {
	parts := strings.SplitN(string(data), " ", 2)

	if _, ok := messages[parts[0]]; !ok {
		return nil, fmt.Errorf("unknown message '%s'", parts[0])
	}
	msg := messages[parts[0]]()

	if err := msg.Unmarshal(data); err != nil {
		return nil, fmt.Errorf("couldn't unmarshal %T: %w", msg, err)
	}

	return msg, nil
}

var (
	gmcpPrefix = []byte{telnet.IAC, telnet.SB, telnet.GMCP}
	gmcpSuffix = []byte{telnet.IAC, telnet.SE}
)

// Wrap embeds a GMCP message in a telnet negotiation sequence.
func Wrap(data []byte) []byte {
	return append(append(gmcpPrefix, data...), gmcpSuffix...)
}

// Unwrap removes telnet control codes from a GMCP message. Returns nil if the
// command actually isn't a GMCP message.
func Unwrap(data []byte) []byte {
	if !bytes.HasPrefix(data, gmcpPrefix) {
		return nil
	}
	if !bytes.HasSuffix(data, gmcpSuffix) {
		return nil
	}

	return data[len(gmcpPrefix) : len(data)-len(gmcpSuffix)]
}

// Marshal converts a Message to JSON data.
func Marshal(msg Message) string {
	data, _ := json.Marshal(msg)
	return fmt.Sprintf("%s %s", msg.ID(), string(data))
}

// Unmarshal hydrates a Message from JSON data.
func Unmarshal(data []byte, msg Message) error {
	data = bytes.TrimSpace(bytes.TrimPrefix(data, []byte(msg.ID())))

	err := json.Unmarshal(data, msg)
	if err != nil {
		return err
	}

	return nil
}

// SplitRank parses values like "Cityname (2)" into separate parts.
func SplitRank(str string) (string, string) {
	parts := strings.SplitN(str, "(", 2)

	if len(parts) == 1 {
		return str, ""
	}

	return strings.TrimSpace(parts[0]), strings.Trim(parts[1], " (%)")
}

// SplitRankInt wraps SplitRank() and converts the rank part to an integer.
func SplitRankInt(str string) (string, int) {
	str, strRank := SplitRank(str)

	rank, err := strconv.Atoi(strRank)
	if err != nil {
		return str, 0
	}

	return str, rank
}

// SplitRankFloat wraps SplitRank() and converts the rank part to a float.
func SplitRankFloat(str string) (string, float64) {
	str, strRank := SplitRank(str)

	rank, err := strconv.ParseFloat(strRank, 64)
	if err != nil {
		return str, 0
	}

	return str, rank
}
