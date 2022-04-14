package pkg_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tobiassjosten/nogfx/pkg"
)

func TestGenericWorld(t *testing.T) {
	assert := assert.New(t)

	ui := &pkg.UIMock{}
	client := &pkg.ClientMock{}

	world := pkg.NewGenericWorld(ui, client)

	input := []byte("input")
	assert.Equal(input, world.Input(input))

	output := []byte("output")
	assert.Equal(output, world.Output(output))

	assert.Nil(world.Command([]byte("command")))
}
