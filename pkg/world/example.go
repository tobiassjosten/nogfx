package world

import (
	"github.com/icza/gox/gox"
	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/navigation"
	"github.com/tobiassjosten/nogfx/pkg/tui"
)

// ExampleWorld is a mock implementation of the pkg.World interface.
type ExampleWorld struct {
	ui pkg.UI
}

// NewExampleWorld creates a new Imperian-specific pkg.World.
func NewExampleWorld(_ pkg.Client, ui pkg.UI) pkg.World {
	go func() {
		ui.Outputs() <- []byte("One lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum et nunc in dui efficitur commodo sed ut lectus. Etiam a urna nec augue gravida imperdiet. Aenean luctus ut augue in laoreet. Nunc ut dui sem. Maecenas id leo purus. Maecenas enim purus, finibus sit amet aliquam sit amet, commodo ut velit. Nunc nunc lectus, pulvinar ut metus quis, laoreet dapibus dolor. In rhoncus quis ligula sit amet pharetra. Aliquam nunc velit, pharetra nec imperdiet nec, iaculis nec ex. Vestibulum porta tristique dignissim. Pellentesque ac maximus lorem, ut viverra risus.")
		ui.Outputs() <- []byte("Two velit ligula, gravida at enim a, interdum egestas massa. Maecenas feugiat commodo velit, et tristique erat dictum ut. Phasellus vel pulvinar nunc, eu convallis sem. In elit eros, fringilla sed porta eu, pharetra at leo. Integer pretium, odio a venenatis elementum, enim elit porttitor purus, eget fermentum est ex eu est. Duis quis euismod velit, eleifend lacinia sapien. Cras non ullamcorper turpis. Pellentesque elementum dui risus, ut fringilla tortor cursus sed. Pellentesque ipsum dui, fringilla vel tellus in, imperdiet semper ligula.")
		ui.Outputs() <- []byte("Three tempus orci nibh, at consectetur nisi sagittis eu. Nam elementum, erat sed mattis malesuada, est neque elementum diam, blandit iaculis ipsum nunc eget diam. Nulla ut dolor sed elit faucibus mattis. Nam elementum urna vel interdum volutpat. Curabitur convallis imperdiet mi, at feugiat enim interdum interdum. In ornare pellentesque lectus, in euismod libero volutpat vitae. Proin hendrerit in est at faucibus. Proin massa leo, pellentesque non aliquam non, lobortis molestie risus. Curabitur vel justo in augue rutrum convallis.")
		ui.Outputs() <- []byte("Four cursus massa sit amet tortor blandit, sit amet volutpat nisi venenatis. Ut laoreet congue tempus. Nulla eu enim in ligula pulvinar tincidunt quis ac tortor. Morbi varius scelerisque mauris in laoreet. Proin sit amet nisl felis. Etiam a venenatis ante. Etiam tempus arcu diam, ac laoreet risus auctor vitae. Praesent facilisis ligula ante. Maecenas eget ipsum varius, pulvinar elit vel, sagittis purus. Etiam auctor enim sit amet blandit egestas. Quisque a volutpat leo, et semper eros. Nullam sodales eleifend nisl quis auctor. Etiam auctor dui non nisl malesuada, vel ornare erat rhoncus. Integer sapien tortor, sagittis vitae lacus nec, ultricies tincidunt sapien.")
		ui.Outputs() <- []byte("Five proin vel nibh placerat arcu consequat efficitur ac  quis orci. Phasellus et metus ligula. Nulla vestibulum varius mi, ut lobortis eros ullamcorper ac. Integer fringilla, tortor in eleifend condimentum, nulla ipsum imperdiet lorem, sed sollicitudin nisl mauris at massa. Morbi lobortis quis lectus et dictum. Praesent ut urna porta, maximus neque vitae, pharetra libero. Integer eu metus non mi vehicula ultrices pellentesque sit amet felis. Aliquam sed risus lectus. Cras accumsan sagittis leo vel sagittis. Sed lacinia feugiat est, a tincidunt odio egestas at.")
		ui.Outputs() <- []byte("Six lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum et nunc in dui efficitur commodo sed ut lectus. Etiam a urna nec augue gravida imperdiet. Aenean luctus ut augue in laoreet. Nunc ut dui sem. Maecenas id leo purus. Maecenas enim purus, finibus sit amet aliquam sit amet, commodo ut velit. Nunc nunc lectus, pulvinar ut metus quis, laoreet dapibus dolor. In rhoncus quis ligula sit amet pharetra. Aliquam nunc velit, pharetra nec imperdiet nec, iaculis nec ex. Vestibulum porta tristique dignissim. Pellentesque ac maximus lorem, ut viverra risus.")
		ui.Outputs() <- []byte("Seven velit ligula, gravida at enim a, interdum egestas massa. Maecenas feugiat commodo velit, et tristique erat dictum ut. Phasellus vel pulvinar nunc, eu convallis sem. In elit eros, fringilla sed porta eu, pharetra at leo. Integer pretium, odio a venenatis elementum, enim elit porttitor purus, eget fermentum est ex eu est. Duis quis euismod velit, eleifend lacinia sapien. Cras non ullamcorper turpis. Pellentesque elementum dui risus, ut fringilla tortor cursus sed. Pellentesque ipsum dui, fringilla vel tellus in, imperdiet semper ligula.")
		ui.Outputs() <- []byte("Eight tempus orci nibh, at consectetur nisi sagittis eu. Nam elementum, erat sed mattis malesuada, est neque elementum diam, blandit iaculis ipsum nunc eget diam. Nulla ut dolor sed elit faucibus mattis. Nam elementum urna vel interdum volutpat. Curabitur convallis imperdiet mi, at feugiat enim interdum interdum. In ornare pellentesque lectus, in euismod libero volutpat vitae. Proin hendrerit in est at faucibus. Proin massa leo, pellentesque non aliquam non, lobortis molestie risus. Curabitur vel justo in augue rutrum convallis.")
		ui.Outputs() <- []byte("Nine cursus massa sit amet tortor blandit, sit amet volutpat nisi venenatis. Ut laoreet congue tempus. Nulla eu enim in ligula pulvinar tincidunt quis ac tortor. Morbi varius scelerisque mauris in laoreet. Proin sit amet nisl felis. Etiam a venenatis ante. Etiam tempus arcu diam, ac laoreet risus auctor vitae. Praesent facilisis ligula ante. Maecenas eget ipsum varius, pulvinar elit vel, sagittis purus. Etiam auctor enim sit amet blandit egestas. Quisque a volutpat leo, et semper eros. Nullam sodales eleifend nisl quis auctor. Etiam auctor dui non nisl malesuada, vel ornare erat rhoncus. Integer sapien tortor, sagittis vitae lacus nec, ultricies tincidunt sapien.")
		ui.Outputs() <- []byte("Ten proin vel nibh placerat arcu consequat efficitur ac  quis orci. Phasellus et metus ligula. Nulla vestibulum varius mi, ut lobortis eros ullamcorper ac. Integer fringilla, tortor in eleifend condimentum, nulla ipsum imperdiet lorem, sed sollicitudin nisl mauris at massa. Morbi lobortis quis lectus et dictum. Praesent ut urna porta, maximus neque vitae, pharetra libero. Integer eu metus non mi vehicula ultrices pellentesque sit amet felis. Aliquam sed risus lectus. Cras accumsan sagittis leo vel sagittis. Sed lacinia feugiat est, a tincidunt odio egestas at.")
	}()

	x := &navigation.Room{
		ID:        1,
		Name:      "1",
		Known:     true,
		HasPlayer: true,
		X:         gox.NewInt(5),
		Y:         gox.NewInt(3),
		Exits:     map[string]*navigation.Room{},
	}
	in := &navigation.Room{
		ID:    2,
		Name:  "2",
		Known: true,
		X:     gox.NewInt(6),
		Y:     gox.NewInt(3),
		Exits: map[string]*navigation.Room{
			"ne":  {ID: 7, Name: "7"},
			"out": x,
		},
	}
	in.Exits["e"] = &navigation.Room{
		ID:   8,
		Name: "8",
		Exits: map[string]*navigation.Room{
			"w": in,
			"d": {ID: 16, Name: "G"},
		},
	}
	inse := &navigation.Room{
		ID:    15,
		Name:  "F",
		Known: true,
		X:     gox.NewInt(8),
		Y:     gox.NewInt(1),
		Exits: map[string]*navigation.Room{
			"nw": in,
		},
	}
	in.Exits["se"] = inse
	n := &navigation.Room{
		ID:    3,
		Name:  "3",
		Known: true,
		X:     gox.NewInt(5),
		Y:     gox.NewInt(5),
		Exits: map[string]*navigation.Room{
			"n": {ID: 9, Name: "9"},
			"s": x,
		},
	}
	s := &navigation.Room{ // hur markera target?
		ID:   4,
		Name: "4",
		Exits: map[string]*navigation.Room{
			"n": x,
		},
	}
	sw := &navigation.Room{
		ID:    5,
		Name:  "5",
		Known: true,
		Exits: map[string]*navigation.Room{
			"ne":  x,
			"u":   {ID: 10, Name: "A"},
			"d":   {ID: 11, Name: "B"},
			"out": {ID: 12, Name: "C"},
		},
	}
	w := &navigation.Room{
		ID:    6,
		Name:  "6",
		Known: true,
		Exits: map[string]*navigation.Room{
			"e":  x,
			"s":  sw,
			"nw": {ID: 13, Name: "D"},
			"u":  {ID: 14, Name: "E"},
		},
	}
	sw.Exits["n"] = w

	x.Exits["n"] = w
	x.Exits["in"] = in
	x.Exits["n"] = n
	x.Exits["s"] = s
	x.Exits["sw"] = sw
	x.Exits["w"] = w

	ui.SetRoom(x)

	ui.AddVital("health", tui.NewHealthVital())
	ui.UpdateVital("health", 800, 1000)
	ui.AddVital("mana", tui.NewManaVital())
	ui.UpdateVital("mana", 950, 1000)

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
