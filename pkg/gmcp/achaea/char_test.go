package achaea_test

import (
	"strings"
	"testing"

	"github.com/tobiassjosten/nogfx/pkg/gmcp"
	agmcp "github.com/tobiassjosten/nogfx/pkg/gmcp/achaea"

	"github.com/icza/gox/gox"
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
			msg:         &agmcp.CharStatus{},
			data:        "Char.Status {}",
			unmarshaled: &agmcp.CharStatus{},
			marshaled:   "Char.Status {}",
		},

		"Char.Status hydrated": {
			msg: &agmcp.CharStatus{},
			data: makeGMCP("Char.Status", map[string]interface{}{
				"age":              "6",
				"bank":             "1",
				"boundcredits":     "7",
				"boundmayancrowns": "8",
				"city":             "Mhaldor (2)",
				"class":            "Monk",
				"explorerrank":     "an Itinerant",
				"fullname":         "Mason Durak",
				"gender":           "male",
				"gold":             "3",
				"house":            "The Dread Legates(9)",
				"lessons":          "10",
				"level":            "69 (23%)",
				"mayancrowns":      "11",
				"name":             "Durak",
				"order":            "Something (12)",
				"race":             "dwarf",
				"specialisation":   "Brawler",
				"target":           "someone",
				"unboundcredits":   "13",
				"unread_msgs":      "4",
				"unread_news":      "5",
				"xprank":           "14",
			}),
			unmarshaled: &agmcp.CharStatus{
				Age:              gox.NewInt(6),
				Bank:             gox.NewInt(1),
				BoundCredits:     gox.NewInt(7),
				BoundMayanCrowns: gox.NewInt(8),
				City:             gox.NewString("Mhaldor"),
				CityRank:         gox.NewInt(2),
				Class:            gox.NewString("Monk"),
				ExplorerRank:     gox.NewString("an Itinerant"),
				Fullname:         gox.NewString("Mason Durak"),
				Gender:           gox.NewString("male"),
				Gold:             gox.NewInt(3),
				House:            gox.NewString("The Dread Legates"),
				HouseRank:        gox.NewInt(9),
				Lessons:          gox.NewInt(10),
				Level:            gox.NewFloat64(69.23),
				MayanCrowns:      gox.NewInt(11),
				Name:             gox.NewString("Durak"),
				Order:            gox.NewString("Something"),
				OrderRank:        gox.NewInt(12),
				Race:             gox.NewString("dwarf"),
				Specialisation:   gox.NewString("Brawler"),
				Target:           gox.NewString("someone"),
				UnboundCredits:   gox.NewInt(13),
				UnreadMsgs:       gox.NewInt(4),
				UnreadNews:       gox.NewInt(5),
				XPRank:           gox.NewInt(14),
			},
			marshaled: makeGMCP("Char.Status", map[string]interface{}{
				"age":              "6",
				"bank":             "1",
				"boundcredits":     "7",
				"boundmayancrowns": "8",
				"city":             "Mhaldor (2)",
				"class":            "Monk",
				"explorerrank":     "an Itinerant",
				"fullname":         "Mason Durak",
				"gender":           "male",
				"gold":             "3",
				"house":            "The Dread Legates (9)",
				"lessons":          "10",
				"level":            "69 (23%)",
				"mayancrowns":      "11",
				"name":             "Durak",
				"order":            "Something (12)",
				"race":             "dwarf",
				"specialisation":   "Brawler",
				"target":           "someone",
				"unboundcredits":   "13",
				"unread_msgs":      "4",
				"unread_news":      "5",
				"xprank":           "14",
			}),
		},

		"Char.Status empty city": {
			msg: &agmcp.CharStatus{},
			data: makeGMCP("Char.Status", map[string]interface{}{
				"city": "(None)",
			}),
			unmarshaled: &agmcp.CharStatus{
				City: gox.NewString(""),
			},
			marshaled: makeGMCP("Char.Status", map[string]interface{}{
				"city": "(None)",
			}),
		},

		"Char.Status fractal-progress level": {
			msg: &agmcp.CharStatus{},
			data: makeGMCP("Char.Status", map[string]interface{}{
				"level": "69 (23.45%)",
			}),
			unmarshaled: &agmcp.CharStatus{
				Level: gox.NewFloat64(69.2345),
			},
			marshaled: makeGMCP("Char.Status", map[string]interface{}{
				"level": "69 (23.45%)",
			}),
		},

		"Char.Status single-part level": {
			msg: &agmcp.CharStatus{},
			data: makeGMCP("Char.Status", map[string]interface{}{
				"level": "69",
			}),
			unmarshaled: &agmcp.CharStatus{
				Level: gox.NewFloat64(69),
			},
			marshaled: makeGMCP("Char.Status", map[string]interface{}{
				"level": "69 (0%)",
			}),
		},

		"Char.Status non-number level": {
			msg: &agmcp.CharStatus{},
			data: makeGMCP("Char.Status", map[string]interface{}{
				"level": "asdf",
			}),
			err: `failed parsing level: strconv.ParseFloat: parsing "asdf": invalid syntax`,
		},

		"Char.Status non-number level progress": {
			msg: &agmcp.CharStatus{},
			data: makeGMCP("Char.Status", map[string]interface{}{
				"level": "69 (xy%)",
			}),
			err: `failed parsing level progress: strconv.ParseFloat: parsing "xy": invalid syntax`,
		},

		"Char.Status invalid JSON": {
			msg:  &agmcp.CharStatus{},
			data: "asdf",
			err:  "invalid character 'a' looking for beginning of value",
		},

		"Char.Status parent invalid JSON": {
			msg: &agmcp.CharStatus{},
			data: makeGMCP("Char.Status", map[string]interface{}{
				"level": "69 (xy%)",
			}),
			err: `failed parsing level progress: strconv.ParseFloat: parsing "xy": invalid syntax`,
		},

		"Char.Vitals empty": {
			msg:         &agmcp.CharVitals{},
			data:        `Char.Vitals {}`,
			unmarshaled: &agmcp.CharVitals{},
			marshaled: makeGMCP("Char.Vitals", map[string]interface{}{
				"hp":     "0",
				"maxhp":  "0",
				"mp":     "0",
				"maxmp":  "0",
				"ep":     "0",
				"maxep":  "0",
				"wp":     "0",
				"maxwp":  "0",
				"nl":     "0",
				"bal":    "0",
				"eq":     "0",
				"string": "",
				"charstats": []string{
					"Bleed: 0",
					"Rage: 0",
				},
			}),
		},

		"Char.Vitals hydrated": {
			msg: &agmcp.CharVitals{},
			data: makeGMCP("Char.Vitals", map[string]interface{}{
				"hp":     "2",
				"maxhp":  "3",
				"mp":     "4",
				"maxmp":  "5",
				"ep":     "6",
				"maxep":  "7",
				"wp":     "8",
				"maxwp":  "9",
				"nl":     "10",
				"bal":    "1",
				"eq":     "1",
				"string": "asdf> ",
				"charstats": []string{
					"Bleed: 1",
					"Rage: 2",
					"Ferocity: 3",
					"Kai: 4%",
					"Karma: 5%",
					"Spec: Sword and Shield",
					"Stance: Scorpion",
				},
			}),
			unmarshaled: &agmcp.CharVitals{
				HP:     2,
				MaxHP:  3,
				MP:     4,
				MaxMP:  5,
				EP:     6,
				MaxEP:  7,
				WP:     8,
				MaxWP:  9,
				NL:     10,
				Bal:    true,
				Eq:     true,
				Prompt: "asdf> ",
				Stats: agmcp.CharVitalsStats{
					Bleed:    1,
					Rage:     2,
					Ferocity: gox.NewInt(3),
					Kai:      gox.NewInt(4),
					Karma:    gox.NewInt(5),
					Spec:     gox.NewString("Sword and Shield"),
					Stance:   gox.NewString("Scorpion"),
				},
			},
			marshaled: makeGMCP("Char.Vitals", map[string]interface{}{
				"hp":     "2",
				"maxhp":  "3",
				"mp":     "4",
				"maxmp":  "5",
				"ep":     "6",
				"maxep":  "7",
				"wp":     "8",
				"maxwp":  "9",
				"nl":     "10",
				"bal":    "1",
				"eq":     "1",
				"string": "asdf> ",
				"charstats": []string{
					"Bleed: 1",
					"Rage: 2",
					"Ferocity: 3",
					"Kai: 4%",
					"Karma: 5%",
					"Spec: Sword and Shield",
					"Stance: Scorpion",
				},
			}),
		},

		"Char.Vitals invalid JSON": {
			msg:  &agmcp.CharVitals{},
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

			if tcdata == "" {
				assert.Equal(t, tcdata, data)
				return
			}

			assert.JSONEq(t, tcdata, data, "marshaling maintains data integrity")

			require.Equal(t, tc.unmarshaled, tc.msg, "marshaling doesn't mutate")
		})
	}
}
