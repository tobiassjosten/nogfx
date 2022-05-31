package gmcp_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tobiassjosten/nogfx/pkg/gmcp"
)

func TestRoomInfoDetails(t *testing.T) {
	msg := &gmcp.RoomInfo{}

	tcs := map[string]struct {
		details []string
		falsers []func() bool
		truther func() bool
	}{
		"base": {
			details: []string{},
			falsers: []func() bool{
				msg.IsBank,
				msg.IsIndoors,
				msg.IsOutdoors,
				msg.IsSewer,
				msg.IsShop,
				msg.IsSubdivision,
				msg.IsWilderness,
			},
			truther: func() bool { return true },
		},

		"bank": {
			details: []string{"bank"},
			truther: msg.IsBank,
		},

		"indoors": {
			details: []string{"bank"},
			truther: msg.IsBank,
		},

		"outdoors": {
			details: []string{"bank"},
			truther: msg.IsBank,
		},

		"sewer": {
			details: []string{"bank"},
			truther: msg.IsBank,
		},

		"shop": {
			details: []string{"bank"},
			truther: msg.IsBank,
		},

		"subdivision": {
			details: []string{"bank"},
			truther: msg.IsBank,
		},

		"wilderness": {
			details: []string{"bank"},
			truther: msg.IsBank,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			msg.Details = tc.details
			for _, falser := range tc.falsers {
				assert.False(t, falser())
			}
			assert.True(t, tc.truther())
		})
	}
}

