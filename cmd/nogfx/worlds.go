package main

import (
	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/tui"
	"github.com/tobiassjosten/nogfx/pkg/world/achaea"
)

var worlds = map[string]func(*tui.TUI, pkg.Client) pkg.World{
	"achaea.com:23":  NewAchaeaWorld,
	"50.31.100.8:23": NewAchaeaWorld,
}

// NewWorld creates a World specific to the game being played.
func NewWorld(tui *tui.TUI, client pkg.Client, address string) pkg.World {
	var world pkg.World = pkg.NewGenericWorld(tui, client)
	if constructor, ok := worlds[address]; ok {
		world = constructor(tui, client)
	}

	return world
}

// NewAchaeaWorld wraps the specific constructor to create an interfaced World.
func NewAchaeaWorld(tui *tui.TUI, client pkg.Client) pkg.World {
	return achaea.NewWorld(tui, client)
}
