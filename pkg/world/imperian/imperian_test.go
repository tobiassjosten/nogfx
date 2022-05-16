package imperian

import (
	"errors"
	"fmt"
	"testing"

	"github.com/tobiassjosten/nogfx/pkg/mock"
	"github.com/tobiassjosten/nogfx/pkg/telnet"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func wrapGMCP(msgs []string) []byte {
	var bs []byte
	for _, msg := range msgs {
		bs = append(bs, telnet.IAC, telnet.SB, telnet.GMCP)
		bs = append(bs, []byte(msg)...)
		bs = append(bs, telnet.IAC, telnet.SE)
	}
	return bs
}

func TestWorldBasics(t *testing.T) {
	assert := assert.New(t)

	client := &mock.ClientMock{}
	ui := &mock.UIMock{
		AddVitalFunc: func(_ string, _ interface{}) error {
			return nil
		},
	}

	world := NewWorld(client, ui)

	input := []byte("input")
	assert.Equal(input, world.ProcessInput(input))

	output := []byte("output")
	assert.Equal(output, world.ProcessOutput(output))
}

func TestCommandsReply(t *testing.T) {
	tcs := []struct {
		command []byte
		sent    []byte
		errs    []bool
		err     string
	}{
		{
			command: []byte{telnet.IAC, telnet.WILL, telnet.GMCP},
			sent: wrapGMCP([]string{
				`Core.Supports.Set ["Char 1","Char.Skills 1","Char.Items 1","Comm.Channel 1","Room 1","IRE.Rift 1"]`,
			}),
		},
		{
			command: []byte{telnet.IAC, telnet.WILL, telnet.GMCP},
			errs:    []bool{true},
			err:     "failed GMCP: ooops",
		},

		{
			command: wrapGMCP([]string{`Asdf.Qwer`}),
			err:     "failed parsing GMCP: unknown message 'Asdf.Qwer'",
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
				AddVitalFunc: func(_ string, _ interface{}) error {
					return nil
				},
			}

			world := NewWorld(client, ui)

			err := world.ProcessCommand(tc.command)

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

func TestCommandsMutateVitals(t *testing.T) {
	tcs := []struct {
		command []byte
		vitals  map[string][][]int
	}{
		{
			command: wrapGMCP([]string{`Char.Vitals { "hp": "3904", "maxhp": "3905", "mp": "3845", "maxmp": "3846" }`}),
			vitals: map[string][][]int{
				"health": [][]int{{3904 / 11, 3905 / 11}},
				"mana":   [][]int{{3845 / 11, 3846 / 11}},
			},
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			vitals := map[string][][]int{}
			ui := &mock.UIMock{
				AddVitalFunc: func(_ string, _ interface{}) error {
					return nil
				},
				UpdateVitalFunc: func(name string, value, max int) error {
					if _, ok := vitals[name]; !ok {
						vitals[name] = [][]int{}
					}
					vitals[name] = append(vitals[name], []int{value, max})
					return nil
				},
			}

			world := NewWorld(&mock.ClientMock{}, ui)

			err := world.ProcessCommand(tc.command)
			require.Nil(t, err)

			assert.Equal(t, tc.vitals, vitals)
		})
	}
}
