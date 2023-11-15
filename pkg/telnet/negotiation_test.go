package telnet_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/tobiassjosten/nogfx/pkg/telnet"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsCommand(t *testing.T) {
	tcs := map[string]struct {
		data    []byte
		verdict bool
	}{
		"empty": {
			data:    []byte{},
			verdict: false,
		},

		"a": {
			data:    []byte("a"),
			verdict: false,
		},

		"as": {
			data:    []byte("as"),
			verdict: false,
		},

		"asd": {
			data:    []byte("asd"),
			verdict: false,
		},

		"asdf": {
			data:    []byte("asdf"),
			verdict: false,
		},

		"iac iac": {
			data:    []byte{telnet.IAC, telnet.IAC},
			verdict: true,
		},

		"iac ga": {
			data:    []byte{telnet.IAC, telnet.GA},
			verdict: true,
		},

		"iac will echo": {
			data:    []byte{telnet.IAC, telnet.Will, telnet.Echo},
			verdict: true,
		},

		"iac will a": {
			data:    []byte{telnet.IAC, telnet.Will, 'a'},
			verdict: true,
		},

		"iac will echo a": {
			data:    []byte{telnet.IAC, telnet.Will, telnet.Echo, 'a'},
			verdict: false,
		},

		"iac wont echo": {
			data:    []byte{telnet.IAC, telnet.Wont, telnet.Echo},
			verdict: true,
		},

		"iac wont a": {
			data:    []byte{telnet.IAC, telnet.Wont, 'a'},
			verdict: true,
		},

		"iac wont echo a": {
			data:    []byte{telnet.IAC, telnet.Wont, telnet.Echo, 'a'},
			verdict: false,
		},

		"iac do echo": {
			data:    []byte{telnet.IAC, telnet.Do, telnet.Echo},
			verdict: true,
		},

		"iac do a": {
			data:    []byte{telnet.IAC, telnet.Do, 'a'},
			verdict: true,
		},

		"iac do echo a": {
			data:    []byte{telnet.IAC, telnet.Do, telnet.Echo, 'a'},
			verdict: false,
		},

		"iac dont echo": {
			data:    []byte{telnet.IAC, telnet.Dont, telnet.Echo},
			verdict: true,
		},

		"iac dont a": {
			data:    []byte{telnet.IAC, telnet.Dont, 'a'},
			verdict: true,
		},

		"iac dont echo a": {
			data:    []byte{telnet.IAC, telnet.Dont, telnet.Echo, 'a'},
			verdict: false,
		},

		"iac 239 echo": {
			data:    []byte{telnet.IAC, 239, telnet.Echo},
			verdict: false,
		},

		"iac 239 a": {
			data:    []byte{telnet.IAC, 239, 'a'},
			verdict: false,
		},

		"sub-negotiation empty": {
			data:    []byte{telnet.IAC, telnet.SB, telnet.IAC, telnet.SE},
			verdict: true,
		},

		"sub-negotiation complete": {
			data:    []byte{telnet.IAC, telnet.SB, 'a', telnet.IAC, telnet.SE},
			verdict: true,
		},

		"sub-negotiation unterminated": {
			data:    []byte{telnet.IAC, telnet.SB, 'a', telnet.IAC},
			verdict: false,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.verdict, telnet.IsCommand(tc.data))
		})
	}
}

