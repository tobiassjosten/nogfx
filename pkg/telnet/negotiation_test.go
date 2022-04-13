package telnet_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/tobiassjosten/nogfx/pkg/telnet"

	"github.com/stretchr/testify/assert"
)

func TestAccepts(t *testing.T) {
	assert := assert.New(t)

	tcs := []struct {
		serverWill byte
		acceptWill byte
		serverDo   byte
		acceptDo   byte
		err        error
	}{
		{
			serverWill: 123,
			acceptWill: 123,
		},
		{
			serverWill: 124,
		},
		{
			serverDo: 125,
			acceptDo: 125,
		},
		{
			serverDo: 126,
		},
		{
			serverWill: 127,
			acceptWill: 127,
			serverDo:   128,
			acceptDo:   128,
		},
		{
			serverWill: 127,
			serverDo:   128,
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			output := []byte{}
			if tc.serverWill > 0 {
				output = append(output, telnet.IAC, telnet.WILL, tc.serverWill)
			}
			if tc.serverDo > 0 {
				output = append(output, telnet.IAC, telnet.DO, tc.serverDo)
			}

			input := []byte{}
			if tc.acceptWill > 0 {
				input = append(input, telnet.IAC, telnet.DO, tc.acceptWill)
			} else if tc.serverWill > 0 {
				input = append(input, telnet.IAC, telnet.DONT, tc.serverWill)
			}
			if tc.acceptDo > 0 {
				input = append(input, telnet.IAC, telnet.WILL, tc.acceptDo)
			} else if tc.serverDo > 0 {
				input = append(input, telnet.IAC, telnet.WONT, tc.serverDo)
			}

			builder := &strings.Builder{}
			stream := &MockStream{
				bytes.NewReader(output),
				builder,
			}

			client, commandChan := telnet.NewClient(stream)

			if tc.acceptWill > 0 {
				client.AcceptWill(tc.acceptWill)
			}
			if tc.acceptDo > 0 {
				client.AcceptDo(tc.acceptDo)
			}

			var commands [][]byte
			go func(commandChan <-chan []byte) {
				for command := range commandChan {
					commands = append(commands, command)
				}
			}(commandChan)

			output, err := ioutil.ReadAll(client)

			if tc.err != nil {
				assert.Equal(tc.err, err)
				return
			}

			assert.Equal(string(input), builder.String())
		})
	}
}
