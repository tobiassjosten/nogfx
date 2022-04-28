package gmcp_test

import (
	"fmt"
	"testing"

	"github.com/tobiassjosten/nogfx/pkg/gmcp"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIREClientMessages(t *testing.T) {
	tcs := []struct {
		message gmcp.ClientMessage
		output  string
	}{
		{
			message: gmcp.IRERiftRequest{},
			output:  "IRE.Rift.Request",
		},

		{
			message: gmcp.IRETargetSet(""),
			output:  "IRE.Target.Set",
		},
		{
			message: gmcp.IRETargetSet("asdf"),
			output:  `IRE.Target.Set "asdf"`,
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			assert := assert.New(t)

			assert.Equal(tc.output, tc.message.String())
		})
	}
}

func TestIREServerMessages(t *testing.T) {
	tcs := []struct {
		command []byte
		message gmcp.ServerMessage
		err     string
	}{
		{
			command: []byte("IRE.Rift.Change"),
			err:     "failed hydrating gmcp.IRERiftChange: unexpected end of JSON input",
		},
		{
			command: []byte(`IRE.Rift.Change { "name": "rawstone", "amount": "1", "desc": "rawstone" }`),
			message: gmcp.IRERiftChange{
				Name:        "rawstone",
				Amount:      1,
				Description: "rawstone",
			},
		},

		{
			command: []byte("IRE.Rift.List"),
			err:     "failed hydrating gmcp.IRERiftList: unexpected end of JSON input",
		},
		{
			command: []byte(`IRE.Rift.List [ { "name": "rawstone", "amount": "1", "desc": "rawstone" } ]`),
			message: gmcp.IRERiftList{{
				Name:        "rawstone",
				Amount:      1,
				Description: "rawstone",
			}},
		},

		{
			command: []byte("IRE.Target.Set"),
			err:     "failed hydrating gmcp.IRETargetSet: unexpected end of JSON input",
		},
		{
			command: []byte(`IRE.Target.Set "asdf"`),
			message: gmcp.IRETargetSet("asdf"),
		},

		{
			command: []byte("IRE.Target.Info"),
			err:     "failed hydrating *gmcp.IRETargetInfo: unexpected end of JSON input",
		},
		{
			command: []byte(`IRE.Target.Info { "hpperc": "asdf" }`),
			err:     `failed hydrating *gmcp.IRETargetInfo: strconv.Atoi: parsing "asdf": invalid syntax`,
		},
		{
			command: []byte(`IRE.Target.Info { "id": "266744", "hpperc": "79%", "short_desc": "a practice dummy" }`),
			message: gmcp.IRETargetInfo{
				ID:          "266744",
				Health:      79,
				Description: "a practice dummy",
			},
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			message, err := gmcp.Parse(tc.command)

			if tc.err != "" {
				require.NotNil(err, fmt.Sprintf(
					"wanted: %s", tc.err,
				))
				assert.Equal(tc.err, err.Error())
				return
			}

			require.Nil(err)
			assert.Equal(tc.message, message)
		})
	}
}
