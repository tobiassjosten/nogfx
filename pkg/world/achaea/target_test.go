package achaea_test

import (
	"testing"

	"github.com/tobiassjosten/nogfx/pkg/gmcp"
	"github.com/tobiassjosten/nogfx/pkg/mock"
	"github.com/tobiassjosten/nogfx/pkg/world/achaea"
)

func TestTargeting(t *testing.T) {
	tcs := map[string]struct {
		messages []gmcp.Message
	}{
		"my case": {
			messages: []gmcp.Message{},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			client := &mock.ClientMock{}
			ui := &mock.UIMock{}

			world := achaea.NewWorld(client, ui)

			for _, message := range tc.messages {
				data := []byte(message.Marshal())
				world.ProcessCommand(gmcp.Wrap(data))
			}
		})
	}
}
