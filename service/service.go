// Package service provides functionalities related to directory listing.
package service

import (
	"fmt"
	"github.com/oussamaM1/treels/module"
	"log"
	"os"
	"path/filepath"
)

const (
	whiteSpaces              = "    " // Whitespaces
	boxUpAndRight            = "└── " // BOX DRAWINGS HEAVY UP AND RIGHT
	boxLightVertical         = "│   " // BOX DRAWINGS LIGHT VERTICAL
	boxLightVerticalAndRight = "├── " // BOX DRAWINGS LIGHT VERTICAL AND RIGHT
	dot                      = "."
)

// Dispatcher func - executes function based on flags
func Dispatcher(options module.Options) {
	fmt.Println(dot)
	if options.Flags.ShowTreeView {
		treeDirectory(options, "", true)
	} else {
		listDirectory(options)
	}
}

// listDirectory func - lists the content of the directory.
func listDirectory(options module.Options) {
	CheckDefaultDirectory(&options.Directory)

	// Open and read the directory
	files, d, err := readDirectory(options.Directory)
	defer closeDirectory(d)
	if err != nil {
		log.Fatalf("Error reading directory: %s\n", err)
	}

	// sort files by name
	sortSlice(files)

	// Print files and directories
	for _, file := range files {
		if !isHidden(file.Name()) || options.Flags.ShowHidden {
			if !options.Flags.HideIcon {
				fmt.Println(printWithIconAndPrefix("", file))
			} else {
				fmt.Println(printFilesAndFolderWithoutIcons("", file))
			}
		}
	}
}

// treeDirectory func - displays a tree view of the directory.
func treeDirectory(options module.Options, indent string, isLastFolder bool) {
	CheckDefaultDirectory(&options.Directory)

	// Open and read the directory
	files, d, err := readDirectory(options.Directory)
	defer closeDirectory(d)
	if err != nil {
		log.Fatalf("Error reading directory: %s\n", err)
	}

	// Sort files by name
	sortSlice(files)

	// Print files and directories
	printFilesAndDirectoriesTreeFormat(files, options, indent, isLastFolder)
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
func printFilesAndDirectoriesTreeFormat(files []os.FileInfo, options module.Options, indent string, isLastFolder bool) {
	lastVisibleIndex := getLastVisibleIndex(files, options.Flags.ShowHidden)
	for i, file := range files {
		if !shouldShowFile(file, options.Flags.ShowHidden) {
			continue
		}

		isLast := i == lastVisibleIndex
		prefix, childIndent := calculateIndent(indent, isLast)

		printFileWithPrefix(prefix, file, options.Flags.HideIcon)

		if file.IsDir() {
			processDirectory(file, options, childIndent, isLast && isLastFolder)
		}
	}
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
func printFileWithPrefix(prefix string, file os.FileInfo, hideIcon bool) {
	if hideIcon {
		fmt.Println(printFilesAndFolderWithoutIcons(prefix, file))
	} else {
		fmt.Println(printWithIconAndPrefix(prefix, file))
	}
}

// processDirectory recursively processes a subdirectory
func processDirectory(file os.FileInfo, options module.Options, childIndent string, isLastFolder bool) {
	newOpts := options
	newOpts.Directory = filepath.Join(options.Directory, file.Name())
	treeDirectory(newOpts, childIndent, isLastFolder)
}
