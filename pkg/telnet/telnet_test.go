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

type MockData struct {
	reader io.Reader
	writer io.Writer
}

func (mock *MockData) Read(p []byte) (int, error) {
	return mock.reader.Read(p)
}

func (mock MockData) Write(p []byte) (int, error) {
	return mock.writer.Write(p)
}

func TestAsdf(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	tcs := []struct {
		data     []byte
		output   []byte
		commands [][]byte
		err      error
	}{
		{
			data:   []byte("xyz\n"),
			output: []byte("xyz\n"),
		},
		{
			data:   []byte{'x', telnet.IAC, telnet.IAC, 'y'},
			output: []byte{'x', telnet.IAC, 'y'},
		},
		{
			data:   []byte{'x', telnet.IAC, telnet.GA, 'y'},
			output: []byte{'x', telnet.GA, 'y'},
		},
		{
			data:     []byte{'x', telnet.IAC, telnet.WILL, telnet.GMCP, 'y'},
			output:   []byte{'x', 'y'},
			commands: [][]byte{{telnet.IAC, telnet.WILL, telnet.GMCP}},
		},
		{
			data:     []byte{'x', telnet.IAC, telnet.SB, 'z', telnet.IAC, telnet.SE, 'y'},
			output:   []byte{'x', 'y'},
			commands: [][]byte{{telnet.IAC, telnet.SB, 'z', telnet.IAC, telnet.SE}},
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			data := &MockData{
				bytes.NewReader(tc.data),
				&strings.Builder{},
			}

			stream, commandChan := telnet.NewStream(data)

			var commands [][]byte
			go func(commandChan <-chan []byte) {
				for command := range commandChan {
					commands = append(commands, command)
				}
			}(commandChan)

			output, err := ioutil.ReadAll(stream)

			if tc.err != nil {
				assert.Equal(tc.err, err)
				return
			}

			require.Nil(err)
			assert.Equal(tc.output, output)
			assert.Equal(commands, tc.commands)
		})
	}
}