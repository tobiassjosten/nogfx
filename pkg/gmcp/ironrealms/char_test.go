package ironrealms_test

import (
	"strings"
	"testing"

	"github.com/icza/gox/gox"
	"github.com/tobiassjosten/nogfx/pkg/gmcp"
	"github.com/tobiassjosten/nogfx/pkg/gmcp/ironrealms"

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
		"Char.Status empty": {
			msg:  &ironrealms.CharStatus{},
			data: `Char.Status {}`,
			unmarshaled: &ironrealms.CharStatus{
				CharStatus: &gmcp.CharStatus{},
			},
			marshaled: `Char.Status {}`,
		},

		"Char.Status hydrated": {
			msg: &ironrealms.CharStatus{},
			data: makeGMCP("Char.Status", map[string]interface{}{
				"bank":        "1",
				"city":        "Mhaldor (2)",
				"class":       "Monk",
				"gold":        "3",
				"unread_msgs": "4",
				"unread_news": "5",
				// Base:
				"name":     "Durak",
				"fullname": "Mason Durak",
				"race":     "dwarf",
				"gender":   "male",
				"level":    "69 (23%)",
			}),
			unmarshaled: &ironrealms.CharStatus{
				CharStatus: &gmcp.CharStatus{
					Name:     gox.NewString("Durak"),
					Fullname: gox.NewString("Mason Durak"),
					Race:     gox.NewString("dwarf"),
					Gender:   gox.NewString("male"),
					Level:    gox.NewFloat64(69.23),
				},
				Bank:       gox.NewInt(1),
				City:       gox.NewString("Mhaldor"),
				CityRank:   gox.NewInt(2),
				Class:      gox.NewString("Monk"),
				Gold:       gox.NewInt(3),
				UnreadMsgs: gox.NewInt(4),
				UnreadNews: gox.NewInt(5),
			},
			marshaled: makeGMCP("Char.Status", map[string]interface{}{
				"bank":        "1",
				"city":        "Mhaldor (2)",
				"class":       "Monk",
				"gold":        "3",
				"unread_msgs": "4",
				"unread_news": "5",
				// Base:
				"name":     "Durak",
				"fullname": "Mason Durak",
				"race":     "dwarf",
				"gender":   "male",
				"level":    "69 (23%)",
			}),
		},

		"Char.Status empty city": {
			msg: &ironrealms.CharStatus{},
			data: makeGMCP("Char.Status", map[string]interface{}{
				"city": "(None)",
			}),
			unmarshaled: &ironrealms.CharStatus{
				CharStatus: &gmcp.CharStatus{},
				City:       gox.NewString(""),
			},
			marshaled: makeGMCP("Char.Status", map[string]interface{}{
				"city": "(None)",
			}),
		},

		"Char.Status invalid JSON": {
			msg:  &ironrealms.CharStatus{},
			data: `Char.Status asdf`,
			err:  "invalid character 'a' looking for beginning of value",
		},

		"Char.Status parent invalid JSON": {
			msg: &ironrealms.CharStatus{},
			data: makeGMCP("Char.Status", map[string]interface{}{
				"level": "69 (xy%)",
			}),
			err: `failed parsing level progress: strconv.ParseFloat: parsing "xy": invalid syntax`,
		},

		"Char.Vitals empty": {
			msg:         &ironrealms.CharVitals{},
			data:        `Char.Vitals {}`,
			unmarshaled: &ironrealms.CharVitals{},
			marshaled: makeGMCP("Char.Vitals", map[string]interface{}{
				"bal": "0",
				"eq":  "0",
				"nl":  "0",
			}),
		},

		"Char.Vitals hydrated": {
			msg: &ironrealms.CharVitals{},
			data: makeGMCP("Char.Vitals", map[string]interface{}{
				"bal": "1",
				"eq":  "1",
				"nl":  "23",
			}),
			unmarshaled: &ironrealms.CharVitals{
				Bal: true,
				Eq:  true,
				NL:  23,
			},
			marshaled: makeGMCP("Char.Vitals", map[string]interface{}{
				"bal": "1",
				"eq":  "1",
				"nl":  "23",
			}),
		},

		"Char.Vitals invalid JSON": {
			msg:  &ironrealms.CharVitals{},
			data: `Char.Vitals asdf`,
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
