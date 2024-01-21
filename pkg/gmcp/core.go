package gmcp

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// CoreGoodbye is a GMCP message finishing a game session.
type CoreGoodbye struct{}

// ID is the prefix before the message's data.
func (msg *CoreGoodbye) ID() string {
	return "Core.Goodbye"
}

// Marshal converts the message to a string.
func (msg *CoreGoodbye) Marshal() string {
	return msg.ID()
}

// Unmarshal populates the message with data.
func (msg *CoreGoodbye) Unmarshal(_ []byte) error {
	return nil
}

// CoreHello is a GMCP message used to identify the client. It has to be the
// first message sent.
type CoreHello struct {
	Client  string `json:"client"`
	Version string `json:"version"`
}

// ID is the prefix before the message's data.
func (msg *CoreHello) ID() string {
	return "Core.Hello"
}

// Marshal converts the message to a string.
func (msg *CoreHello) Marshal() string {
	return Marshal(msg)
}

// Unmarshal populates the message with data.
func (msg *CoreHello) Unmarshal(data []byte) error {
	return Unmarshal(data, msg)
}

// CoreKeepAlive is a GMCP message resetting the timeout counter.
type CoreKeepAlive struct{}

// ID is the prefix before the message's data.
func (msg *CoreKeepAlive) ID() string {
	return "Core.KeepAlive"
}

// Marshal converts the message to a string.
func (msg *CoreKeepAlive) Marshal() string {
	return msg.ID()
}

// Unmarshal populates the message with data.
func (msg *CoreKeepAlive) Unmarshal(_ []byte) error {
	return nil
}

// CorePing is a GMCP message measuring latency.
type CorePing struct {
	Latency *int
}

// ID is the prefix before the message's data.
func (msg *CorePing) ID() string {
	return "Core.Ping"
}

// Marshal converts the message to a string.
func (msg *CorePing) Marshal() string {
	if msg.Latency != nil {
		return fmt.Sprintf("%s %d", msg.ID(), *msg.Latency)
	}

	return msg.ID()
}

// Unmarshal populates the message with data.
func (msg *CorePing) Unmarshal(data []byte) error {
	data = bytes.TrimSpace(bytes.TrimPrefix(data, []byte(msg.ID())))

	if len(data) == 0 {
		return nil
	}

	err := json.Unmarshal(data, &msg.Latency)
	if err != nil {
		return err
	}

	return nil
}
