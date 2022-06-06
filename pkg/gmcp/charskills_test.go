package gmcp_test

import (
	"strings"
	"testing"

	"github.com/icza/gox/gox"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tobiassjosten/nogfx/pkg/gmcp"
)

func TestCharSkillsMessages(t *testing.T) {
	tcs := map[string]struct {
		msg         gmcp.Message
		data        string
		unmarshaled gmcp.Message
		marshaled   string
		err         string
	}{
		"Char.Skills.Get empty": {
			msg:         &gmcp.CharSkillsGet{},
			data:        "Char.Skills.Get {}",
			unmarshaled: &gmcp.CharSkillsGet{},
			marshaled:   "Char.Skills.Get {}",
		},

		"Char.Skills.Get group": {
			msg: &gmcp.CharSkillsGet{},
			data: makeGMCP("Char.Skills.Get", map[string]interface{}{
				"group": "tekura",
			}),
			unmarshaled: &gmcp.CharSkillsGet{
				Group: "tekura",
			},
			marshaled: makeGMCP("Char.Skills.Get", map[string]interface{}{
				"group": "tekura",
			}),
		},

		"Char.Skills.Get name": {
			msg: &gmcp.CharSkillsGet{},
			data: makeGMCP("Char.Skills.Get", map[string]interface{}{
				"name": "sidekick",
			}),
			unmarshaled: &gmcp.CharSkillsGet{
				Name: "sidekick",
			},
			marshaled: "Char.Skills.Get {}",
		},

		"Char.Skills.Get group and name": {
			msg: &gmcp.CharSkillsGet{},
			data: makeGMCP("Char.Skills.Get", map[string]interface{}{
				"group": "tekura",
				"name":  "sidekick",
			}),
			unmarshaled: &gmcp.CharSkillsGet{
				Group: "tekura",
				Name:  "sidekick",
			},
			marshaled: makeGMCP("Char.Skills.Get", map[string]interface{}{
				"group": "tekura",
				"name":  "sidekick",
			}),
		},

		"Char.Skills.Group empty": {
			msg:         &gmcp.CharSkillsGroups{},
			data:        "Char.Skills.Groups []",
			unmarshaled: &gmcp.CharSkillsGroups{},
			marshaled:   "Char.Skills.Groups []",
		},

		"Char.Skills.Group hydrated": {
			msg: &gmcp.CharSkillsGroups{},
			data: makeGMCP("Char.Skills.Groups", []map[string]interface{}{
				{
					"name": "Tekura",
					"rank": "Adept (40%)",
				},
			}),
			unmarshaled: &gmcp.CharSkillsGroups{
				{
					Name:     "Tekura",
					Rank:     "Adept",
					Progress: gox.NewInt(40),
				},
			},
			marshaled: makeGMCP("Char.Skills.Groups", []map[string]interface{}{
				{
					"name": "Tekura",
					"rank": "Adept (40%)",
				},
			}),
		},

		"Char.Skills.Group no progress": {
			msg: &gmcp.CharSkillsGroups{},
			data: makeGMCP("Char.Skills.Groups", []map[string]interface{}{
				{
					"name": "Tekura",
					"rank": "Adept",
				},
			}),
			unmarshaled: &gmcp.CharSkillsGroups{
				{
					Name: "Tekura",
					Rank: "Adept",
				},
			},
			marshaled: makeGMCP("Char.Skills.Groups", []map[string]interface{}{
				{
					"name": "Tekura",
					"rank": "Adept",
				},
			}),
		},

		"Char.Skills.Group invalid JSON": {
			msg:  &gmcp.CharSkillsGroups{},
			data: "Char.Skills.Groups asdf",
			err:  "invalid character 'a' looking for beginning of value",
		},

		"Char.Skills.Group invalid progress": {
			msg: &gmcp.CharSkillsGroups{},
			data: makeGMCP("Char.Skills.Groups", []map[string]interface{}{
				{
					"name": "Tekura",
					"rank": "Adept (xy%)",
				},
			}),
			err: `failed parsing rank progress: strconv.Atoi: parsing "xy": invalid syntax`,
		},

		"Char.Skills.List empty": {
			msg:         &gmcp.CharSkillsList{},
			data:        "Char.Skills.List {}",
			unmarshaled: &gmcp.CharSkillsList{},
			marshaled: makeGMCP("Char.Skills.List", map[string]interface{}{
				"group": "",
				"list":  []string{},
				"descs": []string{},
			}),
		},

		"Char.Skills.List hydrated": {
			msg: &gmcp.CharSkillsList{},
			data: makeGMCP("Char.Skills.List", map[string]interface{}{
				"group": "Tekura",
				"list":  []string{"sidekick"},
				"descs": []string{"A kick to the side"},
			}),
			unmarshaled: &gmcp.CharSkillsList{
				Group:        "Tekura",
				List:         []string{"sidekick"},
				Descriptions: []string{"A kick to the side"},
			},
			marshaled: makeGMCP("Char.Skills.List", map[string]interface{}{
				"group": "Tekura",
				"list":  []string{"sidekick"},
				"descs": []string{"A kick to the side"},
			}),
		},

		"Char.Skills.Info empty": {
			msg:         &gmcp.CharSkillsInfo{},
			data:        "Char.Skills.Info {}",
			unmarshaled: &gmcp.CharSkillsInfo{},
			marshaled: makeGMCP("Char.Skills.Info", map[string]interface{}{
				"group": "",
				"skill": "",
				"info":  "",
			}),
		},

		"Char.Skills.Info hydrated": {
			msg: &gmcp.CharSkillsInfo{},
			data: makeGMCP("Char.Skills.Info", map[string]interface{}{
				"group": "Tekura",
				"skill": "sidekick",
				"info":  "A kick to the side",
			}),
			unmarshaled: &gmcp.CharSkillsInfo{
				Group: "Tekura",
				Skill: "sidekick",
				Info:  "A kick to the side",
			},
			marshaled: makeGMCP("Char.Skills.Info", map[string]interface{}{
				"group": "Tekura",
				"skill": "sidekick",
				"info":  "A kick to the side",
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
