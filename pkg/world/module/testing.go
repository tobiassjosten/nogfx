package module

import (
	"bytes"
	"testing"

	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/mock"

	"github.com/stretchr/testify/assert"
)

// TestEvent represents a player input or server output.
type TestEvent struct {
	Input bool
	Data  []byte
}

// NewTestEvent creates a new TestEvent.
func NewTestEvent(input bool, data []byte) TestEvent {
	return TestEvent{
		Input: input,
		Data:  data,
	}
}

// TestCase is a sequence of inputs/outputs and desired states.
type TestCase struct {
	Events  []TestEvent
	Inputs  [][]byte
	Outputs [][]byte
	Sent    [][]byte
}

// Eval plays the TestCase's inputs/outputs and asserts its desired states.
func (tc TestCase) Eval(t *testing.T, mod pkg.Module, client *mock.ClientMock) {
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
