package gmcp

import (
	"bytes"
	"fmt"
)

// ClientMessage is a GMCP message sent from the client.
type ClientMessage interface {
	String() string
}

// ServerMessage is a GMCP message sent from the server.
type ServerMessage interface {
	Hydrate([]byte) (ServerMessage, error)
}

// ClientServerMessage is a GMCP message sent from both server and client.
type ClientServerMessage interface {
	ClientMessage
	ServerMessage
}

// Parse converts a byte slice into a GMCP message.
func Parse(command []byte) (ServerMessage, error) {
	parts := bytes.SplitN(command, []byte{' '}, 2)

	var hydrator ServerMessage

	switch string(parts[0]) {
	case "Core.Goodbye":
		return CoreGoodbye{}, nil

	case "Core.Ping":
		return CorePing{}, nil

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
