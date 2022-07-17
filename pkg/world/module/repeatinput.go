package module

import (
	"log"
	"strconv"

	"github.com/tobiassjosten/nogfx/pkg"
)

// RepeatInput is a module that lets players repeat commands by inputting them
// in the format of `3 command`, to send "command" thrice.
type RepeatInput struct {
	world pkg.World
}

// NewRepeatInput creates a new RepeatInput module.
func NewRepeatInput(world pkg.World) pkg.Module {
	log.Println("instantiating repeat input")
	return &RepeatInput{
		world: world,
	}
}

func (mod RepeatInput) InputTriggers() []pkg.Trigger[pkg.Input] {
	return []pkg.Trigger[pkg.Input]{
		{
			Pattern:  []byte("{^} {*}"),
			Callback: mod.onRepeat,
		},
	}
}

func (mod RepeatInput) OutputTriggers() []pkg.Trigger[pkg.Output] {
	return []pkg.Trigger[pkg.Output]{}
}

func (mod *RepeatInput) onRepeat(match pkg.TriggerMatch[pkg.Input]) pkg.Input {
	input := match.Content

	number, err := strconv.Atoi(string(match.Captures[0]))
	if err != nil {
		return input
	}

	input = input.Replace(match.Index, match.Captures[1])
	for i := 0; i < number-1; i++ {
		input = input.Insert(match.Index, match.Captures[1])
	}

	return input
}
