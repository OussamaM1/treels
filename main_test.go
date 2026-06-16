package main

import (
	"io"
	"os"
	"testing"
)

func TestMainVersionFlag(t *testing.T) {
	originalArgs := os.Args
	os.Args = []string{"treels", "--version"}
	defer func() {
		os.Args = originalArgs
	}()

	output := captureStdout(t, func() {
		main()
	})

	if output != "treels v1.4.0\n" {
		t.Fatalf("main() output = %q, want version output", output)
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
