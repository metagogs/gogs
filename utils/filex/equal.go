package filex

import (
	"bytes"
	"os"
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

	return bytes.Equal(f1, f2)
}
