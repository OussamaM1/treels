// Package utils - utils/validator.go
package utils

// ValidateDirectoryArgs func - Validates command-line arguments for directory operations.
func ValidateDirectoryArgs(args []string) bool {
	return len(args) == 1
}
