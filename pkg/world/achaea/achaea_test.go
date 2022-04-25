package achaea_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/telnet"
	tn "github.com/tobiassjosten/nogfx/pkg/telnet"
	"github.com/tobiassjosten/nogfx/pkg/world/achaea"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func wrapGMCP(msgs []string) []byte {
	var bs []byte
	for _, msg := range msgs {
		bs = append(bs, tn.IAC, tn.SB, tn.GMCP)
		bs = append(bs, []byte(msg)...)
		bs = append(bs, tn.IAC, tn.SE)
	}
	return bs
}

func TestWorldBasics(t *testing.T) {
	assert := assert.New(t)

	client := &pkg.ClientMock{}
	ui := &pkg.UIMock{
		AddVitalFunc: func(_ string, _ interface{}) {},
	}

	world := achaea.NewWorld(ui, client)

	input := []byte("input")
	assert.Equal(input, world.Input(input))

	output := []byte("output")
	assert.Equal(output, world.Output(output))
}

func TestCommandsReply(t *testing.T) {
	tcs := []struct {
		command []byte
		sent    []byte
		errs    []bool
		err     string
	}{
		{
			command: []byte{tn.IAC, tn.WILL, tn.GMCP},
			sent: wrapGMCP([]string{
				`Core.Hello {"client":"nogfx","version":"0.0.0"}`,
				`Core.Supports.Set ["Char 1","Char.Skills 1","Char.Items 1","Comm.Channel 1","Room 1","IRE.Rift 1"]`,
			}),
		},
		{
			command: []byte{telnet.IAC, telnet.WILL, telnet.GMCP},
			errs:    []bool{true},
			err:     "failed GMCP: ooops",
		},
		{
			command: []byte{telnet.IAC, telnet.WILL, telnet.GMCP},
			errs:    []bool{false, true},
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
				`IRE.Rift.Request`,
				`Comm.Channel.Players`,
				`Char.Items.Inv`,
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

			client := &pkg.ClientMock{
				WriteFunc: func(data []byte) (int, error) {
					defer func() { calls++ }()
					if len(tc.errs) > calls && tc.errs[calls] {
						return 0, errors.New("ooops")
					}

					sent = append(sent, data...)

					return len(data), nil
				},
			}

			ui := &pkg.UIMock{
				AddVitalFunc: func(_ string, _ interface{}) {},
			}

			world := achaea.NewWorld(ui, client)

			err := world.Command(tc.command)

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

func TestCommandsMutate(t *testing.T) {
	tcs := []struct {
		command []byte
		vitals  map[string][][]int
	}{
		// @todo Add Char.Name once that pane exists.
		{
			command: wrapGMCP([]string{`Char.Vitals { "hp": "3904", "maxhp": "3904", "mp": "3845", "maxmp": "3845", "ep": "15020", "maxep": "15020", "wp": "12980", "maxwp": "12980", "nl": "19", "bal": "1", "eq": "1", "vote": "1", "string": "H:3904/3904 M:3845/3845 E:15020/15020 W:12980/12980 NL:19/100 ", "charstats": [ "Bleed: 1", "Rage: 2", "Kai: 4%", "Karma: 5%", "Stance: Crane", "Ferocity: 3", "Spec: Sword and Shield" ] }`}),
			vitals: map[string][][]int{
				"health":    [][]int{{3904, 3904}},
				"mana":      [][]int{{3845, 3845}},
				"endurance": [][]int{{15020, 15020}},
				"willpower": [][]int{{12980, 12980}},
			},
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			vitals := map[string][][]int{}
			ui := &pkg.UIMock{
				AddVitalFunc: func(_ string, _ interface{}) {},
				UpdateVitalFunc: func(name string, value, max int) {
					if _, ok := vitals[name]; !ok {
						vitals[name] = [][]int{}
					}
					vitals[name] = append(vitals[name], []int{value, max})
				},
			}

			world := achaea.NewWorld(ui, &pkg.ClientMock{})

			err := world.Command(tc.command)
			require.Nil(t, err)

			assert.Equal(t, vitals, tc.vitals)
		})
	}
}

func TestMasking(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	client := &pkg.ClientMock{}

	var masked bool

	ui := &pkg.UIMock{
		AddVitalFunc: func(_ string, _ interface{}) {},
		MaskInputFunc: func() {
			masked = true
		},
		UnmaskInputFunc: func() {
			masked = false
		},
	}

	world := achaea.NewWorld(ui, client)

	err := world.Command([]byte{tn.IAC, tn.WILL, tn.ECHO})
	require.Nil(err)

	assert.Equal(true, masked)

	err = world.Command([]byte{tn.IAC, tn.WONT, tn.ECHO})
	require.Nil(err)

	assert.Equal(false, masked)
}
