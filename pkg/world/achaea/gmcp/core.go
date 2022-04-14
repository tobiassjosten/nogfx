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
	Char        bool
	CharSkills  bool
	CharItems   bool
	CommChannel bool
	Room        bool
	IRERift     bool
}

// String is the message's string representation.
func (msg CoreSupportsSet) String() string {
	list := []string{}
	if msg.Char {
		list = append(list, "Char 1")
	}
	if msg.CharSkills {
		list = append(list, "Char.Skills 1")
	}
	if msg.CharItems {
		list = append(list, "Char.Items 1")
	}
	if msg.CommChannel {
		list = append(list, "Comm.Channel 1")
	}
	if msg.Room {
		list = append(list, "Room 1")
	}
	if msg.IRERift {
		list = append(list, "IRE.Rift 1")
	}

	data, err := json.Marshal(list)
	if err != nil {
		data = []byte("[]")
	}

	return fmt.Sprintf("Core.Supports.Set %s", data)
}
