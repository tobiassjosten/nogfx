package ironrealms

import "github.com/tobiassjosten/nogfx/pkg/gmcp"

// IRERiftItem is an item in rift storage.
type IRERiftItem struct {
	Amount      int    `json:"amount,string"`
	Description string `json:"desc"`
	Name        string `json:"name"`
}

// IRERiftChange is a GMCP message lists items in rift storage.
type IRERiftChange IRERiftItem

// ID is the prefix before the message's data.
func (*IRERiftChange) ID() string {
	return "IRE.Rift.Change"
}

// Marshal converts the message to a string.
func (msg *IRERiftChange) Marshal() string {
	return gmcp.Marshal(msg)
}

// Unmarshal populates the message with data.
func (msg *IRERiftChange) Unmarshal(data []byte) error {
	return gmcp.Unmarshal(data, msg)
}

// IRERiftList is a GMCP message lists items in rift storage.
type IRERiftList []IRERiftItem

// ID is the prefix before the message's data.
func (*IRERiftList) ID() string {
	return "IRE.Rift.List"
}

// Marshal converts the message to a string.
func (msg *IRERiftList) Marshal() string {
	return gmcp.Marshal(msg)
}

// Unmarshal populates the message with data.
func (msg *IRERiftList) Unmarshal(data []byte) error {
	return gmcp.Unmarshal(data, msg)
}

// IRERiftRequest is a GMCP message to request a list of items in the player's
// inventory.
type IRERiftRequest struct{}

// ID is the prefix before the message's data.
func (*IRERiftRequest) ID() string {
	return "IRE.Rift.Request"
}

// Marshal converts the message to a string.
func (*IRERiftRequest) Marshal() string {
	return "IRE.Rift.Request"
}

// Unmarshal populates the message with data.
func (*IRERiftRequest) Unmarshal(_ []byte) error {
	return nil
}
