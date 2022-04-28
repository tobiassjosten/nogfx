package igmcp_test

import (
	"fmt"
	"testing"

	"github.com/icza/gox/gox"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tobiassjosten/nogfx/pkg/gmcp"
	"github.com/tobiassjosten/nogfx/pkg/world/imperian"
	"github.com/tobiassjosten/nogfx/pkg/world/imperian/igmcp"
)

func TestCharServerMessages(t *testing.T) {
	tcs := []struct {
		command []byte
		message gmcp.ServerMessage
		err     string
	}{
		{
			command: []byte("Char.Status"),
			err:     "failed hydrating igmcp.CharStatus: unexpected end of JSON input",
		},
		{
			command: []byte(`Char.Status { "nonexistant": {"invalid"] }`),
			err:     `failed hydrating igmcp.CharStatus: invalid character ']' after object key`,
		},
		{
			command: []byte(`Char.Status { "name": "Durak", "fullname": "Mason Durak", "age": "184", "race": "Dwarf", "level": "68 (19%)", "xp": "19%", "class": "Monk", "gender": "male" }`),
			message: igmcp.CharStatus{
				CharStatus: gmcp.CharStatus{
					Name:     gox.NewString("Durak"),
					Fullname: gox.NewString("Mason Durak"),
					Age:      gox.NewInt(184),
					Race:     gox.NewString("Dwarf"),
					Level:    gox.NewFloat64(68),
					XP:       gox.NewInt(19),
					Gender:   gox.NewInt(1),
				},
				Class: gox.NewString("Monk"),
			},
		},

		{
			command: []byte("Char.Vitals"),
			err:     "failed hydrating igmcp.CharVitals: unexpected end of JSON input",
		},
		{
			command: []byte(`Char.Vitals { "charstats": [ "Invalid" ] }`),
			err:     `failed hydrating igmcp.CharVitals: misformed charstat 'Invalid'`,
		},
		{
			command: []byte(`Char.Vitals { "charstats": [ "Unknown: 123" ] }`),
			err:     `failed hydrating igmcp.CharVitals: invalid charstat 'Unknown: 123'`,
		},
		{
			command: []byte(`Char.Vitals { "charstats": [ "Bleed: invalid" ] }`),
			err:     `failed hydrating igmcp.CharVitals: invalid charstat 'Bleed: invalid'`,
		},
		{
			command: []byte(`Char.Vitals { "charstats": [ "Rage: invalid" ] }`),
			err:     `failed hydrating igmcp.CharVitals: invalid charstat 'Rage: invalid'`,
		},
		{
			command: []byte(`Char.Vitals { "charstats": [ "Ferocity: invalid" ] }`),
			err:     `failed hydrating igmcp.CharVitals: invalid charstat 'Ferocity: invalid'`,
		},
		{
			command: []byte(`Char.Vitals { "charstats": [ "Kai: invalid" ] }`),
			err:     `failed hydrating igmcp.CharVitals: invalid charstat 'Kai: invalid'`,
		},
		{
			command: []byte(`Char.Vitals { "charstats": [ "Kai: 1" ] }`),
			err:     `failed hydrating igmcp.CharVitals: invalid charstat 'Kai: 1'`,
		},
		{
			command: []byte(`Char.Vitals { "charstats": [ "Karma: invalid" ] }`),
			err:     `failed hydrating igmcp.CharVitals: invalid charstat 'Karma: invalid'`,
		},
		{
			command: []byte(`Char.Vitals { "charstats": [ "Karma: 1" ] }`),
			err:     `failed hydrating igmcp.CharVitals: invalid charstat 'Karma: 1'`,
		},
		{
			command: []byte(`Char.Vitals { "charstats": [ "Stance: None" ] }`),
			message: igmcp.CharVitals{
				Stats: igmcp.CharVitalsStats{Stance: nil},
			},
		},
		{
			command: []byte(`Char.Vitals { "hp": "3904", "maxhp": "3904", "mp": "3845", "maxmp": "3845", "ep": "15020", "maxep": "15020", "wp": "12980", "maxwp": "12980", "nl": "19", "bal": "1", "eq": "1", "vote": "1", "string": "H:3904/3904 M:3845/3845 E:15020/15020 W:12980/12980 NL:19/100 ", "charstats": [ "Bleed: 1", "Rage: 2", "Kai: 4%", "Karma: 5%", "Stance: Crane", "Ferocity: 3", "Spec: Sword and Shield" ] }`),
			message: igmcp.CharVitals{
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
				Stats: igmcp.CharVitalsStats{
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

			message, err := gmcp.Parse(tc.command, imperian.ServerMessages)

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
