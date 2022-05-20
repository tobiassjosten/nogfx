package gmcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

var (
	_ ClientMessage = &IRERiftRequest{}
	_ ServerMessage = &IRERiftChange{}
	_ ServerMessage = &IRERiftList{}

	_ ClientMessage = IRETargetSet{}
	_ ServerMessage = IRETargetSet{}
	_ ServerMessage = &IRETargetInfo{}
)

// IRERiftItem is an item in rift storage.
type IRERiftItem struct {
	Name        string `json:"name"`
	Amount      int    `json:"amount,string"`
	Description string `json:"desc"`
}

// IRERiftChange is a server-sent GMCP message lists items in rift storage.
type IRERiftChange IRERiftItem

// Hydrate populates the message with data.
func (msg IRERiftChange) Hydrate(data []byte) (ServerMessage, error) {
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// IRERiftList is a server-sent GMCP message lists items in rift storage.
type IRERiftList []IRERiftItem

// Hydrate populates the message with data.
func (msg IRERiftList) Hydrate(data []byte) (ServerMessage, error) {
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// IRERiftRequest is a client-sent GMCP message to request a list of items in
// the player's inventory.
type IRERiftRequest struct{}

// String is the message's string representation.
func (msg IRERiftRequest) String() string {
	return "IRE.Rift.Request"
}

// IRETargetSet is both a a client- and server-sent GMCP message to either set
// or verify the setting of the in-game target variable.
type IRETargetSet struct {
	Target string
}

// String is the message's string representation.
func (msg IRETargetSet) String() string {
	return strings.TrimSpace(fmt.Sprintf(`IRE.Target.Set "%s"`, msg.Target))
}

// Hydrate populates the message with data.
func (msg IRETargetSet) Hydrate(data []byte) (ServerMessage, error) {
	err := json.Unmarshal(data, &msg.Target)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// IRETargetInfo is both a a client- and server-sent GMCP message with
// additional information about the current active server side target.
type IRETargetInfo struct {
	ID          string `json:"id"` // @todo Check if ever a non-number.
	Health      int
	Description string `json:"short_desc"`
}

// Hydrate populates the message with data.
func (msg IRETargetInfo) Hydrate(data []byte) (ServerMessage, error) {
	if bytes.Equal(data, []byte(`""`)) {
		return msg, nil
	}

	type IRETargetInfoAlias IRETargetInfo
	var child struct {
		IRETargetInfoAlias
		CHealth string `json:"hpperc"`
	}

	err := json.Unmarshal(data, &child)
	if err != nil {
		return nil, err
	}

	msg = (IRETargetInfo)(child.IRETargetInfoAlias)

	health, err := strconv.Atoi(strings.Trim(child.CHealth, "%)"))
	if err != nil {
		return nil, err
	}
	msg.Health = health

	return msg, nil
}
