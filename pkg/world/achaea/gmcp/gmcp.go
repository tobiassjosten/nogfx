package gmcp

import (
	"bytes"
	"fmt"
)

// @todo Implement the full set:
// - https://nexus.ironrealms.com/GMCP
// - https://nexus.ironrealms.com/GMCP_Data
// - https://github.com/keneanung/GMCPAdditions

// Message is a GMCP message.
type Message interface {
	Hydrate([]byte) (Message, error)
	String() string
}

// Parse converts a byte slice into a GMCP message.
func Parse(command []byte) (Message, error) {
	parts := bytes.SplitN(command, []byte{' '}, 2)

	var hydrator Message

	switch string(parts[0]) {
	case "Char.Items.Inv":
		return CharItemsInv{}, nil

	case "Char.Name":
		hydrator = CharName{}

	case "Char.Status":
		hydrator = CharStatus{}

	case "Char.Vitals":
		hydrator = CharVitals{}

	default:
		return nil, fmt.Errorf("unknown message '%s'", parts[0])
	}

	if len(parts) == 1 {
		return nil, fmt.Errorf("missing '%T' data", hydrator)
	}

	return hydrator.Hydrate(parts[1])
}
