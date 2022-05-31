package ironrealms

import (
	"fmt"
	"strings"

	"github.com/tobiassjosten/nogfx/pkg/gmcp"
)

func msger(msg gmcp.Message) func() gmcp.Message {
	return func() gmcp.Message { return msg }
}

var messages = map[string]func() gmcp.Message{
	(&IRERiftChange{}).ID():  msger(&IRERiftChange{}),
	(&IRERiftList{}).ID():    msger(&IRERiftList{}),
	(&IRERiftRequest{}).ID(): msger(&IRERiftRequest{}),

	(&IRETargetSet{}).ID():  msger(&IRETargetSet{}),
	(&IRETargetInfo{}).ID(): msger(&IRETargetInfo{}),
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
