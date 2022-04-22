package gmcp_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tobiassjosten/nogfx/pkg/world/achaea/gmcp"
)

func TestRoomServerMessages(t *testing.T) {
	tcs := []struct {
		command []byte
		message gmcp.ServerMessage
		err     string
	}{
		{
			command: []byte("Room.Info"),
			err:     "failed hydrating gmcp.RoomInfo (): unexpected end of JSON input",
		},
		{
			command: []byte(`Room.Info {"details": [ "asdf" ] }`),
			err:     `failed hydrating gmcp.RoomInfo ({"details": [ "asdf" ] }): unknown room detail 'asdf'`,
		},
		{
			command: []byte(`Room.Info {"coords": "123"}`),
			err:     `failed hydrating gmcp.RoomInfo ({"coords": "123"}): failed parsing coords '[123]'`,
		},
		{
			command: []byte(`Room.Info {"coords": "a,5,4,3"}`),
			err:     `failed hydrating gmcp.RoomInfo ({"coords": "a,5,4,3"}): failed parsing area number from coords '[a 5 4 3]': strconv.Atoi: parsing "a": invalid syntax`,
		},
		{
			command: []byte(`Room.Info {"coords": "45,b,4,3"}`),
			err:     `failed hydrating gmcp.RoomInfo ({"coords": "45,b,4,3"}): failed parsing x from coords '[45 b 4 3]': strconv.Atoi: parsing "b": invalid syntax`,
		},
		{
			command: []byte(`Room.Info {"coords": "45,5,c,3"}`),
			err:     `failed hydrating gmcp.RoomInfo ({"coords": "45,5,c,3"}): failed parsing y from coords '[45 5 c 3]': strconv.Atoi: parsing "c": invalid syntax`,
		},
		{
			command: []byte(`Room.Info {"coords": "45,5,4,d"}`),
			err:     `failed hydrating gmcp.RoomInfo ({"coords": "45,5,4,d"}): failed parsing building from coords '[45 5 4 d]': strconv.Atoi: parsing "d": invalid syntax`,
		},
		{
			command: []byte(`Room.Info {"coords": "45,5,4"}`),
			message: gmcp.RoomInfo{
				AreaNumber: 45,
				X:          5,
				Y:          4,
			},
		},
		{
			command: []byte(`Room.Info {"num": 12345, "name": "On a hill", "area": "Barren hills", "environment": "Hills", "coords": "45,5,4,3", "map": "www.example.com/map.php?45,5,4,3", "exits": { "n": 12344, "se": 12336 }, "details": [ "shop", "bank" ] }`),
			message: gmcp.RoomInfo{
				Number:      12345,
				Name:        "On a hill",
				AreaName:    "Barren hills",
				AreaNumber:  45,
				Environment: "Hills",
				X:           5,
				Y:           4,
				Building:    3,
				Map:         "www.example.com/map.php?45,5,4,3",
				Exits: map[string]int{
					"n":  12344,
					"se": 12336,
				},
				Details: gmcp.RoomDetails{
					Shop: true,
					Bank: true,
				},
			},
		},

		{
			command: []byte("Room.Players"),
			err:     "failed hydrating gmcp.RoomPlayers (): unexpected end of JSON input",
		},
		{
			command: []byte(`Room.Players [{ "name": "Tecton", "fullname": "Tecton the Terraformer" }]`),
			message: gmcp.RoomPlayers{{
				Name:     "Tecton",
				Fullname: "Tecton the Terraformer",
			}},
		},

		{
			command: []byte("Room.AddPlayer"),
			err:     "failed hydrating gmcp.RoomAddPlayer (): unexpected end of JSON input",
		},
		{
			command: []byte(`Room.AddPlayer { "name": "Tecton", "fullname": "Tecton the Terraformer" }`),
			message: gmcp.RoomAddPlayer{
				Name:     "Tecton",
				Fullname: "Tecton the Terraformer",
			},
		},

		{
			command: []byte("Room.RemovePlayer"),
			err:     "failed hydrating gmcp.RoomRemovePlayer (): unexpected end of JSON input",
		},
		{
			command: []byte(`Room.RemovePlayer "Tecton"`),
			message: gmcp.RoomRemovePlayer{
				Name: "Tecton",
			},
		},

		{
			command: []byte("Room.WrongDir"),
			err:     "failed hydrating gmcp.RoomWrongDir (): unexpected end of JSON input",
		},
		{
			command: []byte(`Room.WrongDir "ne"`),
			message: gmcp.RoomWrongDir("ne"),
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			message, err := gmcp.Parse(tc.command)

			if tc.err != "" {
				require.NotNil(err, fmt.Sprintf(
					"wanted: %s", tc.err,
				))
				assert.Equal(tc.err, err.Error())
				return
			}

			require.Nil(err)
			assert.Equal(tc.message, message)
		})
	}
}
