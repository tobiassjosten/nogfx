package world_test

import (
	"testing"

	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/mock"
	"github.com/tobiassjosten/nogfx/pkg/world"

	"github.com/stretchr/testify/assert"
)

func TestGenericWorld(t *testing.T) {
	world := world.NewGenericWorld(&mock.ClientMock{}, &mock.UIMock{})

	data := []byte("input")

	input := pkg.Input{pkg.NewCommand(data)}
	assert.Equal(t, data, world.ProcessInput(input)[0])

	output := pkg.Output{pkg.NewLine(data)}
	assert.Equal(t, data, world.ProcessOutput(output)[0])
}
