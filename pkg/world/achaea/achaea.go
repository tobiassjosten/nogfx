package achaea

import (
	"bytes"
	"fmt"

	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/gmcp"
	"github.com/tobiassjosten/nogfx/pkg/telnet"
	"github.com/tobiassjosten/nogfx/pkg/tui"
	"github.com/tobiassjosten/nogfx/pkg/world/achaea/agmcp"

	"github.com/icza/gox/gox"
)

// World is an Achaea-specific implementation of the pkg.World interface.
type World struct {
	client pkg.Client

	ui       pkg.UI
	uiVitals map[string]struct{}

	character *Character
}

// NewWorld creates a new Achaea-specific pkg.World.
func NewWorld(client pkg.Client, ui pkg.UI) pkg.World {
	return &World{
		client: client,

		ui:       ui,
		uiVitals: map[string]struct{}{},

		character: &Character{},
	}
}

// ProcessInput processes player input.
func (world *World) ProcessInput(input []byte) []byte {
	return input
}

// ProcessOutput processes game output.
func (world *World) ProcessOutput(output []byte) []byte {
	return output
}

// ProcessCommand processes telnet commands.
func (world *World) ProcessCommand(command []byte) error {
	if data := gmcp.Unwrap(command); data != nil {
		return world.ProcessGMCP(data)
	}

	switch {
	case bytes.Equal(command, []byte{telnet.IAC, telnet.WILL, telnet.GMCP}):
		err := world.SendGMCP(agmcp.CoreSupportsSet{
			CoreSupports: agmcp.CoreSupports{
				CoreSupports: gmcp.CoreSupports{
					Char:        gox.NewInt(1),
					CharSkills:  gox.NewInt(1),
					CharItems:   gox.NewInt(1),
					CommChannel: gox.NewInt(1),
					Room:        gox.NewInt(1),
				},
				IRERift: gox.NewInt(1),
			},
		})
		if err != nil {
			return fmt.Errorf("failed GMCP: %w", err)
		}
	}

	return nil
}

// ProcessGMCP processes GMCP messages.
func (world *World) ProcessGMCP(data []byte) error {
	message, err := agmcp.Parse(data)
	if err != nil {
		return fmt.Errorf("failed parsing GMCP: %w", err)
	}

	switch msg := message.(type) {
	case gmcp.CharName:
		world.character.FromCharName(msg)

		if err := world.SendGMCP(gmcp.CharItemsInv{}); err != nil {
			return fmt.Errorf("failed GMCP: %w", err)
		}

		if err := world.SendGMCP(gmcp.CommChannelPlayers{}); err != nil {
			return fmt.Errorf("failed GMCP: %w", err)
		}

		if err := world.SendGMCP(gmcp.IRERiftRequest{}); err != nil {
			return fmt.Errorf("failed GMCP: %w", err)
		}

	case agmcp.CharStatus:
		world.character.FromCharStatus(msg)

	case agmcp.CharVitals:
		world.character.FromCharVitals(msg)

		err := world.UpdateVitals()
		if err != nil {
			return fmt.Errorf("failed updating vitals: %w", err)
		}
	}

	return nil
}

// SendGMCP writes a GMCP message to the client.
func (world *World) SendGMCP(message gmcp.ClientMessage) error {
	data := []byte(message.String())
	if _, err := world.client.Write(gmcp.Wrap(data)); err != nil {
		return err
	}

	return nil
}

// UpdateVitals creates sends new current and max values to UI's VitalPanes.
func (world *World) UpdateVitals() error {
	order := []string{"health", "mana", "endurance", "willpower"}

	vitals := map[string]tui.Vital{
		"health":    tui.HealthVital,
		"mana":      tui.ManaVital,
		"endurance": tui.EnduranceVital,
		"willpower": tui.WillpowerVital,
	}

	values := map[string][]int{
		"health":    {world.character.Health, world.character.MaxHealth},
		"mana":      {world.character.Mana, world.character.MaxMana},
		"endurance": {world.character.Endurance, world.character.MaxEndurance},
		"willpower": {world.character.Willpower, world.character.MaxWillpower},
	}

	for _, name := range order {
		value, ok := values[name]
		if !ok || len(value) != 2 {
			return fmt.Errorf("invalid vital data for '%s'", name)
		}

		if _, ok := world.uiVitals[name]; !ok {
			world.ui.AddVital(name, vitals[name])
			world.uiVitals[name] = struct{}{}
		}

		world.ui.UpdateVital(name, value[0], value[1])
	}

	return nil
}
