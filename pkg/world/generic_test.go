package world_test

import (
	"testing"

	"github.com/tobiassjosten/nogfx/pkg/mock"
	"github.com/tobiassjosten/nogfx/pkg/world"

	"github.com/stretchr/testify/assert"
)

func TestGenericWorld(t *testing.T) {
	client := &mock.ClientMock{}

	ui := &mock.UIMock{
		AddVitalFunc: func(_ string, _ interface{}) {},
	}

	world := world.NewGenericWorld(client, ui)

	in := []byte("asdf")

	assert.Equal(t, in, world.ProcessInput(in))
	assert.Equal(t, in, world.ProcessOutput(in))
	assert.Nil(t, world.ProcessCommand(in))
}
