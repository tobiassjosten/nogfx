package main

import (
	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/world/achaea"
)

var worlds = map[string]func(pkg.UI, pkg.Client) pkg.World{
	"achaea.com:23":  NewAchaeaWorld,
	"50.31.100.8:23": NewAchaeaWorld,
}

// NewWorld creates a World specific to the game being played.
func NewWorld(ui pkg.UI, client pkg.Client, address string) pkg.World {
	// @todo Make this actually configurable.
	var world pkg.World = pkg.NewGenericWorld(ui, client)
	if constructor, ok := worlds[address]; ok {
		world = constructor(ui, client)
	}

	return world
}

// NewAchaeaWorld wraps the specific constructor to create an interfaced World.
func NewAchaeaWorld(ui pkg.UI, client pkg.Client) pkg.World {
	return achaea.NewWorld(ui, client)
}
