package process

import (
	"fmt"
	"os"
	"path/filepath"
)

// LogProcessor writes input and output to the given file path. The parent
// directories are created if they don't already exist.
func LogProcessor(dir, filename string) (Processor, error) {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to create logs directory %q: %w", dir, err)
	}

	path := filepath.Join(dir, filename)

	file, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file %q: %s", path, err)
	}

	return func(ins, outs [][]byte) ([][]byte, [][]byte, error) {
		for _, text := range append(ins, outs...) {
			if _, err := file.Write([]byte(text)); err != nil {
				return nil, nil, fmt.Errorf(
					"failed to write to log: %s", err,
				)
			}
		}

		return ins, outs, nil
	}, nil
}
