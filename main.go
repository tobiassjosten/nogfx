package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/process"
	"github.com/tobiassjosten/nogfx/pkg/telnet"
	"github.com/tobiassjosten/nogfx/pkg/tui"
	"github.com/tobiassjosten/nogfx/pkg/world/achaea"

	"github.com/gdamore/tcell/v2"
	"golang.org/x/net/idna"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: nogfx example.com:23")
	}

	f, err := os.OpenFile(
		filepath.Join(pkg.Directory, "errors.log"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)

	address, err := parseAddress(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	if err := run(address); err != nil {
		log.Fatal(err)
	}
}

func parseAddress(address string) (string, error) {
	if !strings.Contains(address, ":") {
		address += ":23"
	}

	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return "", fmt.Errorf("invalid server address %q: %w", address, err)
	}

	if _, err := strconv.ParseFloat(port, 64); err != nil {
		return "", fmt.Errorf("invalid server port %q: %w", port, err)
	}

	host, err = idna.Lookup.ToASCII(host)
	if err != nil {
		return "", fmt.Errorf("invalid server host: %w", err)
	}

	return net.JoinHostPort(host, port), nil
}

func run(address string) error {
	ctx := context.Background()

	client, err := telnet.Dial(address)
	if err != nil {
		return err
	}

	ui, err := ui()
	if err != nil {
		return err
	}

	logProcessor, err := process.LogProcessor(
		filepath.Join(pkg.Directory, "logs"),
		fmt.Sprintf(
			"%s-%s.log",
			strings.Split(address, ":")[0],
			time.Now().Format("20060102-150405"),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create log processor: %w", err)
	}

	engine := &pkg.Engine{
		Client: client,
		UI:     ui,
		Processor: process.ChainProcessor(
			process.RepeatInputProcessor(),
			logProcessor,
		),
	}

	switch address {
	case "achaea.com:23", "50.31.100.8:23":
		processor, err := achaea.Processor(client, ui)
		if err != nil {
			return fmt.Errorf("failed to create Acahea processor: %w", err)
		}

		engine.Processor = processor
	}

	return engine.Run(ctx)
}

func ui() (*tui.TUI, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	return tui.NewTUI(screen), nil
}
