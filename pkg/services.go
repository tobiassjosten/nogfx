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
	SetTarget(*Target)
}

// World represents a game and hooks into all their various specific logic.
type World interface {
	ProcessInput(Input) Input
	ProcessOutput(Output) Output
	ProcessCommand([]byte)

	Send([]byte)
	Print([]byte)
}
