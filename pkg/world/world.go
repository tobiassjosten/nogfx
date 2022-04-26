package world

import (
	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/tui"
	"github.com/tobiassjosten/nogfx/pkg/world/achaea"
	"github.com/tobiassjosten/nogfx/pkg/world/imperian"
)

var worlds = map[string]func(*tui.TUI, pkg.Client) pkg.World{
	"achaea.com:23":  NewAchaeaWorld,
	"50.31.100.8:23": NewAchaeaWorld,

	"imperian.com:23":  NewImperianWorld,
	"67.202.121.44:23": NewImperianWorld,

	// @todo Extend this when we support more games. For now, we list these
	// two so as to force more general, shared functionality.
}

// New creates a World specific to the game being played.
func New(client pkg.Client, tui *tui.TUI, address string) pkg.World {
	// var world pkg.World = pkg.NewGenericWorld(tui, client)
	var world pkg.World
	if constructor, ok := worlds[address]; ok {
		world = constructor(tui, client)
	}

	return world
}

// NewAchaeaWorld wraps the specific constructor to create an interfaced World.
func NewAchaeaWorld(tui *tui.TUI, client pkg.Client) pkg.World {
	return achaea.NewWorld(tui, client)
}

// NewImperianWorld wraps the specific constructor to create an interfaced World.
func NewImperianWorld(tui *tui.TUI, client pkg.Client) pkg.World {
	return imperian.NewWorld(tui, client)
}

/*
if address == "example.com:23" {
	ui.AddVital("health", tui.HealthVital)
	ui.UpdateVital("health", 123, 234)
	ui.AddVital("mana", tui.ManaVital)
	ui.UpdateVital("mana", 100, 200)
	ui.AddVital("endurance", tui.EnduranceVital)
	ui.UpdateVital("endurance", 1000, 1200)
	ui.AddVital("willpower", tui.WillpowerVital)
	ui.UpdateVital("willpower", 1000, 2000)
}
*/

/*
type MockClient struct {
	reader io.Reader
	writer io.Writer
}

func (mock *MockClient) Read(p []byte) (int, error) {
	return mock.reader.Read(p)
}

func (mock MockClient) Write(p []byte) (int, error) {
	return mock.writer.Write(p)
}

func mockReadWriter() io.ReadWriter {
	return &MockClient{
		strings.NewReader("trololol\nqweqwrreqr\none two \033[33mthree \033[39mfour five six seven eight nine ten eleven twelve thirteen fourteen fifteen sixteen seventeen eighteen nineteen twenty twentyone twentytwo twentythree twentyfour twentyfive twentysix twentyseven twentyeight twentynine thirty thirtyone thirtytwo\nzxcxzvzxcxcxzc"),
		&strings.Builder{},
	}
}
*/
