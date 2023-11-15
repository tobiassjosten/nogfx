package process_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tobiassjosten/nogfx/pkg/process"
)

func TestSplitInputProcessor(t *testing.T) {
	tcs := map[string]struct {
		sep     []byte
		preins  [][]byte
		postins [][]byte
	}{
		"empty": {
			preins:  [][]byte{},
			postins: [][]byte{},
		},

		"split semicolon": {
			sep: []byte{';'},
			preins: [][]byte{
				[]byte("asdf;qwer"),
			},
			postins: [][]byte{
				[]byte("asdf"),
				[]byte("qwer"),
			},
		},

		"split pipes": {
			sep: []byte("||"),
			preins: [][]byte{
				[]byte("asdf||qwer"),
			},
			postins: [][]byte{
				[]byte("asdf"),
				[]byte("qwer"),
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			proc := process.SplitInputProcessor(tc.sep)

			ins, outs, _ := proc(tc.preins, nil)
			assert.Equal(t, tc.postins, ins)
			assert.Empty(t, outs)
		})
	}
}
