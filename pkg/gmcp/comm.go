package gmcp

import (
	"encoding/json"
	"fmt"
	"strings"
)

var (
	_ ClientMessage = CommChannelEnable("")
	_ ClientMessage = &CommChannelPlayers{}

	_ ServerMessage = CommChannelEnd("")
	_ ServerMessage = &CommChannelList{}
	_ ServerMessage = &CommChannelPlayers{}
	_ ServerMessage = CommChannelStart("")
	_ ServerMessage = &CommChannelText{}
)

// CommChannelEnable is a client-sent GMCP message used to tell the game to
// turn on a character channel without typing in a command line command
type CommChannelEnable string

// String is the message's string representation.
func (msg CommChannelEnable) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Comm.Channel.Enable "%s"`, string(msg)))
}

// CommChannelEnd is a server-sent GMCP message informing the client that text
// that follows is something said over a communication channel.
type CommChannelEnd string

// Hydrate populates the message with data.
func (msg CommChannelEnd) Hydrate(data []byte) (ServerMessage, error) {
	var channel string

	err := json.Unmarshal(data, &channel)
	if err != nil {
		return nil, err
	}

	return CommChannelEnd(channel), nil
}

// CommChannel contains information about an in-game channel.
type CommChannel struct {
	Name    string `json:"name"`
	Caption string `json:"caption"`
	Command string `json:"command"`
}

// CommChannelList is a server-sent GMCP message listing communication
// channels available to the player.
type CommChannelList []CommChannel

// Hydrate populates the message with data.
func (msg CommChannelList) Hydrate(data []byte) (ServerMessage, error) {
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// CommChannelPlayer shows which channels are shared with a specific player.
type CommChannelPlayer struct {
	Name     string   `json:"name"`
	Channels []string `json:"channels"`
}

// CommChannelPlayers is both a client-sent and server-sent GMCP message, to
// either request data or lists players and which channels (if any) they share
// with the player's character.
type CommChannelPlayers []CommChannelPlayer

// Hydrate populates the message with data.
func (msg CommChannelPlayers) Hydrate(data []byte) (ServerMessage, error) {
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// String is the message's string representation.
func (msg CommChannelPlayers) String() string {
	return "Comm.Channel.Players"
}

// CommChannelStart is a server-sent GMCP message informing the client that text
// that follows is something said over a communication channel.
type CommChannelStart string

// Hydrate populates the message with data.
func (msg CommChannelStart) Hydrate(data []byte) (ServerMessage, error) {
	var channel string

	err := json.Unmarshal(data, &channel)
	if err != nil {
		return nil, err
	}

	return CommChannelStart(channel), nil
}

// CommChannelText is both a client-sent and server-sent GMCP message, to
// either request data or lists players and which channels (if any) they share
// with the player's character.
type CommChannelText struct {
	Channel string `json:"channel"`
	Talker  string `json:"talker"`
	Text    string `json:"text"`
}

// Hydrate populates the message with data.
func (msg CommChannelText) Hydrate(data []byte) (ServerMessage, error) {
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
