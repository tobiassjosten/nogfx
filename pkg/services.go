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
	SplitFunc() bufio.SplitFunc

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
