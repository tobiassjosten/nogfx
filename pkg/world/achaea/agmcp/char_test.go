package agmcp_test

import (
	"fmt"
	"testing"

	"github.com/icza/gox/gox"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tobiassjosten/nogfx/pkg/gmcp"
	"github.com/tobiassjosten/nogfx/pkg/world/achaea"
	"github.com/tobiassjosten/nogfx/pkg/world/achaea/agmcp"
)

func TestCharServerMessages(t *testing.T) {
	tcs := []struct {
		command []byte
		message gmcp.ServerMessage
		err     string
	}{
		{
			command: []byte("Char.Status"),
			err:     "failed hydrating agmcp.CharStatus: unexpected end of JSON input",
		},
		{
			command: []byte(`Char.Status { "nonexistant": {"invalid"] }`),
			err:     `failed hydrating agmcp.CharStatus: invalid character ']' after object key`,
		},
		{
			command: []byte(`Char.Status { "city": "invalid" }`),
			err:     `failed hydrating agmcp.CharStatus: failed parsing city 'invalid'`,
		},
		{
			command: []byte(`Char.Status { "house": "invalid" }`),
			err:     `failed hydrating agmcp.CharStatus: failed parsing house 'invalid'`,
		},
		{
			command: []byte(`Char.Status { "order": "invalid" }`),
			err:     `failed hydrating agmcp.CharStatus: failed parsing order 'invalid'`,
		},
		{
			command: []byte(`Char.Status { "city": "(None)", "house": "(None)", "order": "(None)", "target": "None" }`),
			message: agmcp.CharStatus{
				Target: gox.NewString(""),
			},
		},
		{
			command: []byte(`Char.Status { "target": "Someone (player)" }`),
			message: agmcp.CharStatus{
				Target: gox.NewString("Someone"),
			},
		},
		{
			command: []byte(`Char.Status { "name": "Durak", "fullname": "Mason Durak", "age": "184", "race": "Dwarf", "specialisation": "Brawler", "level": "68 (19%)", "xp": "19%", "xprank": "999", "class": "Monk", "city": "Hashan (1)", "house": "The Somatikos(1)", "order": "Blabla (1)", "boundcredits": "20", "unboundcredits": "1", "lessons": "4073", "explorerrank": "an Itinerant", "mayancrowns": "1", "boundmayancrowns": "2", "gold": "35", "bank": "1590", "unread_news": "3751", "unread_msgs": "1", "target": "123456", "gender": "male" }`),
			message: agmcp.CharStatus{
				CharStatus: gmcp.CharStatus{
					Name:     gox.NewString("Durak"),
					Fullname: gox.NewString("Mason Durak"),
					Age:      gox.NewInt(184),
					Race:     gox.NewString("Dwarf"),
					Level:    gox.NewFloat64(68),
					XP:       gox.NewInt(19),
					Gender:   gox.NewInt(1),
				},
				Specialisation:   gox.NewString("Brawler"),
				XPRank:           gox.NewInt(999),
				Class:            gox.NewString("Monk"),
				City:             gox.NewString("Hashan"),
				CityRank:         gox.NewInt(1),
				House:            gox.NewString("The Somatikos"),
				HouseRank:        gox.NewInt(1),
				Order:            gox.NewString("Blabla"),
				OrderRank:        gox.NewInt(1),
				BoundCredits:     gox.NewInt(20),
				UnboundCredits:   gox.NewInt(1),
				Lessons:          gox.NewInt(4073),
				ExplorerRank:     gox.NewString("an Itinerant"),
				MayanCrowns:      gox.NewInt(1),
				BoundMayanCrowns: gox.NewInt(2),
				Gold:             gox.NewInt(35),
				Bank:             gox.NewInt(1590),
				UnreadNews:       gox.NewInt(3751),
				UnreadMessages:   gox.NewInt(1),
				Target:           gox.NewString("123456"),
			},
		},

		{
			command: []byte("Char.Vitals"),
			err:     "failed hydrating agmcp.CharVitals: unexpected end of JSON input",
		},
		{
			command: []byte(`Char.Vitals { "charstats": [ "Invalid" ] }`),
			err:     `failed hydrating agmcp.CharVitals: misformed charstat 'Invalid'`,
		},
		{
			command: []byte(`Char.Vitals { "charstats": [ "Unknown: 123" ] }`),
			err:     `failed hydrating agmcp.CharVitals: invalid charstat 'Unknown: 123'`,
		},
		{
			command: []byte(`Char.Vitals { "charstats": [ "Bleed: invalid" ] }`),
			err:     `failed hydrating agmcp.CharVitals: invalid charstat 'Bleed: invalid'`,
		},
		{
			command: []byte(`Char.Vitals { "charstats": [ "Rage: invalid" ] }`),
			err:     `failed hydrating agmcp.CharVitals: invalid charstat 'Rage: invalid'`,
		},
		{
			command: []byte(`Char.Vitals { "charstats": [ "Ferocity: invalid" ] }`),
			err:     `failed hydrating agmcp.CharVitals: invalid charstat 'Ferocity: invalid'`,
		},
		{
			command: []byte(`Char.Vitals { "charstats": [ "Kai: invalid" ] }`),
			err:     `failed hydrating agmcp.CharVitals: invalid charstat 'Kai: invalid'`,
		},
		{
			command: []byte(`Char.Vitals { "charstats": [ "Kai: 1" ] }`),
			err:     `failed hydrating agmcp.CharVitals: invalid charstat 'Kai: 1'`,
		},
		{
			command: []byte(`Char.Vitals { "charstats": [ "Karma: invalid" ] }`),
			err:     `failed hydrating agmcp.CharVitals: invalid charstat 'Karma: invalid'`,
		},
		{
			command: []byte(`Char.Vitals { "charstats": [ "Karma: 1" ] }`),
			err:     `failed hydrating agmcp.CharVitals: invalid charstat 'Karma: 1'`,
		},
		{
			command: []byte(`Char.Vitals { "charstats": [ "Stance: None" ] }`),
			message: agmcp.CharVitals{
				Stats: agmcp.CharVitalsStats{Stance: nil},
			},
		},
		{
			command: []byte(`Char.Vitals { "hp": "3904", "maxhp": "3904", "mp": "3845", "maxmp": "3845", "ep": "15020", "maxep": "15020", "wp": "12980", "maxwp": "12980", "nl": "19", "bal": "1", "eq": "1", "vote": "1", "string": "H:3904/3904 M:3845/3845 E:15020/15020 W:12980/12980 NL:19/100 ", "charstats": [ "Bleed: 1", "Rage: 2", "Kai: 4%", "Karma: 5%", "Stance: Crane", "Ferocity: 3", "Spec: Sword and Shield" ] }`),
			message: agmcp.CharVitals{
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
				Stats: agmcp.CharVitalsStats{
					Bleed: 1,
					Rage:  2,

					Ferocity: gox.NewInt(3),
					Kai:      gox.NewInt(4),
					Karma:    gox.NewInt(5),
					Spec:     gox.NewString("Sword and Shield"),
					Stance:   gox.NewString("Crane"),
				},
			},
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			message, err := gmcp.Parse(tc.command, achaea.ServerMessages)

			if tc.err != "" {
				require.NotNil(err, fmt.Sprintf(
					"wanted: %s", tc.err,
				))
				assert.Equal(tc.err, err.Error())
				return
			}

			require.Nil(err)
			assert.Equal(tc.message, message)
		})
	}
}
