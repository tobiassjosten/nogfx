package gmcp

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// CharAffliction is an affliction ailing a character.
type CharAffliction struct {
	Name        string `json:"name"`
	Cure        string `json:"cure"`
	Description string `json:"desc"`
}

// CharAfflictionsList is a GMCP message listing current afflictions.
type CharAfflictionsList []CharAffliction

// ID is the prefix before the message's data.
func (msg *CharAfflictionsList) ID() string {
	return "Char.Afflictions.List"
}

// Marshal converts the message to a string.
func (msg *CharAfflictionsList) Marshal() string {
	return Marshal(msg)
}

// Unmarshal populates the message with data.
func (msg *CharAfflictionsList) Unmarshal(data []byte) error {
	return Unmarshal(data, msg)
}

// CharAfflictionsAdd is a GMCP message detailing an additional affliction.
type CharAfflictionsAdd CharAffliction

// ID is the prefix before the message's data.
func (msg *CharAfflictionsAdd) ID() string {
	return "Char.Afflictions.Add"
}

// Marshal converts the message to a string.
func (msg *CharAfflictionsAdd) Marshal() string {
	return Marshal(msg)
}

// Unmarshal populates the message with data.
func (msg *CharAfflictionsAdd) Unmarshal(data []byte) error {
	return Unmarshal(data, msg)
}

// CharAfflictionsRemove is a GMCP message detailing a cured affliction.
type CharAfflictionsRemove []CharAffliction

// ID is the prefix before the message's data.
func (msg *CharAfflictionsRemove) ID() string {
	return "Char.Afflictions.Remove"
}

// Marshal converts the message to a string.
func (msg *CharAfflictionsRemove) Marshal() string {
	list := []string{}

	for _, affliction := range *msg {
		list = append(list, affliction.Name)
	}

	data, _ := json.Marshal(list)
	return fmt.Sprintf("%s %s", msg.ID(), string(data))
}

// Unmarshal populates the message with data.
func (msg *CharAfflictionsRemove) Unmarshal(data []byte) error {
	data = bytes.TrimPrefix(data, []byte(msg.ID()+" "))

	list := []string{}

	err := json.Unmarshal(data, &list)
	if err != nil {
		return err
	}

	for _, item := range list {
		*msg = append(*msg, CharAffliction{Name: item})
	}

	return nil
}
