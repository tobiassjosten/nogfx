package main

import (
	"context"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/telnet"
	"github.com/tobiassjosten/nogfx/pkg/tui"
)

func main() {
	log.SetOutput(ioutil.Discard)

	fileFlags := os.O_APPEND | os.O_CREATE | os.O_WRONLY
	f, err := os.OpenFile("nogfx.log", fileFlags, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "mock",
				Usage: "mock connection",
			},
		},
		Action: func(c *cli.Context) error {
			return run(c.Bool("mock"))
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(mock bool) error {
	ctx := context.Background()

	ui, err := tui.NewTUI()
	if err != nil {
		log.Fatal(err)
	}

	address := "achaea.com:23"

	connection := mockReadWriter()
	if !mock {
		connection, err = net.Dial("tcp", address)
		if err != nil {
			return err
		}
	}

	client := telnet.NewClient(connection)

	world := NewWorld(ui, client, address)

	return pkg.Run(ctx, client, ui, world)
}
