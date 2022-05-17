package module

import (
	"bytes"
	"testing"

	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/mock"

	"github.com/stretchr/testify/assert"
)

type TestEvent struct {
	Input bool
	Data  []byte
}

func NewTestEvent(input bool, data []byte) TestEvent {
	return TestEvent{
		Input: input,
		Data:  data,
	}
}

type TestCase struct {
	Events  []TestEvent
	Inputs  [][]byte
	Outputs [][]byte
	Sent    [][]byte
}

func Test(t *testing.T, mod pkg.Module, tc TestCase, client *mock.ClientMock) {
	var inputs, outputs, sent [][]byte

	for _, event := range tc.Events {
		if event.Input {
			inputs = append(inputs, mod.ProcessInput(event.Data)...)
		} else {
			outputs = append(outputs, mod.ProcessOutput(event.Data)...)
		}
	}

	for _, call := range client.SendCalls() {
		sent = append(sent, call.Bytes)
	}

	assert.Equal(t, tc.Inputs, inputs, "inputs was "+string(bytes.Join(
		inputs, []byte(" | "),
	)))
	assert.Equal(t, tc.Outputs, outputs, "outputs was: "+string(bytes.Join(
		outputs, []byte(" | "),
	)))
	assert.Equal(t, tc.Sent, sent, "sent was: "+string(bytes.Join(
		sent, []byte(" | "),
	)))
}
