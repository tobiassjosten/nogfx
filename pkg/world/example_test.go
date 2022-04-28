package world_test

import (
	"testing"

	"github.com/tobiassjosten/nogfx/pkg/mock"
	"github.com/tobiassjosten/nogfx/pkg/world"

	"github.com/stretchr/testify/assert"
)

func TestExampleWorld(t *testing.T) {
	client := &mock.ClientMock{}

	printed := [][]byte{}

	ui := &mock.UIMock{
		AddVitalFunc: func(_ string, _ interface{}) {},
		PrintFunc: func(data []byte) {
			printed = append(printed, data)
		},
	}

	world := world.NewExampleWorld(client, ui)

	in := []byte("asdf")

	assert.Equal(t, [][]byte{}, printed)
	assert.Equal(t, in, world.ProcessInput(in))
	assert.Equal(t, [][]byte{append([]byte("> "), in...)}, printed)

	assert.Equal(t, in, world.ProcessOutput(in))

	assert.Nil(t, world.ProcessCommand(in))
}
