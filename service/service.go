// Package service provides functionalities related to directory listing.
package service

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/oussamaM1/treels/module"
)

const (
	whiteSpaces              = "    " // Whitespaces
	boxUpAndRight            = "└── " // BOX DRAWINGS HEAVY UP AND RIGHT
	boxLightVertical         = "│   " // BOX DRAWINGS LIGHT VERTICAL
	boxLightVerticalAndRight = "├── " // BOX DRAWINGS LIGHT VERTICAL AND RIGHT
)

// Dispatcher func - executes function based on flags
func Dispatcher(options module.Options) {
	if options.Flags.ShowTreeView {
		TreeDirectory(options, "", true)
	} else {
		ListDirectory(options)
	}
}

// ListDirectory func - lists the content of the directory.
func ListDirectory(options module.Options) {
	CheckDefaultDirectory(&options.Directory)

	// Open and read the directory
	files, d, err := readDirectory(options.Directory)
	defer closeDirectory(d)
	if err != nil {
		log.Fatalf("Error reading directory: %s\n", err)
	}

	// Print files and directories
	fmt.Printf("Files and directories in %s:\n", options.Directory)
	for _, file := range files {
		if !isHidden(file.Name()) || options.Flags.ShowHidden {
			fmt.Println(file.Name())
		}
	}
}

// TreeDirectory func - displays a tree view of the directory.
func TreeDirectory(options module.Options, indent string, isLastFolder bool) {
	CheckDefaultDirectory(&options.Directory)

	// Open and read the directory
	files, d, err := readDirectory(options.Directory)
	defer closeDirectory(d)
	if err != nil {
		log.Fatalf("Error reading directory: %s\n", err)
	}

	// Sort files by name
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	// Print files and directories
	lastVisibleIndex := getLastVisibleIndex(files, options.Flags.ShowHidden)
	for i, file := range files {
		if !isHidden(file.Name()) || options.Flags.ShowHidden {
			var prefix, childIndentNext string
			if i == lastVisibleIndex && isLastFolder {
				prefix = indent + boxUpAndRight
				childIndentNext = indent + whiteSpaces
			} else if i == lastVisibleIndex {
				prefix = indent + boxUpAndRight
				childIndentNext = indent + whiteSpaces
			} else {
				prefix = indent + boxLightVerticalAndRight
				childIndentNext = indent + boxLightVertical
			}

			fmt.Println(prefix + file.Name())

			if file.IsDir() {
				newDirectory := filepath.Join(options.Directory, file.Name())
				newIsLastFolder := i == lastVisibleIndex && isLastFolder
				TreeDirectory(module.Options{Directory: newDirectory, Flags: module.Flags{ShowHidden: options.Flags.ShowHidden}}, childIndentNext, newIsLastFolder)
			}
		}
	}
}

func getLastVisibleIndex(files []os.FileInfo, showHidden bool) int {
	for i := len(files) - 1; i >= 0; i-- {
		if !isHidden(files[i].Name()) || showHidden {
			return i
		}
	}
	return -1
}
