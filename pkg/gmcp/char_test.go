package gmcp_test

import (
	"strings"
	"testing"

	"github.com/tobiassjosten/nogfx/pkg/gmcp"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCharMessages(t *testing.T) {
	tcs := map[string]struct {
		msg         gmcp.Message
		data        string
		unmarshaled gmcp.Message
		marshaled   string
		err         string
	}{
		"Char.Login empty": {
			msg:         &gmcp.CharLogin{},
			data:        "Char.Login {}",
			unmarshaled: &gmcp.CharLogin{},
			marshaled: makeGMCP("Char.Login", map[string]interface{}{
				"name":     "",
				"password": "",
			}),
		},

		"Char.Login hydrated": {
			msg: &gmcp.CharLogin{},
			data: makeGMCP("Char.Login", map[string]interface{}{
				"name":     "durak",
				"password": "secret",
			}),
			unmarshaled: &gmcp.CharLogin{
				Name:     "durak",
				Password: "secret",
			},
			marshaled: makeGMCP("Char.Login", map[string]interface{}{
				"name":     "durak",
				"password": "secret",
			}),
		},

		"Char.Name empty": {
			msg:         &gmcp.CharName{},
			data:        "Char.Name {}",
			unmarshaled: &gmcp.CharName{},
			marshaled:   `Char.Name {"name":"","fullname":""}`,
		},

		"Char.Name hydrated": {
			msg: &gmcp.CharName{},
			data: makeGMCP("Char.Name", map[string]interface{}{
				"name":     "Durak",
				"fullname": "Mason Durak",
			}),
			unmarshaled: &gmcp.CharName{
				Name:     "Durak",
				Fullname: "Mason Durak",
			},
			marshaled: makeGMCP("Char.Name", map[string]interface{}{
				"name":     "Durak",
				"fullname": "Mason Durak",
			}),
		},

		"Char.StatusVars empty": {
			msg:         &gmcp.CharStatusVars{},
			data:        "Char.StatusVars {}",
			unmarshaled: &gmcp.CharStatusVars{},
			marshaled:   "Char.StatusVars {}",
		},

		"Char.StatusVars hydrated": {
			msg:         &gmcp.CharStatusVars{},
			data:        `Char.StatusVars {"this": "That"}`,
			unmarshaled: &gmcp.CharStatusVars{"this": "That"},
			marshaled:   `Char.StatusVars {"this": "That"}`,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			err := tc.msg.Unmarshal([]byte(tc.data))

			if tc.err != "" {
				require.NotNil(t, err)
				assert.Equal(t, tc.err, err.Error())
				return
			}
			require.Nil(t, err)

			assert.Equal(t, tc.unmarshaled, tc.msg, "unmarshaling hydrates message")

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
