package module_test

import (
	"testing"

	"github.com/tobiassjosten/nogfx/pkg"
	tst "github.com/tobiassjosten/nogfx/pkg/testing"
	amodule "github.com/tobiassjosten/nogfx/pkg/world/achaea/module"
)

func TestLearnMultipleLessons(t *testing.T) {
	// @todo Figure out how best to mock time, so we can properly test the
	// timeout functionality.

	tcs := map[string]tst.IOTestCase{
		"learn 20": {
			Events: []tst.IOEvent{
				tst.IOEIn("learn 35 x from y"),
				tst.IOEOut("Y bows to you - the lesson in X is over."),
				tst.IOEOut("Y bows to you - the lesson in X is over."),
				tst.IOEOut("Y bows to you - the lesson in X is over."),
			},
			Inoutputs: []pkg.Inoutput{
				tst.IOIn("learn 15 x from y"),
				tst.IO(
					"learn 15 x from y",
					"15 of 35 lessons learned, 0 seconds remaining.",
				),
				tst.IO(
					"learn 5 x from y",
					"30 of 35 lessons learned, 0 seconds remaining.",
				),
				tst.IOOut(
					"Y bows to you - the lesson in X is over.",
				).AddAfterOutput(0, []byte("35 of 35 lessons learned.")),
			},
		},

		"incomplete": {
			Events: []tst.IOEvent{
				tst.IOEIn("learn 20"),
			},
			Inoutputs: []pkg.Inoutput{
				tst.IOIn("learn 20"),
			},
		},

		"non-number": {
			Events: []tst.IOEvent{
				tst.IOEIn("learn z x from y"),
			},
			Inoutputs: []pkg.Inoutput{
				tst.IOIn("learn z x from y"),
			},
		},

		"uninitiated": {
			Events: []tst.IOEvent{
				tst.IOEOut("Y bows to you - the lesson in X is over."),
			},
			Inoutputs: []pkg.Inoutput{
				tst.IOOut("Y bows to you - the lesson in X is over."),
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			mod := amodule.NewLearnMultipleLessons()
			tc.Eval(t, mod)
		})
	}
}
