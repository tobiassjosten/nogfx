package tui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapLines(test *testing.T) {
	events := [][]rune{[]rune("asdfg")}
	lines := wrapLines(events, 2)

	expected := [][]rune{[]rune("as"), []rune("df"), []rune("g")}
	assert.Equal(test, expected, lines)
}
