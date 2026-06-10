package service

import (
	"os"
	"path/filepath"
	"strings"
)

func hasIncludePatterns(options directoryOptions) bool {
	return len(options.Flags.IncludePatterns) > 0
}

func passesEntryFilters(file os.FileInfo, options directoryOptions) bool {
	if isHidden(file.Name()) && !options.Flags.ShowHidden {
		return false
	}

	if options.Flags.ShowDirsOnly && !file.IsDir() {
		return false
	}

	filePath := filepath.Join(options.Directory, file.Name())
	if options.gitIgnore != nil && options.gitIgnore.ignores(filePath, file.IsDir()) {
		return false
	}

	return !matchesAnyFilterPattern(options, filePath, file)
}

func matchesIncludePattern(options directoryOptions, filePath string, file os.FileInfo) bool {
	return matchesAnyPattern(options.root, options.Flags.IncludePatterns, filePath, file.Name(), file.IsDir())
}

func matchesAnyFilterPattern(options directoryOptions, filePath string, file os.FileInfo) bool {
	return matchesAnyPattern(options.root, options.Flags.ExcludePatterns, filePath, file.Name(), file.IsDir())
}

func matchesAnyPattern(root string, patterns []string, filePath, name string, isDir bool) bool {
	for _, pattern := range patterns {
		if matchesFilterPattern(root, pattern, filePath, name, isDir) {
			return true
		}
	}
	return false
}

func matchesFilterPattern(root, pattern, filePath, name string, isDir bool) bool {
	pattern = strings.TrimSpace(pattern)
	if pattern == "" {
		return false
	}

	pattern = filepath.ToSlash(pattern)
	relPath, err := filepath.Rel(root, filePath)
	if err != nil {
		return false
	}
	relPath = filepath.ToSlash(relPath)
	if relPath == "." || strings.HasPrefix(relPath, "../") {
		return false
	}

	if strings.HasSuffix(pattern, "/") {
		if !isDir {
			return false
		}
		pattern = strings.TrimRight(pattern, "/")
	}

	if strings.Contains(pattern, "/") || strings.Contains(pattern, "**") {
		if matchSlashPattern(pattern, relPath) {
			return true
		}

		if strings.HasSuffix(pattern, "/**") {
			prefix := strings.TrimSuffix(pattern, "/**")
			return relPath == prefix || strings.HasPrefix(relPath, prefix+"/")
		}
		return false
	}

	return matchSlashPattern(pattern, name)
}
