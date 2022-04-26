package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/telnet"
	"github.com/tobiassjosten/nogfx/pkg/tui"
	"github.com/tobiassjosten/nogfx/pkg/world"

	"github.com/gdamore/tcell/v2"
	"github.com/urfave/cli/v2"
)

const (
	defaultPort = 23
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
		Action: func(c *cli.Context) error {
			address, err := address(c.Args().Get(0))
			if err != nil {
				return err
			}

			return run(address)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func address(host string) (string, error) {
	if strings.Contains(host, ":") {
		parts := strings.Split(host, ":")
		// @todo Add support for IPv6 addresses.
		if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
			return "", fmt.Errorf("invalid address '%s'", host)
		}

		if _, err := strconv.ParseFloat(parts[1], 64); err != nil {
			return "", fmt.Errorf("invalid port '%s'", parts[1])
		}

		return host, nil
	}

	if host == "" {
		host = "example.com"
	}

	return fmt.Sprintf("%s:%d", host, defaultPort), nil
}

func run(address string) error {
	ctx := context.Background()

	client, err := client(address)
	if err != nil {
		return err
	}

	ui, err := ui()
	if err != nil {
		return err
	}

	world := world.New(client, ui, address)

	return pkg.Run(ctx, client, ui, world)
}

func client(address string) (*telnet.Client, error) {
	connection, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	return telnet.NewClient(connection), nil
}

func ui() (*tui.TUI, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	return tui.NewTUI(screen, tui.NewPanes()), nil
}
