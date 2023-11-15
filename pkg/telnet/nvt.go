package telnet

import (
	"bufio"
	"net"
)

// Telnet is a symmetric protocol, with no distinct server and client. Options
// can be negotiated separately for each side and we need to keep track of who
// has enabled what. This type helps us more clearly differentiate the sides.
type optionSide bool

const (
	ourside   optionSide = true
	theirside            = false
)

type optionState int

const (
	// StateDisabled is the default state for any option, meaning its
	// activation hasn't yet been agreed upon.
	StateDisabled optionState = iota

	StateDisabling
	StateEnabled
	StateEnabling
)

func (s optionState) On() bool {
	return s == StateEnabled || s == StateEnabling
}

func (s optionState) Off() bool {
	return s == StateDisabled || s == StateDisabling
}

// CommandFunc is a callback function for when Telnet commands are found, used
// for negotiation and custom logic.
type CommandFunc func(cmd []byte, conn net.Conn) error

// NVT (Network Virtual Terminal) represents a bi-directional character device
// and is a fundamental concept in the Telnet protocol (RFC 854). It acts as
// both "server" and "client", with both ends of a connection being equal, and
// state requiring negotiation and unanimous agreement.
type NVT struct {
	net.Conn

	buffer *bufio.Reader

	options map[optionSide]map[byte]optionState

	outBuffer []byte
	cmdBuffer []byte

	ourCoulds   []byte
	theirCoulds []byte

	CommandFunc CommandFunc
}

// NewNVT creates a NVT with some sane defaults.
func NewNVT(conn net.Conn) *NVT {
	return &NVT{
		Conn: conn,

		buffer: bufio.NewReader(conn),

		options: map[optionSide]map[byte]optionState{
			ourside:   {},
			theirside: {},
		},

		ourCoulds:   []byte{Echo},
		theirCoulds: []byte{SuppressGoAhead, GMCP},
	}
}
