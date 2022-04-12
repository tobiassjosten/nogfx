package gmcp

import (
	"strconv"
	"strings"
)

func splitRank(str string) (string, int) {
	parts := strings.SplitN(str, "(", 2)
	name := strings.Trim(parts[0], " ")
	rank, _ := strconv.Atoi(strings.Trim(parts[1], "%)"))

	return name, rank
}

func splitLevelRank(str string) (int, int) {
	name, rank := splitRank(str)
	level, _ := strconv.Atoi(name)

	return level, rank
}
