package main

import (
	"github.com/tobiassjosten/nogfx-cli"
	"os"
)

func address() (address string) {
	host := "achaea.com"
	port := "23"

	switch {
	case len(os.Args) >= 3:
		port = os.Args[2]
		fallthrough
	case len(os.Args) >= 2:
		host = os.Args[1]
	}

	return host + ":" + port
}

func main() {
	engine := nogfx.NewEngine()
	engine.Run(address())
}
