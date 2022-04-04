package telnet_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

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
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			data := &MockData{
				bytes.NewReader(tc.data),
				&strings.Builder{},
			}

			stream, commands := telnet.NewStream(data)

			output, err := ioutil.ReadAll(stream)

			if assert.Equal(tc.err, err) {
				assert.Equal(tc.output, output)
			}

			if tc.commands == nil {
				return
			}

		main:
			for i := 0; ; {
				select {
				case command := <-commands:
					if assert.LessOrEqual(i+1, len(tc.commands)) {
						assert.Equal(tc.commands[i], command)
					}

				case <-time.After(1 * time.Second):
					assert.Fail("channel failed to close")
					break main
				}

				if commands == nil {
					break
				}
			}
		})
	}
}
