package utils

import (
	"os"
)

func DirExists(d string) bool {
	if _, err := os.Stat(d); os.IsNotExist(err) {
		return false
	}

	return true
}
