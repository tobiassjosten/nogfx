package gmcp_test

import (
	"fmt"
	"testing"

	"github.com/tobiassjosten/nogfx/pkg/gmcp"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommClientMessages(t *testing.T) {
	tcs := []struct {
		message gmcp.ClientMessage
		output  string
	}{
		{
			message: gmcp.CommChannelEnable("newbie"),
			output:  `Comm.Channel.Enable "newbie"`,
		},
		{
			message: gmcp.CommChannelPlayers{},
			output:  "Comm.Channel.Players",
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			assert := assert.New(t)

			assert.Equal(tc.output, tc.message.String())
		})
	}
}

func TestCommServerMessages(t *testing.T) {
	tcs := []struct {
		command []byte
		message gmcp.ServerMessage
		err     string
	}{
		{
			command: []byte("Comm.Channel.End"),
			err:     "failed hydrating gmcp.CommChannelEnd: unexpected end of JSON input",
		},
		{
			command: []byte(`Comm.Channel.End "tell Jeremy"`),
			message: gmcp.CommChannelEnd("tell Jeremy"),
		},

		{
			command: []byte(`Comm.Channel.List`),
			err:     "failed hydrating gmcp.CommChannelList: unexpected end of JSON input",
		},
		{
			command: []byte(`Comm.Channel.List [}`),
			err:     `failed hydrating gmcp.CommChannelList: invalid character '}' looking for beginning of value`,
		},
		{
			command: []byte(`Comm.Channel.List ["asdf"]`),
			err:     `failed hydrating gmcp.CommChannelList: json: cannot unmarshal string into Go value of type gmcp.CommChannel`,
		},
		{
			command: []byte(`Comm.Channel.List [{"name":"ct", "caption":"Some city", "command":"ct"}]`),
			message: gmcp.CommChannelList{gmcp.CommChannel{
				Name:    "ct",
				Caption: "Some city",
				Command: "ct",
			}},
		},

		{
			command: []byte("Comm.Channel.Players"),
			err:     "failed hydrating gmcp.CommChannelPlayers: unexpected end of JSON input",
		},
		{
			command: []byte(`Comm.Channel.Players [}`),
			err:     `failed hydrating gmcp.CommChannelPlayers: invalid character '}' looking for beginning of value`,
		},
		{
			command: []byte(`Comm.Channel.Players ["asdf"]`),
			err:     `failed hydrating gmcp.CommChannelPlayers: json: cannot unmarshal string into Go value of type gmcp.CommChannelPlayer`,
		},
		{
			command: []byte(`Comm.Channel.Players [{"name": "Player1", "channels": ["Some city", "Some guild"]}, {"name": "Player2"}]`),
			message: gmcp.CommChannelPlayers{
				gmcp.CommChannelPlayer{
					Name:     "Player1",
					Channels: []string{"Some city", "Some guild"},
				},
				gmcp.CommChannelPlayer{
					Name: "Player2",
				},
			},
		},

		{
			command: []byte("Comm.Channel.Text"),
			err:     "failed hydrating gmcp.CommChannelText: unexpected end of JSON input",
		},
		{
			command: []byte(`Comm.Channel.Text [}`),
			err:     `failed hydrating gmcp.CommChannelText: invalid character '}' looking for beginning of value`,
		},
		{
			command: []byte(`Comm.Channel.Text []`),
			err:     `failed hydrating gmcp.CommChannelText: json: cannot unmarshal array into Go value of type gmcp.CommChannelText`,
		},
		{
			command: []byte(`Comm.Channel.Text ""`),
			err:     `failed hydrating gmcp.CommChannelText: json: cannot unmarshal string into Go value of type gmcp.CommChannelText`,
		},
		{
			command: []byte(`Comm.Channel.Text 1234`),
			err:     `failed hydrating gmcp.CommChannelText: json: cannot unmarshal number into Go value of type gmcp.CommChannelText`,
		},
		{
			command: []byte(`Comm.Channel.Text { "channel": "newbie", "talker": "Olad", "text": "\u001b[0;1;32m(Newbie): You say, \"Hello.\"\u001b[0;37m" }`),
			message: gmcp.CommChannelText{
				Channel: "newbie",
				Talker:  "Olad",
				Text:    "\u001b[0;1;32m(Newbie): You say, \"Hello.\"\u001b[0;37m",
			},
		},

		{
			command: []byte("Comm.Channel.Start"),
			err:     "failed hydrating gmcp.CommChannelStart: unexpected end of JSON input",
		},
		{
			command: []byte(`Comm.Channel.Start "ct"`),
			message: gmcp.CommChannelStart("ct"),
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			message, err := gmcp.Parse(tc.command, gmcp.ServerMessages)

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
