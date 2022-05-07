package tui

import "math"

func max(a, b int) int {
	if b > a {
		return b
	}
	return a
}

func min(a, b int) int {
	if b < a {
		return b
	}
	return a
}

func abs(a int) int {
	return int(math.Abs(float64(a)))
}

func rel(a int, b int) int {
	return abs(a) * max(-1, min(1, b))
}
