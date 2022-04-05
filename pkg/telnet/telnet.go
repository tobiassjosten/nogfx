package telnet

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
)

const (
	ECHO  byte = 1
	LF    byte = 10
	CR    byte = 13
	TTYPE byte = 24
	MCCP  byte = 85
	MCCP2 byte = 86
	ATCP  byte = 200
	GMCP  byte = 201
	SE    byte = 240
	GA    byte = 249
	SB    byte = 250
	WILL  byte = 251
	WONT  byte = 252
	DO    byte = 253
	DONT  byte = 254
	IAC   byte = 255
)

type Stream struct {
	data     io.ReadWriter
	reader   *bufio.Reader
	commands chan []byte
}

// NewStream wraps a given reader and returns a new Stream.
func NewStream(data io.ReadWriter) (*Stream, <-chan []byte) {
	commands := make(chan []byte)

	stream := &Stream{
		data:     data,
		reader:   bufio.NewReader(data),
		commands: commands,
	}

	return stream, commands
}

func (stream *Stream) Write(buffer []byte) (int, error) {
	return stream.data.Write(buffer)
}

func (stream *Stream) will(command byte) error {
	_, err := stream.data.Write([]byte{IAC, WILL, command})
	return err
}

func (stream *Stream) wont(command byte) error {
	_, err := stream.data.Write([]byte{IAC, WONT, command})
	return err
}

func (stream *Stream) do(command byte) error {
	_, err := stream.data.Write([]byte{IAC, DO, command})
	return err
}

func (stream *Stream) dont(command byte) error {
	_, err := stream.data.Write([]byte{IAC, DONT, command})
	return err
}

func (stream *Stream) subneg(b byte, value []byte) error {
	var v byte = 0
	if len(value) > 0 {
		v = 1
	}

	_, err := stream.data.Write(append(append(
		[]byte{IAC, SB, b, v},
		value...,
	), IAC, SE))
	return err
}

func (stream *Stream) gmcp(value []byte) error {
	_, err := stream.data.Write(append(append(
		[]byte{IAC, SB, GMCP},
		value...,
	), IAC, SE))
	return err
}

func (stream *Stream) Read(buffer []byte) (count int, err error) {
	command := []byte{}

	for bufferlen := len(buffer); count < bufferlen; {
		b, err := stream.reader.ReadByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				close(stream.commands)
			}

			return count, err
		}

		if b == IAC || len(command) > 0 {
			command, b = stream.processCommand(append(command, b))
			if b == 0 {
				continue
			}
		}

		buffer[count] = b
		count++

		// Achaea skickar ibland CR CR LN, ibland CR LN. Kanske kan vi
		// skapa någon konfigurerbar mekanism för att normalisera sånt?
		if b == '\n' || b == GA {
			break
		}
	}

	return count, nil
}

func (stream *Stream) processCommand(command []byte) ([]byte, byte) {
	if bytes.Equal(command, []byte{IAC, IAC}) {
		return []byte{}, IAC
	}

	// More reliable than newline to mark the end of a message, so we relay
	// it upstream for processing in the game logic.
	if bytes.Equal(command, []byte{IAC, GA}) {
		return []byte{}, GA
	}

	if len(command) < 3 {
		return command, 0
	}

	switch command[1] {
	case WILL:
		if command[2] == GMCP {
			if err := stream.do(command[2]); err != nil {
				log.Printf("failed accepting WILL %d", command[2])
			}
			stream.gmcp([]byte(`Core.Hello { "client": "NoGFX", "version": "0.0.1" }`))
			stream.gmcp([]byte(`Core.Supports.Set [ "Char 1", "Char.Skills 1", "Char.Items 1", "Comm.Channel 1", "Room 1", "IRE.Rift 1"]`))
		} else {
			if err := stream.dont(command[2]); err != nil {
				log.Printf("failed rejecting WILL %d", command[2])
			}
		}

		stream.commands <- command
		return []byte{}, 0

	case WONT:
		stream.commands <- command
		return []byte{}, 0

	case DO:
		if err := stream.wont(command[2]); err != nil {
			log.Printf("failed rejecting DO %d", command[2])
		}

		stream.commands <- command
		return []byte{}, 0

	case DONT:
		stream.commands <- command
		return []byte{}, 0

	default:
		// Noop.
	}

	if len(command) < 5 {
		return command, 0
	}

	if bytes.Equal(command[:2], []byte{IAC, SB}) {
		if bytes.Equal(command[len(command)-2:], []byte{IAC, SE}) {
			stream.commands <- command
			return []byte{}, 0
		}
	}

	return command, 0
}