func TestNegotiate(t *testing.T) {
	tcs := map[string]struct {
		serverOutput []byte
		clientInput  []byte
		writeErr     func([]byte) error
	}{
		"enable our echo": {
			serverOutput: []byte{telnet.IAC, telnet.Do, telnet.Echo},
			clientInput:  []byte{telnet.IAC, telnet.Will, telnet.Echo},
		},

		"enable-fail our echo": {
			serverOutput: []byte{telnet.IAC, telnet.Do, telnet.Echo},
			writeErr:     func([]byte) error { return errMock },
		},

		"disable our echo": {
			serverOutput: []byte{telnet.IAC, telnet.Dont, telnet.Echo},
			clientInput:  nil,
		},

		"enable disable our echo": {
			serverOutput: []byte{
				telnet.IAC, telnet.Do, telnet.Echo,
				telnet.IAC, telnet.Dont, telnet.Echo,
			},
			clientInput: []byte{
				telnet.IAC, telnet.Will, telnet.Echo,
				telnet.IAC, telnet.Wont, telnet.Echo,
			},
		},

		"enable disable-fail our echo": {
			serverOutput: []byte{
				telnet.IAC, telnet.Do, telnet.Echo,
				telnet.IAC, telnet.Dont, telnet.Echo,
			},
			clientInput: []byte{
				telnet.IAC, telnet.Will, telnet.Echo,
			},
			writeErr: func(cmd []byte) error {
				if bytes.Equal(cmd, []byte{telnet.IAC, telnet.Wont, telnet.Echo}) {
					return errMock
				}
				return nil
			},
		},

		"enable their echo": {
			serverOutput: []byte{telnet.IAC, telnet.Will, telnet.Echo},
			clientInput:  []byte{telnet.IAC, telnet.Dont, telnet.Echo},
		},

		"enable-fail their echo": {
			serverOutput: []byte{telnet.IAC, telnet.Will, telnet.Echo},
			writeErr:     func([]byte) error { return errMock },
		},

		"disable their echo": {
			serverOutput: []byte{telnet.IAC, telnet.Wont, telnet.Echo},
			clientInput:  nil,
		},

		"enable our suppress-go-ahead": {
			serverOutput: []byte{telnet.IAC, telnet.Do, telnet.SuppressGoAhead},
			clientInput:  []byte{telnet.IAC, telnet.Wont, telnet.SuppressGoAhead},
		},

		"enable-fail our suppress-go-ahead": {
			serverOutput: []byte{telnet.IAC, telnet.Do, telnet.SuppressGoAhead},
			writeErr:     func([]byte) error { return errMock },
		},

		"disable our suppress-go-ahead": {
			serverOutput: []byte{telnet.IAC, telnet.Dont, telnet.SuppressGoAhead},
			clientInput:  nil,
		},

		"enable their suppress-go-ahead": {
			serverOutput: []byte{telnet.IAC, telnet.Will, telnet.SuppressGoAhead},
			clientInput:  []byte{telnet.IAC, telnet.Do, telnet.SuppressGoAhead},
		},

		"enable-fail their suppress-go-ahead": {
			serverOutput: []byte{telnet.IAC, telnet.Will, telnet.SuppressGoAhead},
			writeErr:     func([]byte) error { return errMock },
		},

		"disable their suppress-go-ahead": {
			serverOutput: []byte{telnet.IAC, telnet.Wont, telnet.SuppressGoAhead},
			clientInput:  nil,
		},

		"enable disable their suppress-go-ahead": {
			serverOutput: []byte{
				telnet.IAC, telnet.Will, telnet.SuppressGoAhead,
				telnet.IAC, telnet.Wont, telnet.SuppressGoAhead,
			},
			clientInput: []byte{
				telnet.IAC, telnet.Do, telnet.SuppressGoAhead,
				telnet.IAC, telnet.Dont, telnet.SuppressGoAhead,
			},
		},

		"enable disable-fail their suppress-go-ahead": {
			serverOutput: []byte{
				telnet.IAC, telnet.Will, telnet.SuppressGoAhead,
				telnet.IAC, telnet.Wont, telnet.SuppressGoAhead,
			},
			clientInput: []byte{
				telnet.IAC, telnet.Do, telnet.SuppressGoAhead,
			},
			writeErr: func(cmd []byte) error {
				if bytes.Equal(cmd, []byte{telnet.IAC, telnet.Dont, telnet.SuppressGoAhead}) {
					return errMock
				}
				return nil
			},
		},

		"enable our nonsense": {
			serverOutput: []byte{telnet.IAC, telnet.Do, 'a'},
			clientInput:  []byte{telnet.IAC, telnet.Wont, 'a'},
		},

		"enable-fail our nonsense": {
			serverOutput: []byte{telnet.IAC, telnet.Do, 'a'},
			writeErr:     func([]byte) error { return errMock },
		},

		"disable our nonsense": {
			serverOutput: []byte{telnet.IAC, telnet.Dont, 'a'},
			clientInput:  nil,
		},

		"enable their nonsense": {
			serverOutput: []byte{telnet.IAC, telnet.Will, 'a'},
			clientInput:  []byte{telnet.IAC, telnet.Dont, 'a'},
		},

		"enable-fail their nonsense": {
			serverOutput: []byte{telnet.IAC, telnet.Will, 'a'},
			writeErr:     func([]byte) error { return errMock },
		},

		"disable their nonsense": {
			serverOutput: []byte{telnet.IAC, telnet.Wont, 'a'},
			clientInput:  nil,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			conn := NewMockConn(tc.serverOutput)
			conn.WriteErr = tc.writeErr

			client := telnet.NewNVT(conn)

			_, err := io.ReadAll(client)

			if tc.writeErr != nil {
				require.ErrorIs(t, err, errMock)
			} else {
				require.Nil(t, err, err)
			}

			assert.Equal(t, tc.clientInput, conn.Written)
		})
	}
}
