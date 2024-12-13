package utils

import (
	"fmt"
	"os"
)

// EnsureDir checks if a directory exists and creates it if not
func EnsureDir(dir string) error {
	// Check if directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Attempt to create the directory
		if err = os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return nil
}
