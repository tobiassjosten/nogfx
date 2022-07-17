package module

import (
	"fmt"
	"strings"
	"testing"

	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/mock"
	"github.com/tobiassjosten/nogfx/pkg/simpex"

	"github.com/stretchr/testify/assert"
)

var (
	Input  = "input"
	Output = "output"
)

// TestEvent represents a player input or server output.
type TestEvent struct {
	Type string
	Data []string
}

// TestCase is a sequence of inputs/outputs and desired states.
type TestCase struct {
	Events  []TestEvent
	Inputs  []pkg.Input
	Outputs []pkg.Output
}

// Eval plays the TestCase's inputs/outputs and asserts its desired states.
func (tc TestCase) Eval(t *testing.T, constructor pkg.ModuleConstructor) {
	mod := constructor(&mock.WorldMock{})

	var inputs []pkg.Input
	var tcinputs []pkg.Input
	var outputs []pkg.Output
	var tcoutputs []pkg.Output

	for _, event := range tc.Events {
		if event.Type == Input {
			input := pkg.Input{}
			for _, data := range event.Data {
				input = input.Add([]byte(data))
			}
			tcinputs = append(tcinputs, input)

			for _, trigger := range mod.InputTriggers() {
				for i, command := range input {
					match := simpex.Match(trigger.Pattern, command.Text)
					if match == nil {
						continue
					}

					input = trigger.Callback(pkg.TriggerMatch[pkg.Input]{
						Captures: match,
						Content:  input,
						Index:    i,
					})
				}
			}

			inputs = append(inputs, input)
		} else if event.Type == Output {
			output := pkg.Output{}
			for _, data := range event.Data {
				output = output.Add([]byte(data))
			}
			tcoutputs = append(tcoutputs, output)

			for _, trigger := range mod.OutputTriggers() {
				for i, line := range output {
					match := simpex.Match(trigger.Pattern, line.Text)
					if match == nil {
						continue
					}

					output = trigger.Callback(pkg.TriggerMatch[pkg.Output]{
						Captures: match,
						Content:  output,
						Index:    i,
					})
				}
			}

			outputs = append(outputs, output)
		} else {
			t.Logf("invalid event type '%s'", event.Type)
			t.FailNow()
			return
		}
	}

	inexpected := []string{}
	for _, input := range tcinputs {
		for _, command := range input {
			inexpected = append(inexpected, string(command.Text))
		}
	}
	inactual := []string{}
	for _, input := range inputs {
		for _, command := range input {
			inactual = append(inactual, string(command.Text))
		}
	}
	assert.Equal(t, tc.Inputs, inputs, fmt.Sprintf(
		"'%s' vs '%s'",
		strings.Join(inexpected, " | "),
		strings.Join(inactual, " | "),
	))

	outexpected := []string{}
	for _, output := range tcoutputs {
		for _, line := range output {
			outexpected = append(outexpected, string(line.Text))
		}
	}
	outactual := []string{}
	for _, output := range outputs {
		for _, line := range output {
			outactual = append(outactual, string(line.Text))
		}
	}
	assert.Equal(t, tc.Outputs, outputs, fmt.Sprintf(
		"'%s' vs '%s'",
		strings.Join(outexpected, " | "),
		strings.Join(outactual, " | "),
	))
}
