package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddress(t *testing.T) {
	tcs := []struct {
		in  string
		out string
		err string
	}{
		{
			in:  "",
			out: "example.com:23",
		},
		{
			in:  "achaea.com",
			out: "achaea.com:23",
		},
		{
			in:  "achaea.com:23",
			out: "achaea.com:23",
		},
		{
			in:  "imperianea.com",
			out: "imperianea.com:23",
		},
		{
			in:  "imperianea.com:23",
			out: "imperianea.com:23",
		},
		{
			in:  ":",
			err: "invalid address ':'",
		},
		{
			in:  "a:",
			err: "invalid address 'a:'",
		},
		{
			in:  "a:a",
			err: "invalid port 'a'",
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			out, err := address(tc.in)

			if tc.err != "" {
				assert.Equal(t, tc.err, err.Error())
				return
			}

			assert.Equal(t, tc.out, out)
		})
	}
}
