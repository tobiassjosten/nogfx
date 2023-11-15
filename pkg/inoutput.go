package pkg

/*

import (
	"bytes"

	"github.com/tobiassjosten/nogfx/pkg/gmcp"
	"github.com/tobiassjosten/nogfx/pkg/telnet"
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
	Tags   Tags
	Input  Exput
	Output Exput
}

// NewInoutput creates a new Inoutput.
func NewInoutput(input [][]byte, output [][]byte) (inout Inoutput) {
	inout.Tags = Tags{}

	for _, data := range input {
		inout.Input = inout.Input.Append(data)
	}

	for _, data := range output {
		inout.Output = inout.Output.Append(data)
	}

	return
}

// Commands returns all telnet commnds from the Output.
func (inout Inoutput) Commands() (commands [][]byte) {
	for _, line := range inout.Output {
		if line.Text[0] != telnet.IAC {
			continue
		}
		commands = append(commands, line.Text)
	}

	return
}

func (inout Inoutput) GMCPs() (gmcps [][]byte) {
	for _, cmd := range inout.Commands() {
		if data := gmcp.Unwrap(cmd); data != nil {
			gmcps = append(gmcps, data)
		}
	}

	return
}

func (inout Inoutput) AddTag(tag Tag, values ...any) Inoutput {
	inout.Tags = inout.Tags.Add(tag, values)
	return inout
}

// AppendInput is a shortcut for Inoutput.Input.Append().
func (inout Inoutput) AppendInput(data []byte) Inoutput {
	inout.Input = inout.Input.Append(data)
	return inout
}

// InsertInput is a shortcut for Inoutput.Input.Insert().
func (inout Inoutput) InsertInput(i int, data []byte) Inoutput {
	inout.Input = inout.Input.Insert(i, data)
	return inout
}

// RemoveInputMatching is a shortcut for Inoutput.Input.RemoveMatching().
func (inout Inoutput) RemoveInputMatching(data []byte) Inoutput {
	inout.Input = inout.Input.RemoveMatching(data)
	return inout
}

// RemoveInput is a shortcut for Inoutput.Input.Remove().
func (inout Inoutput) RemoveInput(i int) Inoutput {
	inout.Input = inout.Input.Remove(i)
	return inout
}

// ReplaceInput is a shortcut for Inoutput.Input.Replace().
func (inout Inoutput) ReplaceInput(i int, data []byte) Inoutput {
	inout.Input = inout.Input.Replace(i, data)
	return inout
}

// HasOutput checks if the Output has a text matching the given byte slice.
func (inout Inoutput) HasOutput(data []byte) bool {
	for _, out := range inout.Output {
		if bytes.Equal(data, out.Text) {
			return true
		}
	}
	return false
}

// AppendOutput is a shortcut for Inoutput.Output.Append().
func (inout Inoutput) AppendOutput(data []byte) Inoutput {
	inout.Output = inout.Output.Append(data)
	return inout
}

// InsertOutput is a shortcut for Inoutput.Output.Insert().
func (inout Inoutput) InsertOutput(i int, data []byte) Inoutput {
	inout.Output = inout.Output.Insert(i, data)
	return inout
}

// RemoveOutputMatching is a shortcut for Inoutput.Output.RemoveMatching().
func (inout Inoutput) RemoveOutputMatching(data []byte) Inoutput {
	inout.Output = inout.Output.RemoveMatching(data)
	return inout
}

// RemoveOutput is a shortcut for Inoutput.Output.Remove().
func (inout Inoutput) RemoveOutput(i int) Inoutput {
	inout.Output = inout.Output.Remove(i)
	return inout
}

// ReplaceOutput is a shortcut for Inoutput.Output.Replace().
func (inout Inoutput) ReplaceOutput(i int, data []byte) Inoutput {
	inout.Output = inout.Output.Replace(i, data)
	return inout
}

// AddOutputTag is a shortcut for Inoutput.Output.AddOutputTag().
func (inout Inoutput) AddOutputTag(i int, tag Tag, values ...any) Inoutput {
	inout.Output = inout.Output.AddTag(i, tag, values...)
	return inout
}

// Text is a byte slice, with some utility methods, used for Input and Output.
type Text []byte

// Clean removes ANSI colors from the Text.
func (txt Text) Clean() []byte {
	clean := []byte{}

	var escape byte = 27
	escaped, escaping := false, false
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

	Tags Tags
}

// Exput represents both Input and Output.
type Exput []Line

// NewExput creates a new Exput.
func NewExput(data []byte) Exput {
	return Exput{Line{Text: data, Tags: Tags{}}}
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

// Append appends a new Line to the Exput
func (ex Exput) Append(data []byte) Exput {
	return append(ex, Line{Text: data})
}

// Insert creates a Line from the given data and inserts it at the designated
// position in the Exput, shifting the existing Line and all following one
// position back.
func (ex Exput) Insert(i int, data []byte) Exput {
	newex := append(Exput{}, ex...)
	newex = append(append(newex[:i], NewExput(data)...), newex[i:]...)
	return newex
}

// RemoveMatching removes lines matching the given data.
func (ex Exput) RemoveMatching(data []byte) Exput {
	for i := 0; i < len(ex); {
		if bytes.Equal(data, ex[i].Text.Clean()) {
			ex = append(ex[:i], ex[i+1:]...)
			continue
		}
		i++
	}
	return ex
}

// Remove deletes a given Line.
func (ex Exput) Remove(i int) Exput {
	return append(ex[:i], ex[i+1:]...)
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

// AddTag appends the given tag to the Exput's slice of tags.
func (ex Exput) AddTag(i int, tag Tag, values ...any) Exput {
	newex := append(Exput{}, ex...)
	newex[i].Tags = newex[i].Tags.Add(tag, values...)
	return newex
}

// Bytes assembles the Exput into a slice of byte slices.
func (ex Exput) Bytes() (bs [][]byte) {
	for _, ln := range ex {
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

*/
