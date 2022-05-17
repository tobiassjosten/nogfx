package module_test

import (
	"testing"

	"github.com/tobiassjosten/nogfx/pkg/mock"
	amodule "github.com/tobiassjosten/nogfx/pkg/world/achaea/module"
	"github.com/tobiassjosten/nogfx/pkg/world/module"
)

func TestLearnMultipleLessons(t *testing.T) {
	// @todo Figure out how best to mock time, so we can properly test the
	// timeout functionality.

	tcs := map[string]module.TestCase{
		"learn 20": {
			Events: []module.TestEvent{
				module.NewTestEvent(true, []byte(
					"learn 35 x from y",
				)),
				module.NewTestEvent(false, []byte(
					"Y bows to you - the lesson in X is over.",
				)),
				module.NewTestEvent(false, []byte(
					"Y bows to you - the lesson in X is over.",
				)),
				module.NewTestEvent(false, []byte(
					"Y bows to you - the lesson in X is over.",
				)),
			},
			Inputs: [][]byte{
				[]byte("learn 15 x from y"),
			},
			Outputs: [][]byte{
				[]byte("Y bows to you - the lesson in X is over. [15/35]"),
				[]byte("Y bows to you - the lesson in X is over. [30/35]"),
				[]byte("Y bows to you - the lesson in X is over. [35/35]"),
			},
			Sent: [][]byte{
				[]byte("learn 15 x from y"),
				[]byte("learn 5 x from y"),
			},
		},

		"incomplete": {
			Events: []module.TestEvent{
				module.NewTestEvent(true, []byte(
					"learn 20",
				)),
			},
		},

		"non-number": {
			Events: []module.TestEvent{
				module.NewTestEvent(true, []byte(
					"learn z x from y",
				)),
			},
		},

		"uninitiated": {
			Events: []module.TestEvent{
				module.NewTestEvent(false, []byte(
					"Y bows to you - the lesson in X is over.",
				)),
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			client := &mock.ClientMock{
				SendFunc: func(data []byte) {},
			}
			ui := &mock.UIMock{}

			mod := amodule.NewLearnMultipleLessons(client, ui)

			module.Test(t, mod, tc, client)
		})
	}
}
