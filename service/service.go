// Package service provides functionalities related to directory listing.
package service

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/oussamaM1/treels/module"
	"github.com/oussamaM1/treels/utils"
)

const (
	whiteSpaces              = "    " // Whitespaces
	boxUpAndRight            = "└── " // BOX DRAWINGS HEAVY UP AND RIGHT
	boxLightVertical         = "│   " // BOX DRAWINGS LIGHT VERTICAL
	boxLightVerticalAndRight = "├── " // BOX DRAWINGS LIGHT VERTICAL AND RIGHT
	dot                      = "."
)

type directoryOptions struct {
	module.Options
	root      string
	gitIgnore *gitIgnoreMatcher
	gitStatus *gitStatusMatcher
}

// Dispatcher func - executes function based on flags
func Dispatcher(options module.Options) error {
	return dispatcher(options, os.Stdout)
}

func dispatcher(options module.Options, output io.Writer) error {
	if err := CheckDefaultDirectory(&options.Directory); err != nil {
		return err
	}
	if err := utils.ValidateDirectory(options.Directory); err != nil {
		return err
	}

	traversalOptions := directoryOptions{Options: options, root: options.Directory}
	if options.Flags.RespectGitIgnore {
		gitIgnore, err := newGitIgnoreMatcher(options.Directory)
		if err != nil {
			return err
		}
		traversalOptions.gitIgnore = gitIgnore
	}
	if options.Flags.ShowGitStatus {
		traversalOptions.gitStatus = newGitStatusMatcher(options.Directory)
	}

	if options.Flags.ShowJSON {
		return printJSONDirectory(traversalOptions, output)
	}

	var fileCount, dirCount int
	if _, err := fmt.Fprintln(output, dot); err != nil {
		return err
	}
	if options.Flags.ShowTreeView {
		var err error
		fileCount, dirCount, err = treeDirectory(traversalOptions, output, "", true, 0)
		if err != nil {
			return err
		}
	} else {
		var err error
		fileCount, dirCount, err = listDirectory(traversalOptions, output)
		if err != nil {
			return err
		}
	}
	if options.Flags.HideSummary {
		return nil
	}

	return printNumberOfFilesAndDirectories(output, fileCount, dirCount, options.Flags)
}

