package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateDirectory(t *testing.T) {
	t.Run("existing directory", func(t *testing.T) {
		dir := t.TempDir()

		if err := ValidateDirectory(dir); err != nil {
			t.Fatalf("ValidateDirectory() error = %v, want nil", err)
		}
	})

	t.Run("missing directory", func(t *testing.T) {
		missing := filepath.Join(t.TempDir(), "missing")

		if err := ValidateDirectory(missing); err == nil {
			t.Fatal("ValidateDirectory() error = nil, want error")
		}
	})

	t.Run("path is a file", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "file.txt")
		if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}

		err := ValidateDirectory(path)
		if err == nil {
			t.Fatal("ValidateDirectory() error = nil, want error")
		}
	})
}
