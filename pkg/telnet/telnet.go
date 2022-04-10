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

type Client struct {
	data       io.ReadWriter
	reader     *bufio.Reader
	commands   chan []byte
	acceptWill map[byte]struct{}
	acceptDo   map[byte]struct{}
}

// NewClient wraps a given reader and returns a new Client.
func NewClient(data io.ReadWriter) (*Client, <-chan []byte) {
	commands := make(chan []byte)

	client := &Client{
		data:       data,
		reader:     bufio.NewReader(data),
		commands:   commands,
		acceptWill: map[byte]struct{}{},
		acceptDo:   map[byte]struct{}{},
	}

	return client, commands
}

func (client *Client) Write(buffer []byte) (int, error) {
	log.Printf("> '%s'", string(buffer))
	return client.data.Write(buffer)
}

func (client *Client) Read(buffer []byte) (count int, err error) {
	command := []byte{}

	for bufferlen := len(buffer); count < bufferlen; {
		b, err := client.reader.ReadByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				close(client.commands)
			}

			return count, err
		}

		if b == IAC || len(command) > 0 {
			command, b = client.processCommand(append(command, b))
			if b == 0 {
				continue
			}
		}

		buffer[count] = b
		count++

		if b == GA {
			break
		}
	}

	return count, nil
}

var ScanGA bufio.SplitFunc = func(data []byte, atEOF bool) (int, []byte, error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.IndexByte(data, GA); i >= 0 {
		return i + 1, data[0:i], nil
	}

	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}
