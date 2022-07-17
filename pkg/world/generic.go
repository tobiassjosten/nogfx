package world

import (
	"github.com/tobiassjosten/nogfx/pkg"
)

// GenericWorld is a generic implementation of the pkg.World interface.
type GenericWorld struct {
	client pkg.Client
	ui     pkg.UI
}

// NewGenericWorld creates a new Imperian-specific pkg.World.
func NewGenericWorld(client pkg.Client, ui pkg.UI) pkg.World {
	return &GenericWorld{
		client: client,
		ui:     ui,
	}
}

// Print passes data onto the configured UI.
func (world *GenericWorld) Print(data []byte) {
	world.ui.Print(data)
}

// Send passes data onto the configured Client.
func (world *GenericWorld) Send(data []byte) {
	world.client.Send(data)
}

// ProcessInput processes player input.
func (world *GenericWorld) ProcessInput(input pkg.Input) pkg.Input {
	return input
}

// ProcessOutput processes game output.
func (world *GenericWorld) ProcessOutput(output pkg.Output) pkg.Output {
	return output
}

// ProcessCommand processes telnet commands.
func (world *GenericWorld) ProcessCommand(command []byte) {
}
