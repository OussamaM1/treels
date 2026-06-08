package service

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const gitignoreFileName = ".gitignore"

// gitIgnoreMatcher holds parsed .gitignore rules for a root directory.
type gitIgnoreMatcher struct {
	root  string
	rules []gitIgnoreRule
}

// gitIgnoreRule represents one parsed .gitignore pattern.
type gitIgnoreRule struct {
	pattern  string
	negated  bool
	anchored bool
	dirOnly  bool
	hasSlash bool
}

// newGitIgnoreMatcher creates a matcher from the root directory's .gitignore file.
func newGitIgnoreMatcher(root string) (*gitIgnoreMatcher, error) {
	gitignorePath := filepath.Join(root, gitignoreFileName)
	file, err := os.Open(gitignorePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read %s: %w", gitignorePath, err)
	}
	defer func() {
		_ = file.Close()
	}()

	rules := parseGitIgnoreRules(file)
	if len(rules) == 0 {
		return nil, nil
	}

	return &gitIgnoreMatcher{
		root:  root,
		rules: rules,
	}, nil
}

// parseGitIgnoreRules parses all supported .gitignore rules from a file.
func parseGitIgnoreRules(file *os.File) []gitIgnoreRule {
	var rules []gitIgnoreRule
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if rule, ok := parseGitIgnoreRule(scanner.Text()); ok {
			rules = append(rules, rule)
		}
	}
	return rules
}

// parseGitIgnoreRule parses a single .gitignore line into a rule.
func parseGitIgnoreRule(line string) (gitIgnoreRule, bool) {
	line = strings.TrimSpace(line)
	if line == "" {
		return gitIgnoreRule{}, false
	}

	if strings.HasPrefix(line, `\#`) {
		line = strings.TrimPrefix(line, `\`)
	} else if strings.HasPrefix(line, "#") {
		return gitIgnoreRule{}, false
	}

	negated := false
	if strings.HasPrefix(line, `\!`) {
		line = strings.TrimPrefix(line, `\`)
	} else if strings.HasPrefix(line, "!") {
		negated = true
		line = strings.TrimPrefix(line, "!")
	}

	line = strings.TrimSpace(line)
	if line == "" {
		return gitIgnoreRule{}, false
	}

	anchored := strings.HasPrefix(line, "/")
	line = strings.TrimPrefix(line, "/")

	dirOnly := strings.HasSuffix(line, "/")
	line = strings.TrimRight(line, "/")
	if line == "" {
		return gitIgnoreRule{}, false
	}

	return gitIgnoreRule{
		pattern:  filepath.ToSlash(line),
		negated:  negated,
		anchored: anchored,
		dirOnly:  dirOnly,
		hasSlash: strings.Contains(line, "/"),
	}, true
}

// ignores reports whether a path should be excluded by the parsed rules.
func (m *gitIgnoreMatcher) ignores(filePath string, isDir bool) bool {
	if m == nil {
		return false
	}

	relPath, err := filepath.Rel(m.root, filePath)
	if err != nil {
		return false
	}

	relPath = filepath.ToSlash(relPath)
	if relPath == "." || strings.HasPrefix(relPath, "../") {
		return false
	}

	ignored := false
	for _, rule := range m.rules {
		if rule.matches(relPath, isDir) {
			ignored = !rule.negated
		}
	}
	return ignored
}

// matches reports whether a relative path matches this .gitignore rule.
func (r gitIgnoreRule) matches(relPath string, isDir bool) bool {
	if r.dirOnly && !isDir {
		return false
	}

	if r.anchored {
		return matchSlashPattern(r.pattern, relPath)
	}

	if !r.hasSlash {
		for _, part := range strings.Split(relPath, "/") {
			if matchSlashPattern(r.pattern, part) {
				return true
			}
		}
		return false
	}

	return matchSlashPattern(r.pattern, relPath)
}

// matchSlashPattern matches slash-separated paths with glob support.
func matchSlashPattern(pattern, name string) bool {
	if !strings.Contains(pattern, "**") {
		matched, err := path.Match(pattern, name)
		return err == nil && matched
	}

	return matchPatternSegments(strings.Split(pattern, "/"), strings.Split(name, "/"))
}

// matchPatternSegments matches path segments and supports the ** wildcard.
func matchPatternSegments(patternSegments, nameSegments []string) bool {
	if len(patternSegments) == 0 {
		return len(nameSegments) == 0
	}

	if patternSegments[0] == "**" {
		for i := 0; i <= len(nameSegments); i++ {
			if matchPatternSegments(patternSegments[1:], nameSegments[i:]) {
				return true
			}
		}
		return false
	}

	if len(nameSegments) == 0 {
		return false
	}

	matched, err := path.Match(patternSegments[0], nameSegments[0])
	if err != nil || !matched {
		return false
	}

	return matchPatternSegments(patternSegments[1:], nameSegments[1:])
}
