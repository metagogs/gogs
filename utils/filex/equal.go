package filex

import (
	"bytes"
	"os"
	"strings"
)

func IsFileEqual(a, b string) bool {
	f1, err := os.ReadFile(a)
	if err != nil {
		return false
	}

	f2, err := os.ReadFile(b)
	if err != nil {
		return false
	}

	normF1 := strings.ReplaceAll(string(f1), "\r\n", "\n")
	normF2 := strings.ReplaceAll(string(f2), "\r\n", "\n")

	return bytes.Equal([]byte(normF1), []byte(normF2))
}
