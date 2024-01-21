package main

import (
	"bufio"
	"context"
	_ "embed"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/mock"
	"github.com/tobiassjosten/nogfx/pkg/telnet"
	"github.com/tobiassjosten/nogfx/pkg/tui"
	"github.com/tobiassjosten/nogfx/pkg/world"

	"github.com/gdamore/tcell/v2"
	"github.com/urfave/cli/v2"
)

const (
	defaultPort = 23
)

//go:embed help.tmpl
var appHelpTemplate string

func main() {
log.SetOutput(io.Discard)

	//en kommentar
	log.Println("trololol")

	fileFlags := os.O_APPEND | os.O_CREATE | os.O_WRONLY
	f, err := os.OpenFile("nogfx.log", fileFlags, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)

	cli.AppHelpTemplate = appHelpTemplate

	app := &cli.App{
		Name:      "nogfx",
		Usage:     "because the book is always better",
		ArgsUsage: "<hostname>",
		HideHelp:  true,

		Authors: []*cli.Author{
			{
				Name:  "Tobias Sjösten",
				Email: "tobias@nogfx.org",
			},
		},

		Action: func(c *cli.Context) error {
			if c.Args().Len() == 0 {
				return cli.ShowAppHelp(c)
			}

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

	ctx, err := ctxDirs(ctx)
	if err != nil {
		return err
	}

	client, err := client(address)
	if err != nil {
		return err
	}

	ui, err := ui()
	if err != nil {
		return err
	}

	engine := world.NewEngine(client, ui, address)
	return engine.Run(ctx)
}

func client(address string) (pkg.Client, error) {
	if address == "example.com:23" {
		return &mock.ClientMock{
			ScannerFunc: func() *bufio.Scanner {
				return bufio.NewScanner(strings.NewReader("asdf"))
			},
			CommandsFunc: func() <-chan []byte {
				return make(chan []byte)
			},
			WriteFunc: func(data []byte) (int, error) {
				return len(data), nil
			},
		}, nil
	}

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

	return tui.NewTUI(screen), nil
}

func ctxDirs(ctx context.Context) (context.Context, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed acquiring home directory: %w", err)
	}

	dir += "/nogfx/logs"

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed creating gamelog directory: %w", err)
	}

	ctx = context.WithValue(ctx, pkg.CtxLogdir, dir)

	return ctx, nil
}
