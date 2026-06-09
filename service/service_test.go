package service

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

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

func TestDispatcher_PublicWrapper(t *testing.T) {
	dir := t.TempDir()
	mustWriteFile(t, filepath.Join(dir, "main.go"), "package main")

	if err := Dispatcher(module.Options{Directory: dir, Flags: module.Flags{HideIcon: true, HideSummary: true}}); err != nil {
		t.Fatalf("Dispatcher() error = %v, want nil", err)
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

func TestListAndTreeDirectory_ReadErrors(t *testing.T) {
	path := filepath.Join(t.TempDir(), "regular.txt")
	mustWriteFile(t, path, "content")

	var output bytes.Buffer
	options := directoryOptions{Options: module.Options{Directory: path}}

	if _, _, err := listDirectory(options, &output); err == nil {
		t.Fatal("listDirectory() error = nil, want error")
	}
	if _, _, err := treeDirectory(options, &output, "", true, 0); err == nil {
		t.Fatal("treeDirectory() error = nil, want error")
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

func TestDispatcher_TreeDirectoryDepth(t *testing.T) {
	dir := t.TempDir()
	mustMkdir(t, filepath.Join(dir, "subpkg"))
	mustWriteFile(t, filepath.Join(dir, "subpkg", "nested.go"), "package nested")
	mustWriteFile(t, filepath.Join(dir, "main.go"), "package main")

	tests := []struct {
		name         string
		depth        int
		wantContains []string
		wantMissing  []string
	}{
		{
			name:         "zero shows only root",
			depth:        0,
			wantContains: []string{".", "0 directories, 0 files"},
			wantMissing:  []string{"main.go", "subpkg", "nested.go"},
		},
		{
			name:         "one shows direct children",
			depth:        1,
			wantContains: []string{"main.go", "subpkg", "1 directories, 1 files"},
			wantMissing:  []string{"nested.go"},
		},
		{
			name:         "two shows grandchildren",
			depth:        2,
			wantContains: []string{"main.go", "subpkg", "nested.go", "1 directories, 2 files"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output bytes.Buffer
			err := dispatcher(module.Options{
				Directory: dir,
				Flags: module.Flags{
					HideIcon:       true,
					ShowTreeView:   true,
					TreeDepth:      tt.depth,
					LimitTreeDepth: true,
				},
			}, &output)
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

func TestDispatcher_NoSummary(t *testing.T) {
	dir := t.TempDir()
	mustMkdir(t, filepath.Join(dir, "subpkg"))
	mustWriteFile(t, filepath.Join(dir, "main.go"), "package main")

	tests := []struct {
		name  string
		flags module.Flags
	}{
		{
			name:  "flat mode",
			flags: module.Flags{HideIcon: true, HideSummary: true},
		},
		{
			name:  "tree mode",
			flags: module.Flags{HideIcon: true, HideSummary: true, ShowTreeView: true},
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
			for _, want := range []string{".", "main.go", "subpkg"} {
				if !strings.Contains(got, want) {
					t.Fatalf("dispatcher() output = %q, want to contain %q", got, want)
				}
			}
			if strings.Contains(got, "directories,") {
				t.Fatalf("dispatcher() output = %q, want no summary", got)
			}
		})
	}
}

func TestDispatcher_DirsOnly(t *testing.T) {
	dir := t.TempDir()
	mustMkdir(t, filepath.Join(dir, "cmd"))
	mustMkdir(t, filepath.Join(dir, "service"))
	mustWriteFile(t, filepath.Join(dir, "README.md"), "readme")
	mustWriteFile(t, filepath.Join(dir, "main.go"), "package main")

	tests := []struct {
		name  string
		flags module.Flags
	}{
		{
			name:  "flat mode",
			flags: module.Flags{HideIcon: true, ShowDirsOnly: true},
		},
		{
			name:  "tree mode",
			flags: module.Flags{HideIcon: true, ShowDirsOnly: true, ShowTreeView: true},
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
			for _, want := range []string{"cmd", "service", "2 directories"} {
				if !strings.Contains(got, want) {
					t.Fatalf("dispatcher() output = %q, want to contain %q", got, want)
				}
			}
			if strings.Contains(got, "0 files") {
				t.Fatalf("dispatcher() output = %q, want no file count in dirs-only summary", got)
			}
			for _, missing := range []string{"README.md", "main.go"} {
				if strings.Contains(got, missing) {
					t.Fatalf("dispatcher() output = %q, want not to contain %q", got, missing)
				}
			}
		})
	}
}

func TestDispatcher_DirsOnlyWithDepth(t *testing.T) {
	dir := t.TempDir()
	mustMkdir(t, filepath.Join(dir, "cmd"))
	mustMkdir(t, filepath.Join(dir, "cmd", "internal"))
	mustWriteFile(t, filepath.Join(dir, "cmd", "main.go"), "package main")
	mustMkdir(t, filepath.Join(dir, "service"))

	var output bytes.Buffer
	err := dispatcher(module.Options{
		Directory: dir,
		Flags: module.Flags{
			HideIcon:       true,
			ShowTreeView:   true,
			ShowDirsOnly:   true,
			TreeDepth:      1,
			LimitTreeDepth: true,
		},
	}, &output)
	if err != nil {
		t.Fatalf("dispatcher() error = %v, want nil", err)
	}

	got := stripANSI(output.String())
	for _, want := range []string{"cmd", "service", "2 directories"} {
		if !strings.Contains(got, want) {
			t.Fatalf("dispatcher() output = %q, want to contain %q", got, want)
		}
	}
	if strings.Contains(got, "0 files") {
		t.Fatalf("dispatcher() output = %q, want no file count in dirs-only summary", got)
	}
	for _, missing := range []string{"internal", "main.go"} {
		if strings.Contains(got, missing) {
			t.Fatalf("dispatcher() output = %q, want not to contain %q", got, missing)
		}
	}
}

func TestDispatcher_DirsOnlyWithGitIgnore(t *testing.T) {
	dir := t.TempDir()
	mustWriteFile(t, filepath.Join(dir, ".gitignore"), "ignored-dir/\n")
	mustMkdir(t, filepath.Join(dir, "ignored-dir"))
	mustMkdir(t, filepath.Join(dir, "visible-dir"))
	mustWriteFile(t, filepath.Join(dir, "main.go"), "package main")

	var output bytes.Buffer
	err := dispatcher(module.Options{
		Directory: dir,
		Flags: module.Flags{
			HideIcon:         true,
			ShowTreeView:     true,
			ShowDirsOnly:     true,
			RespectGitIgnore: true,
		},
	}, &output)
	if err != nil {
		t.Fatalf("dispatcher() error = %v, want nil", err)
	}

	got := stripANSI(output.String())
	for _, want := range []string{"visible-dir", "1 directories"} {
		if !strings.Contains(got, want) {
			t.Fatalf("dispatcher() output = %q, want to contain %q", got, want)
		}
	}
	if strings.Contains(got, "0 files") {
		t.Fatalf("dispatcher() output = %q, want no file count in dirs-only summary", got)
	}
	for _, missing := range []string{"ignored-dir", "main.go"} {
		if strings.Contains(got, missing) {
			t.Fatalf("dispatcher() output = %q, want not to contain %q", got, missing)
		}
	}
}

func TestDispatcher_ListDirectoryGitIgnore(t *testing.T) {
	dir := t.TempDir()
	mustWriteFile(t, filepath.Join(dir, ".gitignore"), "*.log\nignored-dir/\n!keep.log\n")
	mustWriteFile(t, filepath.Join(dir, "main.go"), "package main")
	mustWriteFile(t, filepath.Join(dir, "debug.log"), "debug")
	mustWriteFile(t, filepath.Join(dir, "keep.log"), "keep")
	mustMkdir(t, filepath.Join(dir, "ignored-dir"))
	mustWriteFile(t, filepath.Join(dir, "ignored-dir", "ignored.go"), "package ignored")

	var output bytes.Buffer
	err := dispatcher(module.Options{
		Directory: dir,
		Flags: module.Flags{
			HideIcon:         true,
			RespectGitIgnore: true,
		},
	}, &output)
	if err != nil {
		t.Fatalf("dispatcher() error = %v, want nil", err)
	}

	got := stripANSI(output.String())
	for _, want := range []string{"main.go", "keep.log", "0 directories, 2 files"} {
		if !strings.Contains(got, want) {
			t.Fatalf("dispatcher() output = %q, want to contain %q", got, want)
		}
	}
	for _, missing := range []string{"debug.log", "ignored-dir"} {
		if strings.Contains(got, missing) {
			t.Fatalf("dispatcher() output = %q, want not to contain %q", got, missing)
		}
	}
}

func TestDispatcher_TreeDirectoryGitIgnore(t *testing.T) {
	dir := t.TempDir()
	mustWriteFile(t, filepath.Join(dir, ".gitignore"), "node_modules/\ndist/*.js\n*.log\n!keep.log\n")
	mustMkdir(t, filepath.Join(dir, "node_modules"))
	mustWriteFile(t, filepath.Join(dir, "node_modules", "package.js"), "module")
	mustMkdir(t, filepath.Join(dir, "dist"))
	mustWriteFile(t, filepath.Join(dir, "dist", "app.js"), "app")
	mustWriteFile(t, filepath.Join(dir, "dist", "style.css"), "style")
	mustMkdir(t, filepath.Join(dir, "src"))
	mustWriteFile(t, filepath.Join(dir, "src", "main.go"), "package main")
	mustWriteFile(t, filepath.Join(dir, "debug.log"), "debug")
	mustWriteFile(t, filepath.Join(dir, "keep.log"), "keep")

	var output bytes.Buffer
	err := dispatcher(module.Options{
		Directory: dir,
		Flags: module.Flags{
			HideIcon:         true,
			ShowTreeView:     true,
			RespectGitIgnore: true,
		},
	}, &output)
	if err != nil {
		t.Fatalf("dispatcher() error = %v, want nil", err)
	}

	got := stripANSI(output.String())
	for _, want := range []string{"dist", "style.css", "src", "main.go", "keep.log", "2 directories, 3 files"} {
		if !strings.Contains(got, want) {
			t.Fatalf("dispatcher() output = %q, want to contain %q", got, want)
		}
	}
	for _, missing := range []string{"node_modules", "app.js", "debug.log"} {
		if strings.Contains(got, missing) {
			t.Fatalf("dispatcher() output = %q, want not to contain %q", got, missing)
		}
	}
}

func TestDispatcher_GitIgnoreMissingFile(t *testing.T) {
	dir := t.TempDir()
	mustWriteFile(t, filepath.Join(dir, "main.go"), "package main")

	var output bytes.Buffer
	err := dispatcher(module.Options{
		Directory: dir,
		Flags: module.Flags{
			HideIcon:         true,
			RespectGitIgnore: true,
		},
	}, &output)
	if err != nil {
		t.Fatalf("dispatcher() error = %v, want nil", err)
	}

	got := stripANSI(output.String())
	if !strings.Contains(got, "main.go") {
		t.Fatalf("dispatcher() output = %q, want to contain main.go", got)
	}
}

func TestDispatcher_GitIgnoreWithHiddenFiles(t *testing.T) {
	dir := t.TempDir()
	mustWriteFile(t, filepath.Join(dir, ".gitignore"), ".env\n")
	mustWriteFile(t, filepath.Join(dir, ".env"), "secret")
	mustWriteFile(t, filepath.Join(dir, ".hidden"), "hidden")
	mustWriteFile(t, filepath.Join(dir, "main.go"), "package main")

	var output bytes.Buffer
	err := dispatcher(module.Options{
		Directory: dir,
		Flags: module.Flags{
			HideIcon:         true,
			ShowHidden:       true,
			RespectGitIgnore: true,
		},
	}, &output)
	if err != nil {
		t.Fatalf("dispatcher() error = %v, want nil", err)
	}

	got := stripANSI(output.String())
	for _, want := range []string{".gitignore", ".hidden", "main.go", "0 directories, 3 files"} {
		if !strings.Contains(got, want) {
			t.Fatalf("dispatcher() output = %q, want to contain %q", got, want)
		}
	}
	if strings.Contains(got, ".env") {
		t.Fatalf("dispatcher() output = %q, want not to contain ignored .env", got)
	}
}

func TestGitIgnoreRuleParsingAndMatching(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		path     string
		isDir    bool
		wantRule bool
		want     bool
		negated  bool
	}{
		{
			name:     "blank line",
			line:     "   ",
			wantRule: false,
		},
		{
			name:     "comment line",
			line:     "# generated files",
			wantRule: false,
		},
		{
			name:     "escaped comment",
			line:     `\#literal`,
			path:     "#literal",
			wantRule: true,
			want:     true,
		},
		{
			name:     "escaped negation",
			line:     `\!literal`,
			path:     "!literal",
			wantRule: true,
			want:     true,
		},
		{
			name:     "negated rule",
			line:     "!keep.log",
			path:     "keep.log",
			wantRule: true,
			want:     true,
			negated:  true,
		},
		{
			name:     "anchored rule",
			line:     "/dist",
			path:     "dist",
			wantRule: true,
			want:     true,
		},
		{
			name:     "directory only does not match file",
			line:     "build/",
			path:     "build",
			isDir:    false,
			wantRule: true,
			want:     false,
		},
		{
			name:     "directory only matches directory",
			line:     "build/",
			path:     "build",
			isDir:    true,
			wantRule: true,
			want:     true,
		},
		{
			name:     "double star path",
			line:     "logs/**/*.tmp",
			path:     "logs/app/archive/debug.tmp",
			wantRule: true,
			want:     true,
		},
		{
			name:     "invalid glob",
			line:     "[",
			path:     "[",
			wantRule: true,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule, ok := parseGitIgnoreRule(tt.line)
			if ok != tt.wantRule {
				t.Fatalf("parseGitIgnoreRule(%q) ok = %v, want %v", tt.line, ok, tt.wantRule)
			}
			if !ok {
				return
			}
			if rule.negated != tt.negated {
				t.Fatalf("parseGitIgnoreRule(%q) negated = %v, want %v", tt.line, rule.negated, tt.negated)
			}
			if got := rule.matches(tt.path, tt.isDir); got != tt.want {
				t.Fatalf("rule.matches(%q, %v) = %v, want %v", tt.path, tt.isDir, got, tt.want)
			}
		})
	}
}

func TestGitIgnoreMatcherIgnores(t *testing.T) {
	root := t.TempDir()
	matcher := &gitIgnoreMatcher{
		root: root,
		rules: []gitIgnoreRule{
			{pattern: "*.log"},
			{pattern: "keep.log", negated: true},
			{pattern: "cache", dirOnly: true},
		},
	}

	tests := []struct {
		name  string
		path  string
		isDir bool
		want  bool
	}{
		{name: "ignored file", path: filepath.Join(root, "debug.log"), want: true},
		{name: "negated file", path: filepath.Join(root, "keep.log"), want: false},
		{name: "ignored directory", path: filepath.Join(root, "cache"), isDir: true, want: true},
		{name: "outside root", path: filepath.Join(t.TempDir(), "debug.log"), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matcher.ignores(tt.path, tt.isDir); got != tt.want {
				t.Fatalf("matcher.ignores(%q, %v) = %v, want %v", tt.path, tt.isDir, got, tt.want)
			}
		})
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

func TestResolveFileIconStyle(t *testing.T) {
	tests := []struct {
		name      string
		fileName  string
		wantIcon  string
		wantColor string
	}{
		{
			name:      "filename specific icon takes precedence",
			fileName:  "README.md",
			wantIcon:  module.ReadmeIcon,
			wantColor: module.Cyan,
		},
		{
			name:      "known extension is case insensitive",
			fileName:  "main.GO",
			wantIcon:  module.GoLangIcon,
			wantColor: module.LightBlue,
		},
		{
			name:      "multi-part config filename",
			fileName:  ".eslintrc.json",
			wantIcon:  module.ESLintIcon,
			wantColor: module.Purple,
		},
		{
			name:      "exact filename is case insensitive",
			fileName:  "PACKAGE.JSON",
			wantIcon:  module.NPMIcon,
			wantColor: module.Red,
		},
		{
			name:      "unknown extension falls back to file icon",
			fileName:  "archive.unknown",
			wantIcon:  module.FileIcon,
			wantColor: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveFileIconStyle(tt.fileName)
			if got.icon != tt.wantIcon || got.color != tt.wantColor {
				t.Fatalf("resolveFileIconStyle(%q) = %+v, want icon %q and color %q", tt.fileName, got, tt.wantIcon, tt.wantColor)
			}
		})
	}
}

func TestResolveFolderIconStyle(t *testing.T) {
	tests := []struct {
		name       string
		folderName string
		wantIcon   string
		wantColor  string
	}{
		{
			name:       "git folder",
			folderName: ".git",
			wantIcon:   module.GitIcon,
			wantColor:  module.Orange,
		},
		{
			name:       "folder match is case insensitive",
			folderName: ".GIT",
			wantIcon:   module.GitIcon,
			wantColor:  module.Orange,
		},
		{
			name:       "default folder",
			folderName: "pkg",
			wantIcon:   module.FolderIcon,
			wantColor:  module.Pink,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveFolderIconStyle(tt.folderName)
			if got.icon != tt.wantIcon || got.color != tt.wantColor {
				t.Fatalf("resolveFolderIconStyle(%q) = %+v, want icon %q and color %q", tt.folderName, got, tt.wantIcon, tt.wantColor)
			}
		})
	}
}

func TestFormatFileWithOptions_HideIcon(t *testing.T) {
	file := fakeFileInfo{name: "main.go"}

	got := formatFileWithOptions("", file, module.Flags{HideIcon: true})
	if strings.Contains(got, module.GoLangIcon) {
		t.Fatalf("formatFileWithOptions() = %q, want no file icon", got)
	}
	if got != "main.go" {
		t.Fatalf("formatFileWithOptions() = %q, want plain file name", got)
	}
}

func TestFormatFileWithOptions_WithIcons(t *testing.T) {
	tests := []struct {
		name string
		file fakeFileInfo
		want []string
	}{
		{
			name: "known file icon",
			file: fakeFileInfo{name: "main.go"},
			want: []string{module.GoLangIcon, "main.go"},
		},
		{
			name: "default file icon",
			file: fakeFileInfo{name: "unknownfile"},
			want: []string{module.FileIcon, "unknownfile"},
		},
		{
			name: "known folder icon",
			file: fakeFileInfo{name: ".git", isDir: true},
			want: []string{module.GitIcon, ".git"},
		},
		{
			name: "default folder icon",
			file: fakeFileInfo{name: "src", isDir: true},
			want: []string{module.FolderIcon, "src"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatFileWithOptions("├── ", tt.file, module.Flags{})
			for _, want := range tt.want {
				if !strings.Contains(got, want) {
					t.Fatalf("formatFileWithOptions() = %q, want to contain %q", got, want)
				}
			}
		})
	}
}

func TestCheckDefaultDirectory(t *testing.T) {
	var directory string
	if err := CheckDefaultDirectory(&directory); err != nil {
		t.Fatalf("CheckDefaultDirectory() error = %v, want nil", err)
	}
	if directory == "" {
		t.Fatal("CheckDefaultDirectory() directory = empty, want current directory")
	}

	existing := "/tmp"
	if err := CheckDefaultDirectory(&existing); err != nil {
		t.Fatalf("CheckDefaultDirectory() error = %v, want nil", err)
	}
	if existing != "/tmp" {
		t.Fatalf("CheckDefaultDirectory() directory = %q, want /tmp", existing)
	}
}

func TestCloseDirectoryNil(t *testing.T) {
	if err := closeDirectory(nil); err != nil {
		t.Fatalf("closeDirectory(nil) error = %v, want nil", err)
	}
}

func TestPrintGridEmpty(t *testing.T) {
	var output bytes.Buffer
	if err := printGrid(&output, nil, 0); err != nil {
		t.Fatalf("printGrid() error = %v, want nil", err)
	}
	if output.Len() != 0 {
		t.Fatalf("printGrid() output = %q, want empty", output.String())
	}
}

func TestPrintGridWriteErrors(t *testing.T) {
	tests := []struct {
		name      string
		failAfter int
		entries   []string
		maxLen    int
	}{
		{
			name:      "entry write",
			failAfter: 0,
			entries:   []string{"alpha"},
			maxLen:    5,
		},
		{
			name:      "padding write",
			failAfter: 1,
			entries:   []string{"alpha", "beta"},
			maxLen:    5,
		},
		{
			name:      "newline write",
			failAfter: 1,
			entries:   []string{"alpha"},
			maxLen:    200,
		},
		{
			name:      "final newline write",
			failAfter: 4,
			entries:   []string{"alpha", "beta"},
			maxLen:    5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &failingWriter{failAfter: tt.failAfter}
			if err := printGrid(writer, tt.entries, tt.maxLen); err == nil {
				t.Fatal("printGrid() error = nil, want error")
			}
		})
	}
}

func TestMatchPatternSegmentsFailures(t *testing.T) {
	tests := []struct {
		name            string
		patternSegments []string
		nameSegments    []string
	}{
		{
			name:            "name exhausted before pattern",
			patternSegments: []string{"src"},
		},
		{
			name:            "invalid segment glob",
			patternSegments: []string{"["},
			nameSegments:    []string{"["},
		},
		{
			name:            "segment mismatch",
			patternSegments: []string{"src"},
			nameSegments:    []string{"cmd"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if matchPatternSegments(tt.patternSegments, tt.nameSegments) {
				t.Fatalf("matchPatternSegments(%v, %v) = true, want false", tt.patternSegments, tt.nameSegments)
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

func TestDispatcher_JSONFlatDirectory(t *testing.T) {
	dir := t.TempDir()
	mustWriteFile(t, filepath.Join(dir, "main.go"), "package main")
	mustMkdir(t, filepath.Join(dir, "service"))
	mustWriteFile(t, filepath.Join(dir, ".hidden"), "hidden")

	var output bytes.Buffer
	err := dispatcher(module.Options{
		Directory: dir,
		Flags:     module.Flags{ShowJSON: true},
	}, &output)
	if err != nil {
		t.Fatalf("dispatcher() error = %v, want nil", err)
	}

	var got jsonOutput
	if err := json.Unmarshal(output.Bytes(), &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v, output = %q", err, output.String())
	}

	if got.Root != dir {
		t.Fatalf("json root = %q, want %q", got.Root, dir)
	}
	if got.Tree {
		t.Fatal("json tree = true, want false")
	}
	if got.Summary.Directories != 1 || got.Summary.Files != 1 {
		t.Fatalf("json summary = %+v, want 1 directory and 1 file", got.Summary)
	}
	if len(got.Entries) != 2 {
		t.Fatalf("json entries length = %d, want 2", len(got.Entries))
	}
	if got.Entries[0].Name != "main.go" || got.Entries[0].Type != "file" || got.Entries[0].Size != 12 {
		t.Fatalf("first json entry = %+v, want main.go file", got.Entries[0])
	}
	if got.Entries[1].Name != "service" || got.Entries[1].Type != "directory" {
		t.Fatalf("second json entry = %+v, want service directory", got.Entries[1])
	}
	if strings.Contains(output.String(), "directories,") || strings.Contains(output.String(), "├──") {
		t.Fatalf("json output = %q, want no human formatted output", output.String())
	}
}

func TestDispatcher_JSONTreeDirectoryWithDepth(t *testing.T) {
	dir := t.TempDir()
	mustMkdir(t, filepath.Join(dir, "cmd"))
	mustMkdir(t, filepath.Join(dir, "cmd", "internal"))
	mustWriteFile(t, filepath.Join(dir, "cmd", "root.go"), "package cmd")
	mustWriteFile(t, filepath.Join(dir, "README.md"), "readme")

	var output bytes.Buffer
	err := dispatcher(module.Options{
		Directory: dir,
		Flags: module.Flags{
			ShowJSON:       true,
			ShowTreeView:   true,
			TreeDepth:      1,
			LimitTreeDepth: true,
		},
	}, &output)
	if err != nil {
		t.Fatalf("dispatcher() error = %v, want nil", err)
	}

	var got jsonOutput
	if err := json.Unmarshal(output.Bytes(), &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if !got.Tree {
		t.Fatal("json tree = false, want true")
	}
	if got.Summary.Directories != 1 || got.Summary.Files != 1 {
		t.Fatalf("json summary = %+v, want depth-limited 1 directory and 1 file", got.Summary)
	}
	if len(got.Entries) != 2 {
		t.Fatalf("json entries length = %d, want 2", len(got.Entries))
	}
	cmdEntry := got.Entries[1]
	if cmdEntry.Name != "cmd" {
		t.Fatalf("second json entry = %+v, want cmd directory", cmdEntry)
	}
	if len(cmdEntry.Children) != 0 {
		t.Fatalf("cmd children = %+v, want none at depth 1", cmdEntry.Children)
	}
}

func TestDispatcher_JSONWithGitIgnoreAndDirsOnly(t *testing.T) {
	dir := t.TempDir()
	mustWriteFile(t, filepath.Join(dir, ".gitignore"), "ignored/\n")
	mustMkdir(t, filepath.Join(dir, "ignored"))
	mustMkdir(t, filepath.Join(dir, "visible"))
	mustWriteFile(t, filepath.Join(dir, "main.go"), "package main")

	var output bytes.Buffer
	err := dispatcher(module.Options{
		Directory: dir,
		Flags: module.Flags{
			ShowJSON:         true,
			ShowTreeView:     true,
			ShowDirsOnly:     true,
			RespectGitIgnore: true,
		},
	}, &output)
	if err != nil {
		t.Fatalf("dispatcher() error = %v, want nil", err)
	}

	var got jsonOutput
	if err := json.Unmarshal(output.Bytes(), &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if got.Summary.Directories != 1 || got.Summary.Files != 0 {
		t.Fatalf("json summary = %+v, want 1 directory and 0 files", got.Summary)
	}
	if len(got.Entries) != 1 || got.Entries[0].Name != "visible" {
		t.Fatalf("json entries = %+v, want only visible directory", got.Entries)
	}
}

func TestPrintJSONDirectoryErrors(t *testing.T) {
	path := filepath.Join(t.TempDir(), "regular.txt")
	mustWriteFile(t, path, "content")

	var output bytes.Buffer
	err := printJSONDirectory(directoryOptions{Options: module.Options{Directory: path}}, &output)
	if err == nil {
		t.Fatal("printJSONDirectory() error = nil, want read error")
	}

	_, _, err = collectJSONFlatEntries(directoryOptions{Options: module.Options{Directory: path}})
	if err == nil {
		t.Fatal("collectJSONFlatEntries() error = nil, want read error")
	}

	_, _, err = collectJSONTreeEntries(directoryOptions{Options: module.Options{Directory: path}}, 0)
	if err == nil {
		t.Fatal("collectJSONTreeEntries() error = nil, want read error")
	}

	entries, summary, err := collectJSONTreeEntries(directoryOptions{
		Options: module.Options{
			Directory: path,
			Flags: module.Flags{
				TreeDepth:      0,
				LimitTreeDepth: true,
			},
		},
	}, 0)
	if err != nil {
		t.Fatalf("collectJSONTreeEntries() error = %v, want nil at max depth", err)
	}
	if len(entries) != 0 || summary.Directories != 0 || summary.Files != 0 {
		t.Fatalf("collectJSONTreeEntries() = entries %+v summary %+v, want empty", entries, summary)
	}

	dir := t.TempDir()
	mustWriteFile(t, filepath.Join(dir, "main.go"), "package main")
	err = printJSONDirectory(directoryOptions{Options: module.Options{Directory: dir}}, &failingWriter{})
	if err == nil {
		t.Fatal("printJSONDirectory() error = nil, want write error")
	}
}

func TestJSONEntryAndSummaryHelpers(t *testing.T) {
	file := fakeFileInfo{name: "main.go", size: 12}
	fileEntry := newJSONEntry("/tmp/project", file)
	if fileEntry.Name != "main.go" || fileEntry.Type != "file" || fileEntry.Size != 12 {
		t.Fatalf("newJSONEntry() = %+v, want file entry", fileEntry)
	}

	dir := fakeFileInfo{name: "cmd", isDir: true}
	dirEntry := newJSONEntry("/tmp/project", dir)
	if dirEntry.Name != "cmd" || dirEntry.Type != "directory" {
		t.Fatalf("newJSONEntry() = %+v, want directory entry", dirEntry)
	}

	var summary jsonSummary
	addJSONSummaryCount(&summary, file)
	addJSONSummaryCount(&summary, dir)
	if summary.Files != 1 || summary.Directories != 1 {
		t.Fatalf("summary = %+v, want 1 file and 1 directory", summary)
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

type fakeFileInfo struct {
	name  string
	size  int64
	isDir bool
}

func (f fakeFileInfo) Name() string {
	return f.name
}

func (f fakeFileInfo) Size() int64 {
	return f.size
}

func (f fakeFileInfo) Mode() os.FileMode {
	if f.isDir {
		return os.ModeDir
	}
	return 0
}

func (f fakeFileInfo) ModTime() time.Time {
	return time.Time{}
}

func (f fakeFileInfo) IsDir() bool {
	return f.isDir
}

func (f fakeFileInfo) Sys() interface{} {
	return nil
}

type failingWriter struct {
	writes    int
	failAfter int
}

func (w *failingWriter) Write(p []byte) (int, error) {
	if w.writes >= w.failAfter {
		return 0, os.ErrPermission
	}
	w.writes++
	return len(p), nil
}