// listDirectory func - lists the content of the directory.
func listDirectory(options directoryOptions, output io.Writer) (fileCount, dirCount int, err error) {
	// Open and read the directory
	files, d, err := readDirectory(options.Directory)
	if err != nil {
		return 0, 0, fmt.Errorf("read directory %q: %w", options.Directory, err)
	}
	defer func() {
		closeErr := closeDirectory(d)
		if err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	files = options.gitStatus.appendDeletedFiles(options.Directory, files)

	// sort files by requested order
	sortSlice(files, options.Flags)

	// Collect formatted file entries
	var entries []string
	var maxLen int

	for _, file := range files {
		if shouldShowFile(file, options) {
			if file.IsDir() {
				dirCount++
			} else {
				fileCount++
			}

			formatted := formatFileWithOptions(gitStatusPrefix("", file, options), file, options.Flags)

			entries = append(entries, formatted)

			// Track max visible length for column width
			visibleLen := getVisibleLength(formatted)
			if visibleLen > maxLen {
				maxLen = visibleLen
			}
		}
	}

	if options.Flags.ShowLongFormat {
		for _, entry := range entries {
			if _, err := fmt.Fprintln(output, entry); err != nil {
				return 0, 0, err
			}
		}
		return fileCount, dirCount, nil
	}

	// Print in grid format
	if err := printGrid(output, entries, maxLen); err != nil {
		return 0, 0, err
	}

	return fileCount, dirCount, nil
}

// treeDirectory func - displays a tree view of the directory.
func treeDirectory(options directoryOptions, output io.Writer, indent string, isLastFolder bool, depth int) (fileCount, dirCount int, err error) {
	if reachedMaxDepth(options.Flags, depth) {
		return 0, 0, nil
	}

	// Open and read the directory
	files, d, err := readDirectory(options.Directory)
	if err != nil {
		return 0, 0, fmt.Errorf("read directory %q: %w", options.Directory, err)
	}
	defer func() {
		closeErr := closeDirectory(d)
		if err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	files = options.gitStatus.appendDeletedFiles(options.Directory, files)

	// Sort files by requested order
	sortSlice(files, options.Flags)

	visibleFiles, err := visibleTreeFiles(files, options, depth)
	if err != nil {
		return 0, 0, err
	}

	// Print files and directories
	fc, dc, err := printFilesAndDirectoriesTreeFormat(visibleFiles, options, output, indent, isLastFolder, depth)
	if err != nil {
		return 0, 0, err
	}
	return fc, dc, nil
}

func reachedMaxDepth(flags module.Flags, depth int) bool {
	return flags.LimitTreeDepth && depth >= flags.TreeDepth
}

func getLastVisibleIndex(files []os.FileInfo, _ directoryOptions) int {
	return len(files) - 1
}

// printFilesAndDirectoriesTreeFormat - prints files and directories in tree format
func printFilesAndDirectoriesTreeFormat(files []os.FileInfo, options directoryOptions, output io.Writer, indent string, isLastFolder bool, depth int) (fileCount, dirCount int, err error) {
	lastVisibleIndex := getLastVisibleIndex(files, options)
	for i, file := range files {
		isLast := i == lastVisibleIndex
		prefix, childIndent := calculateIndent(indent, isLast)

		if err := printFileWithPrefix(output, prefix, file, options); err != nil {
			return 0, 0, err
		}

		if file.IsDir() {
			dirCount++
			fc, dc, err := processDirectory(file, options, output, childIndent, isLast && isLastFolder, depth+1)
			if err != nil {
				return 0, 0, err
			}
			fileCount += fc
			dirCount += dc
		} else {
			fileCount++
		}
	}
	return fileCount, dirCount, nil
}

// shouldShowFile determines if a file should be displayed based on visibility settings
func shouldShowFile(file os.FileInfo, options directoryOptions) bool {
	if !passesEntryFilters(file, options) {
		return false
	}
	if !hasIncludePatterns(options) {
		return true
	}

	filePath := filepath.Join(options.Directory, file.Name())
	return matchesIncludePattern(options, filePath, file)
}

func visibleTreeFiles(files []os.FileInfo, options directoryOptions, depth int) ([]os.FileInfo, error) {
	visible := make([]os.FileInfo, 0, len(files))
	for _, file := range files {
		show, err := shouldShowTreeFile(file, options, depth)
		if err != nil {
			return nil, err
		}
		if show {
			visible = append(visible, file)
		}
	}
	return visible, nil
}

func shouldShowTreeFile(file os.FileInfo, options directoryOptions, depth int) (bool, error) {
	if !passesEntryFilters(file, options) {
		return false, nil
	}
	if !hasIncludePatterns(options) {
		return true, nil
	}

	filePath := filepath.Join(options.Directory, file.Name())
	if matchesIncludePattern(options, filePath, file) {
		return true, nil
	}
	if !file.IsDir() {
		return false, nil
	}

	childOptions := options
	childOptions.Directory = filePath
	return directoryHasIncludedDescendant(childOptions, depth+1)
}

func directoryHasIncludedDescendant(options directoryOptions, depth int) (bool, error) {
	if reachedMaxDepth(options.Flags, depth) {
		return false, nil
	}

	files, d, err := readDirectory(options.Directory)
	if err != nil {
		return false, fmt.Errorf("read directory %q: %w", options.Directory, err)
	}
	defer func() {
		_ = closeDirectory(d)
	}()

	for _, file := range files {
		show, err := shouldShowTreeFile(file, options, depth)
		if err != nil {
			return false, err
		}
		if show {
			return true, nil
		}
	}
	return false, nil
}

// calculateIndent returns the appropriate prefix and child indent strings
func calculateIndent(indent string, isLast bool) (prefix, childIndent string) {
	if isLast {
		return indent + boxUpAndRight, indent + whiteSpaces
	}
	return indent + boxLightVerticalAndRight, indent + boxLightVertical
}

// printFileWithPrefix prints the file with the given prefix and icon settings
func printFileWithPrefix(output io.Writer, prefix string, file os.FileInfo, options directoryOptions) error {
	_, err := fmt.Fprintln(output, formatFileWithOptions(gitStatusPrefix(prefix, file, options), file, options.Flags))
	return err
}

func gitStatusPrefix(prefix string, file os.FileInfo, options directoryOptions) string {
	if !options.Flags.ShowGitStatus {
		return prefix
	}

	filePath := filepath.Join(options.Directory, file.Name())
	symbol := options.gitStatus.statusFor(filePath)
	if symbol == "" {
		symbol = " "
	}
	return prefix + colorGitStatusSymbol(symbol) + " "
}

func colorGitStatusSymbol(symbol string) string {
	switch symbol {
	case "M":
		return module.Yellow + symbol + module.Reset
	case "A":
		return module.Green + symbol + module.Reset
	case "D":
		return module.Red + symbol + module.Reset
	case "?":
		return module.Cyan + symbol + module.Reset
	case "!":
		return module.Grey + symbol + module.Reset
	default:
		return symbol
	}
}

func formatFileWithOptions(prefix string, file os.FileInfo, flags module.Flags) string {
	if flags.ShowLongFormat {
		return formatLongFileWithOptions(prefix, file, flags)
	}

	formatted := formatFileNameWithOptions(prefix, file, flags)
	if flags.ShowReadableSize {
		formatted = appendReadableSize(formatted, file.Size())
	}

	return formatted
}

func formatFileNameWithOptions(prefix string, file os.FileInfo, flags module.Flags) string {
	if flags.HideIcon {
		return printFilesAndFolderWithoutIcons(prefix, file)
	}

	return printWithIconAndPrefix(prefix, file)
}

func formatLongFileWithOptions(prefix string, file os.FileInfo, flags module.Flags) string {
	name := formatFileNameWithOptions("", file, flags)
	return fmt.Sprintf("%s%s  %10s  %s  %s", prefix, file.Mode().String(), formatLongSize(file.Size(), flags.ShowReadableSize), file.ModTime().Format(longDateFormat), name)
}

func formatLongSize(size int64, readable bool) string {
	if readable {
		return humanReadableSize(size)
	}
	if size < 0 {
		size = 0
	}

	return fmt.Sprintf("%d", size)
}

// processDirectory recursively processes a subdirectory
func processDirectory(file os.FileInfo, options directoryOptions, output io.Writer, childIndent string, isLastFolder bool, depth int) (fileCount, dirCount int, err error) {
	newOpts := options
	newOpts.Directory = filepath.Join(options.Directory, file.Name())
	return treeDirectory(newOpts, output, childIndent, isLastFolder, depth)
}

var errGetwd = errors.New("get current working directory")

// printNumberOfFilesAndDirectories returns number of files and directories
func printNumberOfFilesAndDirectories(output io.Writer, fileCount, dirCount int, flags module.Flags) error {
	if flags.ShowDirsOnly {
		_, err := fmt.Fprintf(output, "\n%d directories\n", dirCount)
		return err
	}

	_, err := fmt.Fprintf(output, "\n%d directories, %d files\n", dirCount, fileCount)
	return err
}
