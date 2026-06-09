package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
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

	t.Run("directory is not readable", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("permission bits are not reliable for this test on windows")
		}

		dir := filepath.Join(t.TempDir(), "restricted")
		if err := os.Mkdir(dir, 0o755); err != nil {
			t.Fatalf("Mkdir(%q) error = %v", dir, err)
		}
		if err := os.Chmod(dir, 0o000); err != nil {
			t.Fatalf("Chmod(%q) error = %v", dir, err)
		}
		t.Cleanup(func() {
			_ = os.Chmod(dir, 0o755)
		})

		err := ValidateDirectory(dir)
		if err == nil {
			t.Fatal("ValidateDirectory() error = nil, want error")
		}
		if !strings.Contains(err.Error(), "open directory") {
			t.Fatalf("ValidateDirectory() error = %q, want open directory context", err)
		}
	})
}
