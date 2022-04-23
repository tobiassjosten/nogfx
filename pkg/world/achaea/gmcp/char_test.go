package gmcp_test

import (
	"fmt"
	"testing"

	"github.com/icza/gox/gox"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tobiassjosten/nogfx/pkg/world/achaea/gmcp"
)

func TestCharClientMessages(t *testing.T) {
	tcs := []struct {
		message gmcp.ClientMessage
		output  string
	}{
		{
			message: gmcp.CharItemsContents{},
			output:  "Char.Items.Contents",
		},
		{
			message: gmcp.CharItemsContents{ID: "123"},
			output:  "Char.Items.Contents 123",
		},

		{
			message: gmcp.CharItemsInv{},
			output:  "Char.Items.Inv",
		},

		{
			message: gmcp.CharItemsRoom{},
			output:  "Char.Items.Room",
		},

		{
			message: gmcp.CharLogin{},
			output:  "Char.Login {}",
		},
		{
			message: gmcp.CharLogin{Name: "durak", Password: "s3cr3t"},
			output:  `Char.Login {"name":"durak","password":"s3cr3t"}`,
		},

		{
			message: gmcp.CharSkillsGet{},
			output:  "Char.Skills.Get {}",
		},
		{
			message: gmcp.CharSkillsGet{
				Group: "Perception",
			},
			output: `Char.Skills.Get {"group":"Perception"}`,
		},
		{
			message: gmcp.CharSkillsGet{
				Group: "Perception",
				Name:  "Deathsight",
			},
			output: `Char.Skills.Get {"group":"Perception","name":"Deathsight"}`,
		},
		{
			message: gmcp.CharSkillsGet{
				Name: "Deathsight",
			},
			output: `Char.Skills.Get {}`,
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			assert := assert.New(t)
			assert.Equal(tc.output, tc.message.String())
		})
	}
}

