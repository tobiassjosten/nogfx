package telnet_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tobiassjosten/nogfx/pkg/telnet"
)

type mockStream struct {
	reader   io.Reader
	writer   io.Writer
	writeErr error
}

func (mock *mockStream) Read(p []byte) (int, error) {
	return mock.reader.Read(p)
}

func (mock mockStream) Write(p []byte) (int, error) {
	if mock.writeErr != nil {
		return 0, mock.writeErr
	}
	return mock.writer.Write(p)
}

func TestReader(t *testing.T) {
	tcs := []struct {
		data      []byte
		output    []byte
		commands  [][]byte
		response  string
		writerErr error
		errMsg    string
	}{
		{
			data:   []byte("xyz\n"),
			output: []byte("xyz\n"),
		},
		{
			data:   []byte{'x', telnet.IAC, telnet.IAC, 'y'},
			output: []byte{'x', 'y'},
		},
		{
			data:   []byte{'x', telnet.IAC, telnet.GA, 'y'},
			output: []byte{'x', telnet.GA, 'y'},
		},
		{
			data:     []byte{'x', telnet.IAC, telnet.WILL, telnet.ECHO, 'y'},
			output:   []byte{'x', 'y'},
			commands: [][]byte{{telnet.IAC, telnet.WILL, telnet.ECHO}},
			response: string([]byte{telnet.IAC, telnet.DO, telnet.ECHO}),
		},
		{
			data:     []byte{'y', telnet.IAC, telnet.WILL, telnet.GMCP, 'x'},
			output:   []byte{'y', 'x'},
			commands: [][]byte{{telnet.IAC, telnet.WILL, telnet.GMCP}},
			response: string([]byte{telnet.IAC, telnet.DO, telnet.GMCP}),
		},
		{
			data:      []byte{telnet.IAC, telnet.WILL, telnet.GMCP},
			writerErr: fmt.Errorf("x"),
			errMsg:    "failed sending response (\xff\xfb\xc9 -> \xff\xfd\xc9): x",
		},
		{
			data:     []byte{'x', telnet.IAC, telnet.SB, 'z', telnet.IAC, telnet.SE, 'y'},
			output:   []byte{'x', 'y'},
			commands: [][]byte{{telnet.IAC, telnet.SB, 'z', telnet.IAC, telnet.SE}},
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			reader := bytes.NewReader(tc.data)
			writer := bytes.NewBuffer([]byte{})
			stream := &mockStream{reader, writer, tc.writerErr}

			client := telnet.NewClient(stream)

			var commands [][]byte
			go func() {
				for command := range client.Commands() {
					commands = append(commands, command)
				}
			}()

			output, err := ioutil.ReadAll(client)

			if tc.errMsg != "" {
				assert.Equal(tc.errMsg, err.Error())
				return
			}

			require.Nil(err)
			assert.Equal(tc.output, output)
			assert.Equal(commands, tc.commands)
			assert.Equal(tc.response, writer.String())
		})
	}
}

func TestScanner(t *testing.T) {
	tcs := []struct {
		data   []byte
		output []byte
		err    error
	}{
		{
			data:   []byte("xyz\n"),
			output: []byte("xyz\n"),
		},
		{
			data:   append([]byte("xyz\n"), telnet.GA),
			output: []byte("xyz\n"),
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			reader := bytes.NewReader(tc.data)
			writer := &strings.Builder{}
			stream := &mockStream{reader, writer, nil}

			client := telnet.NewClient(stream)

			var commands [][]byte
			go func() {
				for command := range client.Commands() {
					commands = append(commands, command)
				}
			}()

			scanner := client.Scanner()

			buf := make([]byte, 2)
			scanner.Buffer(buf, 4096)

			output := []byte{}
			for scanner.Scan() {
				output = append(output, scanner.Bytes()...)
			}

			err := scanner.Err()

			if tc.err != nil {
				assert.Equal(tc.err, err)
				return
			}

			require.Nil(err)
			assert.Equal(tc.output, output)
		})
	}
}

func TestWriter(t *testing.T) {
	tcs := []struct {
		data   []byte
		length int
		err    error
	}{
		{
			data:   []byte("xyz\n"),
			length: 4,
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			reader := bytes.NewReader(tc.data)
			writer := &strings.Builder{}
			stream := &mockStream{reader, writer, nil}

			client := telnet.NewClient(stream)

			length, err := client.Write(tc.data)

			if tc.err != nil {
				assert.Equal(tc.err, err)
				return
			}

			require.Nil(err)
			assert.Equal(tc.length, length)
			assert.Equal(string(tc.data), writer.String())
		})
	}
}

func TestCommandToString(t *testing.T) {
	tcs := []struct {
		command []byte
		output  string
	}{
		{
			command: []byte{telnet.IAC, telnet.GA},
			output:  "IAC GA",
		},
		{
			command: []byte{telnet.IAC, telnet.WILL, telnet.GMCP},
			output:  "IAC WILL GMCP",
		},
		{
			command: []byte{telnet.IAC, telnet.WONT, telnet.MCCP2},
			output:  "IAC WONT MCCP2",
		},
		{
			command: []byte{telnet.IAC, telnet.DO, telnet.ECHO},
			output:  "IAC DO ECHO",
		},
		{
			command: []byte{telnet.IAC, telnet.DONT, telnet.ATCP},
			output:  "IAC DONT ATCP",
		},
		{
			command: []byte{telnet.IAC, telnet.SB, telnet.TTYPE, 65, telnet.IAC, telnet.SE},
			output:  "IAC SB TTYPE 65 IAC SE",
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			assert := assert.New(t)

			assert.Equal(tc.output, telnet.CommandToString(tc.command))
		})
	}
}
