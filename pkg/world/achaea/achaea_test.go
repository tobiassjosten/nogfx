package achaea_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/gmcp"
	"github.com/tobiassjosten/nogfx/pkg/mock"
	"github.com/tobiassjosten/nogfx/pkg/telnet"
	tst "github.com/tobiassjosten/nogfx/pkg/testing"
	"github.com/tobiassjosten/nogfx/pkg/world/achaea"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func wrapGMCP(msg string, data any) (bs []byte) {
	content := []byte(msg)

	if data != nil {
		mdata, _ := json.Marshal(data)
		content = append(content, append([]byte{' '}, mdata...)...)
	}

	bs = append(bs, telnet.IAC, telnet.SB, telnet.GMCP)
	bs = append(bs, content...)
	bs = append(bs, telnet.IAC, telnet.SE)

	return bs
}

func TestInputOutput(t *testing.T) {
	tcs := map[string]struct {
		Events    []tst.IOEvent
		Inoutputs []pkg.Inoutput
	}{
		"separated repeated input": {
			Events: []tst.IOEvent{
				tst.IOEIn("qwer;2 asdf;zxcv"),
			},
			Inoutputs: []pkg.Inoutput{
				tst.IOIns([]string{"qwer", "asdf", "zxcv"}).
					AddAfterInput(1, []byte("asdf")),
			},
		},

		"extranous ga newline": {
			Events: []tst.IOEvent{
				tst.IOEOuts([]string{
					"\033[35m",
					"asdf",
					"123h 234m\0371",
				}),
			},
			Inoutputs: []pkg.Inoutput{
				tst.IOOuts([]string{
					"\033[35m",
					"\033[35masdf",
					"123h 234m\0371",
				}).
					OmitOutput(0),
			},
		},

		"single prompts": {
			Events: []tst.IOEvent{
				tst.IOEOut("123h 234m\0371"),
			},
			Inoutputs: []pkg.Inoutput{
				{Output: pkg.Exput{}},
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			world, ok := achaea.NewWorld(
				&mock.ClientMock{},
				&mock.UIMock{},
			).(*achaea.World)
			require.True(t, ok)

			var inouts []pkg.Inoutput

			for _, event := range tc.Events {
				inout := world.OnInoutput(event.Inoutput())
				inouts = append(inouts, inout)
			}

			assert.Equal(t, len(tc.Inoutputs), len(inouts))
			for i, inout := range tc.Inoutputs {
				if i >= len(inouts) {
					break
				}
				assert.Equal(t, inout, inouts[i], fmt.Sprintf("index %d", i))
			}
		})
	}
}

func TestCommandsReply(t *testing.T) {
	tcs := []struct {
		command []byte
		sent    []byte
		errs    []bool
	}{
		{
			command: []byte{telnet.IAC, telnet.WILL, telnet.GMCP},
			sent: wrapGMCP("Core.Supports.Set", []string{
				"Char 1", "Char.Items 1", "Char.Skills 1",
				"Comm.Channel 1", "IRE.Rift 1", "IRE.Target 1",
				"Room 1",
			}),
		},
		{
			command: []byte{telnet.IAC, telnet.WILL, telnet.GMCP},
			errs:    []bool{true},
		},

		{
			command: wrapGMCP("Char.Name", map[string]string{
				"name":     "Durak",
				"fullname": "Mason Durak",
			}),
			sent: append(wrapGMCP("Char.Items.Inv", nil),
				append(wrapGMCP("Comm.Channel.Players", nil),
					wrapGMCP("IRE.Rift.Request", nil)...)...),
		},
		{
			command: wrapGMCP("Char.Name", map[int]int{}),
			errs:    []bool{true},
		},
		{
			command: wrapGMCP("Char.Name", map[int]int{}),
			errs:    []bool{false, true},
		},
		{
			command: wrapGMCP("Char.Name", map[int]int{}),
			errs:    []bool{false, false, true},
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			var calls int
			var sent []byte

			client := &mock.ClientMock{
				WriteFunc: func(data []byte) (int, error) {
					defer func() { calls++ }()
					if len(tc.errs) > calls && tc.errs[calls] {
						return 0, errors.New("ooops")
					}

					sent = append(sent, data...)

					return len(data), nil
				},
			}

			ui := &mock.UIMock{}

			world := achaea.NewWorld(client, ui)

			world.OnCommand(tc.command)

			if len(tc.sent) > 0 {
				assert.Equal(t, tc.sent, sent, string(sent))
			}
		})
	}
}

