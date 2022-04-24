package achaea

import (
	"bytes"
	"fmt"

	"github.com/icza/gox/gox"
	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/telnet"
	"github.com/tobiassjosten/nogfx/pkg/tui"
	"github.com/tobiassjosten/nogfx/pkg/world/achaea/gmcp"
)

// World is an Achaea-specific implementation of the pkg.World interface.
type World struct {
	ui     pkg.UI
	client pkg.Client

	character Character
}

// NewWorld creates a new Achaea-specific pkg.World.
func NewWorld(ui pkg.UI, client pkg.Client) *World {
	ui.AddVital("health", tui.HealthVital)
	ui.AddVital("mana", tui.ManaVital)
	ui.AddVital("endurance", tui.EnduranceVital)
	ui.AddVital("willpower", tui.WillpowerVital)

	return &World{
		ui:     ui,
		client: client,

		character: Character{},
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
	willEcho := []byte{telnet.IAC, telnet.WILL, telnet.ECHO}
	wontEcho := []byte{telnet.IAC, telnet.WONT, telnet.ECHO}

	willGMCP := []byte{telnet.IAC, telnet.WILL, telnet.GMCP}
	prefixGMCP := []byte{telnet.IAC, telnet.SB, telnet.GMCP}
	suffixGMCP := []byte{telnet.IAC, telnet.SE}

	switch {
	case bytes.Equal(command, willEcho):
		world.ui.MaskInput()

	case bytes.Equal(command, wontEcho):
		world.ui.UnmaskInput()

	case bytes.Equal(command, willGMCP):
		err := world.gmcp(gmcp.CoreHello{
			Client:  "nogfx",
			Version: pkg.Version,
		})
		if err != nil {
			return fmt.Errorf("failed GMCP: %w", err)
		}

		err = world.gmcp(gmcp.CoreSupportsSet{
			CoreSupports: gmcp.CoreSupports{
				Char:        gox.NewInt(1),
				CharSkills:  gox.NewInt(1),
				CharItems:   gox.NewInt(1),
				CommChannel: gox.NewInt(1),
				Room:        gox.NewInt(1),
				IRERift:     gox.NewInt(1),
			},
		})
		if err != nil {
			return fmt.Errorf("failed GMCP: %w", err)
		}

	case bytes.HasPrefix(command, prefixGMCP):
		data := command[len(prefixGMCP) : len(command)-len(suffixGMCP)]
		message, err := gmcp.Parse(data)
		if err != nil {
			world.ui.Print([]byte(fmt.Sprintf("[GMCP error: %s]", err)))
			return nil
		}

		switch msg := message.(type) {
		case gmcp.CharName:
			world.character.FromCharName(msg)

			// We have just logged in, so let's do an inventory.
			err := world.gmcp(gmcp.IRERiftRequest{})
			if err != nil {
				return fmt.Errorf("failed GMCP: %w", err)
			}

			err = world.gmcp(gmcp.CommChannelPlayers{})
			if err != nil {
				return fmt.Errorf("failed GMCP: %w", err)
			}

			err = world.gmcp(gmcp.CharItemsInv{})
			if err != nil {
				return fmt.Errorf("failed GMCP: %w", err)
			}

		case gmcp.CharVitals:
			world.character.FromCharVitals(msg)

			world.ui.UpdateVital("health",
				world.character.Health,
				world.character.MaxHealth,
			)
			world.ui.UpdateVital("mana",
				world.character.Mana,
				world.character.MaxMana,
			)
			world.ui.UpdateVital("endurance",
				world.character.Endurance,
				world.character.MaxEndurance,
			)
			world.ui.UpdateVital("willpower",
				world.character.Willpower,
				world.character.MaxWillpower,
			)
		}
	}

	return nil
}

func (world *World) gmcp(value gmcp.ClientMessage) error {
	_, err := world.client.Write(append(append(
		[]byte{telnet.IAC, telnet.SB, telnet.GMCP},
		[]byte(value.String())...,
	), telnet.IAC, telnet.SE))
	return err
}
