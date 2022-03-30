package telnet

import "bytes"

func isCompleteCommand(command []byte) bool {
	if len(command) < 3 {
		return false
	}

	if command[1] == SB && !bytes.Equal(command[len(command)-2:], []byte{IAC, SE}) {
		return false
	}

	return true
}

func isValidCommand(command []byte) bool {
	if command[0] != IAC {
		return false
	}

	negotiators := []byte{WILL, WONT, DO, DONT}
	for _, negotiator := range negotiators {
		if command[1] == negotiator && len(command) != 3 {
			return false
		}
	}

	return true
}
