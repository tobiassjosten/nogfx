package telnet

import (
	"bufio"
	"bytes"
	"errors"
	"io"
)

// från achaea
// WILL 25
// WILL 200
// WILL 201 (GMCP)
// WILL 86
// WILL 24 ("terminal type"?)

const (
	ECHO byte = 1
	LF   byte = 10
	CR   byte = 13
	GMCP byte = 201
	SE   byte = 240
	GA   byte = 249
	SB   byte = 250
	WILL byte = 251
	WONT byte = 252
	DO   byte = 253
	DONT byte = 254
	IAC  byte = 255
)

// okända kommandon DO/DONT -> WONT?

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

// http://mud-dev.wikidot.com/telnet:negotiation
// https://blog.ikeran.org/?p=129
// http://pcmicro.com/NetFoss/telnet.html
// https://www.ironrealms.com/gmcp-doc
// https://tintin.sourceforge.io/protocols/mssp/
// MCCP?
// http://www.mushclient.com/mushclient/mxp.htm
// https://wiki.mudlet.org/w/Manual:Supported_Protocols

type Stream struct {
	data     io.ReadWriter
	reader   *bufio.Reader
	commands chan []byte
}

// NewStream wraps a given reader and returns a new Stream.
func NewStream(data io.ReadWriter) (*Stream, <-chan []byte) {
	commands := make(chan []byte)

	return &Stream{
		data:     data,
		reader:   bufio.NewReader(data),
		commands: commands,
	}, commands
}

func (stream *Stream) Will(command byte) error {
	_, err := stream.data.Write([]byte{IAC, WILL, command})
	return err
}

func (stream *Stream) Wont(command byte) error {
	_, err := stream.data.Write([]byte{IAC, WONT, command})
	return err
}

func (stream *Stream) Do(command byte) error {
	_, err := stream.data.Write([]byte{IAC, DO, command})
	return err
}

func (stream *Stream) Dont(command byte) error {
	_, err := stream.data.Write([]byte{IAC, DONT, command})
	return err
}

func (stream *Stream) Write(buffer []byte) (int, error) {
	return stream.data.Write(buffer)
}

func (stream *Stream) Read(readBuffer []byte) (count int, err error) {
	readLength := len(readBuffer)

	command := []byte{}

	for count < readLength {
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

		readBuffer[count] = b
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
	if len(command) == 0 {
		return []byte{}, 0
	}

	if command[0] != IAC {
		// @todo Log warning about invalid command sequence.
		return []byte{}, 0
	}

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
		if command[2] == 123 {
			if err := stream.Wont(132); err != nil {
				// log.Warnf("invalid command sequence: %s", command)
			}
		}

		stream.relayCommand(command)
		return []byte{}, 0

	case WONT:
		stream.relayCommand(command)
		return []byte{}, 0

	case DO:
		stream.relayCommand(command)
		return []byte{}, 0

	case DONT:
		stream.relayCommand(command)
		return []byte{}, 0

	default:
		// Noop.
	}

	if len(command) < 5 {
		return command, 0
	}

	if bytes.Equal(command[:2], []byte{IAC, SB}) {
		if bytes.Equal(command[len(command)-2:], []byte{IAC, SE}) {
			stream.relayCommand(command)
			return []byte{}, GA
		}
	}

	return command, 0
}

func (stream *Stream) relayCommand(command []byte) {
	// We can't trust that consumers are listening to our channel, so we
	// spawn goroutines so as to let Go handle buffering for us.
	go func(command []byte) {
		stream.commands <- command
	}(command)
}
