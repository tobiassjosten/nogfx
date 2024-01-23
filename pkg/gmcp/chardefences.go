package gmcp

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// CharDefence is a defence protecting a character.
type CharDefence struct {
	Name        string `json:"name"`
	Cure        string `json:"cure"`
	Description string `json:"desc"`
}

// CharDefencesList is a GMCP message listing current character defenses.
type CharDefencesList []CharDefence

// ID is the prefix before the message's data.
func (*CharDefencesList) ID() string {
	return "Char.Defences.List"
}

// Marshal converts the message to a string.
func (msg *CharDefencesList) Marshal() string {
	return Marshal(msg)
}

// Unmarshal populates the message with data.
func (msg *CharDefencesList) Unmarshal(data []byte) error {
	return Unmarshal(data, msg)
}

// CharDefencesAdd is a GMCP message detailing an additional defense.
type CharDefencesAdd CharDefence

// ID is the prefix before the message's data.
func (*CharDefencesAdd) ID() string {
	return "Char.Defences.Add"
}

// Marshal converts the message to a string.
func (msg *CharDefencesAdd) Marshal() string {
	return Marshal(msg)
}

// Unmarshal populates the message with data.
func (msg *CharDefencesAdd) Unmarshal(data []byte) error {
	return Unmarshal(data, msg)
}

// CharDefencesRemove is a GMCP message detailing a lost defense.
type CharDefencesRemove []CharDefence

// ID is the prefix before the message's data.
func (*CharDefencesRemove) ID() string {
	return "Char.Defences.Remove"
}

// Marshal converts the message to a string.
func (msg *CharDefencesRemove) Marshal() string {
	list := []string{}

	for _, defence := range *msg {
		list = append(list, defence.Name)
	}

	data, _ := json.Marshal(list)

	return fmt.Sprintf("%s %s", msg.ID(), string(data))
}

// Unmarshal populates the message with data.
func (msg *CharDefencesRemove) Unmarshal(data []byte) error {
	data = bytes.TrimPrefix(data, []byte(msg.ID()+" "))

	list := []string{}

	err := json.Unmarshal(data, &list)
	if err != nil {
		return err
	}

	for _, item := range list {
		*msg = append(*msg, CharDefence{Name: item})
	}

	return nil
}
