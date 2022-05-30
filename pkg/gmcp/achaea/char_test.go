package achaea_test

import (
	"strings"
	"testing"

	"github.com/tobiassjosten/nogfx/pkg/gmcp"
	agmcp "github.com/tobiassjosten/nogfx/pkg/gmcp/achaea"
	"github.com/tobiassjosten/nogfx/pkg/gmcp/ironrealms"

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
			msg:  &agmcp.CharStatus{},
			data: "Char.Status {}",
			unmarshaled: &agmcp.CharStatus{
				CharStatus: &ironrealms.CharStatus{
					CharStatus: &gmcp.CharStatus{},
				},
			},
			marshaled: "Char.Status {}",
		},

		"Char.Status hydrated": {
			msg: &agmcp.CharStatus{},
			data: makeGMCP("Char.Status", map[string]interface{}{
				"age":              "6",
				"boundcredits":     "7",
				"boundmayancrowns": "8",
				"explorerrank":     "an Itinerant",
				"house":            "The Dread Legates(9)",
				"lessons":          "10",
				"mayancrowns":      "11",
				"order":            "Something (12)",
				"specialisation":   "Brawler",
				"target":           "someone",
				"unboundcredits":   "13",
				"xprank":           "14",
				// Ironrealms:
				"bank":        "1",
				"city":        "Mhaldor (2)",
				"class":       "Monk",
				"gold":        "3",
				"unread_msgs": "4",
				"unread_news": "5",
				// Base:
				"fullname": "Mason Durak",
				"gender":   "male",
				"level":    "69 (23%)",
				"name":     "Durak",
				"race":     "dwarf",
			}),
			unmarshaled: &agmcp.CharStatus{
				CharStatus: &ironrealms.CharStatus{
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
				Age:              gox.NewInt(6),
				BoundCredits:     gox.NewInt(7),
				BoundMayanCrowns: gox.NewInt(8),
				ExplorerRank:     gox.NewString("an Itinerant"),
				House:            gox.NewString("The Dread Legates"),
				HouseRank:        gox.NewInt(9),
				Lessons:          gox.NewInt(10),
				MayanCrowns:      gox.NewInt(11),
				Order:            gox.NewString("Something"),
				OrderRank:        gox.NewInt(12),
				Specialisation:   gox.NewString("Brawler"),
				Target:           gox.NewString("someone"),
				UnboundCredits:   gox.NewInt(13),
				XPRank:           gox.NewInt(14),
			},
			marshaled: makeGMCP("Char.Status", map[string]interface{}{
				"age":              "6",
				"boundcredits":     "7",
				"boundmayancrowns": "8",
				"explorerrank":     "an Itinerant",
				"house":            "The Dread Legates (9)",
				"lessons":          "10",
				"mayancrowns":      "11",
				"order":            "Something (12)",
				"specialisation":   "Brawler",
				"target":           "someone",
				"unboundcredits":   "13",
				"xprank":           "14",
				// Ironrealms:
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
