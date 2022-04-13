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
		expectDo   byte
		expectDont byte
		expectWill byte
		expectWont byte
		serverDo   byte
		serverWill byte
		err        error
	}{
		{
			serverWill: telnet.ECHO,
			expectDo:   telnet.ECHO,
		},
		{
			serverWill: telnet.GMCP,
			expectDo:   telnet.GMCP,
		},
		{
			serverWill: telnet.ATCP,
			expectDont: telnet.ATCP,
		},
		{
			serverWill: telnet.MCCP,
			expectDont: telnet.MCCP,
		},
		{
			serverWill: telnet.MCCP2,
			expectDont: telnet.MCCP2,
		},
		{
			serverWill: 123,
			expectDont: 123,
		},
		{
			serverDo:   telnet.TTYPE,
			expectWont: telnet.TTYPE,
		},
		{
			serverDo:   124,
			expectWont: 124,
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
			if tc.expectWill > 0 {
				input = append(input, telnet.IAC, telnet.WILL, tc.expectWill)
			}
			if tc.expectWont > 0 {
				input = append(input, telnet.IAC, telnet.WONT, tc.expectWont)
			}
			if tc.expectDo > 0 {
				input = append(input, telnet.IAC, telnet.DO, tc.expectDo)
			}
			if tc.expectDont > 0 {
				input = append(input, telnet.IAC, telnet.DONT, tc.expectDont)
			}

			builder := &strings.Builder{}
			stream := &MockStream{
				bytes.NewReader(output),
				builder,
			}

			client, commandChan := telnet.NewClient(stream)

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
