package module

import (
	"strconv"

	"github.com/tobiassjosten/nogfx/pkg"
)

// RepeatInput is a module that lets players repeat commands by inputting them
// in the format of `3 command`, to send "command" thrice.
type RepeatInput struct {
}

// NewRepeatInput creates a new RepeatInput module.
func NewRepeatInput() pkg.Module {
	return &RepeatInput{}
}

func (mod RepeatInput) Triggers() []pkg.Trigger {
	return []pkg.Trigger{
		{
			Kind:     pkg.Input,
			Pattern:  []byte("{^} {*}"),
			Callback: mod.onRepeat,
		},
	}
}

func (mod *RepeatInput) onRepeat(matches []pkg.Match, inout pkg.Inoutput) pkg.Inoutput {
	for _, match := range matches {
		i := match.Index

		number, err := strconv.Atoi(string(match.Captures[0]))
		if err != nil {
			continue
		}

		inout.Input = inout.Input.Replace(i, match.Captures[1])
		for ii := 0; ii < number-1; ii++ {
			inout.Input = inout.Input.AddAfter(i, match.Captures[1])
		}
	}

	return inout
}
