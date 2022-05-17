package module_test

import (
	"testing"

	"github.com/tobiassjosten/nogfx/pkg/mock"
	"github.com/tobiassjosten/nogfx/pkg/world/module"
)

func TestRepeatInput(t *testing.T) {
	tcs := map[string]module.TestCase{
		"three asdf": {
			Events: []module.TestEvent{
				module.NewTestEvent(true, []byte("3 asdf")),
			},
			Inputs: [][]byte{
				[]byte("asdf"),
				[]byte("asdf"),
				[]byte("asdf"),
			},
		},

		"three empty": {
			Events: []module.TestEvent{
				module.NewTestEvent(true, []byte("3")),
			},
		},

		"non-number": {
			Events: []module.TestEvent{
				module.NewTestEvent(true, []byte("x asdf")),
			},
		},

		"no output processing": {
			Events: []module.TestEvent{
				module.NewTestEvent(false, []byte("asdf")),
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			client := &mock.ClientMock{
				SendFunc: func(data []byte) {},
			}
			ui := &mock.UIMock{}

			mod := module.NewRepeatInput(client, ui)

			tc.Eval(t, mod, client)
		})
	}
}
