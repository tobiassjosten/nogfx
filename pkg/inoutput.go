package pkg

import "bytes"

type Command struct {
	Text []byte
}

func NewCommand(data []byte) Command {
	return Command{
		Text: data,
	}
}

type Input []Command

func NewInput(data []byte) Input {
	return Input{NewCommand(data)}
}

func (in Input) Add(data []byte) Input {
	return append(in, NewCommand(data))
}

func (in Input) Insert(i int, data []byte) Input {
	return append(in[:i], append([]Command{NewCommand(data)}, in[i:]...)...)
}

func (in Input) Replace(i int, data []byte) Input {
	// @todo Check for existance.
	// @todo Avoid mutating, so as to force using return value consistently.
	in[i] = NewCommand(data)
	return in
}

func (in Input) Split(s []byte) Input {
	for i, c := range in {
		parts := bytes.Split(c.Text, s)
		if len(parts) == 1 {
			continue
		}

		in = in.Replace(i, parts[0])
		for _, data := range parts[1:] {
			in = in.Add(data)
		}
	}

	return in
}

func (in Input) Join(s []byte) Input {
	// @todo Implement…
	return in
}

type Line struct {
	Text []byte
	Raw  []byte

	Before []Line
	After  []Line

	replaced bool
}

func NewLine(data []byte) Line {
	text := []byte{}

	escaped, escaping := false, false
	for _, b := range data {
		if b == 27 {
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

		text = append(text, b)
	}

	return Line{
		Text: text,
		Raw:  data,
	}
}

func (l Line) Replace(data []byte) Line {
	// @todo Log error on multiple mutations.
	// @todo Replace while keeping ANSI colors.
	nl := NewLine(data)
	nl.replaced = true

	return nl
}

// @todo Make this work with custom ANSI colors, resetting back to whatever was
// current before the prefix.
func (l Line) Prefix(data []byte) Line {
	nl := NewLine(append(data, l.Raw...))
	nl.Before = l.Before
	nl.After = l.After
	nl.replaced = l.replaced

	return nl
}

// @todo Make this work with custom ANSI colors, resetting back to whatever was
// current before the suffix.
func (l Line) Suffix(data []byte) Line {
	nl := NewLine(append(l.Raw, data...))
	nl.Before = l.Before
	nl.After = l.After
	nl.replaced = l.replaced

	return nl
}

type Output []Line

func (out Output) Add(data []byte) Output {
	return append(out, NewLine(data))
}

func (out Output) AddBefore(i int, data []byte) Output {
	out[i].Before = append(out[i].Before, NewLine(data))
	return out
}

func (out Output) AddAfter(i int, data []byte) Output {
	out[i].After = append(out[i].After, NewLine(data))
	return out
}

func (out Output) Insert(i int, data []byte) Output {
	return append(out[:i], append([]Line{NewLine(data)}, out[i:]...)...)
}

func (out Output) Remove(i int) Output {
	return append(out[:i], out[i+1:]...)
}

func (out Output) Replace(i int, data []byte) Output {
	// @todo Check for existance.
	out[i] = out[i].Replace(data)
	// @todo Mutate also RawLines.
	return out
}

func (out Output) Suffix(i int, data []byte) Output {
	// @todo Check for existance.
	out[i] = out[i].Suffix(data)
	return out
}

func (out Output) Lines() []Line {
	var ls []Line

	for _, line := range out {
		ls = append(ls, line.Before...)
		ls = append(ls, line)
		ls = append(ls, line.After...)
	}

	return ls
}
