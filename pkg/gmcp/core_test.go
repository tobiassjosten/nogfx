package gmcp_test

import (
	"strings"
	"testing"

	"github.com/icza/gox/gox"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tobiassjosten/nogfx/pkg/gmcp"
)

func TestCoreMessages(t *testing.T) {
	tcs := map[string]struct {
		msg         gmcp.Message
		data        string
		unmarshaled gmcp.Message
		marshaled   string
		err         string
	}{
		"Core.Goodbye": {
			msg:         &gmcp.CoreGoodbye{},
			data:        "Core.Goodbye",
			unmarshaled: &gmcp.CoreGoodbye{},
			marshaled:   "Core.Goodbye",
		},

		"Core.Hello empty": {
			msg:         &gmcp.CoreHello{},
			data:        "Core.Hello {}",
			unmarshaled: &gmcp.CoreHello{},
			marshaled: makeGMCP("Core.Hello", map[string]interface{}{
				"client":  "",
				"version": "",
			}),
		},

		"Core.Hello hydrated": {
			msg: &gmcp.CoreHello{},
			data: makeGMCP("Core.Hello", map[string]interface{}{
				"client":  "nogfx",
				"version": "1.0.0",
			}),
			unmarshaled: &gmcp.CoreHello{
				Client:  "nogfx",
				Version: "1.0.0",
			},
			marshaled: makeGMCP("Core.Hello", map[string]interface{}{
				"client":  "nogfx",
				"version": "1.0.0",
			}),
		},

		"Core.KeepAlive": {
			msg:         &gmcp.CoreKeepAlive{},
			data:        "Core.KeepAlive",
			unmarshaled: &gmcp.CoreKeepAlive{},
			marshaled:   "Core.KeepAlive",
		},

		"Core.Ping empty": {
			msg:         &gmcp.CorePing{},
			data:        "Core.Ping",
			unmarshaled: &gmcp.CorePing{},
			marshaled:   "Core.Ping",
		},

		"Core.Ping latency": {
			msg:  &gmcp.CorePing{},
			data: "Core.Ping 1234",
			unmarshaled: &gmcp.CorePing{
				Latency: gox.NewInt(1234),
			},
			marshaled: "Core.Ping 1234",
		},

		"Core.Ping invalid JSON": {
			msg:  &gmcp.CorePing{},
			data: "Core.Ping asdf",
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

			if tcdata == "" {
				assert.Equal(t, tcdata, data)
				return
			}

			assert.JSONEq(t, tcdata, data, "marshaling maintains data integrity")

			require.Equal(t, tc.unmarshaled, tc.msg, "marshaling doesn't mutate")
		})
	}
}
