package service

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type gitStatusMatcher struct {
	root     string
	statuses map[string]string
}

type deletedGitFileInfo struct {
	name string
}

func (f deletedGitFileInfo) Name() string       { return f.name }
func (f deletedGitFileInfo) Size() int64        { return 0 }
func (f deletedGitFileInfo) Mode() os.FileMode  { return 0 }
func (f deletedGitFileInfo) ModTime() time.Time { return time.Time{} }
func (f deletedGitFileInfo) IsDir() bool        { return false }
func (f deletedGitFileInfo) Sys() interface{}   { return nil }

var gitStatusCommand = func(root string) ([]byte, error) {
	return exec.Command("git", "-C", root, "status", "--porcelain=v1", "--ignored=matching", "--untracked-files=all").Output()
}

func newGitStatusMatcher(root string) *gitStatusMatcher {
	output, err := gitStatusCommand(root)
	if err != nil {
		return nil
	}

	statuses := parseGitStatusOutput(string(output))
	if len(statuses) == 0 {
		return nil
	}

	return &gitStatusMatcher{
		root:     root,
		statuses: statuses,
	}
}

func parseGitStatusOutput(output string) map[string]string {
	statuses := make(map[string]string)
	for _, line := range strings.Split(output, "\n") {
		code, path, ok := parseGitStatusLine(line)
		if !ok {
			continue
		}

		if symbol := gitStatusSymbol(code); symbol != "" {
			statuses[normalizeGitStatusPath(path)] = symbol
		}
	}
	return statuses
}

func parseGitStatusLine(line string) (code, path string, ok bool) {
	if len(line) < 4 {
		return "", "", false
	}

	code = line[:2]
	path = strings.TrimSpace(line[3:])
	if path == "" {
		return "", "", false
	}

	if strings.Contains(path, " -> ") {
		parts := strings.Split(path, " -> ")
		path = parts[len(parts)-1]
	}

	return code, path, true
}

func gitStatusSymbol(code string) string {
	if code == "??" {
		return "?"
	}
	if code == "!!" {
		return "!"
	}
	if strings.Contains(code, "D") {
		return "D"
	}
	if strings.Contains(code, "A") {
		return "A"
	}
	if strings.ContainsAny(code, "MRC") {
		return "M"
	}
	return ""
}

func normalizeGitStatusPath(path string) string {
	path = strings.Trim(path, `"`)
	path = filepath.ToSlash(path)
	return strings.TrimRight(path, "/")
}

func (m *gitStatusMatcher) appendDeletedFiles(directory string, files []os.FileInfo) []os.FileInfo {
	if m == nil {
		return files
	}

	existingNames := make(map[string]struct{}, len(files))
	for _, file := range files {
		existingNames[file.Name()] = struct{}{}
	}

	for relPath, status := range m.statuses {
		if status != "D" {
			continue
		}

		fullPath := filepath.Join(m.root, filepath.FromSlash(relPath))
		if filepath.Dir(fullPath) != directory {
			continue
		}

		name := filepath.Base(fullPath)
		if _, exists := existingNames[name]; exists {
			continue
		}
		files = append(files, deletedGitFileInfo{name: name})
		existingNames[name] = struct{}{}
	}

	return files
}

func (m *gitStatusMatcher) statusFor(filePath string) string {
	if m == nil {
		return ""
	}

	relPath, err := filepath.Rel(m.root, filePath)
	if err != nil {
		return ""
	}
	relPath = filepath.ToSlash(relPath)
	if relPath == "." || strings.HasPrefix(relPath, "../") {
		return ""
	}

	return m.statuses[normalizeGitStatusPath(relPath)]
}
