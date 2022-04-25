package tui

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMax(t *testing.T) {
	tcs := []struct {
		in  []int
		out int
	}{
		{
			in:  []int{1, 0},
			out: 1,
		},
		{
			in:  []int{0, 1},
			out: 1,
		},
		{
			in:  []int{1, 1},
			out: 1,
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			assert.Equal(t, tc.out, max(tc.in[0], tc.in[1]))
		})
	}
}

func TestMin(t *testing.T) {
	tcs := []struct {
		in  []int
		out int
	}{
		{
			in:  []int{1, 0},
			out: 0,
		},
		{
			in:  []int{0, 1},
			out: 0,
		},
		{
			in:  []int{1, 1},
			out: 1,
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			assert.Equal(t, tc.out, min(tc.in[0], tc.in[1]))
		})
	}
}
