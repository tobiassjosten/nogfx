package gmcp_test

import (
	"fmt"
	"testing"

	"github.com/tobiassjosten/nogfx/pkg/gmcp"
	"github.com/tobiassjosten/nogfx/pkg/telnet"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tcs := []struct {
		command  []byte
		messages map[string]gmcp.ServerMessage
		err      error
	}{
		{
			command:  []byte("Core.Ping"),
			messages: gmcp.ServerMessages,
		},
		{
			command:  []byte("Core.Ping"),
			messages: map[string]gmcp.ServerMessage{},
		},
		{
			command:  []byte("Asdf"),
			messages: gmcp.ServerMessages,
			err:      fmt.Errorf("unknown message 'Asdf'"),
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			_, err := gmcp.Parse(tc.command, tc.messages)

			if tc.err != nil {
				assert.Equal(tc.err, err)
				return
			}

			require.Nil(err)
		})
	}
}

func TestWrapUnwrap(t *testing.T) {
	tcs := []struct {
		unwrapped []byte
		wrapped   []byte
		out       []byte
	}{
		{
			unwrapped: []byte("Core.Ping"),
			out: append(
				[]byte{telnet.IAC, telnet.SB, telnet.GMCP},
				append(
					[]byte("Core.Ping"),
					telnet.IAC, telnet.SE,
				)...,
			),
		},
		{
			wrapped: append(
				[]byte{telnet.IAC, telnet.SB, telnet.GMCP},
				append(
					[]byte("Valid"),
					telnet.IAC, telnet.SE,
				)...,
			),
			out: []byte("Valid"),
		},
		{
			wrapped: []byte("Invalid"),
			out:     nil,
		},
		{
			wrapped: append(
				[]byte{telnet.IAC, telnet.SB, telnet.GMCP},
				[]byte("Invalid")...,
			),
			out: nil,
		},
		{
			wrapped: append(
				[]byte("Invalid"),
				telnet.IAC, telnet.SE,
			),
			out: nil,
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			if tc.unwrapped != nil {
				out := gmcp.Wrap(tc.unwrapped)
				assert.Equal(t, tc.out, out)
				assert.Equal(t, tc.unwrapped, gmcp.Unwrap(out))
			}

			if tc.wrapped != nil {
				assert.Equal(t, tc.out, gmcp.Unwrap(tc.wrapped))
			}
		})
	}
}
