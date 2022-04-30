package gmcp_test

import (
	"fmt"
	"testing"

	"github.com/tobiassjosten/nogfx/pkg/gmcp"

	"github.com/icza/gox/gox"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			err:     "failed hydrating gmcp.CharAfflictionsList: unexpected end of JSON input",
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
			err:     "failed hydrating gmcp.CharAfflictionsAdd: unexpected end of JSON input",
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
			err:     "failed hydrating gmcp.CharAfflictionsRemove: unexpected end of JSON input",
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
			err:     "failed hydrating gmcp.CharDefencesList: unexpected end of JSON input",
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
			err:     "failed hydrating gmcp.CharDefencesAdd: unexpected end of JSON input",
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
			err:     "failed hydrating gmcp.CharDefencesRemove: unexpected end of JSON input",
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
			err:     "failed hydrating gmcp.CharItemsList: unexpected end of JSON input",
		},
		{
			command: []byte(`Char.Items.List { "location": "inv", "items": [ { "id": "123", "name": "something", "icon": "commodity", "attrib": "awWlLgcrfemdtx" } ] }`),
			err:     `failed hydrating gmcp.CharItemsList: unknown attribute 'a'`,
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
			err:     "failed hydrating gmcp.CharItemsAdd: unexpected end of JSON input",
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
			err:     "failed hydrating gmcp.CharItemsRemove: unexpected end of JSON input",
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
			err:     "failed hydrating gmcp.CharItemsUpdate: unexpected end of JSON input",
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
			err:     "failed hydrating gmcp.CharName: unexpected end of JSON input",
		},
		{
			command: []byte(`Char.Name { "name": {"invalid"] }`),
			err:     `failed hydrating gmcp.CharName: invalid character ']' after object key`,
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
			err:     "failed hydrating gmcp.CharSkillsGroups: unexpected end of JSON input",
		},
		{
			command: []byte(`Char.Skills.Groups [ { "name": "Perception", "rank": "Adept" }  ]`),
			message: gmcp.CharSkillsGroups{
				{
					Name:  "Perception",
					Level: "Adept",
				},
			},
		},
		{
			command: []byte(`Char.Skills.Groups [ { "name": "Perception", "rank": "Adept (1%)" }  ]`),
			message: gmcp.CharSkillsGroups{
				{
					Name:     "Perception",
					Level:    "Adept",
					Progress: gox.NewInt(1),
				},
			},
		},

		{
			command: []byte("Char.Skills.List"),
			err:     "failed hydrating gmcp.CharSkillsList: unexpected end of JSON input",
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
			err:     "failed hydrating gmcp.CharSkillsInfo: unexpected end of JSON input",
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
			err:     "failed hydrating gmcp.CharStatusVars: unexpected end of JSON input",
		},
		{
			command: []byte(`Char.StatusVars { "level": "Level" }`),
			message: gmcp.CharStatusVars(map[string]string{"level": "Level"}),
		},

		{
			command: []byte("Char.Status"),
			err:     "failed hydrating gmcp.CharStatus: unexpected end of JSON input",
		},
		{
			command: []byte(`Char.Status { "nonexistant": {"invalid"] }`),
			err:     `failed hydrating gmcp.CharStatus: invalid character ']' after object key`,
		},
		{
			command: []byte(`Char.Status { "level": "invalid" }`),
			err:     `failed hydrating gmcp.CharStatus: failed parsing level 'invalid'`,
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
			command: []byte(`Char.Status { "name": "Durak", "fullname": "Mason Durak", "age": "184", "race": "Dwarf", "level": "68 (19%)", "xp": "19%", "gender": "male" }`),
			message: gmcp.CharStatus{
				Name:     gox.NewString("Durak"),
				Fullname: gox.NewString("Mason Durak"),
				Age:      gox.NewInt(184),
				Race:     gox.NewString("Dwarf"),
				Level:    gox.NewFloat64(68),
				XP:       gox.NewInt(19),
				Gender:   gox.NewInt(1),
			},
		},

		{
			command: []byte("Char.Vitals"),
			err:     "failed hydrating gmcp.CharVitals: unexpected end of JSON input",
		},
		{
			command: []byte(`Char.Vitals { "hp": "3904", "maxhp": "3904", "mp": "3845", "maxmp": "3845", "ep": "15020", "maxep": "15020", "wp": "12980", "maxwp": "12980", "nl": "19", "bal": "1", "eq": "1", "vote": "1", "string": "H:3904/3904 M:3845/3845 E:15020/15020 W:12980/12980 NL:19/100 ", "charstats": [] }`),
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
				Stats:  gmcp.CharVitalsStats{},
			},
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			message, err := gmcp.Parse(tc.command, gmcp.ServerMessages)

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