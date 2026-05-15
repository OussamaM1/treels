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

	var fileCount, dirCount int
	if _, err := fmt.Fprintln(output, dot); err != nil {
		return err
	}
	if options.Flags.ShowTreeView {
		var err error
		fileCount, dirCount, err = treeDirectory(options, output, "", true)
		if err != nil {
			return err
		}
	} else {
		var err error
		fileCount, dirCount, err = listDirectory(options, output)
		if err != nil {
			return err
		}
	}
	return printNumberOfFilesAndDirectories(output, fileCount, dirCount)
}

// listDirectory func - lists the content of the directory.
func listDirectory(options module.Options, output io.Writer) (fileCount, dirCount int, err error) {
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
		if !isHidden(file.Name()) || options.Flags.ShowHidden {
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
func treeDirectory(options module.Options, output io.Writer, indent string, isLastFolder bool) (fileCount, dirCount int, err error) {
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
	fc, dc, err := printFilesAndDirectoriesTreeFormat(files, options, output, indent, isLastFolder)
	if err != nil {
		return 0, 0, err
	}
	return fc, dc, nil
}

func getLastVisibleIndex(files []os.FileInfo, showHidden bool) int {
	for i := len(files) - 1; i >= 0; i-- {
		if !isHidden(files[i].Name()) || showHidden {
			return i
		}
	}
	return -1
}

// printFilesAndDirectoriesTreeFormat - prints files and directories in tree format
func printFilesAndDirectoriesTreeFormat(files []os.FileInfo, options module.Options, output io.Writer, indent string, isLastFolder bool) (fileCount, dirCount int, err error) {
	lastVisibleIndex := getLastVisibleIndex(files, options.Flags.ShowHidden)
	for i, file := range files {
		if !shouldShowFile(file, options.Flags.ShowHidden) {
			continue
		}

		isLast := i == lastVisibleIndex
		prefix, childIndent := calculateIndent(indent, isLast)

		if err := printFileWithPrefix(output, prefix, file, options.Flags); err != nil {
			return 0, 0, err
		}

		if file.IsDir() {
			dirCount++
			fc, dc, err := processDirectory(file, options, output, childIndent, isLast && isLastFolder)
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
func shouldShowFile(file os.FileInfo, showHidden bool) bool {
	return !isHidden(file.Name()) || showHidden
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
func processDirectory(file os.FileInfo, options module.Options, output io.Writer, childIndent string, isLastFolder bool) (fileCount, dirCount int, err error) {
	newOpts := options
	newOpts.Directory = filepath.Join(options.Directory, file.Name())
	return treeDirectory(newOpts, output, childIndent, isLastFolder)
}

var errGetwd = errors.New("get current working directory")

// printNumberOfFilesAndDirectories returns number of files and directories
func printNumberOfFilesAndDirectories(output io.Writer, fileCount, dirCount int) error {
	_, err := fmt.Fprintf(output, "\n%d directories, %d files\n", dirCount, fileCount)
	return err
}
