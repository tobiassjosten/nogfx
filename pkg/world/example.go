package world

import (
	"github.com/tobiassjosten/nogfx/pkg"
)

// ExampleWorld is a mock implementation of the pkg.World interface.
type ExampleWorld struct {
	ui pkg.UI
}

// NewExampleWorld creates a new Imperian-specific pkg.World.
func NewExampleWorld(_ pkg.Client, ui pkg.UI) pkg.World {
	return &ExampleWorld{
		ui: ui,
	}
}

// ProcessInput processes player input.
func (world *ExampleWorld) ProcessInput(input []byte) []byte {
	world.ui.Print(append([]byte("> "), input...))
	return input
}

// ProcessOutput processes game output.
func (world *ExampleWorld) ProcessOutput(output []byte) []byte {
	return output
}

// ProcessCommand processes telnet commands.
func (world *ExampleWorld) ProcessCommand(command []byte) error {
	return nil
}

// ProcessGMCP processes GMCP messages.
func (world *ExampleWorld) ProcessGMCP(command []byte) error {
	return nil
}
