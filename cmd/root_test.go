package cmd

import (
	"bytes"
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

func TestRootCmd_IconFlagsRemoved(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "long flag", args: []string{"--icon"}, want: "unknown flag: --icon"},
		{name: "short flag", args: []string{"-i"}, want: "unknown shorthand flag: 'i'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := newRootCmd()
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			if err == nil {
				t.Fatal("Execute() error = nil, want error")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("Execute() error = %q, want to contain %q", err, tt.want)
			}
		})
	}
}

func TestRootCmd_VersionFlag(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "long flag", args: []string{"--version"}},
		{name: "short flag", args: []string{"-v"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output bytes.Buffer
			cmd := newRootCmd()
			cmd.SetOut(&output)
			cmd.SetArgs(tt.args)

			if err := cmd.Execute(); err != nil {
				t.Fatalf("Execute() error = %v, want nil", err)
			}

			if got := output.String(); got != "treels v1.3.1\n" {
				t.Fatalf("Execute() output = %q, want version output", got)
			}
		})
	}
}

func TestRootCmd_ValidPathWithFlags(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	output := captureStdout(t, func() {
		cmd := newRootCmd()
		cmd.SetArgs([]string{"--tree", "--no-icons", dir})

		if err := cmd.Execute(); err != nil {
			t.Fatalf("Execute() error = %v, want nil", err)
		}
	})

	if !strings.Contains(output, "1 files") {
		t.Fatalf("Execute() output = %q, want file count", output)
	}
}

func TestRootCmd_ReadableFlag(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	output := captureStdout(t, func() {
		cmd := newRootCmd()
		cmd.SetArgs([]string{"-r", "--no-icons", dir})

		if err := cmd.Execute(); err != nil {
			t.Fatalf("Execute() error = %v, want nil", err)
		}
	})

	if !strings.Contains(output, "main.go (12 B)") {
		t.Fatalf("Execute() output = %q, want readable file size", output)
	}
}

func TestRootCmd_DepthFlag(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "subpkg"), 0o755); err != nil {
		t.Fatalf("Mkdir() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "subpkg", "nested.go"), []byte("package nested"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	output := captureStdout(t, func() {
		cmd := newRootCmd()
		cmd.SetArgs([]string{"--tree", "--depth", "1", "--no-icons", dir})

		if err := cmd.Execute(); err != nil {
			t.Fatalf("Execute() error = %v, want nil", err)
		}
	})

	if !strings.Contains(output, "subpkg") || !strings.Contains(output, "main.go") {
		t.Fatalf("Execute() output = %q, want direct children", output)
	}
	if strings.Contains(output, "nested.go") {
		t.Fatalf("Execute() output = %q, want nested file omitted", output)
	}
}

func TestRootCmd_DepthFlagRejectsNegativeValue(t *testing.T) {
	cmd := newRootCmd()
	cmd.SetArgs([]string{"--tree", "--depth", "-1"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "--depth must be greater than or equal to 0") {
		t.Fatalf("Execute() error = %q, want depth validation error", err)
	}
}

func TestRootCmd_NoSummaryFlag(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	output := captureStdout(t, func() {
		cmd := newRootCmd()
		cmd.SetArgs([]string{"--no-summary", "--no-icons", dir})

		if err := cmd.Execute(); err != nil {
			t.Fatalf("Execute() error = %v, want nil", err)
		}
	})

	if !strings.Contains(output, "main.go") {
		t.Fatalf("Execute() output = %q, want visible file", output)
	}
	if strings.Contains(output, "directories,") {
		t.Fatalf("Execute() output = %q, want no summary", output)
	}
}

func TestRootCmd_GitIgnoreFlag(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".gitignore"), []byte("ignored.txt\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "ignored.txt"), []byte("ignored"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	output := captureStdout(t, func() {
		cmd := newRootCmd()
		cmd.SetArgs([]string{"--gitignore", "--no-icons", dir})

		if err := cmd.Execute(); err != nil {
			t.Fatalf("Execute() error = %v, want nil", err)
		}
	})

	if strings.Contains(output, "ignored.txt") {
		t.Fatalf("Execute() output = %q, want ignored file to be omitted", output)
	}
	if !strings.Contains(output, "main.go") {
		t.Fatalf("Execute() output = %q, want visible file", output)
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
