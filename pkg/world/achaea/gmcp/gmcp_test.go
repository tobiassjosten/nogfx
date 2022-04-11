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
			err:     fmt.Errorf("invalid message 'Asdf'"),
		},
		{
			command: []byte("Char.Items.Inv"),
			message: gmcp.CharItemsInv{},
		},
		{
			command: []byte("Char.Vitals"),
			err:     fmt.Errorf("missing 'Char.Vitals' data"),
		},
		{
			command: []byte(`Char.Vitals { "hp": "3904", "maxhp": "3904", "mp": "3845", "maxmp": "3845", "ep": "15020", "maxep": "15020", "wp": "12980", "maxwp": "12980", "nl": "19", "bal": "1", "eq": "1", "vote": "1", "string": "H:3904/3904 M:3845/3845 E:15020/15020 W:12980/12980 NL:19/100 ", "charstats": [ "Bleed: 0", "Rage: 0", "Kai: 0%", "Stance: None" ] }`),
			message: gmcp.CharVitals{
				HP:    gox.NewInt(3904),
				MaxHP: gox.NewInt(3904),
				MP:    gox.NewInt(3845),
				MaxMP: gox.NewInt(3845),
				EP:    gox.NewInt(15020),
				MaxEP: gox.NewInt(15020),
				WP:    gox.NewInt(12980),
				MaxWP: gox.NewInt(12980),
				NL:    gox.NewInt(19),
				Bal:   gox.NewBool(true),
				Eq:    gox.NewBool(true),
				Vote:  gox.NewBool(true),
				Prompt: gox.NewString(
					"H:3904/3904 M:3845/3845 E:15020/15020 W:12980/12980 NL:19/100 ",
				),
				Stats: gmcp.CharVitalsStats{
					Bleed:  gox.NewInt(0),
					Kai:    gox.NewInt(0),
					Rage:   gox.NewInt(0),
					Stance: gox.NewString("None"),
				},
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
