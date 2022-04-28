package achaea

import (
	"errors"
	"fmt"
	"testing"

	"github.com/tobiassjosten/nogfx/pkg/mock"
	"github.com/tobiassjosten/nogfx/pkg/telnet"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func wrapGMCP(msgs []string) []byte {
	var bs []byte
	for _, msg := range msgs {
		bs = append(bs, telnet.IAC, telnet.SB, telnet.GMCP)
		bs = append(bs, []byte(msg)...)
		bs = append(bs, telnet.IAC, telnet.SE)
	}
	return bs
}

func TestWorldBasics(t *testing.T) {
	assert := assert.New(t)

	client := &mock.ClientMock{}
	ui := &mock.UIMock{
		AddVitalFunc: func(_ string, _ interface{}) {},
	}

	world := NewWorld(client, ui)

	input := []byte("input")
	assert.Equal(input, world.ProcessInput(input))

	output := []byte("output")
	assert.Equal(output, world.ProcessOutput(output))
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
			sent: wrapGMCP([]string{
				`Core.Supports.Set ["Char 1","Char.Skills 1","Char.Items 1","Comm.Channel 1","Room 1","IRE.Rift 1"]`,
			}),
		},
		{
			command: []byte{telnet.IAC, telnet.WILL, telnet.GMCP},
			errs:    []bool{true},
			err:     "failed GMCP: ooops",
		},

		{
			command: wrapGMCP([]string{`Asdf.Qwer`}),
			err:     "failed parsing GMCP: unknown message 'Asdf.Qwer'",
		},

		{
			command: wrapGMCP([]string{
				`Char.Name {"name":"Durak","fullname":"Mason Durak"}`,
			}),
			sent: wrapGMCP([]string{
				`Char.Items.Inv`,
				`Comm.Channel.Players`,
				`IRE.Rift.Request`,
			}),
		},
		{
			command: wrapGMCP([]string{`Char.Name {}`}),
			errs:    []bool{true},
			err:     "failed GMCP: ooops",
		},
		{
			command: wrapGMCP([]string{`Char.Name {}`}),
			errs:    []bool{false, true},
			err:     "failed GMCP: ooops",
		},
		{
			command: wrapGMCP([]string{`Char.Name {}`}),
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

			ui := &mock.UIMock{
				AddVitalFunc: func(_ string, _ interface{}) {},
			}

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
		character *Character
	}{
		{
			command: wrapGMCP([]string{
				`Char.Name {"name":"Durak","fullname":"Mason Durak"}`,
			}),
			character: &Character{
				Name:  "Durak",
				Title: "Mason Durak",
			},
		},
		{
			command: wrapGMCP([]string{
				`Char.Status {"name":"Durak","fullname":"Mason Durak","class":"Monk","level":"68 (19%)"}`,
			}),
			character: &Character{
				Name:  "Durak",
				Title: "Mason Durak",
				Class: "Monk",
				Level: 68,
			},
		},
		{
			command: wrapGMCP([]string{`Char.Vitals { "hp": "3904", "maxhp": "3905", "mp": "3845", "maxmp": "3846", "ep": "15020", "maxep": "15021", "wp": "12980", "maxwp": "12981", "bal": "1", "eq": "1", "charstats": [ "Bleed: 1", "Rage: 2", "Kai: 4%", "Karma: 5%", "Stance: Crane", "Ferocity: 3", "Spec: Sword and Shield" ] }`}),
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
				AddVitalFunc:    func(_ string, _ interface{}) {},
				UpdateVitalFunc: func(_ string, _, _ int) {},
			}

			aworld := NewWorld(client, ui).(*World)

			err := aworld.ProcessCommand(tc.command)
			require.Nil(t, err)

			if tc.character != nil {
				assert.Equal(t, tc.character, aworld.character)
			}
		})
	}
}

func TestCommandsMutateVitals(t *testing.T) {
	tcs := []struct {
		command []byte
		vitals  map[string][][]int
	}{
		// @todo Add Char.Name once that pane exists.
		{
			command: wrapGMCP([]string{`Char.Vitals { "hp": "3904", "maxhp": "3905", "mp": "3845", "maxmp": "3846", "ep": "15020", "maxep": "15021", "wp": "12980", "maxwp": "12981", "nl": "19", "bal": "1", "eq": "1", "vote": "1", "string": "H:3904/3905 M:3845/3846 E:15020/15021 W:12980/12981 NL:19/100 ", "charstats": [ "Bleed: 1", "Rage: 2", "Kai: 4%", "Karma: 5%", "Stance: Crane", "Ferocity: 3", "Spec: Sword and Shield" ] }`}),
			vitals: map[string][][]int{
				"health":    [][]int{{3904, 3905}},
				"mana":      [][]int{{3845, 3846}},
				"endurance": [][]int{{15020, 15021}},
				"willpower": [][]int{{12980, 12981}},
			},
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			vitals := map[string][][]int{}
			ui := &mock.UIMock{
				AddVitalFunc: func(_ string, _ interface{}) {},
				UpdateVitalFunc: func(name string, value, max int) {
					if _, ok := vitals[name]; !ok {
						vitals[name] = [][]int{}
					}
					vitals[name] = append(vitals[name], []int{value, max})
				},
			}

			world := NewWorld(&mock.ClientMock{}, ui)

			err := world.ProcessCommand(tc.command)
			require.Nil(t, err)

			assert.Equal(t, vitals, tc.vitals)
		})
	}
}
