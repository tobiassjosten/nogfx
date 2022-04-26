package imperian

import (
	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/tui"
)

// World is an Imperian-specific implementation of the pkg.World interface.
type World struct {
	ui     pkg.UI
	client pkg.Client
}

// NewWorld creates a new Imperian-specific pkg.World.
func NewWorld(ui pkg.UI, client pkg.Client) *World {
	ui.AddVital("health", tui.HealthVital)
	ui.AddVital("mana", tui.ManaVital)
	ui.AddVital("endurance", tui.EnduranceVital)
	ui.AddVital("willpower", tui.WillpowerVital)

	return &World{
		ui:     ui,
		client: client,
	}
}

// Input processes player input.
func (world *World) Input(input []byte) []byte {
	return input
}

// Output processes game output.
func (world *World) Output(output []byte) []byte {
	return output
}

// Command processes telnet commands.
func (world *World) Command(command []byte) error {
	return nil
}
