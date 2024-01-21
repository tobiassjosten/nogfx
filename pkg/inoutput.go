package pkg

import (
	"bytes"
)

// IOKind signifies the direction of the IO, whether it's player input or
// server output.
type IOKind string

// These are the known directions of IO.
const (
	Input  = IOKind("input")
	Output = IOKind("output")
)

// Inoutput collects one paragraph of output lines and one list of input
// commands, for processing and dispatching to UI/client.
type Inoutput struct {
	Input  Exput
	Output Exput
}

// NewInoutput creates a new Inoutput.
func NewInoutput(input [][]byte, output [][]byte) (inout Inoutput) {
	for _, data := range input {
		inout.Input = inout.Input.Add(data)
	}

	for _, data := range output {
		inout.Output = inout.Output.Add(data)
	}

	return
}

// AddBeforeInput is a shortcut for Inoutput.Input.AddBefore().
func (inout Inoutput) AddBeforeInput(i int, data []byte) Inoutput {
	inout.Input = inout.Input.AddBefore(i, data)
	return inout
}

// AddAfterInput is a shortcut for Inoutput.Input.AddAfter().
func (inout Inoutput) AddAfterInput(i int, data []byte) Inoutput {
	inout.Input = inout.Input.AddAfter(i, data)
	return inout
}

// OmitInput is a shortcut for Inoutput.Input.Omit().
func (inout Inoutput) OmitInput(i int) Inoutput {
	inout.Input = inout.Input.Omit(i)
	return inout
}

// AddBeforeOutput is a shortcut for Inoutput.Output.AddBefore().
func (inout Inoutput) AddBeforeOutput(i int, data []byte) Inoutput {
	inout.Output = inout.Output.AddBefore(i, data)
	return inout
}

// AddAfterOutput is a shortcut for Inoutput.Output.AddAfter().
func (inout Inoutput) AddAfterOutput(i int, data []byte) Inoutput {
	inout.Output = inout.Output.AddAfter(i, data)
	return inout
}

// OmitOutput is a shortcut for Inoutput.Output.Omit().
func (inout Inoutput) OmitOutput(i int) Inoutput {
	inout.Output = inout.Output.Omit(i)
	return inout
}

// Text is a byte slice, with some utility methods, used for Input and Output.
type Text []byte

// Clean removes ANSI colors from the Text.
func (txt Text) Clean() []byte {
	clean := []byte{}

	var (
		escape   byte = 27
		escaped       = false
		escaping      = false
	)

	for _, b := range txt {
		if b == escape {
			escaped = true
			continue
		}

		if escaped {
			escaped = false

			if b == '[' {
				escaping = true
				continue
			}
		}

		if escaping {
			if b == 'm' {
				escaping = false
			}

			continue
		}

		clean = append(clean, b)
	}

	return clean
}

// Replace changes the visible parts of a Text while retaining ANSI colors.
func (txt Text) Replace(data []byte) Text {
	// @todo Replace while keeping ANSI colors.
	return data
}

// Line bundles Texts with some metadata, for use in Input and Output.
type Line struct {
	Text Text

	Before []Text
	After  []Text

	omitted bool
}

// Exput represents both Input and Output.
type Exput []Line

// NewExput creates a new Exput.
func NewExput(data []byte) Exput {
	return Exput{Line{Text: data}}
}

// Inoutput creates an Inoutput based on the Exput data.
func (ex Exput) Inoutput(kind IOKind) Inoutput {
	switch kind {
	case Input:
		return Inoutput{Input: ex}
	case Output:
		return Inoutput{Output: ex}
	}

	return Inoutput{}
}

// Add appends a new Line to the Exput.
func (ex Exput) Add(data []byte) Exput {
	return append(ex, Line{Text: data})
}

// AddBefore appends a new Line before a given other Line. This is useful for
// data that belongs together, as omitting one line also omits all lines added
// before and after it (but not independently of it, with Add()).
func (ex Exput) AddBefore(i int, data []byte) Exput {
	newex := append(Exput{}, ex...)
	newex[i].Before = append(newex[i].Before, data)

	return newex
}

// AddAfter appends a new Line after a given other Line. This is useful for
// data that belongs together, as omitting one line also omits all lines added
// before and after it (but not independently of it, with Add()).
func (ex Exput) AddAfter(i int, data []byte) Exput {
	newex := append(Exput{}, ex...)
	newex[i].After = append(newex[i].After, data)

	return newex
}

// Omit flags a Line to be omitted from the Bytes() summary.
func (ex Exput) Omit(i int) Exput {
	newex := append(Exput{}, ex...)
	newex[i].omitted = true

	return newex
}

// Replace changes the visible parts of a Line while retaining ANSI colors.
func (ex Exput) Replace(i int, data []byte) Exput {
	newex := append(Exput{}, ex...)
	newex[i].Text = newex[i].Text.Replace(data)

	return newex
}

// Split breaks down Lines by the given separator.
func (ex Exput) Split(s []byte) Exput {
	for i, c := range ex {
		parts := bytes.Split(c.Text, s)
		if len(parts) == 1 {
			continue
		}

		ex = ex.Replace(i, parts[0])
		for ii, data := range parts[1:] {
			ex = append(
				ex[:i+ii+1],
				append(NewExput(data), ex[i+ii+1:]...)...,
			)
		}
	}

	return ex
}

// Bytes assembles the Exput into a slice of byte slices.
func (ex Exput) Bytes() (bs [][]byte) {
	for _, ln := range ex {
		if ln.omitted {
			continue
		}

		for _, text := range ln.Before {
			bs = append(bs, text)
		}

		bs = append(bs, ln.Text)

		for _, text := range ln.After {
			bs = append(bs, text)
		}
	}

	return
}
