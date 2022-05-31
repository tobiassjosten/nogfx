package gmcp_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tobiassjosten/nogfx/pkg/gmcp"
)

func TestCoreSupportsMessages(t *testing.T) {
	tcs := map[string]struct {
		msg         gmcp.Message
		data        string
		unmarshaled gmcp.Message
		marshaled   string
		err         string
	}{
		"Core.Supports.Set empty": {
			msg:         &gmcp.CoreSupportsSet{},
			data:        "Core.Supports.Set []",
			unmarshaled: &gmcp.CoreSupportsSet{},
			marshaled:   "Core.Supports.Set []",
		},

		"Core.Supports.Set hydrated": {
			msg: &gmcp.CoreSupportsSet{},
			data: makeGMCP("Core.Supports.Set", []string{
				"Char 1",
				"Char.Items 2",
				"Char.Skills 3",
				"Comm.Channel 4",
				"Room 5",
			}),
			unmarshaled: &gmcp.CoreSupportsSet{
				"Char":         1,
				"Char.Items":   2,
				"Char.Skills":  3,
				"Comm.Channel": 4,
				"Room":         5,
			},
			marshaled: makeGMCP("Core.Supports.Set", []string{
				"Char 1",
				"Char.Items 2",
				"Char.Skills 3",
				"Comm.Channel 4",
				"Room 5",
			}),
		},

		"Core.Supports.Set invalid JSON": {
			msg:  &gmcp.CoreSupportsSet{},
			data: "Core.Supports.Set asdf",
			err:  "invalid character 'a' looking for beginning of value",
		},

		"Core.Supports.Set invalid list": {
			msg:  &gmcp.CoreSupportsSet{},
			data: "Core.Supports.Set {}",
			err:  "json: cannot unmarshal object into Go value of type []string",
		},

		"Core.Supports.Set invalid module": {
			msg: &gmcp.CoreSupportsSet{},
			data: makeGMCP("Core.Supports.Set", []string{
				"Asdf 1",
			}),
			unmarshaled: &gmcp.CoreSupportsSet{
				"Asdf": 1,
			},
			marshaled: makeGMCP("Core.Supports.Set", []string{
				"Asdf 1",
			}),
		},

		"Core.Supports.Set invalid number": {
			msg: &gmcp.CoreSupportsSet{},
			data: makeGMCP("Core.Supports.Set", []string{
				"Char xy",
			}),
			err: `failed parsing module version: strconv.Atoi: parsing "xy": invalid syntax`,
		},

		"Core.Supports.Add empty": {
			msg:         &gmcp.CoreSupportsAdd{},
			data:        "Core.Supports.Add []",
			unmarshaled: &gmcp.CoreSupportsAdd{},
			marshaled:   "Core.Supports.Add []",
		},

		"Core.Supports.Add hydrated": {
			msg: &gmcp.CoreSupportsAdd{},
			data: makeGMCP("Core.Supports.Add", []string{
				"Char 1",
				"Char.Items 2",
				"Char.Skills 3",
				"Comm.Channel 4",
				"Room 5",
			}),
			unmarshaled: &gmcp.CoreSupportsAdd{
				"Char":         1,
				"Char.Items":   2,
				"Char.Skills":  3,
				"Comm.Channel": 4,
				"Room":         5,
			},
			marshaled: makeGMCP("Core.Supports.Add", []string{
				"Char 1",
				"Char.Items 2",
				"Char.Skills 3",
				"Comm.Channel 4",
				"Room 5",
			}),
		},

		"Core.Supports.Add invalid JSON": {
			msg:  &gmcp.CoreSupportsAdd{},
			data: "Core.Supports.Add asdf",
			err:  "invalid character 'a' looking for beginning of value",
		},

		"Core.Supports.Add invalid list": {
			msg:  &gmcp.CoreSupportsAdd{},
			data: "Core.Supports.Add {}",
			err:  "json: cannot unmarshal object into Go value of type []string",
		},

		"Core.Supports.Add invalid module": {
			msg: &gmcp.CoreSupportsAdd{},
			data: makeGMCP("Core.Supports.Add", []string{
				"Asdf 1",
			}),
			unmarshaled: &gmcp.CoreSupportsAdd{
				"Asdf": 1,
			},
			marshaled: makeGMCP("Core.Supports.Add", []string{
				"Asdf 1",
			}),
		},

		"Core.Supports.Add invalid number": {
			msg: &gmcp.CoreSupportsAdd{},
			data: makeGMCP("Core.Supports.Add", []string{
				"Char xy",
			}),
			err: `failed parsing module version: strconv.Atoi: parsing "xy": invalid syntax`,
		},

		"Core.Supports.Remove empty": {
			msg:         &gmcp.CoreSupportsRemove{},
			data:        "Core.Supports.Remove []",
			unmarshaled: &gmcp.CoreSupportsRemove{},
			marshaled:   "Core.Supports.Remove []",
		},

		"Core.Supports.Remove hydrated": {
			msg: &gmcp.CoreSupportsRemove{},
			data: makeGMCP("Core.Supports.Remove", []string{
				"Char",
				"Char.Items",
				"Char.Skills",
				"Comm.Channel",
				"Room",
			}),
			unmarshaled: &gmcp.CoreSupportsRemove{
				"Char":         1,
				"Char.Items":   1,
				"Char.Skills":  1,
				"Comm.Channel": 1,
				"Room":         1,
			},
			marshaled: makeGMCP("Core.Supports.Remove", []string{
				"Char 1",
				"Char.Items 1",
				"Char.Skills 1",
				"Comm.Channel 1",
				"Room 1",
			}),
		},

		"Core.Supports.Remove invalid JSON": {
			msg:  &gmcp.CoreSupportsRemove{},
			data: "Core.Supports.Remove asdf",
			err:  "invalid character 'a' looking for beginning of value",
		},

		"Core.Supports.Remove invalid list": {
			msg:  &gmcp.CoreSupportsRemove{},
			data: "Core.Supports.Remove {}",
			err:  "json: cannot unmarshal object into Go value of type []string",
		},

		"Core.Supports.Remove invalid module": {
			msg: &gmcp.CoreSupportsRemove{},
			data: makeGMCP("Core.Supports.Remove", []string{
				"Asdf",
			}),
			unmarshaled: &gmcp.CoreSupportsRemove{
				"Asdf": 1,
			},
			marshaled: makeGMCP("Core.Supports.Remove", []string{
				"Asdf 1",
			}),
		},

		"Core.Supports.Remove invalid number": {
			msg: &gmcp.CoreSupportsRemove{},
			data: makeGMCP("Core.Supports.Remove", []string{
				"Char xy",
			}),
			err: `failed parsing module version: strconv.Atoi: parsing "xy": invalid syntax`,
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
