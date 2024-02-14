package ironrealms_test

import (
	"strings"
	"testing"

	"github.com/tobiassjosten/nogfx/pkg/gmcp"
	"github.com/tobiassjosten/nogfx/pkg/gmcp/ironrealms"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIRERiftMessages(t *testing.T) {
	tcs := map[string]struct {
		msg         gmcp.Message
		data        string
		unmarshaled gmcp.Message
		marshaled   string
		err         string
	}{
		"IRE.Rift.Change empty": {
			msg:         &ironrealms.IRERiftChange{},
			data:        `IRE.Rift.Change {}`,
			unmarshaled: &ironrealms.IRERiftChange{},
			marshaled: makeGMCP("IRE.Rift.Change", map[string]any{
				"amount": "0",
				"desc":   "",
				"name":   "",
			}),
		},

		"IRE.Rift.Change hydrated": {
			msg: &ironrealms.IRERiftChange{},
			data: makeGMCP("IRE.Rift.Change", map[string]any{
				"amount": "1234",
				"desc":   "alchemical silver",
				"name":   "asilver",
			}),
			unmarshaled: &ironrealms.IRERiftChange{
				Amount:      1234,
				Description: "alchemical silver",
				Name:        "asilver",
			},
			marshaled: makeGMCP("IRE.Rift.Change", map[string]any{
				"amount": "1234",
				"desc":   "alchemical silver",
				"name":   "asilver",
			}),
		},

		"IRE.Rift.Change invalid JSON": {
			msg:  &ironrealms.IRERiftChange{},
			data: `IRE.Rift.Change asdf`,
			err:  "invalid character 'a' looking for beginning of value",
		},

		"IRE.Rift.List empty": {
			msg:         &ironrealms.IRERiftList{},
			data:        `IRE.Rift.List []`,
			unmarshaled: &ironrealms.IRERiftList{},
			marshaled:   makeGMCP("IRE.Rift.List", []map[string]any{}),
		},

		"IRE.Rift.List hydrated": {
			msg: &ironrealms.IRERiftList{},
			data: makeGMCP("IRE.Rift.List", []map[string]any{
				{
					"amount": "1234",
					"desc":   "alchemical silver",
					"name":   "asilver",
				},
			}),
			unmarshaled: &ironrealms.IRERiftList{
				{
					Amount:      1234,
					Description: "alchemical silver",
					Name:        "asilver",
				},
			},
			marshaled: makeGMCP("IRE.Rift.List", []map[string]any{
				{
					"amount": "1234",
					"desc":   "alchemical silver",
					"name":   "asilver",
				},
			}),
		},

		"IRE.Rift.List invalid JSON": {
			msg:  &ironrealms.IRERiftChange{},
			data: `IRE.Rift.Change asdf`,
			err:  "invalid character 'a' looking for beginning of value",
		},

		"IRE.Rift.Request": {
			msg:         &ironrealms.IRERiftRequest{},
			data:        `IRE.Rift.Request {}`,
			unmarshaled: &ironrealms.IRERiftRequest{},
			marshaled:   `IRE.Rift.Request`,
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
