package service

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/oussamaM1/treels/module"
)

func TestDispatcher_MissingDirectory(t *testing.T) {
	var output bytes.Buffer
	err := dispatcher(module.Options{Directory: filepath.Join(t.TempDir(), "missing")}, &output)
	if err == nil {
		t.Fatal("Dispatcher() error = nil, want error")
	}

	if output.Len() != 0 {
		t.Fatalf("dispatcher() output = %q, want no output on validation error", output.String())
	}
}

func TestDispatcher_ListDirectory(t *testing.T) {
	dir := t.TempDir()
	mustWriteFile(t, filepath.Join(dir, "alpha.go"), "package main")
	mustWriteFile(t, filepath.Join(dir, ".hidden"), "secret")
	mustMkdir(t, filepath.Join(dir, "subpkg"))

	tests := []struct {
		name         string
		flags        module.Flags
		wantContains []string
		wantMissing  []string
	}{
		{
			name:         "excludes hidden files by default",
			flags:        module.Flags{HideIcon: true},
			wantContains: []string{".", "alpha.go", "subpkg", "1 directories, 1 files"},
			wantMissing:  []string{".hidden"},
		},
		{
			name:         "includes hidden files with all flag",
			flags:        module.Flags{HideIcon: true, ShowHidden: true},
			wantContains: []string{".", ".hidden", "alpha.go", "subpkg", "1 directories, 2 files"},
		},
		{
			name:         "shows readable file and directory sizes",
			flags:        module.Flags{HideIcon: true, ShowReadableSize: true},
			wantContains: []string{"alpha.go (12 B)", "subpkg (", "1 directories, 1 files"},
			wantMissing:  []string{".hidden"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output bytes.Buffer
			err := dispatcher(module.Options{Directory: dir, Flags: tt.flags}, &output)
			if err != nil {
				t.Fatalf("dispatcher() error = %v, want nil", err)
			}

			got := stripANSI(output.String())
			for _, want := range tt.wantContains {
				if !strings.Contains(got, want) {
					t.Fatalf("dispatcher() output = %q, want to contain %q", got, want)
				}
			}
			for _, missing := range tt.wantMissing {
				if strings.Contains(got, missing) {
					t.Fatalf("dispatcher() output = %q, want not to contain %q", got, missing)
				}
			}
		})
	}
}

func TestDispatcher_TreeDirectory(t *testing.T) {
	dir := t.TempDir()
	mustMkdir(t, filepath.Join(dir, "subpkg"))
	mustWriteFile(t, filepath.Join(dir, "subpkg", "nested.go"), "package nested")
	mustWriteFile(t, filepath.Join(dir, "main.go"), "package main")

	var output bytes.Buffer
	err := dispatcher(module.Options{
		Directory: dir,
		Flags: module.Flags{
			HideIcon:         true,
			ShowTreeView:     true,
			ShowReadableSize: true,
		},
	}, &output)
	if err != nil {
		t.Fatalf("dispatcher() error = %v, want nil", err)
	}

	got := stripANSI(output.String())
	for _, want := range []string{"└── ", "subpkg (", "nested.go (14 B)", "1 directories, 2 files"} {
		if !strings.Contains(got, want) {
			t.Fatalf("dispatcher() output = %q, want to contain %q", got, want)
		}
	}
}

func TestHumanReadableSize(t *testing.T) {
	tests := []struct {
		name string
		size int64
		want string
	}{
		{name: "zero bytes", size: 0, want: "0 B"},
		{name: "bytes", size: 42, want: "42 B"},
		{name: "kilobytes", size: 2048, want: "2.0 KB"},
		{name: "megabytes", size: 1536 * 1024, want: "1.5 MB"},
		{name: "negative size", size: -1, want: "0 B"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := humanReadableSize(tt.size); got != tt.want {
				t.Fatalf("humanReadableSize(%d) = %q, want %q", tt.size, got, tt.want)
			}
		})
	}
}

func TestReadDirectory_FilePath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "regular.txt")
	mustWriteFile(t, path, "content")

	files, dir, err := readDirectory(path)
	if err == nil {
		t.Fatal("readDirectory() error = nil, want error")
	}
	if files != nil {
		t.Fatalf("readDirectory() files = %v, want nil", files)
	}
	if dir != nil {
		t.Fatalf("readDirectory() directory = %v, want nil after read failure", dir)
	}
}

func TestIsHidden(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{name: "dot file", input: ".gitignore", want: true},
		{name: "regular file", input: "main.go", want: false},
		{name: "empty", input: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isHidden(tt.input); got != tt.want {
				t.Fatalf("isHidden(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func mustWriteFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", path, err)
	}
}

func mustMkdir(t *testing.T, path string) {
	t.Helper()
	if err := os.Mkdir(path, 0o755); err != nil {
		t.Fatalf("Mkdir(%q) error = %v", path, err)
	}
}
