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

var serverMessages = map[string]ServerMessage{
	"Char.Name":    CharName{},
	"Char.Status":  CharStatus{},
	"Char.Vitals":  CharVitals{},
	"Core.Goodbye": CoreGoodbye{},
	"Core.Ping":    CorePing{},
}

// Parse converts a byte slice into a GMCP message.
func Parse(command []byte) (ServerMessage, error) {
	parts := bytes.SplitN(command, []byte{' '}, 2)

	message, ok := serverMessages[string(parts[0])]
	if !ok {
		return nil, fmt.Errorf("unknown message '%s'", parts[0])
	}

	if len(parts) == 1 {
		parts = append(parts, []byte{})
	}

	message, err := message.Hydrate(parts[1])
	if err != nil {
		return nil, fmt.Errorf(
			"failed hydrating %T (%s): %w",
			message, parts[1], err,
		)
	}

	return message, nil
}
