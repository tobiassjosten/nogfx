package world_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/tobiassjosten/nogfx/pkg/mock"
	"github.com/tobiassjosten/nogfx/pkg/telnet"
	"github.com/tobiassjosten/nogfx/pkg/world"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	gmcpPrefix = []byte{telnet.IAC, telnet.SB, telnet.GMCP}
	gmcpSuffix = []byte{telnet.IAC, telnet.SE}
	willGMCP   = []byte{telnet.IAC, telnet.WILL, telnet.GMCP}
	willEcho   = []byte{telnet.IAC, telnet.WILL, telnet.ECHO}
	wontEcho   = []byte{telnet.IAC, telnet.WONT, telnet.ECHO}
)

func wrapGMCP(msgs []string) []byte {
	var bs []byte
	for _, msg := range msgs {
		bs = append(bs, gmcpPrefix...)
		bs = append(bs, []byte(msg)...)
		bs = append(bs, gmcpSuffix...)
	}
	return bs
}

func TestCommandsReply(t *testing.T) {
	tcs := []struct {
		command []byte
		sent    []byte
		errs    []bool
		err     string
	}{
		{
			command: willGMCP,
			sent: wrapGMCP([]string{
				`Core.Hello {"client":"nogfx","version":"0.0.0"}`,
			}),
		},
		{
			command: []byte{telnet.IAC, telnet.WILL, telnet.GMCP},
			errs:    []bool{true},
			err:     "failed GMCP: ooops",
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			var calls int
			var sent []byte

			client := &mock.ClientMock{
				WriteFunc: func(data []byte) (int, error) {
					defer func() { calls++ }()
					if len(tc.errs) > calls && tc.errs[calls] {
						return 0, errors.New("ooops")
					}

					sent = append(sent, data...)

					return len(data), nil
				},
			}

			ui := &mock.UIMock{
				AddVitalFunc: func(_ string, _ interface{}) {},
			}

			engine := world.NewEngine(client, ui, "example.com:23")

			err := engine.ProcessCommand(tc.command)

			if tc.err != "" && assert.NotNil(t, err) {
				assert.Equal(t, tc.err, err.Error())
				return
			}

			require.Nil(t, err)

			if len(tc.sent) > 0 {
				assert.Equal(t, tc.sent, sent, string(sent))
			}
		})
	}
}

func TestMasking(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	client := &mock.ClientMock{}

	var masked bool

	ui := &mock.UIMock{
		AddVitalFunc: func(_ string, _ interface{}) {},
		MaskInputFunc: func() {
			masked = true
		},
		UnmaskInputFunc: func() {
			masked = false
		},
	}

	engine := world.NewEngine(client, ui, "example.com:23")

	err := engine.ProcessCommand(willEcho)
	require.Nil(err)

	assert.Equal(true, masked)

	err = engine.ProcessCommand(wontEcho)
	require.Nil(err)

	assert.Equal(false, masked)
}
