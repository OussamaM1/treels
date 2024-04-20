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
)

// Dispatcher func - executes function based on flags
func Dispatcher(options module.Options) {
	fmt.Println(".")
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
		if !isHidden(file.Name()) || options.Flags.ShowHidden {
			var prefix, childIndentNext string
			if i == lastVisibleIndex && isLastFolder {
				prefix = indent + boxUpAndRight
				childIndentNext = indent + whiteSpaces
			} else {
				prefix = indent + boxLightVerticalAndRight
				childIndentNext = indent + boxLightVertical
			}

			if !options.Flags.HideIcon {
				fmt.Println(printWithIconAndPrefix(prefix, file))
			} else {
				fmt.Println(printFilesAndFolderWithoutIcons(prefix, file))
			}

			if file.IsDir() {
				newDirectory := filepath.Join(options.Directory, file.Name())
				newIsLastFolder := i == lastVisibleIndex && isLastFolder
				treeDirectory(module.Options{Directory: newDirectory, Flags: module.Flags{ShowHidden: options.Flags.ShowHidden, HideIcon: options.Flags.HideIcon}}, childIndentNext, newIsLastFolder)
			}
		}
	}
}
