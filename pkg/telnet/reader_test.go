package telnet_test

import (
	"bufio"
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tobiassjosten/nogfx/pkg/telnet"
)

func TestRead(test *testing.T) {
	output := []byte("Hello from server\n")

	connection := NewMockConnection()
	connection.Output = output
	stream := telnet.NewStream(connection)

	actual, err := ioutil.ReadAll(stream)

	require.NoError(test, err)
	assert.Equal(test, output, actual)
}

func TestReadEOF(test *testing.T) {
	connection := NewMockConnection()
	connection.Output = []byte("A")
	stream := telnet.NewStream(connection)

	reader := bufio.NewReader(stream)
	buffer := make([]byte, 2)
	count, err := reader.Read(buffer)

	require.NoError(test, err)
	assert.Equal(test, 1, count)
	assert.Equal(test, byte('A'), buffer[0])

	count, err = reader.Read(buffer)

	assert.Equal(test, err, io.EOF)
	assert.Equal(test, 0, count)
}

func TestReadIACEscape(test *testing.T) {
	output := []byte{telnet.IAC, telnet.IAC, '\n'}

	connection := NewMockConnection()
	connection.Output = output
	stream := telnet.NewStream(connection)

	actual, err := ioutil.ReadAll(stream)

	require.NoError(test, err)
	assert.Equal(test, output[1:3], actual)
}

func TestReadGA(test *testing.T) {
	output := []byte{telnet.IAC, telnet.GA}

	connection := NewMockConnection()
	connection.Output = output
	stream := telnet.NewStream(connection)

	expected := []byte{'\r', '\n'}
	actual, err := ioutil.ReadAll(stream)

	require.NoError(test, err)
	assert.Equal(test, expected, actual)
}

func TestNegotiationUnknown(test *testing.T) {
	connection := NewMockConnection()
	connection.Output = []byte{telnet.IAC, telnet.DO, 123, '\n'}
	stream := telnet.NewStream(connection)

	reader := bufio.NewReader(stream)
	output, err := reader.ReadString('\n')

	require.NoError(test, err)
	assert.Equal(test, "\n", output)
	assert.Equal(test, []byte{telnet.IAC, telnet.WONT, 123}, connection.Input)
}
