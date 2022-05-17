package module

import (
	"strconv"

	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/simpex"
)

var (
	modRIINumber = []byte("{^} {*}")
)

// RepeatInput is a module that lets players repeat commands by inputting them
// in the format of `3 command`, to send "command" thrice.
type RepeatInput struct {
	client pkg.Client
	ui     pkg.UI
}

// NewRepeatInput creates a new RepeatInput module.
func NewRepeatInput(client pkg.Client, ui pkg.UI) pkg.Module {
	return &RepeatInput{
		client: client,
		ui:     ui,
	}
}

// ProcessInput processes player input.
func (mod *RepeatInput) ProcessInput(input []byte) [][]byte {
	matches := simpex.Match(modRIINumber, input)
	if matches == nil {
		return [][]byte{}
	}

	number, err := strconv.Atoi(string(matches[0]))
	if err != nil {
		return [][]byte{}
	}

	var inputs [][]byte
	for i := 0; i < number; i++ {
		inputs = append(inputs, matches[1])
	}

	return inputs
}

// ProcessOutput processes server output.
func (mod *RepeatInput) ProcessOutput(output []byte) [][]byte {
	return [][]byte{}
}
