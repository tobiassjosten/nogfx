package gmcp_test

import (
	"strings"
	"testing"

	"github.com/icza/gox/gox"
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

		"Char.Status empty": {
			msg:         &gmcp.CharStatus{},
			data:        "Char.Status {}",
			unmarshaled: &gmcp.CharStatus{},
			marshaled:   "Char.Status {}",
		},

		"Char.Status hydrated": {
			msg: &gmcp.CharStatus{},
			data: makeGMCP("Char.Status", map[string]interface{}{
				"name":     "Durak",
				"fullname": "Mason Durak",
				"race":     "dwarf",
				"gender":   "male",
				"level":    "69 (23%)",
			}),
			unmarshaled: &gmcp.CharStatus{
				Name:     gox.NewString("Durak"),
				Fullname: gox.NewString("Mason Durak"),
				Race:     gox.NewString("dwarf"),
				Gender:   gox.NewString("male"),
				Level:    gox.NewFloat64(69.23),
			},
			marshaled: makeGMCP("Char.Status", map[string]interface{}{
				"name":     "Durak",
				"fullname": "Mason Durak",
				"race":     "dwarf",
				"gender":   "male",
				"level":    "69 (23%)",
			}),
		},

		"Char.Status fractal-progress level": {
			msg: &gmcp.CharStatus{},
			data: makeGMCP("Char.Status", map[string]interface{}{
				"level": "69 (23.45%)",
			}),
			unmarshaled: &gmcp.CharStatus{
				Level: gox.NewFloat64(69.2345),
			},
			marshaled: makeGMCP("Char.Status", map[string]interface{}{
				"level": "69 (23.45%)",
			}),
		},

		"Char.Status single-part level": {
			msg: &gmcp.CharStatus{},
			data: makeGMCP("Char.Status", map[string]interface{}{
				"level": "69",
			}),
			unmarshaled: &gmcp.CharStatus{
				Level: gox.NewFloat64(69),
			},
			marshaled: makeGMCP("Char.Status", map[string]interface{}{
				"level": "69 (0%)",
			}),
		},

		"Char.Status non-number level": {
			msg: &gmcp.CharStatus{},
			data: makeGMCP("Char.Status", map[string]interface{}{
				"level": "asdf",
			}),
			err: `failed parsing level: strconv.ParseFloat: parsing "asdf": invalid syntax`,
		},

		"Char.Status non-number level progress": {
			msg: &gmcp.CharStatus{},
			data: makeGMCP("Char.Status", map[string]interface{}{
				"level": "69 (xy%)",
			}),
			err: `failed parsing level progress: strconv.ParseFloat: parsing "xy": invalid syntax`,
		},

		"Char.Status invalid JSON": {
			msg:  &gmcp.CharStatus{},
			data: "asdf",
			err:  "invalid character 'a' looking for beginning of value",
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

		"Char.Vitals empty": {
			msg:         &gmcp.CharVitals{},
			data:        `Char.Vitals {}`,
			unmarshaled: &gmcp.CharVitals{},
			marshaled: makeGMCP("Char.Vitals", map[string]interface{}{
				"hp":        "0",
				"maxhp":     "0",
				"mp":        "0",
				"maxmp":     "0",
				"charstats": []string{},
			}),
		},

		"Char.Vitals hydrated": {
			msg: &gmcp.CharVitals{},
			data: makeGMCP("Char.Vitals", map[string]interface{}{
				"hp":        "1",
				"maxhp":     "2",
				"mp":        "3",
				"maxmp":     "4",
				"charstats": []string{"Asdf 1"},
			}),
			unmarshaled: &gmcp.CharVitals{
				HP:    1,
				MaxHP: 2,
				MP:    3,
				MaxMP: 4,
				Stats: []string{"Asdf 1"},
			},
			marshaled: makeGMCP("Char.Vitals", map[string]interface{}{
				"hp":        "1",
				"maxhp":     "2",
				"mp":        "3",
				"maxmp":     "4",
				"charstats": []string{"Asdf 1"},
			}),
		},

		"Char.Vitals invalid JSON": {
			msg:  &gmcp.CharVitals{},
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
