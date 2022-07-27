package pkg

import (
	"github.com/tobiassjosten/nogfx/pkg/simpex"
)

type Match struct {
	Kind     IOKind
	Captures [][]byte
	Index    int
	// @todo Add the matched text on top of the Captures.
}

type Callback func([]Match, Inoutput) Inoutput

var NoopCallback Callback = func(_ []Match, inout Inoutput) Inoutput {
	return inout
}

type Trigger struct {
	Kind     IOKind
	Pattern  []byte
	Patterns [][]byte
	Callback Callback
}

func (t Trigger) Match(datas [][]byte, inout Inoutput) Inoutput {
	var matches []Match

	patterns := t.Patterns
	if t.Pattern != nil {
		patterns = append([][]byte{t.Pattern}, patterns...)
	}

	for i, data := range datas {
		for _, pattern := range patterns {
			captures := simpex.Match(pattern, data)
			if captures == nil {
				continue
			}

			matches = append(matches, Match{
				Kind:     t.Kind,
				Captures: captures,
				Index:    i,
			})
		}
	}

	if matches != nil {
		// @todo Make sure to recover from panics in the callback, eg
		// when users try accessing a non-existent Captures index
		inout = t.Callback(matches, inout)
	}

	return inout
}
