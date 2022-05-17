package world_test

import (
	"testing"

	"github.com/tobiassjosten/nogfx/pkg/mock"
	"github.com/tobiassjosten/nogfx/pkg/navigation"
	"github.com/tobiassjosten/nogfx/pkg/world"

	"github.com/stretchr/testify/assert"
)

func TestExampleWorld(t *testing.T) {
	client := &mock.ClientMock{}

	printed := [][]byte{}

	outputs := make(chan []byte, 1000)

	ui := &mock.UIMock{
		AddVitalFunc: func(_ string, _ interface{}) error {
			return nil
		},
		OutputsFunc: func() chan<- []byte {
			return outputs
		},
		PrintFunc: func(data []byte) {
			printed = append(printed, data)
		},
		SetRoomFunc: func(_ *navigation.Room) {},
		UpdateVitalFunc: func(_ string, _, _ int) error {
			return nil
		},
	}

	world := world.NewExampleWorld(client, ui)

	in := []byte("asdf")

	assert.Equal(t, [][]byte{}, printed)
	assert.Equal(t, [][]byte{in}, world.ProcessInput(in))
	assert.Equal(t, [][]byte{append([]byte("> "), in...)}, printed)

	assert.Equal(t, [][]byte{in}, world.ProcessOutput(in))

	assert.Nil(t, world.ProcessCommand(in))
}
