package process

import "fmt"

// Processor parses, makes sense of, and reacts to input from the player and
// output from the game.
type Processor func(ins, outs [][]byte) (postins, postouts [][]byte, err error)

// ChainProcessor creates a Processor to run all the given Processors in order.
func ChainProcessor(ps ...Processor) Processor {
	return func(ins, outs [][]byte) ([][]byte, [][]byte, error) {
		for i, p := range ps {
			pins, pouts, err := p(ins, outs)
			if err != nil {
				return nil, nil, fmt.Errorf("processor %d failed: %w", i, err)
			}

			ins, outs = pins, pouts
		}

		return ins, outs, nil
	}
}
