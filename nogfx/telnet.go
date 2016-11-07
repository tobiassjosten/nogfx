package main

import (
	"bufio"
	"net"
)

const (
	IAC  byte = '\xff'
	DONT byte = '\xfe'
	DO   byte = '\xfd'
	WONT byte = '\xfc'
	WILL byte = '\xfb'
	SB   byte = '\xfa'
	GA   byte = '\xf9'
	SE   byte = '\xf0'
	GMCP byte = '\xc9'
)

type Telnet struct {
	connection   net.Conn
	reader       *bufio.Reader
	serverOutput chan string
	gmcpOutput   chan string
}

func NewTelnet() (*Telnet, <-chan string) {
	telnet := new(Telnet)
	telnet.serverOutput = make(chan string)
	return telnet, telnet.serverOutput
}

func (telnet *Telnet) Send(input string) {
	telnet.connection.Write(append(append([]byte(input), '\r'), '\n'))
}

func (telnet *Telnet) readByte() byte {
	b, err := telnet.reader.ReadByte()
	if err != nil {
		panic(err.Error())
	}

	return b
}

func (telnet *Telnet) serverWill(b byte) {
	switch b {
	case GMCP:
		telnet.connection.Write([]byte("\xff\xfa\xc9Core.Hello {\"client\":\"NoGFX\",\"version\":\"0.0.1\"}\xff\xf0"))
		telnet.connection.Write([]byte("\xff\xfa\xc9Core.Supports.Set [\"Char 1\",\"Char.Skills 1\",\"Char.Items 1\",\"Comm.Channel 1\",\"Room 1\",\"IRE.Misc 1\",\"IRE.Rift 1\",\"IRE.Target 1\",\"IRE.Tasks 1\",\"IRE.Time 1\"]\xff\xf0"))
	}
}

func (telnet *Telnet) serverWont(b byte) {
}

func (telnet *Telnet) serverDo(b byte) {
}

func (telnet *Telnet) serverDont(b byte) {
}

func (telnet *Telnet) serverBuffer(optionCode byte) {
	var buffer []byte

	for {
		b := telnet.readByte()

		if IAC == b {
			if SE == telnet.readByte() {
				break
			} else {
				panic("Invalid telnet negotiation")
			}
		}

		buffer = append(buffer, b)
	}

	switch optionCode {
	case GMCP:
		// @todo Parse (examples):
		// IRE.Time.Update { "daynight": "33" }
		// Room.RemovePlayer "Somedude"
		// Char.Defences.Remove [ "mosstattoo" ]
	}
}

func (telnet *Telnet) Main(address string) {
	connection, err := net.Dial("tcp", address)
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	telnet.connection = connection
	telnet.reader = bufio.NewReader(telnet.connection)

	var buffer []byte

	// @todo Can we trigger this from a connection event? Is it safe to assume
	// we can write to the connection here and now?
	telnet.connection.Write([]byte{IAC, DO, GMCP})

	for {
		b := telnet.readByte()

		if IAC == b {
			switch telnet.readByte() {
			case WILL:
				telnet.serverWill(telnet.readByte())

			case WONT:
				telnet.serverWont(telnet.readByte())

			case DO:
				telnet.serverDo(telnet.readByte())

			case DONT:
				telnet.serverDont(telnet.readByte())

			case SB:
				telnet.serverBuffer(telnet.readByte())

			case GA:

			default:
				panic("Invalid telnet negotiation")
			}
		} else {
			buffer = append(buffer, b)

			if '\n' == b {
				telnet.serverOutput <- string(buffer)
				buffer = buffer[:0]
			}
		}
	}
}
