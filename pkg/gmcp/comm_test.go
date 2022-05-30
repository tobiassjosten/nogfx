package gmcp_test

import (
	"strings"
	"testing"

	"github.com/tobiassjosten/nogfx/pkg/gmcp"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommMessages(t *testing.T) {
	tcs := map[string]struct {
		msg         gmcp.Message
		data        string
		unmarshaled gmcp.Message
		marshaled   string
		err         string
	}{
		"Comm.Channel.Enable empty": {
			msg:         &gmcp.CommChannelEnable{},
			data:        `Comm.Channel.Enable ""`,
			unmarshaled: &gmcp.CommChannelEnable{},
			marshaled:   `Comm.Channel.Enable ""`,
		},

		"Comm.Channel.Enable hydrated": {
			msg:  &gmcp.CommChannelEnable{},
			data: `Comm.Channel.Enable "asdf"`,
			unmarshaled: &gmcp.CommChannelEnable{
				Channel: "asdf",
			},
			marshaled: `Comm.Channel.Enable "asdf"`,
		},

		"Comm.Channel.Enable invalid JSON": {
			msg:  &gmcp.CommChannelEnable{},
			data: `Comm.Channel.Enable asdf`,
			err:  "invalid character 'a' looking for beginning of value",
		},

		"Comm.Channel.End empty": {
			msg:         &gmcp.CommChannelEnd{},
			data:        `Comm.Channel.End ""`,
			unmarshaled: &gmcp.CommChannelEnd{},
			marshaled:   `Comm.Channel.End ""`,
		},

		"Comm.Channel.End hydrated": {
			msg:  &gmcp.CommChannelEnd{},
			data: `Comm.Channel.End "asdf"`,
			unmarshaled: &gmcp.CommChannelEnd{
				Channel: "asdf",
			},
			marshaled: `Comm.Channel.End "asdf"`,
		},

		"Comm.Channel.End invalid JSON": {
			msg:  &gmcp.CommChannelEnd{},
			data: `Comm.Channel.End asdf`,
			err:  "invalid character 'a' looking for beginning of value",
		},

		"Comm.Channel.List empty": {
			msg:         &gmcp.CommChannelList{},
			data:        "Comm.Channel.List []",
			unmarshaled: &gmcp.CommChannelList{},
			marshaled:   "Comm.Channel.List []",
		},

		"Comm.Channel.List hydrated": {
			msg: &gmcp.CommChannelList{},
			data: makeGMCP("Comm.Channel.List", []map[string]interface{}{
				{
					"name":    "ct",
					"caption": "Some city",
					"command": "ct",
				},
			}),
			unmarshaled: &gmcp.CommChannelList{
				{
					Name:    "ct",
					Caption: "Some city",
					Command: "ct",
				},
			},
			marshaled: makeGMCP("Comm.Channel.List", []map[string]interface{}{
				{
					"name":    "ct",
					"caption": "Some city",
					"command": "ct",
				},
			}),
		},

		"Comm.Channel.Players empty": {
			msg:         &gmcp.CommChannelPlayers{},
			data:        "Comm.Channel.Players []",
			unmarshaled: &gmcp.CommChannelPlayers{},
			marshaled:   "Comm.Channel.Players []",
		},

		"Comm.Channel.Players hydrated": {
			msg: &gmcp.CommChannelPlayers{},
			data: makeGMCP("Comm.Channel.Players", []map[string]interface{}{
				{
					"name":     "Durak",
					"channels": []string{"Some city"},
				},
			}),
			unmarshaled: &gmcp.CommChannelPlayers{
				{
					Name:     "Durak",
					Channels: []string{"Some city"},
				},
			},
			marshaled: makeGMCP("Comm.Channel.Players", []map[string]interface{}{
				{
					"name":     "Durak",
					"channels": []string{"Some city"},
				},
			}),
		},

		"Comm.Channel.Start empty": {
			msg:         &gmcp.CommChannelStart{},
			data:        `Comm.Channel.Start ""`,
			unmarshaled: &gmcp.CommChannelStart{},
			marshaled:   `Comm.Channel.Start ""`,
		},

		"Comm.Channel.Start hydrated": {
			msg:  &gmcp.CommChannelStart{},
			data: `Comm.Channel.Start "asdf"`,
			unmarshaled: &gmcp.CommChannelStart{
				Channel: "asdf",
			},
			marshaled: `Comm.Channel.Start "asdf"`,
		},

		"Comm.Channel.Start invalid JSON": {
			msg:  &gmcp.CommChannelStart{},
			data: `Comm.Channel.Start asdf`,
			err:  "invalid character 'a' looking for beginning of value",
		},

		"Comm.Channel.Text empty": {
			msg:         &gmcp.CommChannelText{},
			data:        `Comm.Channel.Text {}`,
			unmarshaled: &gmcp.CommChannelText{},
			marshaled: makeGMCP("Comm.Channel.Text", map[string]interface{}{
				"channel": "",
				"talker":  "",
				"text":    "",
			}),
		},

		"Comm.Channel.Text hydrated": {
			msg: &gmcp.CommChannelText{},
			data: makeGMCP("Comm.Channel.Text", map[string]interface{}{
				"channel": "ct",
				"talker":  "Durak",
				"text":    `(Somecity): Durak says, "Yo!"`,
			}),
			unmarshaled: &gmcp.CommChannelText{
				Channel: "ct",
				Talker:  "Durak",
				Text:    `(Somecity): Durak says, "Yo!"`,
			},
			marshaled: makeGMCP("Comm.Channel.Text", map[string]interface{}{
				"channel": "ct",
				"talker":  "Durak",
				"text":    `(Somecity): Durak says, "Yo!"`,
			}),
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			err := tc.msg.Unmarshal([]byte(tc.data))

			if tc.err != "" {
				require.NotNil(t, err)
				assert.Equal(t, tc.err, err.Error())
				return
			} else if err != nil {
				require.Equal(t, "", err.Error())
			}

			require.Equal(t, tc.unmarshaled, tc.msg, "unmarshaling hydrates message")

			if tc.marshaled == "" {
				return
			}

			marshaled := tc.msg.Marshal()
			data := strings.TrimSpace(strings.TrimPrefix(marshaled, tc.msg.ID()))
			tcdata := strings.TrimSpace(strings.TrimPrefix(tc.marshaled, tc.msg.ID()))

			assert.NotEqual(t, marshaled, data, "marshaled data has ID prefix")
			assert.NotEqual(t, tc.marshaled, tcdata, "marshaled data has ID prefix")

			if tcdata == "" {
				assert.Equal(t, tcdata, data)
				return
			}

			assert.JSONEq(t, tcdata, data, "marshaling maintains data integrity")

			require.Equal(t, tc.unmarshaled, tc.msg, "marshaling doesn't mutate")
		})
	}
}
