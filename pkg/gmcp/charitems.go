package gmcp

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// CharItemsContents is a client-sent GMCP message to request a list of items
// located inside another item.
type CharItemsContents struct {
	Container int
}

// ID is the prefix before the message's data.
func (msg *CharItemsContents) ID() string {
	return "Char.Items.Contents"
}

// Marshal converts the message to a string.
func (msg *CharItemsContents) Marshal() string {
	return fmt.Sprintf("%s %d", msg.ID(), msg.Container)
}

// Unmarshal populates the message with data.
func (msg *CharItemsContents) Unmarshal(data []byte) error {
	data = bytes.TrimSpace(bytes.TrimPrefix(data, []byte(msg.ID())))

	var container int
	if err := json.Unmarshal(data, &container); err != nil {
		return err
	}

	msg.Container = container

	return nil
}

// CharItemsInv is a client-sent GMCP message to request a list of items in the
// player's inventory.
type CharItemsInv struct{}

// ID is the prefix before the message's data.
func (msg *CharItemsInv) ID() string {
	return "Char.Items.Inv"
}

// Marshal converts the message to a string.
func (msg *CharItemsInv) Marshal() string {
	return msg.ID()
}

// Unmarshal populates the message with data.
func (msg *CharItemsInv) Unmarshal(_ []byte) error {
	return nil
}

// CharItemsRoom is a client-sent GMCP message to request an updated list of
// items in the current room.
type CharItemsRoom struct{}

// ID is the prefix before the message's data.
func (msg *CharItemsRoom) ID() string {
	return "Char.Items.Room"
}

// Marshal converts the message to a string.
func (msg *CharItemsRoom) Marshal() string {
	return msg.ID()
}

// Unmarshal populates the message with data.
func (msg *CharItemsRoom) Unmarshal(_ []byte) error {
	return nil
}

// CharItem is an item within a player's inventory, the current room, or any
// other container.
type CharItem struct {
	ID         int                `json:"id"`
	Name       string             `json:"name"`
	Attributes CharItemAttributes `json:"attrib"`
	Icon       string             `json:"icon"`
}

// CharItemAttributes is a set of flags denoting how to interact with an item.
type CharItemAttributes struct {
	Container    bool
	Dangerous    bool
	Dead         bool
	Edible       bool
	Fluid        bool
	Groupable    bool
	Monster      bool
	Riftable     bool
	Takeable     bool
	Wearable     bool
	WieldedLeft  bool
	WieldedRight bool
	Worn         bool
}

// MarshalJSON transforms CharItemAttributes to a string.
func (as *CharItemAttributes) MarshalJSON() ([]byte, error) {
	var attribs string

	if as.Container {
		attribs += "c"
	}

	if as.Dead {
		attribs += "d"
	}

	if as.Edible {
		attribs += "e"
	}

	if as.Fluid {
		attribs += "f"
	}

	if as.Groupable {
		attribs += "g"
	}

	if as.WieldedLeft {
		attribs += "l"
	}

	if as.WieldedRight {
		attribs += "L"
	}

	if as.Monster {
		attribs += "m"
	}

	if as.Riftable {
		attribs += "r"
	}

	if as.Takeable {
		attribs += "t"
	}

	if as.Worn {
		attribs += "w"
	}

	if as.Wearable {
		attribs += "W"
	}

	if as.Dangerous {
		attribs += "x"
	}

	return []byte(`"` + attribs + `"`), nil
}

// UnmarshalJSON hydrates CharItemAttributes from a string.
func (as *CharItemAttributes) UnmarshalJSON(data []byte) error {
	for _, char := range bytes.Trim(data, `"`) {
		switch char {
		case 'c':
			as.Container = true

		case 'd':
			as.Dead = true

		case 'e':
			as.Edible = true

		case 'f':
			as.Fluid = true

		case 'g':
			as.Groupable = true

		case 'l':
			as.WieldedLeft = true

		case 'L':
			as.WieldedRight = true

		case 'm':
			as.Monster = true

		case 'r':
			as.Riftable = true

		case 't':
			as.Takeable = true

		case 'w':
			as.Worn = true

		case 'W':
			as.Wearable = true

		case 'x':
			as.Dangerous = true

		default:
			return fmt.Errorf("unknown attribute '%s'", string(char))
		}
	}

	return nil
}

// CharItemsList is a server-sent GMCP message listing items at the specified
// location.
type CharItemsList struct {
	Location string     `json:"location"`
	Items    []CharItem `json:"items"`
}

// ID is the prefix before the message's data.
func (msg *CharItemsList) ID() string {
	return "Char.Items.List"
}

// Marshal converts the message to a string.
func (msg *CharItemsList) Marshal() string {
	proxy := struct {
		*CharItemsList
		Items []CharItem `json:"items"`
	}{
		CharItemsList: msg,
	}

	proxy.Items = msg.Items
	if msg.Items == nil {
		proxy.Items = []CharItem{}
	}

	return Marshal(proxy)
}

// Unmarshal populates the message with data.
func (msg *CharItemsList) Unmarshal(data []byte) error {
	return Unmarshal(data, msg)
}

// CharItemsAdd is a server-sent GMCP message informing the client about an
// item being added to the specified location.
type CharItemsAdd struct {
	Location string   `json:"location"`
	Item     CharItem `json:"item"`
}

// ID is the prefix before the message's data.
func (msg *CharItemsAdd) ID() string {
	return "Char.Items.Add"
}

// Marshal converts the message to a string.
func (msg *CharItemsAdd) Marshal() string {
	return Marshal(msg)
}

// Unmarshal populates the message with data.
func (msg *CharItemsAdd) Unmarshal(data []byte) error {
	return Unmarshal(data, msg)
}

// CharItemsRemove is a server-sent GMCP message informing the client about an
// item being removed from the specified location.
type CharItemsRemove struct {
	Location string   `json:"location"`
	Item     CharItem `json:"item"`
}

// ID is the prefix before the message's data.
func (msg *CharItemsRemove) ID() string {
	return "Char.Items.Remove"
}

// Marshal converts the message to a string.
func (msg *CharItemsRemove) Marshal() string {
	return Marshal(msg)
}

// Unmarshal populates the message with data.
func (msg *CharItemsRemove) Unmarshal(data []byte) error {
	return Unmarshal(data, msg)
}

// CharItemsUpdate is a server-sent GMCP message informing the client about an
// item being removed from the specified location.
type CharItemsUpdate struct {
	Location string   `json:"location"`
	Item     CharItem `json:"item"`
}

// ID is the prefix before the message's data.
func (msg *CharItemsUpdate) ID() string {
	return "Char.Items.Update"
}

// Marshal converts the message to a string.
func (msg *CharItemsUpdate) Marshal() string {
	return Marshal(msg)
}

// Unmarshal populates the message with data.
func (msg *CharItemsUpdate) Unmarshal(data []byte) error {
	return Unmarshal(data, msg)
}
