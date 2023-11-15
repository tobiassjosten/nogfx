package process

import (
	"strconv"
)

// RepeatInputProcessor unfolds input with a numeric prefix into that input
// repeated the specified number of times.
func RepeatInputProcessor() Processor {
	return MatchInput("{^} {*}", onRepeatInput)
}

func onRepeatInput(match Match, ins, outs [][]byte) ([][]byte, [][]byte, error) {
	postins := [][]byte{}

	previous := -1

	for i, capture := range match {
		for ii := previous + 1; ii < i; ii++ {
			postins = append(postins, ins[ii])
		}
		previous = i

		number, err := strconv.Atoi(string(capture[0]))
		if err != nil {
			postins = append(postins, ins[i])
			continue
		}

		for i := 0; i < number; i++ {
			postins = append(postins, capture[1])
		}
	}

	inslen := len(ins)
	for i := inslen - (inslen - previous) + 1; i < inslen; i++ {
		postins = append(postins, ins[i])
	}

	return postins, nil, nil
}
