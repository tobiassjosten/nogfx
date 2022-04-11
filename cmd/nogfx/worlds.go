package main

import (
	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/world/achaea"
)

var worlds = map[string]func(pkg.Client) pkg.World{
	"achaea.com:23":  achaeaWorld,
	"50.31.100.8:23": achaeaWorld,
}

func achaeaWorld(client pkg.Client) pkg.World {
	return achaea.NewWorld(client)
}
