package tui

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountLines(t *testing.T) {
	assert := assert.New(t)

	tcs := []struct {
		text  []rune
		width int
		lines int
	}{
		{
			text:  []rune("asdf"),
			width: 5,
			lines: 1,
		},
		{
			text:  []rune("asdf\n"),
			width: 5,
			lines: 2,
		},
		{
			text:  []rune("asdf"),
			width: 4,
			lines: 1,
		},
		{
			text:  []rune("asdf"),
			width: 2,
			lines: 2,
		},
		{
			text:  []rune("as\ndf"),
			width: 2,
			lines: 2,
		},
		{
			text:  []rune("as\ndfg"),
			width: 2,
			lines: 3,
		},
		{
			text:  []rune("asdf"),
			width: 1,
			lines: 4,
		},
		{
			text:  []rune("as\ndf\n"),
			width: 1,
			lines: 5,
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			assert.Equal(
				tc.lines,
				countLines(tc.text, tc.width),
				fmt.Sprintf("%s / %d", strings.ReplaceAll(string(tc.text), "\n", "\\n"), tc.width),
			)
		})
	}
}
