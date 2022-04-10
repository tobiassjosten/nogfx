package telnet

import (
	"bytes"
	"log"
)

func (client *Client) AcceptWill(command byte) {
	client.acceptWill[command] = struct{}{}
}

func (client *Client) AcceptDo(command byte) {
	client.acceptDo[command] = struct{}{}
}

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

func (client *Client) processCommand(command []byte) ([]byte, byte) {
	log.Printf("CMD: '%s': %s", command, string(command))
	if bytes.Equal(command, []byte{IAC, IAC}) {
		return []byte{}, IAC
	}

	// More reliable than newline to mark the end of a message, so we relay
	// it upstream for processing in the game logic.
	if bytes.Equal(command, []byte{IAC, GA}) {
		return []byte{}, GA
	}

	if len(command) < 3 {
		return command, 0
	}

	switch command[1] {
	case WILL:
		if _, ok := client.acceptWill[command[2]]; ok {
			if err := client.Do(GMCP); err != nil {
				log.Printf(
					"failed accepting WILL %d: %s",
					command[2], err,
				)
			}
		} else if err := client.Dont(command[2]); err != nil {
			log.Printf(
				"failed rejecting WILL %d: %s",
				command[2], err,
			)
		}

		client.commands <- command
		return []byte{}, 0

	case WONT:
		client.commands <- command
		return []byte{}, 0

	case DO:
		if _, ok := client.acceptDo[command[2]]; ok {
			if err := client.Will(GMCP); err != nil {
				log.Printf(
					"failed accepting DO %d: %s",
					command[2], err,
				)
			}
		} else if err := client.Wont(command[2]); err != nil {
			log.Printf(
				"failed rejecting DO %d: %s",
				command[2], err,
			)
		}

		client.commands <- command
		return []byte{}, 0

	case DONT:
		client.commands <- command
		return []byte{}, 0

	default: // Noop.
	}

	if len(command) < 5 {
		return command, 0
	}

	if bytes.Equal(command[:2], []byte{IAC, SB}) {
		if bytes.Equal(command[len(command)-2:], []byte{IAC, SE}) {
			client.commands <- command
			return []byte{}, 0
		}
	}

	return command, 0
}