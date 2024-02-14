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

				tst.IOEOut("Y begins the lesson in X."),
				tst.IOEOut("Y continues your training in X."),
				tst.IOEOut("Y finishes the lesson in X."),

				tst.IOEOut("Y begins the lesson in X."),
				tst.IOEOut("Y continues your training in X."),
				tst.IOEOut("Y finishes the lesson in X."),

				tst.IOEOut("Y begins the lesson in X."),
				tst.IOEOut("Y continues your training in X."),
				tst.IOEOut("Y finishes the lesson in X."),
			},
			Inoutputs: []pkg.Inoutput{
				tst.IOIn("learn 15 x from y"),

				tst.IOOut("Y begins the lesson in X."),
				tst.IOOut("Y continues your training in X.").OmitOutput(0),
				tst.IO(
					"learn 15 x from y",
					"15 of 35 lessons learned, 0 seconds remaining.",
				),

				tst.IOOut("Y begins the lesson in X.").OmitOutput(0),
				tst.IOOut("Y continues your training in X.").OmitOutput(0),
				tst.IO(
					"learn 5 x from y",
					"30 of 35 lessons learned, 0 seconds remaining.",
				),

				tst.IOOut("Y begins the lesson in X.").OmitOutput(0),
				tst.IOOut("Y continues your training in X.").OmitOutput(0),
				tst.IOOut(
					"Y finishes the lesson in X.",
				).AddAfterOutput(0, []byte("35 of 35 lessons learned.")),
			},
		},

		"uninitiated": {
			Events: []tst.IOEvent{
				tst.IOEOut("Y begins the lesson in X."),
				tst.IOEOut("Y continues your training in X."),
				tst.IOEOut("Y finishes the lesson in X."),
			},
			Inoutputs: []pkg.Inoutput{
				tst.IOOut("Y begins the lesson in X."),
				tst.IOOut("Y continues your training in X."),
				tst.IOOut("Y finishes the lesson in X."),
			},
		},

		"unnecessary": {
			Events: []tst.IOEvent{
				tst.IOEIn("learn 15 x from y"),
				tst.IOEOut("Y begins the lesson in X."),
				tst.IOEOut("Y continues your training in X."),
				tst.IOEOut("Y finishes the lesson in X."),
			},
			Inoutputs: []pkg.Inoutput{
				tst.IOIn("learn 15 x from y"),
				tst.IOOut("Y begins the lesson in X."),
				tst.IOOut("Y continues your training in X."),
				tst.IOOut("Y finishes the lesson in X."),
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
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			mod := amodule.NewLearnMultipleLessons()
			tc.Eval(t, mod)
		})
	}
}
