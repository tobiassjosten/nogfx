package telnet

import (
	"bytes"
)

// Convenience constants to make telnet commands more readable.
const (
	ECHO  byte = 1
	TTYPE byte = 24
	MCCP2 byte = 86
	ATCP  byte = 200
	GMCP  byte = 201
	SE    byte = 240
	GA    byte = 249
	SB    byte = 250
	WILL  byte = 251
	WONT  byte = 252
	DO    byte = 253
	DONT  byte = 254
	IAC   byte = 255
)

var (
	acceptWill = map[byte]struct{}{
		ECHO: {},
		GMCP: {},
	}
)

// Will sends the IAC WILL <CMD> sequence.
func (client *Client) Will(command byte) error {
	_, err := client.data.Write([]byte{IAC, WILL, command})
	return err
}

// Wont sends the IAC WONT <CMD> sequence.
func (client *Client) Wont(command byte) error {
	_, err := client.data.Write([]byte{IAC, WONT, command})
	return err
}

// Do sends the IAC DO <CMD> sequence.
func (client *Client) Do(command byte) error {
	_, err := client.data.Write([]byte{IAC, DO, command})
	return err
}

// Dont sends the IAC DONT <CMD> sequence.
func (client *Client) Dont(command byte) error {
	_, err := client.data.Write([]byte{IAC, DONT, command})
	return err
}

// Subneg sends the IAC SB <CMD> 0/1 <DATA> IAC SE sequence.
func (client *Client) Subneg(b byte, value []byte) error {
	var v byte
	if len(value) > 0 {
		v = 1
	}

	_, err := client.data.Write(append(append(
		[]byte{IAC, SB, b, v},
		value...,
	), IAC, SE))

	return err
}

func (*Client) processCommand(command []byte) (bool, [][]byte) {
	var responses [][]byte

	if len(command) < 3 {
		return false, nil
	}

	switch command[1] {
	case WILL:
		if _, ok := acceptWill[command[2]]; ok {
			responses = append(responses, []byte{IAC, DO, command[2]})
			return true, responses
		}

		responses = append(responses, []byte{IAC, DONT, command[2]})

		return true, responses

	case WONT:
		responses = append(responses, []byte{IAC, DONT, command[2]})
		return true, responses

	case DO:
		responses = append(responses, []byte{IAC, WONT, command[2]})
		return true, responses

	case DONT:
		responses = append(responses, []byte{IAC, WONT, command[2]})
		return true, responses
	}

	if len(command) < 5 {
		return false, nil
	}

	if bytes.Equal(command[:2], []byte{IAC, SB}) {
		if bytes.Equal(command[len(command)-2:], []byte{IAC, SE}) {
			return true, nil
		}
	}

	return false, nil
}
