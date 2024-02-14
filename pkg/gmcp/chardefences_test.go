package gmcp_test

import (
	"strings"
	"testing"

	"github.com/tobiassjosten/nogfx/pkg/gmcp"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCharDefencesMessages(t *testing.T) {
	tcs := map[string]struct {
		msg         gmcp.Message
		data        string
		unmarshaled gmcp.Message
		marshaled   string
		err         string
	}{
		"Char.Defences.List empty": {
			msg:         &gmcp.CharDefencesList{},
			data:        "Char.Defences.List []",
			unmarshaled: &gmcp.CharDefencesList{},
			marshaled:   "Char.Defences.List []",
		},

		"Char.Defences.List hydrated": {
			msg: &gmcp.CharDefencesList{},
			data: makeGMCP("Char.Defences.List", []map[string]string{
				{
					"name": "Name1",
					"cure": "Cure1",
					"desc": "Desc1",
				},
				{
					"name": "Name2",
					"cure": "Cure2",
					"desc": "Desc2",
				},
			}),
			unmarshaled: &gmcp.CharDefencesList{
				{
					Name:        "Name1",
					Cure:        "Cure1",
					Description: "Desc1",
				},
				{
					Name:        "Name2",
					Cure:        "Cure2",
					Description: "Desc2",
				},
			},
			marshaled: makeGMCP("Char.Defences.List", []map[string]string{
				{
					"name": "Name1",
					"cure": "Cure1",
					"desc": "Desc1",
				},
				{
					"name": "Name2",
					"cure": "Cure2",
					"desc": "Desc2",
				},
			}),
		},

		"Char.Defences.List invalid JSON": {
			msg:  &gmcp.CharDefencesList{},
			data: "asdf",
			err:  "invalid character 'a' looking for beginning of value",
		},

		"Char.Defences.Add empty": {
			msg:         &gmcp.CharDefencesAdd{},
			data:        "Char.Defences.Add {}",
			unmarshaled: &gmcp.CharDefencesAdd{},
			marshaled: makeGMCP("Char.Defences.Add", map[string]string{
				"name": "",
				"cure": "",
				"desc": "",
			}),
		},

		"Char.Defences.Add hydrated": {
			msg: &gmcp.CharDefencesAdd{},
			data: makeGMCP("Char.Defences.Add", map[string]string{
				"name": "Name1",
				"cure": "Cure1",
				"desc": "Desc1",
			}),
			unmarshaled: &gmcp.CharDefencesAdd{
				Name:        "Name1",
				Cure:        "Cure1",
				Description: "Desc1",
			},
			marshaled: makeGMCP("Char.Defences.Add", map[string]string{
				"name": "Name1",
				"cure": "Cure1",
				"desc": "Desc1",
			}),
		},

		"Char.Defences.Add invalid JSON": {
			msg:  &gmcp.CharDefencesAdd{},
			data: "asdf",
			err:  "invalid character 'a' looking for beginning of value",
		},

		"Char.Defences.Remove empty": {
			msg:         &gmcp.CharDefencesRemove{},
			data:        "Char.Defences.Remove []",
			unmarshaled: &gmcp.CharDefencesRemove{},
			marshaled:   "Char.Defences.Remove []",
		},

		"Char.Defences.Remove hydrated": {
			msg: &gmcp.CharDefencesRemove{},
			data: makeGMCP("Char.Defences.Remove", []string{
				"Name1",
			}),
			unmarshaled: &gmcp.CharDefencesRemove{
				{
					Name:        "Name1",
					Cure:        "",
					Description: "",
				},
			},
			marshaled: makeGMCP("Char.Defences.Remove", []string{
				"Name1",
			}),
		},

		"Char.Defences.Remove invalid JSON": {
			msg:  &gmcp.CharDefencesRemove{},
			data: "asdf",
			err:  "invalid character 'a' looking for beginning of value",
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

			assert.JSONEq(t, tcdata, data, "marshaling maintains data integrity")

			require.Equal(t, tc.unmarshaled, tc.msg, "marshaling doesn't mutate")
		})
	}
}
