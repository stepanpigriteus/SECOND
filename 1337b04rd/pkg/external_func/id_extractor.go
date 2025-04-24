package externalfunc

import (
	"fmt"
	"strconv"
	"strings"
)

func ExtractIDFromPath(path string) (int, error) {
	segments := strings.Split(path, "/")
	if len(segments) < 3 {
		return 0, fmt.Errorf("invalid path")
	}

	idStr := segments[len(segments)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid id")
	}

	return id, nil
}
