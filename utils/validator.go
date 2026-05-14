// Package utils - utils/validator.go
package utils

import (
	"fmt"
	"os"
)

// ValidateDirectory func - validates that a path exists, is a directory, and is readable.
func ValidateDirectory(directory string) error {
	info, err := os.Stat(directory)
	if err != nil {
		return fmt.Errorf("validate directory %q: %w", directory, err)
	}

	if !info.IsDir() {
		return fmt.Errorf("%q is not a directory", directory)
	}

	dir, err := os.Open(directory)
	if err != nil {
		return fmt.Errorf("open directory %q: %w", directory, err)
	}

	if err := dir.Close(); err != nil {
		return fmt.Errorf("close directory %q: %w", directory, err)
	}

	return nil
}
