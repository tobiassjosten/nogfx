package gmcp_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/icza/gox/gox"
	"github.com/stretchr/testify/assert"
	"github.com/tobiassjosten/nogfx/pkg/gmcp"
	"github.com/tobiassjosten/nogfx/pkg/telnet"
)

func makeGMCP(id string, data interface{}) string {
	jsondata, _ := json.Marshal(data)
	return fmt.Sprintf("%s %s", id, string(jsondata))
}

func TestParse(t *testing.T) {
	tcs := map[string]struct {
		datas []string
		msgs  []gmcp.Message
		errs  []string
	}{
		"Char.Login": {
			datas: []string{"Char.Login {}"},
			msgs:  []gmcp.Message{&gmcp.CharLogin{}},
		},

		"Char.Name": {
			datas: []string{"Char.Name {}"},
			msgs:  []gmcp.Message{&gmcp.CharName{}},
		},

		"Char.StatusVars": {
			datas: []string{"Char.StatusVars {}"},
			msgs:  []gmcp.Message{&gmcp.CharStatusVars{}},
		},

		"Char.Afflictions.List": {
			datas: []string{"Char.Afflictions.List []"},
			msgs:  []gmcp.Message{&gmcp.CharAfflictionsList{}},
		},

		"Char.Afflictions.Add": {
			datas: []string{"Char.Afflictions.Add {}"},
			msgs:  []gmcp.Message{&gmcp.CharAfflictionsAdd{}},
		},

		"Char.Afflictions.Remove": {
			datas: []string{"Char.Afflictions.Remove []"},
			msgs:  []gmcp.Message{&gmcp.CharAfflictionsRemove{}},
		},

		"Char.Defences.List": {
			datas: []string{"Char.Defences.List []"},
			msgs:  []gmcp.Message{&gmcp.CharDefencesList{}},
		},

		"Char.Defences.Add": {
			datas: []string{"Char.Defences.Add {}"},
			msgs:  []gmcp.Message{&gmcp.CharDefencesAdd{}},
		},

		"Char.Defences.Remove": {
			datas: []string{"Char.Defences.Remove []"},
			msgs:  []gmcp.Message{&gmcp.CharDefencesRemove{}},
		},

		"Char.Items.Contents": {
			datas: []string{"Char.Items.Contents 0"},
			msgs:  []gmcp.Message{&gmcp.CharItemsContents{}},
		},

		"Char.Items.Inv": {
			datas: []string{"Char.Items.Inv"},
			msgs:  []gmcp.Message{&gmcp.CharItemsInv{}},
		},

		"Char.Items.Room": {
			datas: []string{"Char.Items.Room"},
			msgs:  []gmcp.Message{&gmcp.CharItemsRoom{}},
		},

		"Char.Items.List": {
			datas: []string{"Char.Items.List {}"},
			msgs:  []gmcp.Message{&gmcp.CharItemsList{}},
		},

		"Char.Items.Add": {
			datas: []string{"Char.Items.Add {}"},
			msgs:  []gmcp.Message{&gmcp.CharItemsAdd{}},
		},

		"Char.Items.Remove": {
			datas: []string{"Char.Items.Remove {}"},
			msgs:  []gmcp.Message{&gmcp.CharItemsRemove{}},
		},

		"Char.Items.Update": {
			datas: []string{"Char.Items.Update {}"},
			msgs:  []gmcp.Message{&gmcp.CharItemsUpdate{}},
		},

		"Char.Skills.Get": {
			datas: []string{"Char.Skills.Get {}"},
			msgs:  []gmcp.Message{&gmcp.CharSkillsGet{}},
		},

		"Char.Skills.Groups": {
			datas: []string{"Char.Skills.Groups []"},
			msgs:  []gmcp.Message{&gmcp.CharSkillsGroups{}},
		},

		"Char.Skills.Info": {
			datas: []string{"Char.Skills.Info {}"},
			msgs:  []gmcp.Message{&gmcp.CharSkillsInfo{}},
		},

		"Char.Skills.List": {
			datas: []string{"Char.Skills.List {}"},
			msgs:  []gmcp.Message{&gmcp.CharSkillsList{}},
		},

		"Comm.Channel.Enable": {
			datas: []string{`Comm.Channel.Enable ""`},
			msgs:  []gmcp.Message{&gmcp.CommChannelEnable{}},
		},

		"Comm.Channel.List": {
			datas: []string{"Comm.Channel.List []"},
			msgs:  []gmcp.Message{&gmcp.CommChannelList{}},
		},

		"Comm.Channel.Players": {
			datas: []string{"Comm.Channel.Players []"},
			msgs:  []gmcp.Message{&gmcp.CommChannelPlayers{}},
		},

		"Comm.Channel.Text": {
			datas: []string{"Comm.Channel.Text {}"},
			msgs:  []gmcp.Message{&gmcp.CommChannelText{}},
		},

		"Core.Goodbye": {
			datas: []string{"Core.Goodbye"},
			msgs:  []gmcp.Message{&gmcp.CoreGoodbye{}},
		},

		"Core.Hello": {
			datas: []string{"Core.Hello {}"},
			msgs:  []gmcp.Message{&gmcp.CoreHello{}},
		},

		"Core.KeepAlive": {
			datas: []string{"Core.KeepAlive"},
			msgs:  []gmcp.Message{&gmcp.CoreKeepAlive{}},
		},

		"Core.Ping": {
			datas: []string{"Core.Ping"},
			msgs:  []gmcp.Message{&gmcp.CorePing{}},
		},

		"Core.Supports.Set": {
			datas: []string{"Core.Supports.Set []"},
			msgs:  []gmcp.Message{&gmcp.CoreSupportsSet{}},
		},

		"Core.Supports.Add": {
			datas: []string{"Core.Supports.Add []"},
			msgs:  []gmcp.Message{&gmcp.CoreSupportsAdd{}},
		},

		"Core.Supports.Remove": {
			datas: []string{"Core.Supports.Remove []"},
			msgs:  []gmcp.Message{&gmcp.CoreSupportsRemove{}},
		},

		"Room.Info": {
			datas: []string{"Room.Info {}"},
			msgs:  []gmcp.Message{&gmcp.RoomInfo{}},
		},

		"Room.Players": {
			datas: []string{"Room.Players []"},
			msgs:  []gmcp.Message{&gmcp.RoomPlayers{}},
		},

		"Room.AddPlayer": {
			datas: []string{"Room.AddPlayer {}"},
			msgs:  []gmcp.Message{&gmcp.RoomAddPlayer{}},
		},

		"Room.RemovePlayer": {
			datas: []string{"Room.RemovePlayer {}"},
			msgs:  []gmcp.Message{&gmcp.RoomRemovePlayer{}},
		},

		"non-existant": {
			datas: []string{"Non.Existant"},
			errs:  []string{"unknown message 'Non.Existant'"},
		},

		"invalid JSON": {
			datas: []string{"Char.Login asdf"},
			errs:  []string{"couldn't unmarshal *gmcp.CharLogin: invalid character 'a' looking for beginning of value"},
		},

		"Room.Info x2": {
			datas: []string{
				`Room.Info {"exits":{"n":2,"s":3}}`,
				`Room.Info {"exits":{"w":4,"e":5}}`,
			},
			msgs: []gmcp.Message{
				&gmcp.RoomInfo{Exits: map[string]int{"n": 2, "s": 3}},
				&gmcp.RoomInfo{Exits: map[string]int{"w": 4, "e": 5}},
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			var msgs []gmcp.Message
			var errs []string

			for _, data := range tc.datas {
				msg, err := gmcp.Parse([]byte(data))
				msgs = append(msgs, msg)
				if err != nil {
					errs = append(errs, err.Error())
				}
			}

			if tc.errs != nil {
				if assert.NotNil(t, errs) {
					assert.Equal(t, tc.errs, errs)
				}
				return
			} else if errs != nil {
				assert.Nil(t, errs)
				return
			}

			assert.Equal(t, tc.msgs, msgs)
		})
	}
}

