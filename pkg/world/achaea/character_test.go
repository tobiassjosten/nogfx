package achaea_test

import (
	"fmt"
	"testing"

	"github.com/tobiassjosten/nogfx/pkg/gmcp"
	agmcp "github.com/tobiassjosten/nogfx/pkg/gmcp/achaea"
	"github.com/tobiassjosten/nogfx/pkg/world/achaea"

	"github.com/icza/gox/gox"
	"github.com/stretchr/testify/assert"
)

func TestFromGMCP(t *testing.T) {
	tcs := []struct {
		in      *achaea.Character
		message gmcp.Message
		out     *achaea.Character
	}{
		{
			in: &achaea.Character{},
			message: &gmcp.CharName{
				Name:     "Durak",
				Fullname: "Mason Durak",
			},
			out: &achaea.Character{
				Name:  "Durak",
				Title: "Mason Durak",
			},
		},

		{
			in: &achaea.Character{},
			message: &agmcp.CharStatus{
				Class:    gox.NewString("Monk"),
				Fullname: gox.NewString("Mason Durak"),
				Level:    gox.NewFloat64(68),
				Name:     gox.NewString("Durak"),
			},
			out: &achaea.Character{
				Name:  "Durak",
				Title: "Mason Durak",
				Level: 68,
				Class: "Monk",
			},
		},

		{
			in: &achaea.Character{},
			message: &agmcp.CharVitals{
				HP:    123,
				MaxHP: 124,
				MP:    234,
				MaxMP: 235,
				EP:    345,
				MaxEP: 346,
				WP:    456,
				MaxWP: 457,
				NL:    56,
				Bal:   true,
				Eq:    true,
				Vote:  true,
				Stats: agmcp.CharVitalsStats{
					Bleed:    12,
					Rage:     23,
					Ferocity: gox.NewInt(34),
					Kai:      gox.NewInt(45),
					Karma:    gox.NewInt(56),
					Spec:     gox.NewString("Sword and Shield"),
					Stance:   gox.NewString("Scorpion"),
				},
			},
			out: &achaea.Character{
				XP: 56,

				Balance:     true,
				Equilibrium: true,

				Health:       123,
				MaxHealth:    124,
				Mana:         234,
				MaxMana:      235,
				Endurance:    345,
				MaxEndurance: 346,
				Willpower:    456,
				MaxWillpower: 457,

				Bleed: 12,
				Rage:  23,

				Ferocity: 34,
				Kai:      45,
				Karma:    56,
				Spec:     "Sword and Shield",
				Stance:   "Scorpion",
			},
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			if msg, ok := tc.message.(*gmcp.CharName); ok {
				tc.in.FromCharName(msg)
			}

			if msg, ok := tc.message.(*agmcp.CharStatus); ok {
				tc.in.FromCharStatus(msg)
			}

			if msg, ok := tc.message.(*agmcp.CharVitals); ok {
				tc.in.FromCharVitals(msg)
			}

			assert.Equal(t, tc.out, tc.in)
		})
	}
}
