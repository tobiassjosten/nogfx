package gmcp

import (
	"bytes"
	"fmt"
)

// @todo Implement the full set:
// - https://nexus.ironrealms.com/GMCP
// - https://nexus.ironrealms.com/GMCP_Data
// - https://github.com/keneanung/GMCPAdditions

// ClientMessage is a GMCP message sent from the client.
type ClientMessage interface {
	String() string
}

// ServerMessage is a GMCP message sent from the server.
type ServerMessage interface {
	Hydrate([]byte) (ServerMessage, error)
}

type ClientServerMessage interface {
	ClientMessage
	ServerMessage
}

// Parse converts a byte slice into a GMCP message.
func Parse(command []byte) (ServerMessage, error) {
	parts := bytes.SplitN(command, []byte{' '}, 2)

	var hydrator ServerMessage

	switch string(parts[0]) {
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
