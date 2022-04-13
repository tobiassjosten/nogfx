package telnet

import (
	"bytes"
	"fmt"
)

const (
	ECHO  byte = 1
	LF    byte = 10
	CR    byte = 13
	TTYPE byte = 24
	MCCP  byte = 85
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
	acceptDo = map[byte]struct{}{}
)

func (client *Client) Will(command byte) error {
	_, err := client.data.Write([]byte{IAC, WILL, command})
	return err
}

func (client *Client) Wont(command byte) error {
	_, err := client.data.Write([]byte{IAC, WONT, command})
	return err
}

func (client *Client) Do(command byte) error {
	_, err := client.data.Write([]byte{IAC, DO, command})
	return err
}

func (client *Client) Dont(command byte) error {
	_, err := client.data.Write([]byte{IAC, DONT, command})
	return err
}

func (client *Client) Subneg(b byte, value []byte) error {
	var v byte = 0
	if len(value) > 0 {
		v = 1
	}

	_, err := client.data.Write(append(append(
		[]byte{IAC, SB, b, v},
		value...,
	), IAC, SE))
	return err
}

func (client *Client) processCommand(command []byte) ([]byte, byte, error) {
	if bytes.Equal(command, []byte{IAC, IAC}) {
		return []byte{}, IAC, nil
	}

	if bytes.Equal(command, []byte{IAC, GA}) {
		return []byte{}, GA, nil
	}

	if len(command) < 3 {
		return command, 0, nil
	}

	switch command[1] {
	case WILL:
		if _, ok := acceptWill[command[2]]; ok {
			if err := client.Do(command[2]); err != nil {
				return nil, 0, fmt.Errorf(
					"failed accepting WILL %d: %s",
					command[2], err,
				)
			}
		} else if err := client.Dont(command[2]); err != nil {
			return nil, 0, fmt.Errorf(
				"failed rejecting WILL %d: %s",
				command[2], err,
			)
		}

		client.commands <- command
		return []byte{}, 0, nil

	case WONT:
		if err := client.Dont(command[2]); err != nil {
			return nil, 0, fmt.Errorf(
				"failed rejecting WONT %d: %s",
				command[2], err,
			)
		}

		client.commands <- command
		return []byte{}, 0, nil

	case DO:
		if _, ok := acceptDo[command[2]]; ok {
			if err := client.Will(command[2]); err != nil {
				return nil, 0, fmt.Errorf(
					"failed accepting DO %d: %s",
					command[2], err,
				)
			}
		} else if err := client.Wont(command[2]); err != nil {
			return nil, 0, fmt.Errorf(
				"failed rejecting DO %d: %s",
				command[2], err,
			)
		}

		client.commands <- command
		return []byte{}, 0, nil

	case DONT:
		if err := client.Wont(command[2]); err != nil {
			return nil, 0, fmt.Errorf(
				"failed rejecting DONT %d: %s",
				command[2], err,
			)
		}

		client.commands <- command
		return []byte{}, 0, nil

	default: // Noop.
	}

	if len(command) < 5 {
		return command, 0, nil
	}

	if bytes.Equal(command[:2], []byte{IAC, SB}) {
		if bytes.Equal(command[len(command)-2:], []byte{IAC, SE}) {
			client.commands <- command
			return []byte{}, 0, nil
		}
	}

	return command, 0, nil
}
