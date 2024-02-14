package gmcp

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// CommChannelEnable is a GMCP message used to tell the game to turn on a
// character channel without typing in a command line command.
type CommChannelEnable struct {
	Channel string
}

// ID is the prefix before the message's data.
func (*CommChannelEnable) ID() string {
	return "Comm.Channel.Enable"
}

// Marshal converts the message to a string.
func (msg *CommChannelEnable) Marshal() string {
	return fmt.Sprintf("%s %q", msg.ID(), msg.Channel)
}

// Unmarshal populates the message with data.
func (msg *CommChannelEnable) Unmarshal(data []byte) error {
	data = bytes.TrimSpace(bytes.TrimPrefix(data, []byte(msg.ID())))

	err := json.Unmarshal(data, &msg.Channel)
	if err != nil {
		return err
	}

	return nil
}

// CommChannel contains information about an in-game channel.
type CommChannel struct {
	Name    string `json:"name"`
	Caption string `json:"caption"`
	Command string `json:"command"`
}

// CommChannelList is a GMCP message listing communication channels available
// to the player.
type CommChannelList []CommChannel

// ID is the prefix before the message's data.
func (*CommChannelList) ID() string {
	return "Comm.Channel.List"
}

// Marshal converts the message to a string.
func (msg *CommChannelList) Marshal() string {
	return Marshal(msg)
}

// Unmarshal populates the message with data.
func (msg *CommChannelList) Unmarshal(data []byte) error {
	return Unmarshal(data, msg)
}

// CommChannelPlayer shows which channels are shared with a specific player.
type CommChannelPlayer struct {
	Name     string   `json:"name"`
	Channels []string `json:"channels"`
}

// CommChannelPlayers is a GMCP message to either request data or lists
// players and which channels (if any) they share with the player's character.
type CommChannelPlayers []CommChannelPlayer

// ID is the prefix before the message's data.
func (*CommChannelPlayers) ID() string {
	return "Comm.Channel.Players"
}

// Marshal converts the message to a string.
func (msg *CommChannelPlayers) Marshal() string {
	return Marshal(msg)
}

// Unmarshal populates the message with data.
func (msg *CommChannelPlayers) Unmarshal(data []byte) error {
	return Unmarshal(data, msg)
}

// CommChannelText is a GMCP message with complete information about a
// communication.
type CommChannelText struct {
	Channel string `json:"channel"`
	Talker  string `json:"talker"`
	Text    string `json:"text"`
}

// ID is the prefix before the message's data.
func (*CommChannelText) ID() string {
	return "Comm.Channel.Text"
}

// Marshal converts the message to a string.
func (msg *CommChannelText) Marshal() string {
	return Marshal(msg)
}

// Unmarshal populates the message with data.
func (msg *CommChannelText) Unmarshal(data []byte) error {
	return Unmarshal(data, msg)
}
