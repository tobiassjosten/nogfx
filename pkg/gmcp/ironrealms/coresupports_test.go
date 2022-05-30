package ironrealms_test

import (
	"strings"
	"testing"

	"github.com/icza/gox/gox"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tobiassjosten/nogfx/pkg/gmcp"
	igmcp "github.com/tobiassjosten/nogfx/pkg/gmcp/ironrealms"
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
			msg:  &igmcp.CoreSupportsSet{},
			data: "Core.Supports.Set []",
			unmarshaled: &igmcp.CoreSupportsSet{
				&igmcp.CoreSupports{
					CoreSupports: &gmcp.CoreSupports{},
				},
			},
			marshaled: "Core.Supports.Set []",
		},

		"Core.Supports.Set hydrated": {
			msg: &igmcp.CoreSupportsSet{},
			data: makeGMCP("Core.Supports.Set", []string{
				"Char 1",
				"Char.Skills 2",
				"Char.Items 3",
				"Comm.Channel 4",
				"Room 5",
				"IRE.Rift 6",
				"IRE.Target 7",
			}),
			unmarshaled: &igmcp.CoreSupportsSet{
				CoreSupports: &igmcp.CoreSupports{
					CoreSupports: &gmcp.CoreSupports{
						Char:        gox.NewInt(1),
						CharSkills:  gox.NewInt(2),
						CharItems:   gox.NewInt(3),
						CommChannel: gox.NewInt(4),
						Room:        gox.NewInt(5),
					},
					IRERift:   gox.NewInt(6),
					IRETarget: gox.NewInt(7),
				},
			},
			marshaled: makeGMCP("Core.Supports.Set", []string{
				"Char 1",
				"Char.Skills 2",
				"Char.Items 3",
				"Comm.Channel 4",
				"Room 5",
				"IRE.Rift 6",
				"IRE.Target 7",
			}),
		},

		"Core.Supports.Set invalid JSON": {
			msg:  &igmcp.CoreSupportsSet{},
			data: "Core.Supports.Set asdf",
			err:  "invalid character 'a' looking for beginning of value",
		},

		"Core.Supports.Set invalid list": {
			msg:  &igmcp.CoreSupportsSet{},
			data: "Core.Supports.Set {}",
			err:  "json: cannot unmarshal object into Go value of type []string",
		},

		"Core.Supports.Set invalid module": {
			msg: &igmcp.CoreSupportsSet{},
			data: makeGMCP("Core.Supports.Set", []string{
				"Asdf 1",
			}),
			unmarshaled: &igmcp.CoreSupportsSet{
				&igmcp.CoreSupports{
					CoreSupports: &gmcp.CoreSupports{},
				},
			},
			marshaled: "Core.Supports.Set []",
		},

		"Core.Supports.Set invalid number": {
			msg: &igmcp.CoreSupportsSet{},
			data: makeGMCP("Core.Supports.Set", []string{
				"Char xy",
			}),
			err: `failed parsing module version: strconv.Atoi: parsing "xy": invalid syntax`,
		},

		"Core.Supports.Add empty": {
			msg:  &igmcp.CoreSupportsAdd{},
			data: "Core.Supports.Add []",
			unmarshaled: &igmcp.CoreSupportsAdd{
				&igmcp.CoreSupports{
					CoreSupports: &gmcp.CoreSupports{},
				},
			},
			marshaled: "Core.Supports.Add []",
		},

		"Core.Supports.Add hydrated": {
			msg: &igmcp.CoreSupportsAdd{},
			data: makeGMCP("Core.Supports.Add", []string{
				"Char 1",
				"Char.Skills 2",
				"Char.Items 3",
				"Comm.Channel 4",
				"Room 5",
				"IRE.Rift 6",
				"IRE.Target 7",
			}),
			unmarshaled: &igmcp.CoreSupportsAdd{
				&igmcp.CoreSupports{
					CoreSupports: &gmcp.CoreSupports{
						Char:        gox.NewInt(1),
						CharSkills:  gox.NewInt(2),
						CharItems:   gox.NewInt(3),
						CommChannel: gox.NewInt(4),
						Room:        gox.NewInt(5),
					},
					IRERift:   gox.NewInt(6),
					IRETarget: gox.NewInt(7),
				},
			},
			marshaled: makeGMCP("Core.Supports.Add", []string{
				"Char 1",
				"Char.Skills 2",
				"Char.Items 3",
				"Comm.Channel 4",
				"Room 5",
				"IRE.Rift 6",
				"IRE.Target 7",
			}),
		},

		"Core.Supports.Add invalid JSON": {
			msg:  &igmcp.CoreSupportsAdd{},
			data: "Core.Supports.Add asdf",
			err:  "invalid character 'a' looking for beginning of value",
		},

		"Core.Supports.Add invalid list": {
			msg:  &igmcp.CoreSupportsAdd{},
			data: "Core.Supports.Add {}",
			err:  "json: cannot unmarshal object into Go value of type []string",
		},

		"Core.Supports.Add invalid module": {
			msg: &igmcp.CoreSupportsAdd{},
			data: makeGMCP("Core.Supports.Add", []string{
				"Asdf 1",
			}),
			unmarshaled: &igmcp.CoreSupportsAdd{
				&igmcp.CoreSupports{
					CoreSupports: &gmcp.CoreSupports{},
				},
			},
			marshaled: "Core.Supports.Add []",
		},

		"Core.Supports.Add invalid number": {
			msg: &igmcp.CoreSupportsAdd{},
			data: makeGMCP("Core.Supports.Add", []string{
				"Char xy",
			}),
			err: `failed parsing module version: strconv.Atoi: parsing "xy": invalid syntax`,
		},

		"Core.Supports.Remove empty": {
			msg:  &igmcp.CoreSupportsRemove{},
			data: "Core.Supports.Remove []",
			unmarshaled: &igmcp.CoreSupportsRemove{
				&igmcp.CoreSupports{
					CoreSupports: &gmcp.CoreSupports{},
				},
			},
			marshaled: "Core.Supports.Remove []",
		},

		"Core.Supports.Remove hydrated": {
			msg: &igmcp.CoreSupportsRemove{},
			data: makeGMCP("Core.Supports.Remove", []string{
				"Char",
				"Char.Skills",
				"Char.Items",
				"Comm.Channel",
				"Room",
				"IRE.Rift",
				"IRE.Target",
			}),
			unmarshaled: &igmcp.CoreSupportsRemove{
				&igmcp.CoreSupports{
					CoreSupports: &gmcp.CoreSupports{
						Char:        gox.NewInt(1),
						CharSkills:  gox.NewInt(1),
						CharItems:   gox.NewInt(1),
						CommChannel: gox.NewInt(1),
						Room:        gox.NewInt(1),
					},
					IRERift:   gox.NewInt(1),
					IRETarget: gox.NewInt(1),
				},
			},
			marshaled: makeGMCP("Core.Supports.Remove", []string{
				"Char 1",
				"Char.Skills 1",
				"Char.Items 1",
				"Comm.Channel 1",
				"Room 1",
				"IRE.Rift 1",
				"IRE.Target 1",
			}),
		},

		"Core.Supports.Remove invalid JSON": {
			msg:  &igmcp.CoreSupportsRemove{},
			data: "Core.Supports.Remove asdf",
			err:  "invalid character 'a' looking for beginning of value",
		},

		"Core.Supports.Remove invalid list": {
			msg:  &igmcp.CoreSupportsRemove{},
			data: "Core.Supports.Remove {}",
			err:  "json: cannot unmarshal object into Go value of type []string",
		},

		"Core.Supports.Remove invalid module": {
			msg: &igmcp.CoreSupportsRemove{},
			data: makeGMCP("Core.Supports.Remove", []string{
				"Asdf 1",
			}),
			unmarshaled: &igmcp.CoreSupportsRemove{
				&igmcp.CoreSupports{
					CoreSupports: &gmcp.CoreSupports{},
				},
			},
			marshaled: "Core.Supports.Remove []",
		},

		"Core.Supports.Remove invalid number": {
			msg: &igmcp.CoreSupportsRemove{},
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
