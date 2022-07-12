package world_test

import (
	"testing"

	"github.com/tobiassjosten/nogfx/pkg/mock"
	"github.com/tobiassjosten/nogfx/pkg/world"

	"github.com/stretchr/testify/assert"
)

func TestGenericWorld(t *testing.T) {
	client := &mock.ClientMock{}

	ui := &mock.UIMock{}

	world := world.NewGenericWorld(client, ui)

	in := []byte("asdf")

	inputs := [][]byte{[]byte("input")}
	assert.Equal(t, inputs, world.ProcessInput(inputs[0]))

	outputs := [][]byte{[]byte("output")}
	assert.Equal(t, outputs, world.ProcessOutput(outputs[0]))

	assert.Nil(t, world.ProcessCommand(in))
}
