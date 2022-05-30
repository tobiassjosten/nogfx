package gmcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/icza/gox/gox"
)

// CoreSupports is a list of potentially supported modules.
type CoreSupports struct {
	Char        *int
	CharSkills  *int
	CharItems   *int
	CommChannel *int
	Room        *int
}

// List transforms the struct to a list of strings.
func (cs CoreSupports) List() []string {
	list := []string{}
	if cs.Char != nil {
		list = append(list, fmt.Sprintf("Char %d", *cs.Char))
	}
	if cs.CharSkills != nil {
		list = append(list, fmt.Sprintf("Char.Skills %d", *cs.CharSkills))
	}
	if cs.CharItems != nil {
		list = append(list, fmt.Sprintf("Char.Items %d", *cs.CharItems))
	}
	if cs.CommChannel != nil {
		list = append(list, fmt.Sprintf("Comm.Channel %d", *cs.CommChannel))
	}
	if cs.Room != nil {
		list = append(list, fmt.Sprintf("Room %d", *cs.Room))
	}

	return list
}

// Unlist hydrates the struct from a list of strings.
func (cs *CoreSupports) UnmarshalJSON(data []byte) error {
	var list []string

	err := json.Unmarshal(data, &list)
	if err != nil {
		return err
	}

	for _, module := range list {
		parts := strings.SplitN(module, " ", 2)
		if len(parts) < 2 {
			// With Core.Supports.Remove the module version isn't
			// mandatory, so we use a default to support that.
			parts = append(parts, "1")
		}

		version, err := strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("failed parsing module version: %w", err)
		}

		switch parts[0] {
		case "Char":
			cs.Char = gox.NewInt(version)

		case "Char.Skills":
			cs.CharSkills = gox.NewInt(version)

		case "Char.Items":
			cs.CharItems = gox.NewInt(version)

		case "Comm.Channel":
			cs.CommChannel = gox.NewInt(version)

		case "Room":
			cs.Room = gox.NewInt(version)
		}
	}

	return nil
}

// CoreSupportsSet is a client-sent GMCP message containing supported modules.
type CoreSupportsSet struct {
	CoreSupports
}

// ID is the prefix before the message's data.
func (msg *CoreSupportsSet) ID() string {
	return "Core.Supports.Set"
}

// Marshal converts the message to a string.
func (msg *CoreSupportsSet) Marshal() string {
	data, _ := json.Marshal(msg.CoreSupports.List())
	return fmt.Sprintf("%s %s", msg.ID(), string(data))
}

// Unmarshal populates the message with data.
func (msg *CoreSupportsSet) Unmarshal(data []byte) error {
	data = bytes.TrimSpace(bytes.TrimPrefix(data, []byte(msg.ID())))
	return json.Unmarshal(data, &msg.CoreSupports)
}

// CoreSupportsAdd is a client-sent GMCP message adding supported modules.
type CoreSupportsAdd struct {
	CoreSupports
}

// ID is the prefix before the message's data.
func (msg *CoreSupportsAdd) ID() string {
	return "Core.Supports.Add"
}

// Marshal converts the message to a string.
func (msg *CoreSupportsAdd) Marshal() string {
	data, _ := json.Marshal(msg.CoreSupports.List())
	return fmt.Sprintf("%s %s", msg.ID(), string(data))
}

// Unmarshal populates the message with data.
func (msg *CoreSupportsAdd) Unmarshal(data []byte) error {
	data = bytes.TrimSpace(bytes.TrimPrefix(data, []byte(msg.ID())))
	return json.Unmarshal(data, &msg.CoreSupports)
}

// CoreSupportsRemove is a client-sent GMCP message removing supported modules.
type CoreSupportsRemove struct {
	CoreSupports
}

// ID is the prefix before the message's data.
func (msg *CoreSupportsRemove) ID() string {
	return "Core.Supports.Remove"
}

// Marshal converts the message to a string.
func (msg *CoreSupportsRemove) Marshal() string {
	data, _ := json.Marshal(msg.CoreSupports.List())
	return fmt.Sprintf("%s %s", msg.ID(), string(data))
}

// Unmarshal populates the message with data.
func (msg *CoreSupportsRemove) Unmarshal(data []byte) error {
	data = bytes.TrimSpace(bytes.TrimPrefix(data, []byte(msg.ID())))
	return json.Unmarshal(data, &msg.CoreSupports)
}
