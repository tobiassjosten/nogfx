package achaea

import (
	"errors"
	"fmt"
	"testing"

	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/gmcp"
	"github.com/tobiassjosten/nogfx/pkg/mock"
	"github.com/tobiassjosten/nogfx/pkg/telnet"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type FuncModule func([]byte) [][]byte

func (mod FuncModule) ProcessInput(data []byte) [][]byte {
	return mod(data)
}

func (mod FuncModule) ProcessOutput(data []byte) [][]byte {
	return mod(data)
}

func wrapGMCP(msgs ...string) []byte {
	var bs []byte
	for _, msg := range msgs {
		bs = append(bs, telnet.IAC, telnet.SB, telnet.GMCP)
		bs = append(bs, []byte(msg)...)
		bs = append(bs, telnet.IAC, telnet.SE)
	}
	return bs
}

func TestProcessInput(t *testing.T) {
	tcs := map[string]struct {
		in  []byte
		out [][]byte
		mod pkg.Module
	}{
		"plain asdf": {
			in:  []byte("asdf"),
			out: [][]byte{[]byte("asdf")},
		},

		"separated, repeated input": {
			in:  []byte("asdf;3 qwer;zxcv"),
			out: [][]byte{[]byte("asdf;qwer;qwer;qwer;zxcv")},
		},

		"blocked input": {
			in:  []byte("asdf"),
			out: nil,
			mod: FuncModule(func(_ []byte) [][]byte {
				return nil
			}),
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			client := &mock.ClientMock{}
			ui := &mock.UIMock{}

			world := NewWorld(client, ui).(*World)

			if tc.mod != nil {
				world.modules = append(world.modules, tc.mod)
			}

			assert.Equal(t, tc.out, world.ProcessInput(tc.in))
		})
	}
}

func TestProcessOutput(t *testing.T) {
	tcs := map[string]struct {
		in  []byte
		out [][]byte
		mod pkg.Module
	}{
		"plain asdf": {
			in:  []byte("asdf"),
			out: [][]byte{[]byte("asdf")},
		},

		"blocked output": {
			in:  []byte("asdf"),
			out: nil,
			mod: FuncModule(func(_ []byte) [][]byte {
				return nil
			}),
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			client := &mock.ClientMock{}
			ui := &mock.UIMock{}

			world := NewWorld(client, ui).(*World)

			if tc.mod != nil {
				world.modules = append(world.modules, tc.mod)
			}

			assert.Equal(t, tc.out, world.ProcessOutput(tc.in))
		})
	}
}

func TestCommandsReply(t *testing.T) {
	tcs := []struct {
		command []byte
		sent    []byte
		errs    []bool
		err     string
	}{
		{
			command: []byte{telnet.IAC, telnet.WILL, telnet.GMCP},
			sent: wrapGMCP(
				`Core.Supports.Set ["Char 1","Char.Items 1","Char.Skills 1","Comm.Channel 1","IRE.Rift 1","IRE.Target 1","Room 1"]`,
			),
		},
		{
			command: []byte{telnet.IAC, telnet.WILL, telnet.GMCP},
			errs:    []bool{true},
			err:     "failed GMCP: ooops",
		},

		{
			command: wrapGMCP(`Asdf.Qwer`),
			err:     "failed parsing GMCP: unknown message 'Asdf.Qwer'",
		},

		{
			command: wrapGMCP(
				`Char.Name {"name":"Durak","fullname":"Mason Durak"}`,
			),
			sent: wrapGMCP(
				`Char.Items.Inv`,
				`Comm.Channel.Players`,
				`IRE.Rift.Request`,
			),
		},
		{
			command: wrapGMCP(`Char.Name {}`),
			errs:    []bool{true},
			err:     "failed GMCP: ooops",
		},
		{
			command: wrapGMCP(`Char.Name {}`),
			errs:    []bool{false, true},
			err:     "failed GMCP: ooops",
		},
		{
			command: wrapGMCP(`Char.Name {}`),
			errs:    []bool{false, false, true},
			err:     "failed GMCP: ooops",
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

			world := NewWorld(client, ui)

			err := world.ProcessCommand(tc.command)

			if tc.err != "" && assert.NotNil(t, err) {
				assert.Equal(t, tc.err, err.Error())
				return
			}

			require.Nil(t, err)

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
		character *Character
		target    *pkg.Target
	}{
		{
			messages: []gmcp.Message{
				&gmcp.CharName{
					Name:     "Durak",
					Fullname: "Mason Durak",
				},
			},
			character: &Character{
				Name:  "Durak",
				Title: "Mason Durak",
			},
		},

		{
			// @todo Replace with []gmcp.Message.
			command: wrapGMCP(
				`Char.Status {"name":"Durak","fullname":"Mason Durak","class":"Monk","level":"68 (19%)"}`,
			),
			character: &Character{
				Name:  "Durak",
				Title: "Mason Durak",
				Class: "Monk",
				Level: 68,
			},
		},

		{
			// @todo Replace with []gmcp.Message.
			command: wrapGMCP(`Char.Vitals { "hp": "3904", "maxhp": "3905", "mp": "3845", "maxmp": "3846", "ep": "15020", "maxep": "15021", "wp": "12980", "maxwp": "12981", "bal": "1", "eq": "1", "charstats": [ "Bleed: 1", "Rage: 2", "Kai: 4%", "Karma: 5%", "Stance: Crane", "Ferocity: 3", "Spec: Sword and Shield" ] }`),
			character: &Character{
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

			aworld := NewWorld(client, ui).(*World)

			if len(tc.command) > 0 {
				err := aworld.ProcessCommand(tc.command)
				require.Nil(t, err)
			}

			for _, message := range tc.messages {
				err := aworld.ProcessCommand(wrapGMCP(message.Marshal()))
				require.Nil(t, err)
			}

			if tc.character != nil {
				assert.Equal(t, tc.character, aworld.Character)
			}

			if tc.target != nil {
				assert.Equal(t, *tc.target, aworld.Target.PkgTarget())
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
			command: wrapGMCP(`Char.Vitals { "hp": "3904", "maxhp": "3905", "mp": "3845", "maxmp": "3846", "ep": "15020", "maxep": "15021", "wp": "12980", "maxwp": "12981", "nl": "19", "bal": "1", "eq": "1", "vote": "1", "string": "H:3904/3905 M:3845/3846 E:15020/15021 W:12980/12981 NL:19/100 ", "charstats": [ "Bleed: 1", "Rage: 2", "Kai: 4%", "Karma: 5%", "Stance: Crane", "Ferocity: 3", "Spec: Sword and Shield" ] }`),
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

			world := NewWorld(&mock.ClientMock{}, ui)

			err := world.ProcessCommand(tc.command)
			require.Nil(t, err)

			assert.Equal(t, tc.character, character)
		})
	}
}
