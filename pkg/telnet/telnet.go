package telnet

import (
	"bufio"
	"bytes"
	"errors"
	"io"
)

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

func (client *Client) Write(data []byte) (int, error) {
	return client.data.Write(data)
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
