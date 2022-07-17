package module_test

import (
	"testing"

	"github.com/tobiassjosten/nogfx/pkg"
	amodule "github.com/tobiassjosten/nogfx/pkg/world/achaea/module"
	"github.com/tobiassjosten/nogfx/pkg/world/module"
)

func TestLearnMultipleLessons(t *testing.T) {
	// @todo Figure out how best to mock time, so we can properly test the
	// timeout functionality.

	tcs := map[string]module.TestCase{
		"learn 20": {
			Events: []module.TestEvent{
				{
					Type: module.Input,
					Data: []string{"learn 35 x from y"},
				},
				{
					Type: module.Output,
					Data: []string{
						"Y bows to you - the lesson in X is over.",
					},
				},
				{
					Type: module.Output,
					Data: []string{
						"Y bows to you - the lesson in X is over.",
					},
				},
				{
					Type: module.Output,
					Data: []string{
						"Y bows to you - the lesson in X is over.",
					},
				},
			},
			Inputs: []pkg.Input{
				{pkg.NewCommand([]byte("learn 15 x from y"))},
				{pkg.NewCommand([]byte("learn 15 x from y"))},
				{pkg.NewCommand([]byte("learn 5 x from y"))},
			},
			Outputs: []pkg.Output{
				{pkg.NewLine([]byte("Y bows to you - the lesson in X is over."))},
			},
		},

		"incomplete": {
			Events: []module.TestEvent{
				{
					Type: module.Input,
					Data: []string{"learn 20"},
				},
			},
		},

		"non-number": {
			Events: []module.TestEvent{
				{
					Type: module.Input,
					Data: []string{"learn z x from y"},
				},
			},
		},

		"uninitiated": {
			Events: []module.TestEvent{
				{
					Type: module.Output,
					Data: []string{"Y bows to you - the lesson in X is over."},
				},
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			tc.Eval(t, amodule.NewLearnMultipleLessons)
		})
	}
}
