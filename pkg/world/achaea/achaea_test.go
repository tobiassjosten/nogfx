package achaea_test

import (
	"fmt"
	"testing"

	"github.com/tobiassjosten/nogfx/pkg"
	tn "github.com/tobiassjosten/nogfx/pkg/telnet"
	"github.com/tobiassjosten/nogfx/pkg/tui"
	"github.com/tobiassjosten/nogfx/pkg/world/achaea"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorldBasics(t *testing.T) {
	assert := assert.New(t)

	screen := &pkg.ScreenMock{
		SetStyleFunc: func(_ tcell.Style) {
		},
		SetCursorStyleFunc: func(_ tcell.CursorStyle) {
		},
	}

	ui := tui.NewTUI(screen, tui.NewPanes())

	client := &pkg.ClientMock{}

	world := achaea.NewWorld(ui, client)

	input := []byte("input")
	assert.Equal(input, world.Input(input))

	output := []byte("output")
	assert.Equal(output, world.Output(output))
}

func TestCommands(t *testing.T) {
	var sent []byte
	client := &pkg.ClientMock{
		WriteFunc: func(data []byte) (int, error) {
			sent = append(sent, data...)
			return len(data), nil
		},
	}

	screen := &pkg.ScreenMock{
		SetStyleFunc: func(_ tcell.Style) {
		},
		SetCursorStyleFunc: func(_ tcell.CursorStyle) {
		},
	}

	ui := tui.NewTUI(screen, tui.NewPanes())
	world := achaea.NewWorld(ui, client)

	wrapGMCP := func(msgs []string) []byte {
		var bs []byte
		for _, msg := range msgs {
			bs = append(bs, tn.IAC, tn.SB, tn.GMCP)
			bs = append(bs, []byte(msg)...)
			bs = append(bs, tn.IAC, tn.SE)
		}
		return bs
	}

	tcs := []struct {
		command []byte
		sent    []byte
	}{
		{
			command: []byte{tn.IAC, tn.WILL, tn.GMCP},
			sent: wrapGMCP([]string{
				`Core.Hello {"client":"nogfx","version":"0.0.0"}`,
				`Core.Supports.Set ["Char 1","Char.Skills 1","Char.Items 1","Comm.Channel 1","Room 1","IRE.Rift 1"]`,
			}),
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
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			sent = []byte{}

			err := world.Command(tc.command)
			require.Nil(err)

			if len(tc.sent) > 0 {
				assert.Equal(tc.sent, sent, string(sent))
			}
		})
	}
}
