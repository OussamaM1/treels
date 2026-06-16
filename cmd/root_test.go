package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExecuteCommandErrorPath(t *testing.T) {
	cmd := newRootCmd()
	cmd.SetArgs([]string{"one", "two"})

	var output bytes.Buffer
	var exitCode int
	execute(cmd, &output, func(code int) {
		exitCode = code
	})

	if exitCode != 1 {
		t.Fatalf("execute() exit code = %d, want 1", exitCode)
	}
	if !strings.Contains(output.String(), "accepts at most 1 arg") {
		t.Fatalf("execute() error output = %q, want too many args error", output.String())
	}
}

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

			if got := output.String(); got != "treels v1.4.0\n" {
				t.Fatalf("Execute() output = %q, want version output", got)
			}
		})
	}
}

func TestExecute_VersionFlag(t *testing.T) {
	originalArgs := os.Args
	os.Args = []string{"treels", "--version"}
	defer func() {
		os.Args = originalArgs
	}()

	output := captureStdout(t, func() {
		Execute()
	})

	if got := output; got != "treels v1.4.0\n" {
		t.Fatalf("Execute() output = %q, want version output", got)
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

func TestRootCmd_LongFlag(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "long flag", args: []string{"--long"}},
		{name: "short flag", args: []string{"-l"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main"), 0o644); err != nil {
				t.Fatalf("WriteFile() error = %v", err)
			}

			output := captureStdout(t, func() {
				cmd := newRootCmd()
				args := append(tt.args, "--no-icons", dir)
				cmd.SetArgs(args)

				if err := cmd.Execute(); err != nil {
					t.Fatalf("Execute() error = %v, want nil", err)
				}
			})

			for _, want := range []string{"-rw", "main.go", "12", "0 directories, 1 files"} {
				if !strings.Contains(output, want) {
					t.Fatalf("Execute() output = %q, want to contain %q", output, want)
				}
			}
			if strings.Contains(output, "main.go (12 B)") {
				t.Fatalf("Execute() output = %q, want size in metadata column only", output)
			}
		})
	}
}

func TestRootCmd_IncludeExcludeFlags(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("readme"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "debug.log"), []byte("debug"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	output := captureStdout(t, func() {
		cmd := newRootCmd()
		cmd.SetArgs([]string{"--include", "*.go", "--include", "*.md", "--exclude", "debug.log", "--no-icons", dir})

		if err := cmd.Execute(); err != nil {
			t.Fatalf("Execute() error = %v, want nil", err)
		}
	})

	for _, want := range []string{"main.go", "README.md", "0 directories, 2 files"} {
		if !strings.Contains(output, want) {
			t.Fatalf("Execute() output = %q, want to contain %q", output, want)
		}
	}
	if strings.Contains(output, "debug.log") {
		t.Fatalf("Execute() output = %q, want debug.log excluded", output)
	}
}

func TestRootCmd_SortFlags(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "small.txt"), []byte("x"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "large.txt"), []byte("1234567890"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if err := os.Mkdir(filepath.Join(dir, "folder"), 0o755); err != nil {
		t.Fatalf("Mkdir() error = %v", err)
	}

	output := captureStdout(t, func() {
		cmd := newRootCmd()
		cmd.SetArgs([]string{"--long", "--no-icons", "--sort", "size", "--reverse", "--dirs-first", dir})

		if err := cmd.Execute(); err != nil {
			t.Fatalf("Execute() error = %v, want nil", err)
		}
	})

	assertOutputOrder(t, output, []string{"folder", "large.txt", "small.txt"})
}

func TestRootCmd_SortFlagRejectsInvalidValue(t *testing.T) {
	cmd := newRootCmd()
	cmd.SetArgs([]string{"--sort", "unknown"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want invalid sort error")
	}
	if !strings.Contains(err.Error(), "--sort must be one of: name, size, modified, type") {
		t.Fatalf("Execute() error = %q, want invalid sort validation error", err)
	}
}

func TestRootCmd_SortFlagIsCaseInsensitive(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("x"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	output := captureStdout(t, func() {
		cmd := newRootCmd()
		cmd.SetArgs([]string{"--sort", "SIZE", "--no-icons", dir})

		if err := cmd.Execute(); err != nil {
			t.Fatalf("Execute() error = %v, want nil", err)
		}
	})
	if !strings.Contains(output, "a.txt") {
		t.Fatalf("Execute() output = %q, want sorted output", output)
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

func TestRootCmd_DirsOnlyFlag(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "cmd"), 0o755); err != nil {
		t.Fatalf("Mkdir() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	output := captureStdout(t, func() {
		cmd := newRootCmd()
		cmd.SetArgs([]string{"--dirs-only", "--no-icons", dir})

		if err := cmd.Execute(); err != nil {
			t.Fatalf("Execute() error = %v, want nil", err)
		}
	})

	if !strings.Contains(output, "cmd") {
		t.Fatalf("Execute() output = %q, want directory", output)
	}
	if strings.Contains(output, "main.go") {
		t.Fatalf("Execute() output = %q, want file omitted", output)
	}
	if !strings.Contains(output, "1 directories") {
		t.Fatalf("Execute() output = %q, want directory-only summary", output)
	}
	if strings.Contains(output, "0 files") {
		t.Fatalf("Execute() output = %q, want no file count in dirs-only summary", output)
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

func TestRootCmd_JSONFlag(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	output := captureStdout(t, func() {
		cmd := newRootCmd()
		cmd.SetArgs([]string{"--json", dir})

		if err := cmd.Execute(); err != nil {
			t.Fatalf("Execute() error = %v, want nil", err)
		}
	})

	var got struct {
		Tree    bool `json:"tree"`
		Summary struct {
			Files int `json:"files"`
		} `json:"summary"`
		Entries []struct {
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"entries"`
	}
	if err := json.Unmarshal([]byte(output), &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v, output = %q", err, output)
	}
	if got.Tree {
		t.Fatal("json tree = true, want false")
	}
	if got.Summary.Files != 1 {
		t.Fatalf("json summary files = %d, want 1", got.Summary.Files)
	}
	if len(got.Entries) != 1 || got.Entries[0].Name != "main.go" || got.Entries[0].Type != "file" {
		t.Fatalf("json entries = %+v, want main.go file", got.Entries)
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

func assertOutputOrder(t *testing.T, output string, orderedNames []string) {
	t.Helper()
	previousIndex := -1
	for _, name := range orderedNames {
		index := strings.Index(output, name)
		if index == -1 {
			t.Fatalf("output = %q, want to contain %q", output, name)
		}
		if index <= previousIndex {
			t.Fatalf("output = %q, want %q after previous entries %v", output, name, orderedNames)
		}
		previousIndex = index
	}
}
