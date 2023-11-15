package telnet

import (
	"bytes"
	"fmt"
)

// SplitFunc looks looks for string termination based on negotiated options. By
// default, newline and GA is used, but the latter can be negotiated.
func (nvt *NVT) SplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	lastCR := false
	for i, b := range data {
		if lastCR && b == '\n' && nvt.options[theirside][SuppressGoAhead].On() {
			return i + 1, data[:i+1], nil
		}
		lastCR = b == '\r'

		if b == GA {
			return i + 1, data[:i+1], nil
		}
	}

	if atEOF {
		return len(data), data, nil
	}

	// @todo Test this.
	return 0, nil, nil
}

// Read parses and returns data received from the server.
func (nvt *NVT) Read(buffer []byte) (count int, err error) {
	l := len(buffer)
	if l == 0 {
		return 0, nil
	}

	for lastCR := false; count < l; {
		b, err := nvt.buffer.ReadByte()
		if err != nil {
			return count, err
		}

		if b != IAC && len(nvt.cmdBuffer) == 0 {
			buffer[count] = b
			count++

			if lastCR && b == '\n' && nvt.options[theirside][SuppressGoAhead].On() {
				return count, nil
			}
			lastCR = b == '\r'

			continue
		}

		nvt.cmdBuffer = append(nvt.cmdBuffer, b)

		if !IsCommand(nvt.cmdBuffer) {
			continue
		}

		if bytes.Equal(nvt.cmdBuffer, []byte{IAC, IAC}) {
			nvt.cmdBuffer = []byte{}

			continue
		}

		if bytes.Equal(nvt.cmdBuffer, []byte{IAC, GA}) {
			nvt.cmdBuffer = []byte{}

			buffer[count] = GA
			count++

			return count, nil
		}

		err = nvt.negotiate(nvt.cmdBuffer)
		if err != nil {
			return count, fmt.Errorf("failed negotiation: %w", err)
		}

		if nvt.CommandFunc != nil {
			err := nvt.CommandFunc(nvt.cmdBuffer, nvt)
			if err != nil {
				fmt.Println(">", err)
				return count, fmt.Errorf("failed command processing: %w", err)
			}
		}

		nvt.cmdBuffer = []byte{}
	}

	return count, nil
}

// Write sends data to the server.
func (nvt *NVT) Write(data []byte) (int, error) {
	// Telnet specifies <CR><LF> endings, so we make sure we adhere.
	if ld := len(data); len(data) == 0 || data[0] != IAC {
		if ld > 2 && data[ld-2] == '\r' && data[ld-1] == '\n' {
			data = data[0 : ld-2]
		}
		data = append(data, '\r', '\n')
	}

	for i := 0; i < len(data); i++ {
		if data[i] != IAC {
			continue
		}

		// @todo Test this.
		if data[i+1] == SB {
			ii := bytes.IndexByte(data[i+1:], IAC)
			if ii < 0 {
				break
			}
		}

		switch data[i+1] {
		case Do:
			nvt.options[theirside][data[i+2]] = StateEnabling
		case Will:
			nvt.options[ourside][data[i+2]] = StateEnabling
		}

		if bytes.Contains([]byte{Do, Dont, Will, Wont}, []byte{data[i+1]}) {
			i += 2
		}
	}

	// @todo Pick up commands and mutate nvt.ourCoulds and nvt.theirCoulds.
	// @todo Potentially add Do(), Dont(), Will(), Wont() methods.

	return nvt.Conn.Write(data)
}
