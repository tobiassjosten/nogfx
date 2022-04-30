package simpex

import (
	"bytes"
)

// Match a text against a pattern to see if it matches. If it does, captured
// matches are returned. If it doesn't, nil is returned.
func Match(pattern []byte, text []byte) [][]byte {
	captures := [][]byte{}

	var capture []byte

	tick := func() {
		if capture != nil {
			capture = append(capture, text[0])
		}

		pattern = pattern[1:]
		text = text[1:]
	}

main:
	for len(pattern) > 0 && len(text) > 0 {
		switch pattern[0] {
		case '{':
			if indexNot(pattern, pattern[0])%2 == 0 {
				if pattern[0] != text[0] {
					return nil
				}

				pattern = pattern[1:]
				tick()
				continue
			}

			capture = []byte{}
			pattern = pattern[1:]

			continue main

		case '}':
			if len(pattern) > 1 && pattern[0] == pattern[1] {
				if pattern[0] != text[0] {
					return nil
				}

				pattern = pattern[1:]
				tick()
				continue
			}

			captures = append(captures, capture)
			capture = nil
			pattern = pattern[1:]

			continue main

		case '^':
			if len(pattern) > 1 && pattern[0] == pattern[1] {
				if pattern[0] != text[0] {
					return nil
				}

				pattern = pattern[1:]
				tick()
				continue
			}

			edge := indexNonAlphanum(text)
			if edge < 1 {
				return nil
			}

			if capture != nil {
				capture = append(capture, text[:edge]...)
			}

			pattern = pattern[1:]
			text = text[edge:]

			continue main

		case '*':
			if len(pattern) > 1 && pattern[0] == pattern[1] {
				if pattern[0] != text[0] {
					return nil
				}

				pattern = pattern[1:]
				tick()
				continue
			}

			if len(pattern) == 1 {
				if capture != nil {
					return nil
				}
				break main
			}

			if capture != nil && len(pattern) == 2 && pattern[1] == '}' {
				capture = append(capture, text...)
				captures = append(captures, capture)

				break main
			}

			start := indexNonSpecial(pattern)

			end := indexSpecial(pattern[start:]) - 1
			if end < 0 {
				end = len(pattern[start:]) - 1
			}

			segment := pattern[start : start+end+1]

			edge := bytes.Index(text, segment)
			if edge < 0 {
				return nil
			}

			if capture != nil {
				capture = append(capture, text[:edge]...)
			}

			pattern = pattern[1:]
			text = text[edge:]

		case '?':
			if len(pattern) > 1 && pattern[0] == pattern[1] {
				if pattern[0] != text[0] {
					return nil
				}

				pattern = pattern[1:]
				tick()
				continue
			}

			pattern[0] = text[0]

			tick()

		default:
			if pattern[0] != text[0] {
				return nil
			}

			tick()
		}
	}

	return captures
}

func isAlphanum(b byte) bool {
	return (b >= '0' && b <= '9') ||
		(b >= 'A' && b <= 'Z') ||
		(b >= 'a' && b <= 'z')
}

func indexNonAlphanum(text []byte) int {
	for i, b := range text {
		if !isAlphanum(b) {
			return i
		}
	}

	return -1
}

func indexSpecial(text []byte) int {
	if len(text) == 0 {
		return -1
	}

main:
	for i := 0; i < len(text); {
		switch text[i] {
		case '{', '}', '?', '^', '*':
			if in := indexNot(text[i:], text[i]); in%2 == 0 {
				i += in
				continue main
			}
			return i
		}
		i++
	}

	return -1
}

func indexNonSpecial(text []byte) int {
	var previous byte
	for i, b := range text {
		switch b {
		case '{', '}', '?', '^', '*':
			if previous == b {
				return i
			}
			previous = b
		default:
			return i
		}
	}

	return -1
}

func indexNot(text []byte, c byte) int {
	for i, b := range text {
		if b != c {
			return i
		}
	}

	return -1
}
