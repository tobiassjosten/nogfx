package module_test

import (
	"testing"

	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/world/module"
)

func TestRepeatInput(t *testing.T) {
	tcs := map[string]module.TestCase{
		"three asdf": {
			Events: []module.TestEvent{
				{module.Input, []string{"3 asdf"}},
			},
			Inputs: []pkg.Input{
				{
					pkg.NewCommand([]byte("asdf")),
					pkg.NewCommand([]byte("asdf")),
					pkg.NewCommand([]byte("asdf")),
				},
			},
		},

		"three empty": {
			Events: []module.TestEvent{
				{module.Input, []string{"3"}},
			},
			Inputs: []pkg.Input{{pkg.NewCommand([]byte("3"))}},
		},

		"non-number": {
			Events: []module.TestEvent{
				{module.Input, []string{"x asdf"}},
			},
			Inputs: []pkg.Input{{pkg.NewCommand([]byte("x asdf"))}},
		},

		"no output processing": {
			Events: []module.TestEvent{
				{module.Output, []string{"asdf"}},
			},
			Outputs: []pkg.Output{{pkg.NewLine([]byte("asdf"))}},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			tc.Eval(t, module.NewRepeatInput)
		})
	}
}
