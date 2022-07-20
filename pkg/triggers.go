package pkg

import (
	"github.com/tobiassjosten/nogfx/pkg/simpex"
)

type Match struct {
	Kind     IOKind
	Captures [][]byte
	Index    int
}

type Callback func([]Match, Inoutput) Inoutput

var NoopCallback Callback = func(_ []Match, inout Inoutput) Inoutput {
	return inout
}

type Trigger struct {
	Kind     IOKind
	Pattern  []byte
	Callback Callback
}

func (t Trigger) Match(datas [][]byte, inout Inoutput) Inoutput {
	var matches []Match

	for i, data := range datas {
		captures := simpex.Match(t.Pattern, data)
		if captures == nil {
			continue
		}

		matches = append(matches, Match{
			Kind:     t.Kind,
			Captures: captures,
			Index:    i,
		})
	}

	if matches != nil {
		// @todo Make sure to recover from panics in the callback, eg
		// when users try accessing a non-existent Captures index
		inout = t.Callback(matches, inout)
	}

	return inout
}
