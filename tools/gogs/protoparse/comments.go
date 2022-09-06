package protoparse

import "strings"

func CommentsContains(lines []string, target string) bool {
	for _, line := range lines {
		if strings.Contains(line, target) {
			return true
		}
	}

	return false
}
