package imperian

import (
	"bytes"
	"fmt"

	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/gmcp"
	"github.com/tobiassjosten/nogfx/pkg/telnet"
	"github.com/tobiassjosten/nogfx/pkg/tui"
	"github.com/tobiassjosten/nogfx/pkg/world/imperian/igmcp"

	"github.com/icza/gox/gox"
)

// World is an Imperian-specific implementation of the pkg.World interface.
type World struct {
	client pkg.Client

	ui       pkg.UI
	uiVitals map[string]struct{}
}

// NewWorld creates a new Imperian-specific pkg.World.
func NewWorld(client pkg.Client, ui pkg.UI) pkg.World {
	return &World{
		client: client,

		ui:       ui,
		uiVitals: map[string]struct{}{},
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
		err := world.SendGMCP(igmcp.CoreSupportsSet{
			CoreSupports: igmcp.CoreSupports{
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

// ServerMessages maps GMCP messages to associated structs.
var ServerMessages = map[string]gmcp.ServerMessage{
	"Char.Status": igmcp.CharStatus{},
	"Char.Vitals": igmcp.CharVitals{},
}

// ProcessGMCP processes GMCP messages.
func (world *World) ProcessGMCP(data []byte) error {
	message, err := gmcp.Parse(data, ServerMessages)
	if err != nil {
		return fmt.Errorf("failed parsing GMCP: %w", err)
	}

	switch msg := message.(type) {
	case igmcp.CharVitals:
		err := world.UpdateVitals(msg)
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
func (world *World) UpdateVitals(msg igmcp.CharVitals) error {
	order := []string{"health", "mana"}

	vitals := map[string]*tui.Vital{
		"health": tui.NewHealthVital(),
		"mana":   tui.NewManaVital(),
	}

	values := map[string][]int{
		"health": {msg.HP / 11, msg.MaxHP / 11},
		"mana":   {msg.MP / 11, msg.MaxMP / 11},
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
