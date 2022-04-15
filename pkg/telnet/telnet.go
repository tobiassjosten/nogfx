package telnet

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

// Client is a wrapper around a telnet io.ReadWriter stream.
type Client struct {
	data     io.ReadWriter
	reader   *bufio.Reader
	commands chan []byte
}

// NewClient wraps a given reader and returns a new Client.
func NewClient(data io.ReadWriter) (*Client, <-chan []byte) {
	commands := make(chan []byte)

	client := &Client{
		data:     data,
		reader:   bufio.NewReader(data),
		commands: commands,
	}

	return client, commands
}

// Scanner creates a bufio.Scanner to abstract some low-level reading.
func (client *Client) Scanner() *bufio.Scanner {
	scanner := bufio.NewScanner(client)
	scanner.Split(func(data []byte, atEOF bool) (int, []byte, error) {
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
	})

	return scanner
}

// Write sends data to the server.
func (client *Client) Write(data []byte) (int, error) {
	return client.data.Write(data)
}

// Read parses and returns data received from the server.
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

		if b != IAC && len(command) == 0 {
			buffer[count] = b
			count++
			continue
		}

		command = append(command, b)

		if bytes.Equal(command, []byte{IAC, IAC}) {
			command = []byte{}
			continue
		}

		if bytes.Equal(command, []byte{IAC, GA}) {
			buffer[count] = b
			count++
			return count, nil
		}

		processed, responses := client.processCommand(command)

		for _, response := range responses {
			if _, err := client.data.Write(response); err != nil {
				return count, fmt.Errorf(
					"failed sending response (%s -> %s): %w",
					command, response, err,
				)
			}
		}

		if processed {
			client.commands <- command
			command = []byte{}
		}
	}

	return count, nil
}

// CommandToString creates a string representation of a telnet command sequence.
func CommandToString(command []byte) string {
	var chars []string
	for _, b := range command {
		switch b {
		case ECHO:
			chars = append(chars, "ECHO")
		case TTYPE:
			chars = append(chars, "TTYPE")
		case MCCP2:
			chars = append(chars, "MCCP2")
		case ATCP:
			chars = append(chars, "ATCP")
		case GMCP:
			chars = append(chars, "GMCP")
		case SE:
			chars = append(chars, "SE")
		case GA:
			chars = append(chars, "GA")
		case SB:
			chars = append(chars, "SB")
		case WILL:
			chars = append(chars, "WILL")
		case WONT:
			chars = append(chars, "WONT")
		case DO:
			chars = append(chars, "DO")
		case DONT:
			chars = append(chars, "DONT")
		case IAC:
			chars = append(chars, "IAC")
		default:
			chars = append(chars, fmt.Sprintf("%d", b))
		}
	}

	return strings.Join(chars, " ")
}
