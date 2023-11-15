package process_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tobiassjosten/nogfx/pkg/process"
)

func TestRepeatInputProcessor(t *testing.T) {
	tcs := map[string]struct {
		preins  [][]byte
		postins [][]byte
	}{
		"empty": {
			preins:  [][]byte{},
			postins: [][]byte{},
		},

		"repeat": {
			preins: [][]byte{
				[]byte("2 asdf"),
			},
			postins: [][]byte{
				[]byte("asdf"),
				[]byte("asdf"),
			},
		},

		"non-numeric": {
			preins: [][]byte{
				[]byte("x asdf"),
			},
			postins: [][]byte{
				[]byte("x asdf"),
			},
		},

		"zero": {
			preins: [][]byte{
				[]byte("0 asdf"),
			},
			postins: [][]byte{},
		},

		"straggler": {
			preins: [][]byte{
				[]byte("2 asdf"),
				[]byte("qwer"),
			},
			postins: [][]byte{
				[]byte("asdf"),
				[]byte("asdf"),
				[]byte("qwer"),
			},
		},

		"mixed": {
			preins: [][]byte{
				[]byte("qwer"),
				[]byte("3 asdf"),
				[]byte("x zxcv"),
				[]byte("2 fdsa"),
				[]byte("rewq"),
				[]byte("0 vcxz"),
			},
			postins: [][]byte{
				[]byte("qwer"),
				[]byte("asdf"),
				[]byte("asdf"),
				[]byte("asdf"),
				[]byte("x zxcv"),
				[]byte("fdsa"),
				[]byte("fdsa"),
				[]byte("rewq"),
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			proc := process.RepeatInputProcessor()

			ins, outs, _ := proc(tc.preins, nil)
			assert.Equal(t, tc.postins, ins)
			assert.Empty(t, outs)
		})
	}
}
