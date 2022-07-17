package pkg_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tobiassjosten/nogfx/pkg"
)

func TestInput(t *testing.T) {
	tcs := map[string]struct {
		data  []byte
		add   [][]byte
		split []byte
		input pkg.Input
	}{
		"instantiate": {
			data:  []byte("asdf"),
			input: pkg.Input{pkg.NewCommand([]byte("asdf"))},
		},

		"add once": {
			data: []byte("asdf"),
			add:  [][]byte{[]byte("qwer")},
			input: pkg.Input{
				pkg.NewCommand([]byte("asdf")),
				pkg.NewCommand([]byte("qwer")),
			},
		},

		"add twice": {
			data: []byte("asdf"),
			add: [][]byte{
				[]byte("qwer"),
				[]byte("zxcv"),
			},
			input: pkg.Input{
				pkg.NewCommand([]byte("asdf")),
				pkg.NewCommand([]byte("qwer")),
				pkg.NewCommand([]byte("zxcv")),
			},
		},

		"split": {
			data:  []byte("as;df;gh"),
			split: []byte{';'},
			input: pkg.Input{
				pkg.NewCommand([]byte("as")),
				pkg.NewCommand([]byte("df")),
				pkg.NewCommand([]byte("gh")),
			},
		},

		// @todo Make sure operations don't mutate the object (but only
		// return a new instance of it).
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			var input pkg.Input
			if tc.data != nil {
				input = pkg.NewInput(tc.data)
			}

			if tc.add != nil {
				for _, add := range tc.add {
					input = input.Add(add)
				}
			}

			if tc.split != nil {
				input = input.Split(tc.split)
			}

			assert.Equal(t, tc.input, input)
		})
	}
}
