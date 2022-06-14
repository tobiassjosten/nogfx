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

		"entering genji manticore present": {
			messages: []gmcp.Message{
				&gmcp.RoomInfo{Number: 1, AreaNumber: 137},
				&gmcp.CharItemsList{
					Location: "room",
					Items: []gmcp.CharItem{
						{Name: "a ferocious manticore"},
					},
				},
			},
			sent: []string{"settarget manticore"},
		},

		"entering genji manticore shaman present": {
			messages: []gmcp.Message{
				&gmcp.RoomInfo{Number: 1, AreaNumber: 137},
				&gmcp.CharItemsList{
					Location: "room",
					Items: []gmcp.CharItem{
						{Name: "a ferocious manticore"},
						{Name: "an atavian shaman"},
					},
				},
			},
			sent: []string{"settarget shaman"},
		},

		"entering genji manticore present pvping": {
			messages: []gmcp.Message{
				&igmcp.IRETargetSet{Target: "someone"},
				&gmcp.RoomInfo{Number: 1, AreaNumber: 137},
				&gmcp.CharItemsList{
					Location: "room",
					Items: []gmcp.CharItem{
						{Name: "a ferocious manticore"},
						{Name: "an atavian shaman"},
					},
				},
			},
		},

		"entering genji manticore enters": {
			messages: []gmcp.Message{
				&gmcp.RoomInfo{Number: 1, AreaNumber: 137},
				&gmcp.CharItemsAdd{
					Location: "room",
					Item: gmcp.CharItem{
						Name: "a ferocious manticore",
					},
				},
			},
			sent: []string{"settarget manticore"},
		},

		"entering genji manticore shaman enter": {
			messages: []gmcp.Message{
				&gmcp.RoomInfo{Number: 1, AreaNumber: 137},
				&gmcp.CharItemsAdd{
					Location: "room",
					Item: gmcp.CharItem{
						Name: "a ferocious manticore",
					},
				},
				&gmcp.CharItemsAdd{
					Location: "room",
					Item: gmcp.CharItem{
						Name: "an atavian shaman",
					},
				},
			},
			sent: []string{
				"settarget manticore",
				"settarget shaman",
			},
		},

		"entering genji manticore enters pvping": {
			messages: []gmcp.Message{
				&igmcp.IRETargetSet{Target: "someone"},
				&gmcp.RoomInfo{Number: 1, AreaNumber: 137},
				&gmcp.CharItemsAdd{
					Location: "room",
					Item: gmcp.CharItem{
						Name: "a ferocious manticore",
					},
				},
			},
		},

		"entering genji manticore shaman present shaman leaves": {
			messages: []gmcp.Message{
				&gmcp.RoomInfo{Number: 1, AreaNumber: 137},
				&gmcp.CharItemsList{
					Location: "room",
					Items: []gmcp.CharItem{
						{Name: "a ferocious manticore"},
						{Name: "an atavian shaman"},
					},
				},
				&gmcp.CharItemsRemove{
					Location: "room",
					Item: gmcp.CharItem{
						Name: "an atavian shaman",
					},
				},
			},
			sent: []string{
				"settarget shaman",
				"settarget manticore",
			},
		},

		"entering genji manticore present in unknown": {
			messages: []gmcp.Message{
				&gmcp.RoomInfo{Number: 1, AreaNumber: 137},
				&gmcp.CharItemsList{
					Location: "unknown",
					Items: []gmcp.CharItem{
						{Name: "a ferocious manticore"},
					},
				},
			},
		},

		"entering genji manticore enters/leaves in unknown": {
			messages: []gmcp.Message{
				&gmcp.RoomInfo{Number: 1, AreaNumber: 137},
				&gmcp.CharItemsAdd{
					Location: "unknown",
					Item: gmcp.CharItem{
						Name: "a ferocious manticore",
					},
				},
				&gmcp.CharItemsRemove{
					Location: "unknown",
					Item: gmcp.CharItem{
						Name: "a ferocious manticore",
					},
				},
			},
		},

		"entering genji manticore present leaving genji": {
			messages: []gmcp.Message{
				&gmcp.RoomInfo{Number: 1, AreaNumber: 137},
				&gmcp.CharItemsList{
					Location: "room",
					Items: []gmcp.CharItem{
						{Name: "a ferocious manticore"},
					},
				},
				&agmcp.CharStatus{Target: gox.NewString("manticore")},
				&gmcp.RoomInfo{Number: 2, AreaNumber: 731},
				&agmcp.CharStatus{Target: gox.NewString("")},
			},
			sent: []string{
				"settarget manticore",
				"settarget none",
			},
		},

		"entering invalid": {
			messages: []gmcp.Message{
				&gmcp.RoomInfo{},
			},
		},

		"entering genji manaual target shaman leaving genji": {
			messages: []gmcp.Message{
				&gmcp.RoomInfo{Number: 1, AreaNumber: 137},
				&agmcp.CharStatus{Target: gox.NewString("shaman")},
				&gmcp.RoomInfo{Number: 2, AreaNumber: 731},
				&agmcp.CharStatus{Target: gox.NewString("")},
			},
			sent: []string{
				"settarget none",
			},
		},

		"entering genji manaual target unknown leaving genji": {
			messages: []gmcp.Message{
				&gmcp.RoomInfo{Number: 1, AreaNumber: 137},
				&agmcp.CharStatus{Target: gox.NewString("unknown")},
				&gmcp.RoomInfo{Number: 2, AreaNumber: 731},
				&agmcp.CharStatus{Target: gox.NewString("")},
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
