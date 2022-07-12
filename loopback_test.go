package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoopback(t *testing.T) {
	loopback := NewLoopback()

	done := make(chan struct{})

	in := []byte{'a'}
	go func() {
		_, err := loopback.Write(in)
		require.Nil(t, err)
	}()

	out := make([]byte, 1)
	go func() {
		_, err := loopback.Read(out)
		require.Nil(t, err)
		done <- struct{}{}
	}()

	<-done

	assert.Equal(t, in, out)
}