func TestCharServerMessages(t *testing.T) {
	tcs := []struct {
		command []byte
		message gmcp.ServerMessage
		err     string
	}{
		{
			command: []byte("Char.Afflictions.List"),
			err:     "failed hydrating gmcp.CharAfflictionsList (): unexpected end of JSON input",
		},
		{
			command: []byte(`Char.Afflictions.List [ { "name": "weariness", "cure": "eat kelp", "desc": "Decreases cutting and blunt damage that you inflict by 30%." }, { "name": "asthma", "cure": "eat kelp", "desc": "Makes you unable to smoke pipes." } ]`),
			message: gmcp.CharAfflictionsList{
				{
					Name:        "weariness",
					Cure:        "eat kelp",
					Description: "Decreases cutting and blunt damage that you inflict by 30%.",
				},
				{
					Name:        "asthma",
					Cure:        "eat kelp",
					Description: "Makes you unable to smoke pipes.",
				},
			},
		},

		{
			command: []byte("Char.Afflictions.Add"),
			err:     "failed hydrating gmcp.CharAfflictionsAdd (): unexpected end of JSON input",
		},
		{
			command: []byte(`Char.Afflictions.Add { "name": "weariness", "cure": "eat kelp", "desc": "Decreases cutting and blunt damage that you inflict by 30%." }`),
			message: gmcp.CharAfflictionsAdd{
				Name:        "weariness",
				Cure:        "eat kelp",
				Description: "Decreases cutting and blunt damage that you inflict by 30%.",
			},
		},

		{
			command: []byte("Char.Afflictions.Remove"),
			err:     "failed hydrating gmcp.CharAfflictionsRemove (): unexpected end of JSON input",
		},
		{
			command: []byte(`Char.Afflictions.Remove [ "weariness" ]`),
			message: gmcp.CharAfflictionsRemove{
				{
					Name:        "weariness",
					Description: "",
				},
			},
		},

		{
			command: []byte("Char.Defences.List"),
			err:     "failed hydrating gmcp.CharDefencesList (): unexpected end of JSON input",
		},
		{
			command: []byte(`Char.Defences.List [ { "name": "deaf", "desc": "deaf" }, { "name": "blind", "desc": "blind" } ]`),
			message: gmcp.CharDefencesList{
				{
					Name:        "deaf",
					Description: "deaf",
				},
				{
					Name:        "blind",
					Description: "blind",
				},
			},
		},

		{
			command: []byte("Char.Defences.Add"),
			err:     "failed hydrating gmcp.CharDefencesAdd (): unexpected end of JSON input",
		},
		{
			command: []byte(`Char.Defences.Add { "name": "deaf", "desc": "deaf" }`),
			message: gmcp.CharDefencesAdd{
				Name:        "deaf",
				Description: "deaf",
			},
		},

		{
			command: []byte("Char.Defences.Remove"),
			err:     "failed hydrating gmcp.CharDefencesRemove (): unexpected end of JSON input",
		},
		{
			command: []byte(`Char.Defences.Remove [ "deaf" ]`),
			message: gmcp.CharDefencesRemove{
				{
					Name:        "deaf",
					Description: "",
				},
			},
		},

		{
			command: []byte("Char.Items.List"),
			err:     "failed hydrating gmcp.CharItemsList (): unexpected end of JSON input",
		},
		{
			command: []byte(`Char.Items.List { "location": "inv", "items": [ { "id": "123", "name": "something", "icon": "commodity", "attrib": "awWlLgcrfemdtx" } ] }`),
			err:     `failed hydrating gmcp.CharItemsList ({ "location": "inv", "items": [ { "id": "123", "name": "something", "icon": "commodity", "attrib": "awWlLgcrfemdtx" } ] }): unknown attribute 'a'`,
		},
		{
			command: []byte(`Char.Items.List { "location": "inv", "items": [ { "id": "123", "name": "something", "icon": "commodity", "attrib": "wWlLgcrfemdtx" } ] }`),
			message: gmcp.CharItemsList{
				Location: "inv",
				Items: []gmcp.CharItem{{
					ID:   "123",
					Name: "something",
					Icon: "commodity",
					Attributes: gmcp.CharItemAttributes{
						Container:    true,
						Dangerous:    true,
						Dead:         true,
						Edible:       true,
						Fluid:        true,
						Groupable:    true,
						Monster:      true,
						Riftable:     true,
						Takeable:     true,
						Wearable:     true,
						WieldedLeft:  true,
						WieldedRight: true,
						Worn:         true,
					},
				}},
			},
		},

		{
			command: []byte("Char.Items.Add"),
			err:     "failed hydrating gmcp.CharItemsAdd (): unexpected end of JSON input",
		},
		{
			command: []byte(`Char.Items.Add { "location": "room", "item": { "id": "239602", "name": "an elegant white letter", "icon": "container", "attrib": "c" } }`),
			message: gmcp.CharItemsAdd{
				Location: "room",
				Item: gmcp.CharItem{
					ID:         "239602",
					Name:       "an elegant white letter",
					Icon:       "container",
					Attributes: gmcp.CharItemAttributes{Container: true},
				},
			},
		},

		{
			command: []byte("Char.Items.Remove"),
			err:     "failed hydrating gmcp.CharItemsRemove (): unexpected end of JSON input",
		},
		{
			command: []byte(`Char.Items.Remove { "location": "inv", "item": { "id": "303060", "name": "a gold nugget", "icon": "commodity", "attrib": "t" } }`),
			message: gmcp.CharItemsRemove{
				Location: "inv",
				Item: gmcp.CharItem{
					ID:         "303060",
					Name:       "a gold nugget",
					Icon:       "commodity",
					Attributes: gmcp.CharItemAttributes{Takeable: true},
				},
			},
		},

		{
			command: []byte("Char.Items.Update"),
			err:     "failed hydrating gmcp.CharItemsUpdate (): unexpected end of JSON input",
		},
		{
			command: []byte(`Char.Items.Update { "location": "inv", "item": { "id": "60572", "name": "an ornate steel rapier" } }`),
			message: gmcp.CharItemsUpdate{
				Location: "inv",
				Item: gmcp.CharItem{
					ID:   "60572",
					Name: "an ornate steel rapier",
				},
			},
		},

		{
			command: []byte("Char.Name"),
			err:     "failed hydrating gmcp.CharName (): unexpected end of JSON input",
		},
		{
			command: []byte(`Char.Name { "name": {"invalid"] }`),
			err:     `failed hydrating gmcp.CharName ({ "name": {"invalid"] }): invalid character ']' after object key`,
		},
		{
			command: []byte(`Char.Name { "name": "Durak", "fullname": "Mason Durak" }`),
			message: gmcp.CharName{
				Name:     "Durak",
				Fullname: "Mason Durak",
			},
		},

		{
			command: []byte("Char.Skills.Groups"),
			err:     "failed hydrating gmcp.CharSkillsGroups (): unexpected end of JSON input",
		},
		{
			command: []byte(`Char.Skills.Groups [ { "name": "Perception", "rank": "Adept" }  ]`),
			err:     `failed hydrating gmcp.CharSkillsGroups ([ { "name": "Perception", "rank": "Adept" }  ]): failed parsing rank 'Adept'`,
		},
		{
			command: []byte(`Char.Skills.Groups [ { "name": "Perception", "rank": "Adept (x%)" }  ]`),
			err:     `failed hydrating gmcp.CharSkillsGroups ([ { "name": "Perception", "rank": "Adept (x%)" }  ]): failed parsing rank 'Adept (x%)'`,
		},
		{
			command: []byte(`Char.Skills.Groups [ { "name": "Perception", "rank": "Adept (1%)" }  ]`),
			message: gmcp.CharSkillsGroups{
				{
					Name:     "Perception",
					Level:    "Adept",
					Progress: 1,
				},
			},
		},

		{
			command: []byte("Char.Skills.List"),
			err:     "failed hydrating gmcp.CharSkillsList (): unexpected end of JSON input",
		},
		{
			command: []byte(`Char.Skills.List { "group": "perception", "list": [ "Deathsight" ], "descs": [ "Using this ability…" ] }`),
			message: gmcp.CharSkillsList{
				Group:        "perception",
				List:         []string{"Deathsight"},
				Descriptions: []string{"Using this ability…"},
			},
		},

		{
			command: []byte("Char.Skills.Info"),
			err:     "failed hydrating gmcp.CharSkillsInfo (): unexpected end of JSON input",
		},
		{
			command: []byte(`Char.Skills.Info { "group": "Perception", "skill": "deathsight", "info": "Using this ability…" }`),
			message: gmcp.CharSkillsInfo{
				Group:       "Perception",
				Skill:       "deathsight",
				Information: "Using this ability…",
			},
		},

		{
			command: []byte("Char.StatusVars"),
			err:     "failed hydrating gmcp.CharStatusVars (): unexpected end of JSON input",
		},
		{
			command: []byte(`Char.StatusVars { "level": "Level" }`),
			message: gmcp.CharStatusVars(map[string]string{"level": "Level"}),
		},

		{
			command: []byte("Char.Status"),
			err:     "failed hydrating gmcp.CharStatus (): unexpected end of JSON input",
		},
		{
			command: []byte(`Char.Status { "nonexistant": {"invalid"] }`),
			err:     `failed hydrating gmcp.CharStatus ({ "nonexistant": {"invalid"] }): invalid character ']' after object key`,
		},
		{
			command: []byte(`Char.Status { "level": "invalid" }`),
			err:     `failed hydrating gmcp.CharStatus ({ "level": "invalid" }): failed parsing level 'invalid'`,
		},
		{
			command: []byte(`Char.Status { "city": "invalid" }`),
			err:     `failed hydrating gmcp.CharStatus ({ "city": "invalid" }): failed parsing city 'invalid'`,
		},
		{
			command: []byte(`Char.Status { "house": "invalid" }`),
			err:     `failed hydrating gmcp.CharStatus ({ "house": "invalid" }): failed parsing house 'invalid'`,
		},
		{
			command: []byte(`Char.Status { "order": "invalid" }`),
			err:     `failed hydrating gmcp.CharStatus ({ "order": "invalid" }): failed parsing order 'invalid'`,
		},
		{
			command: []byte(`Char.Status { "target": "invalid" }`),
			err:     `failed hydrating gmcp.CharStatus ({ "target": "invalid" }): failed parsing target 'invalid'`,
		},
		{
			command: []byte(`Char.Status { "city": "(None)", "house": "(None)", "order": "(None)", "target": "None" }`),
			message: gmcp.CharStatus{},
		},
		{
			command: []byte(`Char.Status { "gender": "female" }`),
			message: gmcp.CharStatus{Gender: gox.NewInt(2)},
		},
		{
			command: []byte(`Char.Status { "gender": "invalid" }`),
			message: gmcp.CharStatus{Gender: gox.NewInt(9)},
		},
		{
			command: []byte(`Char.Status { "name": "Durak", "fullname": "Mason Durak", "age": "184", "race": "Dwarf", "specialisation": "Brawler", "level": "68 (19%)", "xp": "19%", "xprank": "999", "class": "Monk", "city": "Hashan (1)", "house": "The Somatikos(1)", "order": "Blabla (1)", "boundcredits": "20", "unboundcredits": "1", "lessons": "4073", "explorerrank": "an Itinerant", "mayancrowns": "1", "boundmayancrowns": "2", "gold": "35", "bank": "1590", "unread_news": "3751", "unread_msgs": "1", "target": "123456", "gender": "male" }`),
			message: gmcp.CharStatus{
				Name:             gox.NewString("Durak"),
				Fullname:         gox.NewString("Mason Durak"),
				Age:              gox.NewInt(184),
				Race:             gox.NewString("Dwarf"),
				Specialisation:   gox.NewString("Brawler"),
				Level:            gox.NewInt(68),
				XP:               gox.NewInt(19),
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
				Target:           gox.NewInt(123456),
				Gender:           gox.NewInt(1),
			},
		},

		{
			command: []byte("Char.Vitals"),
			err:     "failed hydrating gmcp.CharVitals (): unexpected end of JSON input",
		},
		{
			command: []byte(`Char.Vitals { "charstats": [ "Invalid" ] }`),
			err:     `failed hydrating gmcp.CharVitals ({ "charstats": [ "Invalid" ] }): misformed charstat 'Invalid'`,
		},
		{
			command: []byte(`Char.Vitals { "charstats": [ "Unknown: 123" ] }`),
			err:     `failed hydrating gmcp.CharVitals ({ "charstats": [ "Unknown: 123" ] }): invalid charstat 'Unknown: 123'`,
		},
		{
			command: []byte(`Char.Vitals { "charstats": [ "Bleed: invalid" ] }`),
			err:     `failed hydrating gmcp.CharVitals ({ "charstats": [ "Bleed: invalid" ] }): invalid charstat 'Bleed: invalid'`,
		},
		{
			command: []byte(`Char.Vitals { "charstats": [ "Rage: invalid" ] }`),
			err:     `failed hydrating gmcp.CharVitals ({ "charstats": [ "Rage: invalid" ] }): invalid charstat 'Rage: invalid'`,
		},
		{
			command: []byte(`Char.Vitals { "charstats": [ "Ferocity: invalid" ] }`),
			err:     `failed hydrating gmcp.CharVitals ({ "charstats": [ "Ferocity: invalid" ] }): invalid charstat 'Ferocity: invalid'`,
		},
		{
			command: []byte(`Char.Vitals { "charstats": [ "Kai: invalid" ] }`),
			err:     `failed hydrating gmcp.CharVitals ({ "charstats": [ "Kai: invalid" ] }): invalid charstat 'Kai: invalid'`,
		},
		{
			command: []byte(`Char.Vitals { "charstats": [ "Kai: 1" ] }`),
			err:     `failed hydrating gmcp.CharVitals ({ "charstats": [ "Kai: 1" ] }): invalid charstat 'Kai: 1'`,
		},
		{
			command: []byte(`Char.Vitals { "charstats": [ "Karma: invalid" ] }`),
			err:     `failed hydrating gmcp.CharVitals ({ "charstats": [ "Karma: invalid" ] }): invalid charstat 'Karma: invalid'`,
		},
		{
			command: []byte(`Char.Vitals { "charstats": [ "Karma: 1" ] }`),
			err:     `failed hydrating gmcp.CharVitals ({ "charstats": [ "Karma: 1" ] }): invalid charstat 'Karma: 1'`,
		},
		{
			command: []byte(`Char.Vitals { "charstats": [ "Stance: None" ] }`),
			message: gmcp.CharVitals{
				Stats: gmcp.CharVitalsStats{Stance: nil},
			},
		},
		{
			command: []byte(`Char.Vitals { "hp": "3904", "maxhp": "3904", "mp": "3845", "maxmp": "3845", "ep": "15020", "maxep": "15020", "wp": "12980", "maxwp": "12980", "nl": "19", "bal": "1", "eq": "1", "vote": "1", "string": "H:3904/3904 M:3845/3845 E:15020/15020 W:12980/12980 NL:19/100 ", "charstats": [ "Bleed: 1", "Rage: 2", "Kai: 4%", "Karma: 5%", "Stance: Crane", "Ferocity: 3", "Spec: Sword and Shield" ] }`),
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

			message, err := gmcp.Parse(tc.command)

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
