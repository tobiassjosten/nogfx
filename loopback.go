package main

import (
	"io"
	"net"
)

// Loopback is a mock connection that reads and writes to itself, for use when
// connecting to "example.com" for testing purposes.
type Loopback struct {
	client io.Writer
	server io.Reader
}

// NewLoopback creates a new Loopback from a network pipe.
func NewLoopback() *Loopback {
	client, server := net.Pipe()
	return &Loopback{client, server}
}

// Read returns data that has previously been written.
func (l Loopback) Read(data []byte) (n int, err error) {
	return l.server.Read(data)
}

// Write adds data to later be read.
func (l Loopback) Write(data []byte) (n int, err error) {
	return l.client.Write(data)
}
