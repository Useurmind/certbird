package utils

import (
	"os"
)

// DoesFileExist checks whether the given file path exists and returns true if so, else false.
func DoesFileExist(filePath string) bool {
	if _, err := os.Stat(filePath); err == nil {
		return true;
	} else {
		return false;
	}
}