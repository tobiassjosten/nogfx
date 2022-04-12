package main

import (
	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/world/achaea"
)

var worlds = map[string]func(pkg.UI, pkg.Client) pkg.World{
	"achaea.com:23":  achaeaWorld,
	"50.31.100.8:23": achaeaWorld,
}

func NewWorld(ui pkg.UI, client pkg.Client) pkg.World {
	var world pkg.World = pkg.NewGenericWorld(ui, client)
	if constructor, ok := worlds["achaea.com:23"]; ok {
		world = constructor(ui, client)
	}

	return world
}

func achaeaWorld(ui pkg.UI, client pkg.Client) pkg.World {
	return achaea.NewWorld(ui, client)
}
