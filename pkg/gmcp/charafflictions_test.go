package gmcp_test

import (
	"strings"
	"testing"

	"github.com/tobiassjosten/nogfx/pkg/gmcp"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCharAfflictionsMessages(t *testing.T) {
	tcs := map[string]struct {
		msg         gmcp.Message
		data        string
		unmarshaled gmcp.Message
		marshaled   string
		err         string
	}{
		"Char.Afflictions.List empty": {
			msg:         &gmcp.CharAfflictionsList{},
			data:        "Char.Afflictions.List []",
			unmarshaled: &gmcp.CharAfflictionsList{},
			marshaled:   "Char.Afflictions.List []",
		},

		"Char.Afflictions.List hydrated": {
			msg: &gmcp.CharAfflictionsList{},
			data: makeGMCP("Char.Afflictions.List", []map[string]string{
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
			unmarshaled: &gmcp.CharAfflictionsList{
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
			marshaled: makeGMCP("Char.Afflictions.List", []map[string]string{
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

		"Char.Afflictions.List invalid JSON": {
			msg:  &gmcp.CharAfflictionsList{},
			data: "asdf",
			err:  "invalid character 'a' looking for beginning of value",
		},

		"Char.Afflictions.Add empty": {
			msg:         &gmcp.CharAfflictionsAdd{},
			data:        "Char.Afflictions.Add {}",
			unmarshaled: &gmcp.CharAfflictionsAdd{},
			marshaled: makeGMCP("Char.Afflictions.Add", map[string]string{
				"name": "",
				"cure": "",
				"desc": "",
			}),
		},

		"Char.Afflictions.Add hydrated": {
			msg: &gmcp.CharAfflictionsAdd{},
			data: makeGMCP("Char.Afflictions.Add", map[string]string{
				"name": "Name1",
				"cure": "Cure1",
				"desc": "Desc1",
			}),
			unmarshaled: &gmcp.CharAfflictionsAdd{
				Name:        "Name1",
				Cure:        "Cure1",
				Description: "Desc1",
			},
			marshaled: makeGMCP("Char.Afflictions.Add", map[string]string{
				"name": "Name1",
				"cure": "Cure1",
				"desc": "Desc1",
			}),
		},

		"Char.Afflictions.Add invalid JSON": {
			msg:  &gmcp.CharAfflictionsAdd{},
			data: "asdf",
			err:  "invalid character 'a' looking for beginning of value",
		},

		"Char.Afflictions.Remove empty": {
			msg:         &gmcp.CharAfflictionsRemove{},
			data:        "Char.Afflictions.Remove []",
			unmarshaled: &gmcp.CharAfflictionsRemove{},
			marshaled:   "Char.Afflictions.Remove []",
		},

		"Char.Afflictions.Remove hydrated": {
			msg: &gmcp.CharAfflictionsRemove{},
			data: makeGMCP("Char.Afflictions.Remove", []string{
				"Name1",
			}),
			unmarshaled: &gmcp.CharAfflictionsRemove{
				{
					Name:        "Name1",
					Cure:        "",
					Description: "",
				},
			},
			marshaled: makeGMCP("Char.Afflictions.Remove", []string{
				"Name1",
			}),
		},

		"Char.Afflictions.Remove invalid JSON": {
			msg:  &gmcp.CharAfflictionsRemove{},
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
