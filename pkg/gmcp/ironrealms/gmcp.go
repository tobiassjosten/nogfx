package ironrealms

import (
	"fmt"
	"strings"

	"github.com/tobiassjosten/nogfx/pkg/gmcp"
)

var messages = map[string]func() gmcp.Message{
	(&IRERiftChange{}).ID():  func() gmcp.Message { return &IRERiftChange{} },
	(&IRERiftList{}).ID():    func() gmcp.Message { return &IRERiftList{} },
	(&IRERiftRequest{}).ID(): func() gmcp.Message { return &IRERiftRequest{} },

	(&IRETargetSet{}).ID():  func() gmcp.Message { return &IRETargetSet{} },
	(&IRETargetInfo{}).ID(): func() gmcp.Message { return &IRETargetInfo{} },
}

// Parse converts a byte slice into a GMCP message.
func Parse(data []byte) (gmcp.Message, error) {
	parts := strings.SplitN(string(data), " ", 2)

	if _, ok := messages[parts[0]]; !ok {
		return gmcp.Parse(data)
	}

	msg := messages[parts[0]]()

	if err := msg.Unmarshal(data); err != nil {
		return nil, fmt.Errorf("couldn't unmarshal %T: %w", msg, err)
	}

	return msg, nil
}
