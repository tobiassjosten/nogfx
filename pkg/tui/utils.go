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