func TestWrap(t *testing.T) {
	var (
		gmcpPrefix = []byte{telnet.IAC, telnet.SB, telnet.GMCP}
		gmcpSuffix = []byte{telnet.IAC, telnet.SE}
	)

	data := []byte("asdf")
	assert.Equal(t, data, gmcp.Unwrap(gmcp.Wrap(data)))

	assert.Nil(t, gmcp.Unwrap(append(data, gmcpSuffix...)))
	assert.Nil(t, gmcp.Unwrap(append(gmcpPrefix, data...)))
	assert.NotNil(t, gmcp.Unwrap(append(append(gmcpPrefix, data...), gmcpSuffix...)))
}

func TestSplitRank(t *testing.T) {
	tcs := map[string]struct {
		input string
		item  string
		rank  *string
		ranki *int
		rankf *float64
	}{
		"string Something (rank)": {
			input: "Something (rank)",
			item:  "Something",
			rank:  gox.NewString("rank"),
		},

		"string Something(rank)": {
			input: "Something(rank)",
			item:  "Something",
			rank:  gox.NewString("rank"),
		},

		"string Something only": {
			input: "Something",
			item:  "Something",
			rank:  gox.NewString(""),
		},

		"int Something (1)": {
			input: "Something (1)",
			item:  "Something",
			ranki: gox.NewInt(1),
		},

		"int Something(1)": {
			input: "Something(1)",
			item:  "Something",
			ranki: gox.NewInt(1),
		},

		"int Something (1%)": {
			input: "Something (1%)",
			item:  "Something",
			ranki: gox.NewInt(1),
		},

		"int Something (x)": {
			input: "Something (x)",
			item:  "Something",
			ranki: gox.NewInt(0),
		},

		"int Something only": {
			input: "Something",
			item:  "Something",
			ranki: gox.NewInt(0),
		},

		"float Something (1.2)": {
			input: "Something (1.2)",
			item:  "Something",
			rankf: gox.NewFloat64(1.2),
		},

		"float Something(1.2)": {
			input: "Something(1.2)",
			item:  "Something",
			rankf: gox.NewFloat64(1.2),
		},

		"float Something (1.2%)": {
			input: "Something (1.2%)",
			item:  "Something",
			rankf: gox.NewFloat64(1.2),
		},

		"float Something (x)": {
			input: "Something (x)",
			item:  "Something",
			rankf: gox.NewFloat64(0),
		},

		"float Something only": {
			input: "Something",
			item:  "Something",
			rankf: gox.NewFloat64(0),
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			if tc.rank != nil {
				item, rank := gmcp.SplitRank(tc.input)
				assert.Equal(t, tc.item, item)
				assert.Equal(t, *tc.rank, rank)
			}

			if tc.ranki != nil {
				item, ranki := gmcp.SplitRankInt(tc.input)
				assert.Equal(t, tc.item, item)
				assert.Equal(t, *tc.ranki, ranki)
			}

			if tc.rankf != nil {
				item, rankf := gmcp.SplitRankFloat(tc.input)
				assert.Equal(t, tc.item, item)
				assert.Equal(t, *tc.rankf, rankf)
			}
		})
	}
}
