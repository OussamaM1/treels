package cmd

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRootCmd_InvalidPath(t *testing.T) {
	cmd := newRootCmd()
	cmd.SetArgs([]string{filepath.Join(t.TempDir(), "missing")})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want error")
	}

	if !strings.Contains(err.Error(), "validate directory") {
		t.Fatalf("Execute() error = %q, want validation context", err)
	}
}

func TestRootCmd_TooManyArgs(t *testing.T) {
	cmd := newRootCmd()
	cmd.SetArgs([]string{"one", "two"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want error")
	}
}

func TestRootCmd_ValidPathWithFlags(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	output := captureStdout(t, func() {
		cmd := newRootCmd()
		cmd.SetArgs([]string{"--tree", "--icon", dir})

		if err := cmd.Execute(); err != nil {
			t.Fatalf("Execute() error = %v, want nil", err)
		}
	})

	if !strings.Contains(output, "1 files") {
		t.Fatalf("Execute() output = %q, want file count", output)
	}
}

func captureStdout(t *testing.T, run func()) string {
	t.Helper()

	originalStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("Pipe() error = %v", err)
	}

	os.Stdout = writer
	defer func() {
		os.Stdout = originalStdout
	}()

	run()

	if err := writer.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	output, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}
	if err := reader.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	return string(output)
}
