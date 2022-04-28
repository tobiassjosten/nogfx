package agmcp

import (
	"strconv"
	"strings"
)

func splitRank(str string) (string, *int) {
	parts := strings.SplitN(str, "(", 2)
	name := strings.Trim(parts[0], " ")

	var rank *int
	if len(parts) > 1 {
		r, err := strconv.Atoi(strings.Trim(parts[1], "%)"))
		if err == nil {
			rank = &r
		}
	}

	return name, rank
}
