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
	gitIgnore *gitIgnoreMatcher
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

	traversalOptions := directoryOptions{Options: options}
	if options.Flags.RespectGitIgnore {
		gitIgnore, err := newGitIgnoreMatcher(options.Directory)
		if err != nil {
			return err
		}
		traversalOptions.gitIgnore = gitIgnore
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

	// sort files by name
	sortSlice(files)

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

			formatted := formatFileWithOptions("", file, options.Flags)

			entries = append(entries, formatted)

			// Track max visible length for column width
			visibleLen := getVisibleLength(formatted)
			if visibleLen > maxLen {
				maxLen = visibleLen
			}
		}
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

	// Sort files by name
	sortSlice(files)

	// Print files and directories
	fc, dc, err := printFilesAndDirectoriesTreeFormat(files, options, output, indent, isLastFolder, depth)
	if err != nil {
		return 0, 0, err
	}
	return fc, dc, nil
}

func reachedMaxDepth(flags module.Flags, depth int) bool {
	return flags.LimitTreeDepth && depth >= flags.TreeDepth
}

func getLastVisibleIndex(files []os.FileInfo, options directoryOptions) int {
	for i := len(files) - 1; i >= 0; i-- {
		if shouldShowFile(files[i], options) {
			return i
		}
	}
	return -1
}

// printFilesAndDirectoriesTreeFormat - prints files and directories in tree format
func printFilesAndDirectoriesTreeFormat(files []os.FileInfo, options directoryOptions, output io.Writer, indent string, isLastFolder bool, depth int) (fileCount, dirCount int, err error) {
	lastVisibleIndex := getLastVisibleIndex(files, options)
	for i, file := range files {
		if !shouldShowFile(file, options) {
			continue
		}

		isLast := i == lastVisibleIndex
		prefix, childIndent := calculateIndent(indent, isLast)

		if err := printFileWithPrefix(output, prefix, file, options.Flags); err != nil {
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
	if isHidden(file.Name()) && !options.Flags.ShowHidden {
		return false
	}

	if options.Flags.ShowDirsOnly && !file.IsDir() {
		return false
	}

	if options.gitIgnore == nil {
		return true
	}

	filePath := filepath.Join(options.Directory, file.Name())
	return !options.gitIgnore.ignores(filePath, file.IsDir())
}

// calculateIndent returns the appropriate prefix and child indent strings
func calculateIndent(indent string, isLast bool) (prefix, childIndent string) {
	if isLast {
		return indent + boxUpAndRight, indent + whiteSpaces
	}
	return indent + boxLightVerticalAndRight, indent + boxLightVertical
}

// printFileWithPrefix prints the file with the given prefix and icon settings
func printFileWithPrefix(output io.Writer, prefix string, file os.FileInfo, flags module.Flags) error {
	_, err := fmt.Fprintln(output, formatFileWithOptions(prefix, file, flags))
	return err
}

func formatFileWithOptions(prefix string, file os.FileInfo, flags module.Flags) string {
	var formatted string
	if flags.HideIcon {
		formatted = printFilesAndFolderWithoutIcons(prefix, file)
	} else {
		formatted = printWithIconAndPrefix(prefix, file)
	}

	if flags.ShowReadableSize {
		formatted = appendReadableSize(formatted, file.Size())
	}

	return formatted
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
