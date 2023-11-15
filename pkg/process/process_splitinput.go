package process

import (
	"bytes"
)

// SplitInputProcessor splits Input at the given separator to create a list of
// individual inputs from a compound one.
// E.g. "one;two;three" -> "one" "two" "three".
func SplitInputProcessor(sep []byte) Processor {
	return func(ins, outs [][]byte) ([][]byte, [][]byte, error) {
		postins := [][]byte{}

		for _, in := range ins {
			for _, bs := range bytes.Split(in, sep) {
				postins = append(postins, bytes.TrimSpace(bs))
			}
		}

		return postins, nil, nil
	}
}
