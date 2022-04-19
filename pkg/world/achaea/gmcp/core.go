package gmcp

import (
	"encoding/json"
	"fmt"
)

var (
	_ ClientMessage = &CoreHello{}
	_ ClientMessage = &CoreSupportsSet{}
)

// CoreHello is a client-sent GMCP message used to identify the client. It has
// to be the first message sent.
type CoreHello struct {
	Client  string `json:"client"`
	Version string `json:"version"`
}

// String is the message's string representation.
func (msg CoreHello) String() string {
	data, err := json.Marshal(msg)
	if err != nil {
		data = []byte("{}")
	}

	return fmt.Sprintf("Core.Hello %s", data)
}

// CoreSupportsSet is a client-sent GMCP message containing packages supported.
type CoreSupportsSet struct {
	Char        *int
	CharSkills  *int
	CharItems   *int
	CommChannel *int
	Room        *int
	IRERift     *int
}

// String is the message's string representation.
func (msg CoreSupportsSet) String() string {
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

	data, err := json.Marshal(list)
	if err != nil {
		data = []byte("[]")
	}

	return fmt.Sprintf("Core.Supports.Set %s", data)
}
