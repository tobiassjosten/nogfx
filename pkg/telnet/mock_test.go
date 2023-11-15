package telnet_test

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"time"
)

type ErrMock struct{}

func (ErrMock) Error() string {
	return "mock error"
}

var errMock error = ErrMock{}

type mockAddr struct{}

func (mockAddr) Network() string { return "mock" }
func (mockAddr) String() string  { return "mock" }

type MockReader func(p []byte) (n int, err error)

func (rd MockReader) Read(p []byte) (n int, err error) {
	return rd(p)
}

// MockConn is a mock of net.Conn, with two byte slices. The false key signifies
// the server output (which we read from) and true holds the data written.
type MockConn struct {
	io.Reader
	Closed   bool
	Written  []byte
	WriteErr func([]byte) error
}

var _ net.Conn = &MockConn{}

func NewMockConn(read []byte) *MockConn {
	return &MockConn{
		Reader: bufio.NewReader(bytes.NewReader(read)),
	}
}

func (mc *MockConn) Read(b []byte) (n int, err error) {
	if mc.Closed {
		return 0, net.ErrClosed
	}

	return mc.Reader.Read(b)
}

func (mc *MockConn) Write(b []byte) (n int, err error) {
	if mc.WriteErr != nil {
		err := mc.WriteErr(b)
		if err != nil {
			return 0, err
		}
	}

	if mc.Closed {
		return 0, net.ErrClosed
	}

	mc.Written = append(mc.Written, b...)

	return len(b), nil
}

func (mc *MockConn) Close() error {
	mc.Closed = true

	return nil
}

func (mc *MockConn) LocalAddr() net.Addr {
	return &mockAddr{}
}

func (mc *MockConn) RemoteAddr() net.Addr {
	return &mockAddr{}
}

func (mc *MockConn) SetDeadline(_ time.Time) error {
	return nil
}

func (mc *MockConn) SetReadDeadline(_ time.Time) error {
	return nil
}

func (mc *MockConn) SetWriteDeadline(_ time.Time) error {
	return nil
}
