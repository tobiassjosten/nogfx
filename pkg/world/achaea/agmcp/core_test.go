package agmcp_test

import (
	"fmt"
	"testing"

	"github.com/tobiassjosten/nogfx/pkg/gmcp"
	"github.com/tobiassjosten/nogfx/pkg/world/achaea/agmcp"

	"github.com/icza/gox/gox"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCoreClientMessages(t *testing.T) {
	tcs := []struct {
		message gmcp.ClientMessage
		output  string
	}{
		{
			message: gmcp.CoreHello{
				Client:  "nogfx",
				Version: "1.2.3",
			},
			output: `Core.Hello {"client":"nogfx","version":"1.2.3"}`,
		},
		{
			message: gmcp.CoreKeepAlive{},
			output:  "Core.KeepAlive",
		},
		{
			message: gmcp.CorePing{},
			output:  "Core.Ping",
		},
		{
			message: gmcp.CorePing{Latency: gox.NewInt(120)},
			output:  "Core.Ping 120",
		},
		{
			message: gmcp.CoreSupportsSet{},
			output:  "Core.Supports.Set []",
		},
		{
			message: agmcp.CoreSupportsSet{
				agmcp.CoreSupports{
					CoreSupports: gmcp.CoreSupports{
						Char:        gox.NewInt(1),
						CharSkills:  gox.NewInt(2),
						CharItems:   gox.NewInt(3),
						CommChannel: gox.NewInt(4),
						Room:        gox.NewInt(5),
					},
					IRERift: gox.NewInt(6),
				},
			},
			output: `Core.Supports.Set ["Char 1","Char.Skills 2","Char.Items 3","Comm.Channel 4","Room 5","IRE.Rift 6"]`,
		},
		{
			message: gmcp.CoreSupportsAdd{},
			output:  "Core.Supports.Add []",
		},
		{
			message: agmcp.CoreSupportsAdd{
				agmcp.CoreSupports{
					CoreSupports: gmcp.CoreSupports{
						Char:        gox.NewInt(1),
						CharSkills:  gox.NewInt(2),
						CharItems:   gox.NewInt(3),
						CommChannel: gox.NewInt(4),
						Room:        gox.NewInt(5),
					},
					IRERift: gox.NewInt(6),
				},
			},
			output: `Core.Supports.Add ["Char 1","Char.Skills 2","Char.Items 3","Comm.Channel 4","Room 5","IRE.Rift 6"]`,
		},
		{
			message: gmcp.CoreSupportsRemove{},
			output:  "Core.Supports.Remove []",
		},
		{
			message: agmcp.CoreSupportsRemove{
				agmcp.CoreSupports{
					CoreSupports: gmcp.CoreSupports{
						Char:        gox.NewInt(1),
						CharSkills:  gox.NewInt(2),
						CharItems:   gox.NewInt(3),
						CommChannel: gox.NewInt(4),
						Room:        gox.NewInt(5),
					},
					IRERift: gox.NewInt(6),
				},
			},
			output: `Core.Supports.Remove ["Char 1","Char.Skills 2","Char.Items 3","Comm.Channel 4","Room 5","IRE.Rift 6"]`,
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			assert := assert.New(t)

			assert.Equal(tc.output, tc.message.String())
		})
	}
}

func TestCoreServerMessages(t *testing.T) {
	tcs := []struct {
		command []byte
		message gmcp.ServerMessage
	}{
		{
			command: []byte("Core.Goodbye"),
			message: gmcp.CoreGoodbye{},
		},
		{
			command: []byte("Core.Ping"),
			message: gmcp.CorePing{},
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			message, err := gmcp.Parse(tc.command)

			require.Nil(err)
			assert.Equal(tc.message, message)
		})
	}
}
