package testing

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tobiassjosten/nogfx/pkg"
)

type IOEvent struct {
	Kind   pkg.IOKind
	Input  []string
	Output []string
}

func IOEIns(inputs []string, output ...string) IOEvent {
	te := IOEvent{
		Kind:  pkg.Input,
		Input: inputs,
	}

	if len(output) > 0 {
		te.Output = output
	}

	return te
}

func IOEIn(input string, output ...string) IOEvent {
	return IOEIns([]string{input}, output...)
}

func IOEOuts(output []string, input ...string) IOEvent {
	te := IOEvent{
		Kind:   pkg.Output,
		Output: output,
	}

	if len(input) > 0 {
		te.Input = input
	}

	return te
}

func IOEOut(output string, input ...string) IOEvent {
	return IOEOuts([]string{output}, input...)
}

// Bytes converts the event data to slices of slice bytes.
func (te IOEvent) Bytes() (bs [][]byte) {
	var datas []string

	switch te.Kind {
	case pkg.Input:
		datas = te.Input
	case pkg.Output:
		datas = te.Output
	}

	for _, data := range datas {
		bs = append(bs, []byte(data))
	}

	return
}

// Inoutput wraps the event data with an Inoutput.
func (te IOEvent) Inoutput() pkg.Inoutput {
	inout := pkg.Inoutput{}

	for _, data := range te.Input {
		inout.Input = inout.Input.Append([]byte(data))
	}
	for _, data := range te.Output {
		inout.Output = inout.Output.Append([]byte(data))
	}

	return inout
}

type IOTestCase struct {
	Events    []IOEvent
	Inoutputs []pkg.Inoutput
}

// Eval plays the TestCase's inputs/outputs and asserts its desired states.
func (tc IOTestCase) Eval(t *testing.T, mod pkg.Module) {
	var inouts []pkg.Inoutput

	for _, event := range tc.Events {
		inout := event.Inoutput()

		for _, trigger := range mod.Triggers() {
			if trigger.Kind == pkg.Input && len(inout.Input) > 0 {
				inout = trigger.Match(event.Bytes(), inout)
			}
			if trigger.Kind == pkg.Output && len(inout.Output) > 0 {
				inout = trigger.Match(event.Bytes(), inout)
			}
		}

		inouts = append(inouts, inout)
	}

	assert.Equal(t, len(tc.Inoutputs), len(inouts))
	for i, inout := range tc.Inoutputs {
		if i >= len(inouts) {
			break
		}
		assert.Equal(t, inout, inouts[i])
	}
}

func IO(input, output string) pkg.Inoutput {
	return pkg.NewInoutput(
		[][]byte{[]byte(input)},
		[][]byte{[]byte(output)},
	)
}

func IOIn(input string) pkg.Inoutput {
	return pkg.NewInoutput([][]byte{[]byte(input)}, nil)
}

func IOIns(inputs []string) pkg.Inoutput {
	var bs [][]byte
	for _, str := range inputs {
		bs = append(bs, []byte(str))
	}

	return pkg.NewInoutput(bs, nil)
}

func IOOut(output string) pkg.Inoutput {
	return pkg.NewInoutput(nil, [][]byte{[]byte(output)})
}

func IOOuts(outputs []string) pkg.Inoutput {
	var bs [][]byte
	for _, str := range outputs {
		bs = append(bs, []byte(str))
	}

	return pkg.NewInoutput(nil, bs)
}
