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

func TestWillWontDoDont(t *testing.T) {
	assert := assert.New(t)

	reader := bytes.NewReader([]byte{})
	writer := bytes.NewBuffer([]byte{})
	stream := &mockStream{reader, writer, nil}

	client, _ := telnet.NewClient(stream)

	tcs := []struct {
		verb byte
		noun byte
		f    func(byte) error
	}{
		{
			verb: telnet.WILL,
			noun: telnet.ECHO,
			f:    client.Will,
		},
		{
			verb: telnet.WONT,
			noun: telnet.ECHO,
			f:    client.Wont,
		},
		{
			verb: telnet.DO,
			noun: telnet.ECHO,
			f:    client.Do,
		},
		{
			verb: telnet.DONT,
			noun: telnet.ECHO,
			f:    client.Dont,
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			writer.Reset()
			_ = tc.f(tc.noun)
			assert.Equal([]byte{telnet.IAC, tc.verb, tc.noun}, writer.Bytes())
		})
	}
}

func TestSubneg(t *testing.T) {
	assert := assert.New(t)

	tcs := []struct {
		b     byte
		value []byte
		out   []byte
	}{
		{
			b:     telnet.TTYPE,
			value: []byte{},
			out: []byte{
				telnet.IAC, telnet.SB, telnet.TTYPE, 0,
				telnet.IAC, telnet.SE,
			},
		},
		{
			b:     telnet.TTYPE,
			value: []byte{'x', 'y'},
			out: []byte{
				telnet.IAC, telnet.SB, telnet.TTYPE, 1, 'x',
				'y', telnet.IAC, telnet.SE,
			},
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			reader := bytes.NewReader([]byte{})
			writer := bytes.NewBuffer([]byte{})
			stream := &mockStream{reader, writer, nil}

			client, _ := telnet.NewClient(stream)

			_ = client.Subneg(tc.b, tc.value)
			assert.Equal(tc.out, writer.Bytes())
		})
	}
}

func TestNegotiation(t *testing.T) {
	assert := assert.New(t)

	tcs := []struct {
		expectDo   byte
		expectDont byte
		expectWill byte
		expectWont byte
		serverDo   byte
		serverDont byte
		serverWill byte
		serverWont byte
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
			serverWill: telnet.MCCP2,
			expectDont: telnet.MCCP2,
		},
		{
			serverWill: 123,
			expectDont: 123,
		},
		{
			serverWont: telnet.ECHO,
			expectDont: telnet.ECHO,
		},
		{
			serverDo:   telnet.TTYPE,
			expectWont: telnet.TTYPE,
		},
		{
			serverDo:   124,
			expectWont: 124,
		},
		{
			serverDont: telnet.TTYPE,
			expectWont: telnet.TTYPE,
		},
		{
			serverDont: 124,
			expectWont: 124,
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			output := []byte{}
			if tc.serverWill > 0 {
				output = append(output, telnet.IAC, telnet.WILL, tc.serverWill)
			}
			if tc.serverWont > 0 {
				output = append(output, telnet.IAC, telnet.WONT, tc.serverWont)
			}
			if tc.serverDo > 0 {
				output = append(output, telnet.IAC, telnet.DO, tc.serverDo)
			}
			if tc.serverDont > 0 {
				output = append(output, telnet.IAC, telnet.DONT, tc.serverDont)
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
			stream := &mockStream{bytes.NewReader(output), builder, nil}

			client, commandChan := telnet.NewClient(stream)

			var commands [][]byte
			go func(commandChan <-chan []byte) {
				for command := range commandChan {
					commands = append(commands, command)
				}
			}(commandChan)

			_, err := ioutil.ReadAll(client)

			if tc.err != nil {
				assert.Equal(tc.err, err)
				return
			}

			assert.Equal(string(input), builder.String())
		})
	}
}
