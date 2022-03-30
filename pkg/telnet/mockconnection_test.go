package telnet_test

import (
	"io"
	"math"
)

type MockConnection struct {
	Output []byte
	Input  []byte
	closed bool
}

func NewMockConnection() *MockConnection {
	return &MockConnection{[]byte{}, []byte{}, false}
}

func (connection *MockConnection) Read(data []byte) (int, error) {
	if connection.closed {
		panic("mock connection closed")
	}

	if len(connection.Output) == 0 {
		return 0, io.EOF
	}

	length := int(math.Min(float64(len(data)), float64(len(connection.Output))))

	for i := 0; i < length; i++ {
		data[i] = connection.Output[i]
	}

	connection.Output = connection.Output[length:]

	return length, nil
}

func (connection *MockConnection) Write(data []byte) (int, error) {
	if connection.closed {
		panic("mock connection closed")
	}

	connection.Input = append(connection.Input, data...)

	return len(data), nil
}

func (connection *MockConnection) Close() error {
	if connection.closed {
		panic("mock connection already closed")
	}

	connection.closed = true
	return nil
}