func TestRoomMessages(t *testing.T) {
	tcs := map[string]struct {
		msg         gmcp.Message
		data        string
		unmarshaled gmcp.Message
		marshaled   string
		err         string
	}{
		"Room.Info empty": {
			msg:         &gmcp.RoomInfo{},
			data:        "Room.Info {}",
			unmarshaled: &gmcp.RoomInfo{},
			marshaled: makeGMCP("Room.Info", map[string]interface{}{
				"num":         0,
				"name":        "",
				"area":        "",
				"environment": "",
				"coords":      "",
				"map":         "",
				"exits":       map[string]int{},
				"details":     []string{},
			}),
		},

		"Room.Info hydrated": {
			msg: &gmcp.RoomInfo{},
			data: makeGMCP("Room.Info", map[string]interface{}{
				"num":         1234,
				"name":        "A room",
				"area":        "An area",
				"environment": "urban",
				"coords":      "1,2,3,4",
				"map":         "www.example.com/map?1,2,3,4",
				"exits": map[string]int{
					"n": 4321,
				},
				"details": []string{"shop", "bank"},
			}),
			unmarshaled: &gmcp.RoomInfo{
				Number:      1234,
				Name:        "A room",
				AreaName:    "An area",
				AreaNumber:  1,
				Environment: "urban",
				X:           2,
				Y:           3,
				Building:    4,
				Map:         "www.example.com/map?1,2,3,4",
				Exits: map[string]int{
					"n": 4321,
				},
				Details: []string{"shop", "bank"},
			},
			marshaled: makeGMCP("Room.Info", map[string]interface{}{
				"num":         1234,
				"name":        "A room",
				"area":        "An area",
				"environment": "urban",
				"coords":      "1,2,3,4",
				"map":         "www.example.com/map?1,2,3,4",
				"exits": map[string]int{
					"n": 4321,
				},
				"details": []string{"shop", "bank"},
			}),
		},

		"Room.Info invalid JSON": {
			msg:  &gmcp.RoomInfo{},
			data: "Room.Info asdf",
			err:  "invalid character 'a' looking for beginning of value",
		},

		"Room.Info no building": {
			msg: &gmcp.RoomInfo{},
			data: makeGMCP("Room.Info", map[string]interface{}{
				"coords": "1,2,3",
			}),
			unmarshaled: &gmcp.RoomInfo{
				AreaNumber: 1,
				X:          2,
				Y:          3,
				Building:   0,
			},
		},

		"Room.Info invalid building": {
			msg: &gmcp.RoomInfo{},
			data: makeGMCP("Room.Info", map[string]interface{}{
				"coords": "1,2,3,b",
			}),
			err: `failed parsing building from coords: strconv.Atoi: parsing "b": invalid syntax`,
		},

		"Room.Info invalid area": {
			msg: &gmcp.RoomInfo{},
			data: makeGMCP("Room.Info", map[string]interface{}{
				"coords": "a,2,3,4",
			}),
			err: `failed parsing area number from coords: strconv.Atoi: parsing "a": invalid syntax`,
		},

		"Room.Info invalid x": {
			msg: &gmcp.RoomInfo{},
			data: makeGMCP("Room.Info", map[string]interface{}{
				"coords": "1,x,3,4",
			}),
			err: `failed parsing x from coords: strconv.Atoi: parsing "x": invalid syntax`,
		},

		"Room.Info invalid y": {
			msg: &gmcp.RoomInfo{},
			data: makeGMCP("Room.Info", map[string]interface{}{
				"coords": "1,2,y,4",
			}),
			err: `failed parsing y from coords: strconv.Atoi: parsing "y": invalid syntax`,
		},

		"Room.Info invalid coords": {
			msg: &gmcp.RoomInfo{},
			data: makeGMCP("Room.Info", map[string]interface{}{
				"coords": "1,2",
			}),
			err: "failed parsing coords '[1 2]'",
		},

		"Room.Players empty": {
			msg:         &gmcp.RoomPlayers{},
			data:        "Room.Players []",
			unmarshaled: &gmcp.RoomPlayers{},
			marshaled:   "Room.Players []",
		},

		"Room.Players hydrated": {
			msg: &gmcp.RoomPlayers{},
			data: makeGMCP("Room.Players", []map[string]interface{}{
				{
					"name":     "Durak",
					"fullname": "Mason Durak",
				},
			}),
			unmarshaled: &gmcp.RoomPlayers{
				{
					Name:     "Durak",
					Fullname: "Mason Durak",
				},
			},
			marshaled: makeGMCP("Room.Players", []map[string]interface{}{
				{
					"name":     "Durak",
					"fullname": "Mason Durak",
				},
			}),
		},

		"Room.AddPlayer empty": {
			msg:         &gmcp.RoomAddPlayer{},
			data:        "Room.AddPlayer {}",
			unmarshaled: &gmcp.RoomAddPlayer{},
			marshaled: makeGMCP("Room.AddPlayer", map[string]interface{}{
				"name":     "",
				"fullname": "",
			}),
		},

		"Room.AddPlayer hydrated": {
			msg: &gmcp.RoomAddPlayer{},
			data: makeGMCP("Room.AddPlayer", map[string]interface{}{
				"name":     "Durak",
				"fullname": "Mason Durak",
			}),
			unmarshaled: &gmcp.RoomAddPlayer{
				Name:     "Durak",
				Fullname: "Mason Durak",
			},
			marshaled: makeGMCP("Room.AddPlayer", map[string]interface{}{
				"name":     "Durak",
				"fullname": "Mason Durak",
			}),
		},

		"Room.RemovePlayer empty": {
			msg:         &gmcp.RoomRemovePlayer{},
			data:        "Room.RemovePlayer {}",
			unmarshaled: &gmcp.RoomRemovePlayer{},
			marshaled: makeGMCP("Room.RemovePlayer", map[string]interface{}{
				"name":     "",
				"fullname": "",
			}),
		},

		"Room.RemovePlayer hydrated": {
			msg: &gmcp.RoomRemovePlayer{},
			data: makeGMCP("Room.RemovePlayer", map[string]interface{}{
				"name":     "Durak",
				"fullname": "Mason Durak",
			}),
			unmarshaled: &gmcp.RoomRemovePlayer{
				Name:     "Durak",
				Fullname: "Mason Durak",
			},
			marshaled: makeGMCP("Room.RemovePlayer", map[string]interface{}{
				"name":     "Durak",
				"fullname": "Mason Durak",
			}),
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			err := tc.msg.Unmarshal([]byte(tc.data))

			if tc.err != "" {
				require.NotNil(t, err)
				assert.Equal(t, tc.err, err.Error())
				return
			} else if err != nil {
				require.Equal(t, "", err.Error())
			}

			require.Equal(t, tc.unmarshaled, tc.msg, "unmarshaling hydrates message")

			if tc.marshaled == "" {
				return
			}

			marshaled := tc.msg.Marshal()
			data := strings.TrimSpace(strings.TrimPrefix(marshaled, tc.msg.ID()))
			tcdata := strings.TrimSpace(strings.TrimPrefix(tc.marshaled, tc.msg.ID()))

			assert.NotEqual(t, marshaled, data, "marshaled data has ID prefix")
			assert.NotEqual(t, tc.marshaled, tcdata, "marshaled data has ID prefix")

			if tcdata == "" {
				assert.Equal(t, tcdata, data)
				return
			}

			assert.JSONEq(t, tcdata, data, "marshaling maintains data integrity")

			require.Equal(t, tc.unmarshaled, tc.msg, "marshaling doesn't mutate")
		})
	}
}
