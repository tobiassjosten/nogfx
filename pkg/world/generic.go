package world

import (
	"github.com/tobiassjosten/nogfx/pkg"
)

// GenericWorld is a generic implementation of the pkg.World interface.
type GenericWorld struct {
}

// NewGenericWorld creates a new Imperian-specific pkg.World.
func NewGenericWorld(_ pkg.Client, _ pkg.UI) pkg.World {
	return &GenericWorld{}
}

// ProcessInput processes player input.
func (world *GenericWorld) ProcessInput(input []byte) [][]byte {
	return [][]byte{input}
}

// ProcessOutput processes game output.
func (world *GenericWorld) ProcessOutput(output []byte) [][]byte {
	return [][]byte{output}
}

// ProcessCommand processes telnet commands.
func (world *GenericWorld) ProcessCommand(command []byte) error {
	return nil
}
