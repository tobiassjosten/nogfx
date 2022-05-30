package ironrealms_test

import (
	"strings"
	"testing"

	"github.com/tobiassjosten/nogfx/pkg/gmcp"
	"github.com/tobiassjosten/nogfx/pkg/gmcp/ironrealms"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIRETargetMessages(t *testing.T) {
	tcs := map[string]struct {
		msg         gmcp.Message
		data        string
		unmarshaled gmcp.Message
		marshaled   string
		err         string
	}{
		"IRE.Target.Set empty": {
			msg:         &ironrealms.IRETargetSet{},
			data:        `IRE.Target.Set ""`,
			unmarshaled: &ironrealms.IRETargetSet{},
			marshaled:   `IRE.Target.Set ""`,
		},

		"IRE.Target.Set hydrated": {
			msg:  &ironrealms.IRETargetSet{},
			data: `IRE.Target.Set "asdf"`,
			unmarshaled: &ironrealms.IRETargetSet{
				Target: "asdf",
			},
			marshaled: `IRE.Target.Set "asdf"`,
		},

		"IRE.Target.Set invalid JSON": {
			msg:  &ironrealms.IRETargetSet{},
			data: `IRE.Target.Set asdf`,
			err:  "invalid character 'a' looking for beginning of value",
		},

		"IRE.Target.Info empty": {
			msg:         &ironrealms.IRETargetInfo{},
			data:        `IRE.Target.Info {}`,
			unmarshaled: &ironrealms.IRETargetInfo{},
			marshaled: makeGMCP("IRE.Target.Info", map[string]interface{}{
				"id":         "",
				"hpperc":     "",
				"short_desc": "",
			}),
		},

		"IRE.Target.Info hydrated": {
			msg: &ironrealms.IRETargetInfo{},
			data: makeGMCP("IRE.Target.Info", map[string]interface{}{
				"id":         "1234",
				"hpperc":     "69%",
				"short_desc": "a target",
			}),
			unmarshaled: &ironrealms.IRETargetInfo{
				Identity:    "1234",
				Health:      69,
				Description: "a target",
			},
			marshaled: makeGMCP("IRE.Target.Info", map[string]interface{}{
				"id":         "1234",
				"hpperc":     "69%",
				"short_desc": "a target",
			}),
		},

		"IRE.Target.Info invalid JSON": {
			msg:  &ironrealms.IRETargetInfo{},
			data: `IRE.Target.Info asdf`,
			err:  "invalid character 'a' looking for beginning of value",
		},

		"IRE.Target.Info invalid health": {
			msg: &ironrealms.IRETargetInfo{},
			data: makeGMCP("IRE.Target.Info", map[string]interface{}{
				"hpperc": "asdf",
			}),
			err: `strconv.Atoi: parsing "asdf": invalid syntax`,
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
