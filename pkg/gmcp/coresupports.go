package gmcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func marshalCoreSupports(m map[string]int) string {
	list := []string{}
	for module, version := range m {
		list = append(list, fmt.Sprintf("%s %d", module, version))
	}
	sort.Strings(list)

	data, _ := json.Marshal(list)

	return string(data)
}

func unmarshalCoreSupports(data []byte, msg Message) (map[string]int, error) {
	data = bytes.TrimSpace(bytes.TrimPrefix(data, []byte(msg.ID())))

	var list []string
	err := json.Unmarshal(data, &list)

	cs := map[string]int{}
	for _, item := range list {
		parts := strings.SplitN(item, " ", 2)

		version := 1
		if len(parts) == 2 {
			v, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf(
					"failed parsing module version: %w",
					err,
				)
			}
			version = v
		}

		cs[parts[0]] = version
	}

	return cs, err
}

// CoreSupportsSet is a client-sent GMCP message containing supported modules.
type CoreSupportsSet map[string]int

// ID is the prefix before the message's data.
func (msg *CoreSupportsSet) ID() string {
	return "Core.Supports.Set"
}

// Marshal converts the message to a string.
func (msg *CoreSupportsSet) Marshal() string {
	data := marshalCoreSupports(map[string]int(*msg))
	return fmt.Sprintf("%s %s", msg.ID(), string(data))
}

// Unmarshal populates the message with data.
func (msg *CoreSupportsSet) Unmarshal(data []byte) error {
	m, err := unmarshalCoreSupports(data, msg)
	*msg = CoreSupportsSet(m)
	return err
}

// CoreSupportsAdd is a client-sent GMCP message adding supported modules.
type CoreSupportsAdd map[string]int

// ID is the prefix before the message's data.
func (msg *CoreSupportsAdd) ID() string {
	return "Core.Supports.Add"
}

// Marshal converts the message to a string.
func (msg *CoreSupportsAdd) Marshal() string {
	data := marshalCoreSupports(map[string]int(*msg))
	return fmt.Sprintf("%s %s", msg.ID(), string(data))
}

// Unmarshal populates the message with data.
func (msg *CoreSupportsAdd) Unmarshal(data []byte) error {
	m, err := unmarshalCoreSupports(data, msg)
	*msg = CoreSupportsAdd(m)
	return err
}

// CoreSupportsRemove is a client-sent GMCP message removing supported modules.
type CoreSupportsRemove map[string]int

// ID is the prefix before the message's data.
func (msg *CoreSupportsRemove) ID() string {
	return "Core.Supports.Remove"
}

// Marshal converts the message to a string.
func (msg *CoreSupportsRemove) Marshal() string {
	data := marshalCoreSupports(map[string]int(*msg))
	return fmt.Sprintf("%s %s", msg.ID(), string(data))
}

// Unmarshal populates the message with data.
func (msg *CoreSupportsRemove) Unmarshal(data []byte) error {
	m, err := unmarshalCoreSupports(data, msg)
	*msg = CoreSupportsRemove(m)
	return err
}