func TestCommandsMutateWorld(t *testing.T) {
	tcs := []struct {
		command   []byte
		messages  []gmcp.Message
		character *achaea.Character
		target    *pkg.Target
	}{
		{
			messages: []gmcp.Message{
				&gmcp.CharName{
					Name:     "Durak",
					Fullname: "Mason Durak",
				},
			},
			character: &achaea.Character{
				Name:  "Durak",
				Title: "Mason Durak",
			},
		},

		{
			command: wrapGMCP("Char.Status", map[string]string{
				"name":     "Durak",
				"fullname": "Mason Durak",
				"class":    "Monk",
				"level":    "68 (19%)",
			}),
			character: &achaea.Character{
				Name:  "Durak",
				Title: "Mason Durak",
				Class: "Monk",
				Level: 68,
			},
		},

		{
			command: wrapGMCP("Char.Vitals", map[string]any{
				"hp":    "3904",
				"maxhp": "3905",
				"mp":    "3845",
				"maxmp": "3846",
				"ep":    "15020",
				"maxep": "15021",
				"wp":    "12980",
				"maxwp": "12981",
				"bal":   "1",
				"eq":    "1",
				"charstats": []string{
					"Bleed: 1", "Rage: 2", "Kai: 4%",
					"Karma: 5%", "Stance: Crane",
					"Ferocity: 3", "Spec: Sword and Shield",
				},
			}),
			character: &achaea.Character{
				Balance:     true,
				Equilibrium: true,

				Health:       3904,
				MaxHealth:    3905,
				Mana:         3845,
				MaxMana:      3846,
				Endurance:    15020,
				MaxEndurance: 15021,
				Willpower:    12980,
				MaxWillpower: 12981,

				Bleed: 1,
				Rage:  2,

				Ferocity: 3,
				Kai:      4,
				Karma:    5,
				Spec:     "Sword and Shield",
				Stance:   "Crane",
			},
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			client := &mock.ClientMock{
				WriteFunc: func(data []byte) (int, error) {
					return 0, nil
				},
			}

			ui := &mock.UIMock{
				SetCharacterFunc: func(_ pkg.Character) {},
				SetTargetFunc:    func(_ *pkg.Target) {},
			}

			world, ok := achaea.NewWorld(client, ui).(*achaea.World)
			require.True(t, ok)

			if len(tc.command) > 0 {
				world.OnCommand(tc.command)
			}

			for _, message := range tc.messages {
				world.OnCommand(wrapGMCP(message.Marshal(), nil))
			}

			if tc.character != nil {
				assert.Equal(t, tc.character, world.Character)
			}

			if tc.target != nil {
				assert.Equal(t, *tc.target, world.Target.PkgTarget())
			}
		})
	}
}

func TestCommandsMutateVitals(t *testing.T) {
	tcs := []struct {
		command   []byte
		character pkg.Character
	}{
		{
			command: wrapGMCP("Char.Vitals", map[string]any{
				"hp":     "3904",
				"maxhp":  "3905",
				"mp":     "3845",
				"maxmp":  "3846",
				"ep":     "15020",
				"maxep":  "15021",
				"wp":     "12980",
				"maxwp":  "12981",
				"nl":     "19",
				"bal":    "1",
				"eq":     "1",
				"vote":   "1",
				"string": "H:3904/3905 M:3845/3846 E:15020/15021 W:12980/12981 NL:19/100 ",
				"charstats": []string{
					"Bleed: 1", "Rage: 2", "Kai: 4%",
					"Karma: 5%", "Stance: Crane",
					"Ferocity: 3", "Spec: Sword and Shield",
				},
			}),
			character: pkg.Character{
				Vitals: map[string]pkg.CharacterVital{
					"health":    {Value: 3904, Max: 3905},
					"mana":      {Value: 3845, Max: 3846},
					"endurance": {Value: 15020, Max: 15021},
					"willpower": {Value: 12980, Max: 12981},
					"ferocity":  {Value: 3, Max: 100},
					"kai":       {Value: 4, Max: 100},
					"karma":     {Value: 5, Max: 100},
				},
			},
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			var character pkg.Character
			ui := &mock.UIMock{
				SetCharacterFunc: func(char pkg.Character) {
					character = char
				},
			}

			world := achaea.NewWorld(&mock.ClientMock{}, ui)

			world.OnCommand(tc.command)

			assert.Equal(t, tc.character, character)
		})
	}
}
