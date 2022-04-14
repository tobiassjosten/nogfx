package achaea_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/world/achaea"
)

func TestWorldBasics(t *testing.T) {
	assert := assert.New(t)

	ui := &pkg.UIMock{}
	client := &pkg.ClientMock{}

	world := achaea.NewWorld(ui, client)

	input := []byte("input")
	assert.Equal(input, world.Input(input))

	output := []byte("output")
	assert.Equal(output, world.Output(output))

	// assert.Nil(world.Command([]byte("command")))
}
