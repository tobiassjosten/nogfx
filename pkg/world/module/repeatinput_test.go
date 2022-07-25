package module_test

import (
	"testing"

	"github.com/tobiassjosten/nogfx/pkg"
	tst "github.com/tobiassjosten/nogfx/pkg/testing"
	"github.com/tobiassjosten/nogfx/pkg/world/module"
)

func TestRepeatInput(t *testing.T) {
	tcs := map[string]tst.IOTestCase{
		"three asdf": {
			Events: []tst.IOEvent{tst.IOEIn("3 asdf")},
			Inoutputs: []pkg.Inoutput{
				tst.IOIn("asdf").
					AddAfterInput(0, []byte("asdf")).
					AddAfterInput(0, []byte("asdf")),
			},
		},

		"three empty": {
			Events:    []tst.IOEvent{tst.IOEIn("3")},
			Inoutputs: []pkg.Inoutput{tst.IOIn("3")},
		},

		"non-number": {
			Events:    []tst.IOEvent{tst.IOEIn("x asdf")},
			Inoutputs: []pkg.Inoutput{tst.IOIn("x asdf")},
		},

		"no output processing": {
			Events:    []tst.IOEvent{tst.IOEOut("asdf")},
			Inoutputs: []pkg.Inoutput{tst.IOOut("asdf")},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			mod := module.NewRepeatInput()
			tc.Eval(t, mod)
		})
	}
}
