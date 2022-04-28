package agmcp

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/tobiassjosten/nogfx/pkg/gmcp"
)

var serverMessages = map[string]gmcp.ServerMessage{
	"Char.Status": CharStatus{},
	"Char.Vitals": CharVitals{},
}

// Parse converts a byte slice into a GMCP message.
func Parse(command []byte) (gmcp.ServerMessage, error) {
	parts := bytes.SplitN(command, []byte{' '}, 2)

	message, ok := serverMessages[string(parts[0])]
	if !ok {
		return gmcp.Parse(command)
	}

	// Some messages don't have a message body but we want each message to
	// be responsible for its own hydration and validation. So we mock
	// missing bodies and proceed with hydration as normal.
	if len(parts) == 1 {
		parts = append(parts, []byte{})
	}

	msg, err := message.Hydrate(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed hydrating %T: %w", message, err)
	}

	return msg, nil
}

func splitRank(str string) (string, *int) {
	parts := strings.SplitN(str, "(", 2)
	name := strings.Trim(parts[0], " ")

	var rank *int
	if len(parts) > 1 {
		r, err := strconv.Atoi(strings.Trim(parts[1], "%)"))
		if err == nil {
			rank = &r
		}
	}

	return name, rank
}

func splitLevelRank(str string) (int, *int) {
	name, rank := splitRank(str)
	level, _ := strconv.Atoi(name)

	return level, rank
}
