package gmcp_test

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/tobiassjosten/nogfx/pkg/gmcp"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoomServerMessages(t *testing.T) {
	tcs := []struct {
		command []byte
		message gmcp.ServerMessage
		err     string
		log     string
	}{
		{
			command: []byte("Room.Info"),
			err:     "failed hydrating gmcp.RoomInfo: unexpected end of JSON input",
		},
		{
			command: []byte(`Room.Info {"details": [ "asdf" ] }`),
			log:     "unknown Room.Info detail 'asdf'\n",
			message: gmcp.RoomInfo{},
		},
		{
			command: []byte(`Room.Info {"coords": "123"}`),
			err:     `failed hydrating gmcp.RoomInfo: failed parsing coords '[123]'`,
		},
		{
			command: []byte(`Room.Info {"coords": "a,5,4,3"}`),
			err:     `failed hydrating gmcp.RoomInfo: failed parsing area number from coords '[a 5 4 3]': strconv.Atoi: parsing "a": invalid syntax`,
		},
		{
			command: []byte(`Room.Info {"coords": "45,b,4,3"}`),
			err:     `failed hydrating gmcp.RoomInfo: failed parsing x from coords '[45 b 4 3]': strconv.Atoi: parsing "b": invalid syntax`,
		},
		{
			command: []byte(`Room.Info {"coords": "45,5,c,3"}`),
			err:     `failed hydrating gmcp.RoomInfo: failed parsing y from coords '[45 5 c 3]': strconv.Atoi: parsing "c": invalid syntax`,
		},
		{
			command: []byte(`Room.Info {"coords": "45,5,4,d"}`),
			err:     `failed hydrating gmcp.RoomInfo: failed parsing building from coords '[45 5 4 d]': strconv.Atoi: parsing "d": invalid syntax`,
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
			err:     "failed hydrating gmcp.RoomPlayers: unexpected end of JSON input",
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
			err:     "failed hydrating gmcp.RoomAddPlayer: unexpected end of JSON input",
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
			err:     "failed hydrating gmcp.RoomRemovePlayer: unexpected end of JSON input",
		},
		{
			command: []byte(`Room.RemovePlayer "Tecton"`),
			message: gmcp.RoomRemovePlayer{
				Name: "Tecton",
			},
		},

		{
			command: []byte("Room.WrongDir"),
			err:     "failed hydrating gmcp.RoomWrongDir: unexpected end of JSON input",
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

			var buf bytes.Buffer
			log.SetOutput(&buf)
			logFlags := log.Flags()
			log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
			defer func() {
				log.SetOutput(os.Stderr)
				log.SetFlags(logFlags)
			}()

			message, err := gmcp.Parse(tc.command, gmcp.ServerMessages)

			if tc.err != "" {
				require.NotNil(err, fmt.Sprintf(
					"wanted: %s", tc.err,
				))
				assert.Equal(tc.err, err.Error())
				return
			}

			if tc.log != "" {
				assert.Equal(tc.log, buf.String())
			}

			require.Nil(err)
			assert.Equal(tc.message, message)
		})
	}
}
