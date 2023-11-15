package telnet

import (
	"bytes"
	"fmt"

	"golang.org/x/exp/slices"
)

// IsCommand determines whether the given sequence is a valid Telnet command.
func IsCommand(cmd []byte) bool {
	l := len(cmd)

	if l < 2 {
		return false
	}

	if bytes.Equal(cmd, []byte{IAC, IAC}) {
		return true
	}

	if bytes.Equal(cmd, []byte{IAC, GA}) {
		return true
	}

	if l == 3 && (cmd[1] == Will || cmd[1] == Wont || cmd[1] == Do || cmd[1] == Dont) {
		return true
	}

	if bytes.Equal(cmd[:2], []byte{IAC, SB}) && bytes.Equal(cmd[l-2:], []byte{IAC, SE}) {
		return true
	}

	return false
}

// @todo Implement support for RFC 885. Achaea sends IAC WILL 25[19 in hex] and
// obviously wants to negotiate something. What will they send if we accept
// that? Should we respond with the same IAC WILL 25?

// @todo Implement support for MCCP2.

func (nvt *NVT) negotiate(cmd []byte) error {
	switch cmd[1] {
	case Do:
		if nvt.options[ourside][cmd[2]].On() {
			return nil
		}

		if !slices.Contains(nvt.ourCoulds, cmd[2]) {
			_, err := nvt.Write([]byte{IAC, Wont, cmd[2]})
			if err != nil {
				return fmt.Errorf("failed to decline DO %d: %w", cmd[2], err)
			}

			return nil
		}

		_, err := nvt.Write([]byte{IAC, Will, cmd[2]})
		if err != nil {
			return fmt.Errorf("failed to accept DO %d: %w", cmd[2], err)
		}

		nvt.options[ourside][cmd[2]] = StateEnabled

	case Will:
		if nvt.options[theirside][cmd[2]].On() {
			return nil
		}

		if !slices.Contains(nvt.theirCoulds, cmd[2]) {
			_, err := nvt.Write([]byte{IAC, Dont, cmd[2]})
			if err != nil {
				return fmt.Errorf("failed to decline WILL %d: %w", cmd[2], err)
			}

			return nil
		}

		_, err := nvt.Write([]byte{IAC, Do, cmd[2]})
		if err != nil {
			return fmt.Errorf("failed to accept Will %d: %w", cmd[2], err)
		}

		nvt.options[theirside][cmd[2]] = StateEnabled

	case Dont:
		if nvt.options[ourside][cmd[2]].Off() {
			return nil
		}

		_, err := nvt.Write([]byte{IAC, Wont, cmd[2]})
		if err != nil {
			return fmt.Errorf("failed to accept DONT %d: %w", cmd[2], err)
		}

		nvt.options[ourside][cmd[2]] = StateDisabled

	case Wont:
		if nvt.options[theirside][cmd[2]].Off() {
			return nil
		}

		_, err := nvt.Write([]byte{IAC, Dont, cmd[2]})
		if err != nil {
			return fmt.Errorf("failed to accept WONT %d: %w", cmd[2], err)
		}

		nvt.options[theirside][cmd[2]] = StateDisabled
	}

	return nil
}
