package pkg

import (
	"bufio"
	"context"
	"io"

	"github.com/tobiassjosten/nogfx/pkg/navigation"
)

// Client is the application's main connection to the game server.
type Client interface {
	io.ReadWriter
	Send([]byte)
	Commands() <-chan []byte
	Scanner() *bufio.Scanner

	// Telnet utilities.
	Will(byte) error
	Wont(byte) error
	Do(byte) error
	Dont(byte) error
	Subneg(byte, []byte) error
}

// Module is an atomic extension of game logic.
type Module interface {
	ProcessInput([]byte) [][]byte
	ProcessOutput([]byte) [][]byte
}

// ModuleConstructor is a function that constructs a Module.
type ModuleConstructor func(Client, UI) Module

// UI is the primary user interface for the application.
type UI interface {
	Inputs() <-chan []byte
	Outputs() chan<- []byte
	Run(context.Context) error

	Print([]byte)

	MaskInput()
	UnmaskInput()

	SetCharacter(Character)
	SetRoom(*navigation.Room)
}

// World represents a game and hooks into all their various specific logic.
// When processors return a slice of slices of bytes, they can use that to
// signal that player input or server output should be omitted (nil) or that it
// should be replaced (a new [][]byte).
type World interface {
	ProcessInput([]byte) [][]byte
	ProcessOutput([]byte) [][]byte
	ProcessCommand([]byte) error
}
