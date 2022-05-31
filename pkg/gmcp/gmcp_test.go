package gmcp_test

import (
	"encoding/json"
	"fmt"
	"testing"

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
		data string
		msg  gmcp.Message
		err  string
	}{
		"Char.Login": {
			data: "Char.Login {}",
			msg:  &gmcp.CharLogin{},
		},

		"Char.Name": {
			data: "Char.Name {}",
			msg:  &gmcp.CharName{},
		},

		"Char.StatusVars": {
			data: "Char.StatusVars {}",
			msg:  &gmcp.CharStatusVars{},
		},

		"Char.Afflictions.List": {
			data: "Char.Afflictions.List []",
			msg:  &gmcp.CharAfflictionsList{},
		},

		"Char.Afflictions.Add": {
			data: "Char.Afflictions.Add {}",
			msg:  &gmcp.CharAfflictionsAdd{},
		},

		"Char.Afflictions.Remove": {
			data: "Char.Afflictions.Remove []",
			msg:  &gmcp.CharAfflictionsRemove{},
		},

		"Char.Defences.List": {
			data: "Char.Defences.List []",
			msg:  &gmcp.CharDefencesList{},
		},

		"Char.Defences.Add": {
			data: "Char.Defences.Add {}",
			msg:  &gmcp.CharDefencesAdd{},
		},

		"Char.Defences.Remove": {
			data: "Char.Defences.Remove []",
			msg:  &gmcp.CharDefencesRemove{},
		},

		"Char.Items.Contents": {
			data: "Char.Items.Contents 0",
			msg:  &gmcp.CharItemsContents{},
		},

		"Char.Items.Inv": {
			data: "Char.Items.Inv",
			msg:  &gmcp.CharItemsInv{},
		},

		"Char.Items.Room": {
			data: "Char.Items.Room",
			msg:  &gmcp.CharItemsRoom{},
		},

		"Char.Items.List": {
			data: "Char.Items.List {}",
			msg:  &gmcp.CharItemsList{},
		},

		"Char.Items.Add": {
			data: "Char.Items.Add {}",
			msg:  &gmcp.CharItemsAdd{},
		},

		"Char.Items.Remove": {
			data: "Char.Items.Remove {}",
			msg:  &gmcp.CharItemsRemove{},
		},

		"Char.Items.Update": {
			data: "Char.Items.Update {}",
			msg:  &gmcp.CharItemsUpdate{},
		},

		"Char.Skills.Get": {
			data: "Char.Skills.Get {}",
			msg:  &gmcp.CharSkillsGet{},
		},

		"Char.Skills.Groups": {
			data: "Char.Skills.Groups []",
			msg:  &gmcp.CharSkillsGroups{},
		},

		"Char.Skills.Info": {
			data: "Char.Skills.Info {}",
			msg:  &gmcp.CharSkillsInfo{},
		},

		"Char.Skills.List": {
			data: "Char.Skills.List {}",
			msg:  &gmcp.CharSkillsList{},
		},

		"Comm.Channel.Enable": {
			data: `Comm.Channel.Enable ""`,
			msg:  &gmcp.CommChannelEnable{},
		},

		"Comm.Channel.List": {
			data: "Comm.Channel.List []",
			msg:  &gmcp.CommChannelList{},
		},

		"Comm.Channel.Players": {
			data: "Comm.Channel.Players []",
			msg:  &gmcp.CommChannelPlayers{},
		},

		"Comm.Channel.Start": {
			data: `Comm.Channel.Start ""`,
			msg:  &gmcp.CommChannelStart{},
		},

		"Comm.Channel.End": {
			data: `Comm.Channel.End ""`,
			msg:  &gmcp.CommChannelEnd{},
		},

		"Comm.Channel.Text": {
			data: "Comm.Channel.Text {}",
			msg:  &gmcp.CommChannelText{},
		},

		"Core.Goodbye": {
			data: "Core.Goodbye",
			msg:  &gmcp.CoreGoodbye{},
		},

		"Core.Hello": {
			data: "Core.Hello {}",
			msg:  &gmcp.CoreHello{},
		},

		"Core.KeepAlive": {
			data: "Core.KeepAlive",
			msg:  &gmcp.CoreKeepAlive{},
		},

		"Core.Ping": {
			data: "Core.Ping",
			msg:  &gmcp.CorePing{},
		},

		"Core.Supports.Set": {
			data: "Core.Supports.Set []",
			msg:  &gmcp.CoreSupportsSet{},
		},

		"Core.Supports.Add": {
			data: "Core.Supports.Add []",
			msg:  &gmcp.CoreSupportsAdd{},
		},

		"Core.Supports.Remove": {
			data: "Core.Supports.Remove []",
			msg:  &gmcp.CoreSupportsRemove{},
		},

		"Room.Info": {
			data: "Room.Info {}",
			msg:  &gmcp.RoomInfo{},
		},

		"Room.Players": {
			data: "Room.Players []",
			msg:  &gmcp.RoomPlayers{},
		},

		"Room.AddPlayer": {
			data: "Room.AddPlayer {}",
			msg:  &gmcp.RoomAddPlayer{},
		},

		"Room.RemovePlayer": {
			data: "Room.RemovePlayer {}",
			msg:  &gmcp.RoomRemovePlayer{},
		},

		"Room.WrongDir": {
			data: `Room.WrongDir ""`,
			msg:  &gmcp.RoomWrongDir{},
		},

		"non-existant": {
			data: "Non.Existant",
			err:  "unknown message 'Non.Existant'",
		},

		"invalid JSON": {
			data: "Char.Login asdf",
			err:  "couldn't unmarshal *gmcp.CharLogin: invalid character 'a' looking for beginning of value",
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			msg, err := gmcp.Parse([]byte(tc.data))

			if tc.err != "" {
				if assert.NotNil(t, err) {
					assert.Equal(t, tc.err, err.Error())
				}
				return
			} else if err != nil {
				assert.Equal(t, "", err.Error())
				return
			}

			assert.Equal(t, tc.msg, msg)
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
