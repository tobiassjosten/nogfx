package gmcp

import (
	"encoding/json"
	"fmt"
)

var (
	_ ClientMessage = &CoreHello{}
	_ ClientMessage = &CoreKeepAlive{}
	_ ClientMessage = &CorePing{}
	_ ClientMessage = &CoreSupportsAdd{}
	_ ClientMessage = &CoreSupportsRemove{}
	_ ClientMessage = &CoreSupportsSet{}

	_ ServerMessage = &CoreGoodbye{}
	_ ServerMessage = &CorePing{}
)

// CoreGoodbye is a server-sent GMCP message finishing a game session.
type CoreGoodbye struct {
}

// Hydrate populates the message with data.
func (msg CoreGoodbye) Hydrate(_ []byte) (ServerMessage, error) {
	return msg, nil
}

// CoreHello is a client-sent GMCP message used to identify the client. It has
// to be the first message sent.
type CoreHello struct {
	Client  string `json:"client"`
	Version string `json:"version"`
}

// String is the message's string representation.
func (msg CoreHello) String() string {
	data, _ := json.Marshal(msg)
	return fmt.Sprintf("Core.Hello %s", data)
}

// CoreKeepAlive is a client-sent GMCP message resetting the timeout counter.
type CoreKeepAlive struct {
}

// String is the message's string representation.
func (msg CoreKeepAlive) String() string {
	return "Core.KeepAlive"
}

// CorePing is a client- and server-sent GMCP message measuring latency.
type CorePing struct {
	Latency *int
}

// Hydrate populates the message with data.
func (msg CorePing) Hydrate(_ []byte) (ServerMessage, error) {
	return msg, nil
}

// String is the message's string representation.
func (msg CorePing) String() string {
	if msg.Latency != nil {
		return fmt.Sprintf("Core.Ping %d", *msg.Latency)
	}
	return "Core.Ping"
}

// CoreSupports is a list of potentially supported modules.
type CoreSupports struct {
	Char        *int
	CharSkills  *int
	CharItems   *int
	CommChannel *int
	Room        *int
	IRERift     *int
}

func coreSupportsList(msg CoreSupports) []byte {
	list := []string{}
	if msg.Char != nil {
		list = append(list, fmt.Sprintf("Char %d", *msg.Char))
	}
	if msg.CharSkills != nil {
		list = append(list, fmt.Sprintf("Char.Skills %d", *msg.CharSkills))
	}
	if msg.CharItems != nil {
		list = append(list, fmt.Sprintf("Char.Items %d", *msg.CharItems))
	}
	if msg.CommChannel != nil {
		list = append(list, fmt.Sprintf("Comm.Channel %d", *msg.CommChannel))
	}
	if msg.Room != nil {
		list = append(list, fmt.Sprintf("Room %d", *msg.Room))
	}
	if msg.IRERift != nil {
		list = append(list, fmt.Sprintf("IRE.Rift %d", *msg.IRERift))
	}

	data, _ := json.Marshal(list)

	return data
}

// CoreSupportsSet is a client-sent GMCP message containing supported modules.
type CoreSupportsSet struct {
	CoreSupports
}

// String is the message's string representation.
func (msg CoreSupportsSet) String() string {
	return fmt.Sprintf("Core.Supports.Set %s", coreSupportsList(msg.CoreSupports))
}

// CoreSupportsAdd is a client-sent GMCP message adding supported modules.
type CoreSupportsAdd struct {
	CoreSupports
}

// String is the message's string representation.
func (msg CoreSupportsAdd) String() string {
	return fmt.Sprintf("Core.Supports.Add %s", coreSupportsList(msg.CoreSupports))
}

// CoreSupportsRemove is a client-sent GMCP message removing supported modules.
type CoreSupportsRemove struct {
	CoreSupports
}

// String is the message's string representation.
func (msg CoreSupportsRemove) String() string {
	return fmt.Sprintf("Core.Supports.Remove %s", coreSupportsList(msg.CoreSupports))
}
