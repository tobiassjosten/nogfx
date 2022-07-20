package pkg

import (
	"bytes"
)

type IOKind string

const (
	Input  = IOKind("input")
	Output = IOKind("output")
)

type Inoutput struct {
	Input  Exput
	Output Exput
}

func NewInoutput(input [][]byte, output [][]byte) (inout Inoutput) {
	for _, data := range input {
		inout.Input = inout.Input.Add(data)
	}
	for _, data := range output {
		inout.Output = inout.Output.Add(data)
	}
	return
}

func (inout Inoutput) AddBeforeInput(i int, data []byte) Inoutput {
	inout.Input = inout.Input.AddBefore(i, data)
	return inout
}

func (inout Inoutput) AddAfterInput(i int, data []byte) Inoutput {
	inout.Input = inout.Input.AddAfter(i, data)
	return inout
}

func (inout Inoutput) AddBeforeOutput(i int, data []byte) Inoutput {
	inout.Output = inout.Output.AddBefore(i, data)
	return inout
}

func (inout Inoutput) AddAfterOutput(i int, data []byte) Inoutput {
	inout.Output = inout.Output.AddAfter(i, data)
	return inout
}

type Text []byte

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

func (txt Text) Replace(data []byte) Text {
	// @todo Replace while keeping ANSI colors.
	return data
}

type Line struct {
	Text Text

	Before []Text
	After  []Text

	omitted bool
}

type Exput []Line

func NewExput(data []byte) Exput {
	return Exput{Line{Text: data}}
}

func (ex Exput) Inoutput(kind IOKind) Inoutput {
	switch kind {
	case Input:
		return Inoutput{Input: ex}
	case Output:
		return Inoutput{Output: ex}
	}
	return Inoutput{}
}

func (ex Exput) Add(data []byte) Exput {
	return append(ex, Line{Text: data})
}

func (ex Exput) AddBefore(i int, data []byte) Exput {
	newex := append(Exput{}, ex...)
	newex[i].Before = append(newex[i].Before, data)
	return newex
}

func (ex Exput) AddAfter(i int, data []byte) Exput {
	newex := append(Exput{}, ex...)
	newex[i].After = append(newex[i].After, data)
	return newex
}

func (ex Exput) Omit(i int) Exput {
	newex := append(Exput{}, ex...)
	newex[i].omitted = true
	return newex
}

func (ex Exput) Replace(i int, data []byte) Exput {
	newex := append(Exput{}, ex...)
	newex[i].Text = newex[i].Text.Replace(data)
	return newex
}

func (ex Exput) Split(s []byte) Exput {
	for i, c := range ex {
		parts := bytes.Split(c.Text, s)
		if len(parts) == 1 {
			continue
		}

		ex = ex.Replace(i, parts[0])
		for _, data := range parts[1:] {
			ex = ex.Add(data)
		}
	}

	return ex
}

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
