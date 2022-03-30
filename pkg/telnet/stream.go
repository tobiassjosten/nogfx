package telnet

import (
	"bytes"
	"io"
	"net"
)

const (
	ECHO = 1
	LF   = 10
	CR   = 13
	GMCP = 201

	SE   = 240
	GA   = 249
	SB   = 250
	WILL = 251
	WONT = 252
	DO   = 253
	DONT = 254
	IAC  = 255
)

// okända kommandon DO/DONT -> WONT

// implementera:
// ECHO (1): http://pcmicro.com/NetFoss/RFC857.html
// SUPRESS GO-AHEAD (3): http://pcmicro.com/NetFoss/RFC858.html
// STATUS (5): http://pcmicro.com/NetFoss/RFC859.html
// NAOCRD (10): https://www.ietf.org/rfc/rfc652.txt:w
// LOGOUT (18)?: https://www.ietf.org/rfc/rfc727.txt
// TERMINAL-TYPE (24): http://pcmicro.com/NetFoss/RFC1091.html
// NAWS (31): http://pcmicro.com/NetFoss/RFC1073.html
// CHARSET (42): https://www.ietf.org/rfc/rfc2066.txt
// Telnet Suppress Local Echo (45)?

type Stream struct {
	connection io.ReadWriteCloser
	buffer     []byte
	command    []byte
	linebreak  []byte
}

func NewStream(connection io.ReadWriteCloser) *Stream {
	return &Stream{
		connection: connection,
		linebreak:  []byte{CR, LF},
	}
}

func Dial(network, address string) (*Stream, error) {
	connection, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}

	stream := NewStream(connection)

	err = stream.Do(GMCP)
	if err != nil {
		return nil, err
	}

	return stream, nil
}

func (stream *Stream) Will(command byte) error {
	_, err := stream.connection.Write([]byte{IAC, WILL, command})
	return err
}

func (stream *Stream) Wont(command byte) error {
	_, err := stream.connection.Write([]byte{IAC, WONT, command})
	return err
}

func (stream *Stream) Do(command byte) error {
	_, err := stream.connection.Write([]byte{IAC, DO, command})
	return err
}

func (stream *Stream) Dont(command byte) error {
	_, err := stream.connection.Write([]byte{IAC, DONT, command})
	return err
}

func (stream *Stream) processCommand() {
	if bytes.Equal(stream.command, []byte{IAC, DO, 123}) {
		_, err := stream.Write([]byte{IAC, WONT, stream.command[2]})
		if err != nil {
			panic(err)
		}
	}

	stream.command = []byte{}
}
