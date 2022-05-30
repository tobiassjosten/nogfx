package ironrealms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/tobiassjosten/nogfx/pkg/gmcp"

	"github.com/icza/gox/gox"
)

// CoreSupports is a list of potentially supported modules.
type CoreSupports struct {
	*gmcp.CoreSupports
	IRERift   *int
	IRETarget *int
}

// List transforms the struct to a list of strings.
func (cs *CoreSupports) List() []string {
	list := cs.CoreSupports.List()
	if cs.IRERift != nil {
		list = append(list, fmt.Sprintf("IRE.Rift %d", *cs.IRERift))
	}
	if cs.IRETarget != nil {
		list = append(list, fmt.Sprintf("IRE.Target %d", *cs.IRETarget))
	}

	return list
}

// UnmarshalJSON hydrates the struct from a string.
func (cs *CoreSupports) UnmarshalJSON(data []byte) error {
	if cs.CoreSupports == nil {
		cs.CoreSupports = &gmcp.CoreSupports{}
	}

	err := json.Unmarshal(data, cs.CoreSupports)
	if err != nil {
		return err
	}

	var list []string

	err = json.Unmarshal(data, &list)
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
		case "IRE.Rift":
			cs.IRERift = gox.NewInt(version)

		case "IRE.Target":
			cs.IRETarget = gox.NewInt(version)
		}
	}

	return nil
}

// CoreSupportsSet is a client-sent GMCP message containing supported modules.
type CoreSupportsSet struct {
	*CoreSupports
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
	*CoreSupports
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
	*CoreSupports
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
