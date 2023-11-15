package process

import (
	"fmt"

	"github.com/tobiassjosten/nogfx/pkg/simpex"
)

// Match is the result of a pattern successfully applied to lines of text. Each
// key matches an index within the lines of texts and its value contains
// potential captures from the pattern.
type Match map[int][][]byte

// Callback is a function to execute upon matching pattern(s).
type Callback func(match Match, ins, outs [][]byte) (postins, postouts [][]byte, err error)

// MatchInput creates a Processor that runs the given callback when the given
// pattern matches line(s) within processing inputs.
func MatchInput(pat string, cb Callback) Processor {
	return MatchInputs([]string{pat}, cb)
}

// MatchInputs creates a Processor that runs the given callback when the given
// patterns matches line(s) within processing inputs.
func MatchInputs(pats []string, cb Callback) Processor {
	// We accept patterns as strings only for convenience reasons.
	var bss [][]byte
	for _, pat := range pats {
		bss = append(bss, []byte(pat))
	}

	return func(ins, outs [][]byte) ([][]byte, [][]byte, error) {
		return match(bss, ins, ins, outs, cb)
	}
}

// MatchOutput creates a Processor that runs the given callback when the given
// pattern matches line(s) within processing outputs.
func MatchOutput(pat string, cb Callback) Processor {
	return MatchOutputs([]string{pat}, cb)
}

// MatchOutputs creates a Processor that runs the given callback when the given
// patterns matches line(s) within processing outputs.
func MatchOutputs(pats []string, cb Callback) Processor {
	// We accept pats as strings only for convenience reasons.
	var bss [][]byte
	for _, pat := range pats {
		bss = append(bss, []byte(pat))
	}

	return func(ins, outs [][]byte) ([][]byte, [][]byte, error) {
		return match(bss, outs, ins, outs, cb)
	}
}

func match(pats, txts, ins, outs [][]byte, cb Callback) (postins, postouts [][]byte, err error) {
	match := Match{}

	for i, text := range txts {
		for _, pat := range pats {
			captures := simpex.Match(pat, text)
			if captures == nil {
				continue
			}

			match[i] = captures

			break
		}
	}

	if len(match) == 0 {
		return ins, outs, nil
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("match callback failed: %s", r)
			if rerr, ok := r.(error); ok {
				err = fmt.Errorf("match callback failed: %w", rerr)
			}
		}
	}()

	return cb(match, ins, outs)
}
