package gmcp_test

import (
	"fmt"
	"testing"

	"github.com/icza/gox/gox"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tobiassjosten/nogfx/pkg/world/achaea/gmcp"
)

func TestParse(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	tcs := []struct {
		command []byte
		message gmcp.Message
		err     error
	}{
		{
			command: []byte("Asdf"),
			err:     fmt.Errorf("unknown message 'Asdf'"),
		},
		{
			command: []byte("Char.Items.Inv"),
			message: gmcp.CharItemsInv{},
		},
		{
			command: []byte("Char.Vitals"),
			err:     fmt.Errorf("missing 'gmcp.CharVitals' data"),
		},
		{
			command: []byte(`Char.Vitals { "hp": "3904", "maxhp": "3904", "mp": "3845", "maxmp": "3845", "ep": "15020", "maxep": "15020", "wp": "12980", "maxwp": "12980", "nl": "19", "bal": "1", "eq": "1", "vote": "1", "string": "H:3904/3904 M:3845/3845 E:15020/15020 W:12980/12980 NL:19/100 ", "charstats": [ "Bleed: 1", "Rage: 2", "Kai: 0%", "Stance: None" ] }`),
			message: gmcp.CharVitals{
				HP:     3904,
				MaxHP:  3904,
				MP:     3845,
				MaxMP:  3845,
				EP:     15020,
				MaxEP:  15020,
				WP:     12980,
				MaxWP:  12980,
				NL:     19,
				Bal:    true,
				Eq:     true,
				Vote:   true,
				Prompt: "H:3904/3904 M:3845/3845 E:15020/15020 W:12980/12980 NL:19/100 ",
				Stats: gmcp.CharVitalsStats{
					Bleed:  1,
					Kai:    gox.NewInt(0),
					Rage:   2,
					Stance: gox.NewString("None"),
				},
			},
		},
		{
			command: []byte(`Char.Status { "name": "Durak", "fullname": "Mason Durak", "age": "184", "race": "Dwarf", "specialisation": "Brawler", "level": "68 (19%)", "xp": "19%", "xprank": "999", "class": "Monk", "city": "Hashan (1)", "house": "The Somatikos(1)", "order": "(None)", "boundcredits": "20", "unboundcredits": "1", "lessons": "4073", "explorerrank": "an Itinerant", "mayancrowns": "1", "boundmayancrowns": "2", "gold": "35", "bank": "1590", "unread_news": "3751", "unread_msgs": "1", "target": "None", "gender": "male" }`),
			message: gmcp.CharStatus{
				Name:             "Durak",
				Fullname:         "Mason Durak",
				Age:              184,
				Race:             "Dwarf",
				Specialisation:   "Brawler",
				Level:            68,
				XP:               19,
				XPRank:           999,
				Class:            "Monk",
				City:             gox.NewString("Hashan"),
				CityRank:         gox.NewInt(1),
				House:            gox.NewString("The Somatikos"),
				HouseRank:        gox.NewInt(1),
				Order:            nil,
				BoundCredits:     20,
				UnboundCredits:   1,
				Lessons:          4073,
				ExplorerRank:     "an Itinerant",
				MayanCrowns:      1,
				BoundMayanCrowns: 2,
				Gold:             35,
				Bank:             1590,
				UnreadNews:       3751,
				UnreadMessages:   1,
				Target:           nil,
				Gender:           1,
			},
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			message, err := gmcp.Parse(tc.command)

			if tc.err != nil {
				assert.Equal(tc.err, err)
				return
			}

			require.Nil(err)
			assert.Equal(tc.message, message)
		})
	}
}
