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

func msger(msg Message) func() Message {
	return func() Message { return msg }
}

var messages = map[string]func() Message{
	(&CharLogin{}).ID():      msger(&CharLogin{}),
	(&CharName{}).ID():       msger(&CharName{}),
	(&CharStatus{}).ID():     msger(&CharStatus{}),
	(&CharStatusVars{}).ID(): msger(&CharStatusVars{}),
	(&CharVitals{}).ID():     msger(&CharVitals{}),

	(&CharAfflictionsList{}).ID():   msger(&CharAfflictionsList{}),
	(&CharAfflictionsAdd{}).ID():    msger(&CharAfflictionsAdd{}),
	(&CharAfflictionsRemove{}).ID(): msger(&CharAfflictionsRemove{}),

	(&CharDefencesList{}).ID():   msger(&CharDefencesList{}),
	(&CharDefencesAdd{}).ID():    msger(&CharDefencesAdd{}),
	(&CharDefencesRemove{}).ID(): msger(&CharDefencesRemove{}),

	(&CharItemsContents{}).ID(): msger(&CharItemsContents{}),
	(&CharItemsInv{}).ID():      msger(&CharItemsInv{}),
	(&CharItemsRoom{}).ID():     msger(&CharItemsRoom{}),
	(&CharItemsList{}).ID():     msger(&CharItemsList{}),
	(&CharItemsAdd{}).ID():      msger(&CharItemsAdd{}),
	(&CharItemsRemove{}).ID():   msger(&CharItemsRemove{}),
	(&CharItemsUpdate{}).ID():   msger(&CharItemsUpdate{}),

	(&CharSkillsGet{}).ID():    msger(&CharSkillsGet{}),
	(&CharSkillsGroups{}).ID(): msger(&CharSkillsGroups{}),
	(&CharSkillsInfo{}).ID():   msger(&CharSkillsInfo{}),
	(&CharSkillsList{}).ID():   msger(&CharSkillsList{}),

	(&CommChannelEnable{}).ID():  msger(&CommChannelEnable{}),
	(&CommChannelEnd{}).ID():     msger(&CommChannelEnd{}),
	(&CommChannelList{}).ID():    msger(&CommChannelList{}),
	(&CommChannelPlayers{}).ID(): msger(&CommChannelPlayers{}),
	(&CommChannelStart{}).ID():   msger(&CommChannelStart{}),
	(&CommChannelText{}).ID():    msger(&CommChannelText{}),

	(&CoreGoodbye{}).ID():        msger(&CoreGoodbye{}),
	(&CoreHello{}).ID():          msger(&CoreHello{}),
	(&CoreKeepAlive{}).ID():      msger(&CoreKeepAlive{}),
	(&CorePing{}).ID():           msger(&CorePing{}),
	(&CoreSupportsSet{}).ID():    msger(&CoreSupportsSet{}),
	(&CoreSupportsAdd{}).ID():    msger(&CoreSupportsAdd{}),
	(&CoreSupportsRemove{}).ID(): msger(&CoreSupportsRemove{}),

	(&RoomInfo{}).ID():         msger(&RoomInfo{}),
	(&RoomPlayers{}).ID():      msger(&RoomPlayers{}),
	(&RoomAddPlayer{}).ID():    msger(&RoomAddPlayer{}),
	(&RoomRemovePlayer{}).ID(): msger(&RoomRemovePlayer{}),
	(&RoomWrongDir{}).ID():     msger(&RoomWrongDir{}),
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

func Marshal(msg Message) string {
	data, _ := json.Marshal(msg)
	return fmt.Sprintf("%s %s", msg.ID(), string(data))
}

func Unmarshal(data []byte, msg Message) error {
	data = bytes.TrimSpace(bytes.TrimPrefix(data, []byte(msg.ID())))

	err := json.Unmarshal(data, msg)
	if err != nil {
		return err
	}

	return nil
}

func SplitRank(str string) (string, string) {
	parts := strings.SplitN(str, "(", 2)

	if len(parts) == 1 {
		return str, ""
	}

	return strings.TrimSpace(parts[0]), strings.Trim(parts[1], " (%)")
}

func SplitRankInt(str string) (string, int) {
	str, strRank := SplitRank(str)

	rank, err := strconv.Atoi(strRank)
	if err != nil {
		return "", 0
	}

	return str, rank
}

func SplitRankFloat(str string) (string, float64) {
	str, strRank := SplitRank(str)

	rank, err := strconv.ParseFloat(strRank, 64)
	if err != nil {
		return "", 0
	}

	return str, rank
}
