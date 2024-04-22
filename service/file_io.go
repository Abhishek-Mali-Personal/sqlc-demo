package service

import (
	"os"
	"path/filepath"
)

// FetchAbsolutePath returns an absolute path if file or directory exists else returns error
func FetchAbsolutePath(relativePath string) (string, error) {
	// GET THE ABSOLUTE PATH
	absolutePath, err := filepath.Abs(relativePath)
	if err != nil {
		return "", err
	}

	// CHECK IF THE FILE EXISTS
	if _, err := os.Stat(absolutePath); os.IsNotExist(err) {
		return "", err
	}
	return absolutePath, nil
}
