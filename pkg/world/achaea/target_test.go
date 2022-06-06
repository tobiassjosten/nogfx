package achaea_test

import (
	"testing"

	"github.com/icza/gox/gox"
	"github.com/tobiassjosten/nogfx/pkg/gmcp"
	agmcp "github.com/tobiassjosten/nogfx/pkg/gmcp/achaea"
	igmcp "github.com/tobiassjosten/nogfx/pkg/gmcp/ironrealms"
	"github.com/tobiassjosten/nogfx/pkg/mock"
	"github.com/tobiassjosten/nogfx/pkg/navigation"
	"github.com/tobiassjosten/nogfx/pkg/world/achaea"

	"github.com/stretchr/testify/assert"
)

/*
	case *gmcp.CharItemsList:
		world.Target.FromCharItemsList(msg)

	case *gmcp.CharItemsAdd:
		world.Target.FromCharItemsAdd(msg)

	case *gmcp.CharItemsRemove:
		world.Target.FromCharItemsRemove(msg)

	case *gmcp.RoomInfo:
		world.Target.FromRoomInfo(msg)
*/

func TestWorldTargeting(t *testing.T) {
	tcs := map[string]struct {
		messages []gmcp.Message
		name     string
		health   int
		sent     []string
	}{
		"initial state": {
			messages: []gmcp.Message{},
			name:     "",
			health:   0,
		},

		"target verified": {
			messages: []gmcp.Message{
				&agmcp.CharStatus{Target: gox.NewString("AsDf")},
			},
			name:   "asdf",
			health: 0,
		},

		"target health": {
			messages: []gmcp.Message{
				&igmcp.IRETargetInfo{Health: 1234},
			},
			name:   "",
			health: 1234,
		},

		"entering unknown area": {
			messages: []gmcp.Message{
				&gmcp.RoomInfo{},
			},
		},

		"entering unsupported area": {
			messages: []gmcp.Message{
				&gmcp.RoomInfo{AreaNumber: 1},
			},
		},

		"entering genji targetless": {
			messages: []gmcp.Message{
				&gmcp.RoomInfo{AreaNumber: 137},
			},
			sent: []string{"settarget shaman"},
		},

		"entering genji pvping": {
			messages: []gmcp.Message{
				&igmcp.IRETargetSet{Target: "durak"},
				&gmcp.RoomInfo{AreaNumber: 137},
			},
		},

		"traversing genji": {
			messages: []gmcp.Message{
				&gmcp.RoomInfo{Number: 1, AreaNumber: 137},
				&gmcp.RoomInfo{Number: 2, AreaNumber: 137},
			},
			sent: []string{"settarget shaman"},
		},

		"exiting genji": {
			messages: []gmcp.Message{
				&gmcp.RoomInfo{Number: 1, AreaNumber: 137},
				&agmcp.CharStatus{Target: gox.NewString("shaman")},
				&gmcp.RoomInfo{Number: 2, AreaNumber: 1},
				&agmcp.CharStatus{Target: gox.NewString("")},
			},
			sent: []string{
				"settarget shaman",
				"settarget none",
			},
		},

		"secondary target found primary": {
			messages: []gmcp.Message{
				&gmcp.RoomInfo{Number: 1, AreaNumber: 137},
				&gmcp.CharItemsList{
					Location: "room",
					Items: []gmcp.CharItem{
						{Name: "a ferocious manticore"},
					},
				},
			},
			sent: []string{
				"settarget shaman",
				"settarget manticore",
			},
		},

		"secondary target found secondarily": {
			messages: []gmcp.Message{
				&gmcp.RoomInfo{Number: 1, AreaNumber: 137},
				&agmcp.CharStatus{Target: gox.NewString("shaman")},
				&gmcp.CharItemsList{
					Location: "room",
					Items: []gmcp.CharItem{
						{Name: "an atavian shaman"},
						{Name: "manticore"},
					},
				},
			},
			name: "shaman",
			sent: []string{"settarget shaman"},
		},

		"secondary target entered": {
			messages: []gmcp.Message{
				&gmcp.RoomInfo{Number: 1, AreaNumber: 137},
				&agmcp.CharStatus{Target: gox.NewString("shaman")},
				&gmcp.CharItemsAdd{
					Location: "room",
					Item:     gmcp.CharItem{Name: "manticore"},
				},
			},
			name: "shaman",
			sent: []string{
				"settarget shaman",
				"settarget manticore",
			},
		},

		"secondary target become primary": {
			messages: []gmcp.Message{
				&gmcp.RoomInfo{Number: 1, AreaNumber: 137},
				&agmcp.CharStatus{Target: gox.NewString("shaman")},
				&gmcp.CharItemsList{
					Location: "room",
					Items: []gmcp.CharItem{
						{Name: "shaman"},
						{Name: "manticore"},
					},
				},
				&gmcp.CharItemsRemove{
					Location: "room",
					Item:     gmcp.CharItem{Name: "shaman"},
				},
				&agmcp.CharStatus{Target: gox.NewString("manticore")},
			},
			name: "manticore",
			sent: []string{
				"settarget shaman",
				"settarget manticore",
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			navigation.Reset()

			var sent []string
			client := &mock.ClientMock{
				SendFunc: func(data []byte) {
					sent = append(sent, string(data))
				},
			}

			target := achaea.NewTarget(client)

			for _, message := range tc.messages {
				switch msg := message.(type) {
				case *gmcp.CharItemsList:
					target.FromCharItemsList(msg)

				case *gmcp.CharItemsAdd:
					target.FromCharItemsAdd(msg)

				case *gmcp.CharItemsRemove:
					target.FromCharItemsRemove(msg)

				case *agmcp.CharStatus:
					target.FromCharStatus(msg)

				case *gmcp.RoomInfo:
					target.FromRoomInfo(msg)

				case *igmcp.IRETargetSet:
					target.FromIRETargetSet(msg)

				case *igmcp.IRETargetInfo:
					target.FromIRETargetInfo(msg)

				default:
					t.Fatalf("unsupported message %T", msg)
				}
			}

			assert.Equal(t, tc.name, target.Name)
			assert.Equal(t, tc.health, target.Health)
			assert.Equal(t, tc.sent, sent)
		})
	}
}
