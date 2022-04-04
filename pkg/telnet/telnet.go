package telnet

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
)

const (
	ECHO  byte = 1
	LF    byte = 10
	CR    byte = 13
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

// implementera:
// ECHO (1): http://pcmicro.com/NetFoss/RFC857.html
// STATUS (5): http://pcmicro.com/NetFoss/RFC859.html
// NAOCRD (10): https://www.ietf.org/rfc/rfc652.txt:w
// LOGOUT (18)?: https://www.ietf.org/rfc/rfc727.txt
// TERMINAL-TYPE (24): http://pcmicro.com/NetFoss/RFC1091.html
// NAWS (31): http://pcmicro.com/NetFoss/RFC1073.html
// CHARSET (42): https://www.ietf.org/rfc/rfc2066.txt
// Telnet Suppress Local Echo (45)?

// http://mud-dev.wikidot.com/telnet:negotiation
// https://blog.ikeran.org/?p=129
// http://pcmicro.com/NetFoss/telnet.html
// https://www.ironrealms.com/gmcp-doc
// https://tintin.sourceforge.io/protocols/mssp/
// http://www.mushclient.com/mushclient/mxp.htm
// https://wiki.mudlet.org/w/Manual:Supported_Protocols

type Stream struct {
	data     io.ReadWriter
	reader   *bufio.Reader
	commands chan []byte

	enabled map[byte]struct{}
	accepts map[byte]struct{}
}

// NewStream wraps a given reader and returns a new Stream.
func NewStream(data io.ReadWriter) (*Stream, <-chan []byte) {
	commands := make(chan []byte)

	return &Stream{
		data:     data,
		reader:   bufio.NewReader(data),
		commands: commands,
		accepts: map[byte]struct{}{
			ATCP: {},
			GMCP: {},
		},
		enabled: map[byte]struct{}{},
	}, commands
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

func (stream *Stream) Write(buffer []byte) (int, error) {
	return stream.data.Write(buffer)
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
		_, enabled := stream.enabled[command[2]]
		_, accepts := stream.accepts[command[2]]

		switch {
		case enabled:
			fmt.Printf("[WILL] %d (%X) enabled\n", command[2], command[2])
			stream.enabled[command[2]] = struct{}{}

		case accepts:
			if err := stream.do(command[2]); err != nil {
				// log.Warnf("failed accepting option: %d (%X)", command[2], command[2])
			}
			fmt.Printf("[WILL] %d (%X) accepted\n", command[2], command[2])
			stream.enabled[command[2]] = struct{}{}

		default:
			if err := stream.dont(command[2]); err != nil {
				// log.Warnf("failed rejecting option: %d (%X)", command[2], command[2])
			}
			fmt.Printf("[WILL] %d (%X) rejected\n", command[2], command[2])
		}

		stream.commands <- command
		return []byte{}, 0

	case WONT:
		fmt.Printf("[WONT] %d (%X)\n", command[2], command[2])
		stream.commands <- command
		return []byte{}, 0

	case DO:
		if err := stream.wont(command[2]); err != nil {
			// log.Warnf("failed rejecting option: %d (%X)", command[2], command[2])
		}
		fmt.Printf("[DO] %d (%X) rejected\n", command[2], command[2])

		stream.commands <- command
		return []byte{}, 0

	case DONT:
		fmt.Printf("[DONT] %d (%X)\n", command[2], command[2])
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
