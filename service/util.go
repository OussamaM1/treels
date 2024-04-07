// Package service - service/util.go
package service

import (
	"fmt"
	"github.com/oussamaM1/treels/module"
	"log"
	"os"
	"strings"
)

// CheckDefaultDirectory func - returns current working directory if no directory is specified
func CheckDefaultDirectory(directory *string) {
	if *directory == "" {
		// Get the current working directory - default value
		var err error
		*directory, err = os.Getwd()
		if err != nil {
			log.Fatalf("Error getting current working directory: %s\n", err)
		}
	}
}

// readDirectory func - opens and reads the directory.
func readDirectory(directory string) ([]os.FileInfo, *os.File, error) {
	// Open the directory
	d, err := os.Open(directory)
	if err != nil {
		return nil, nil, err
	}
	// Read the directory contents
	files, err := d.Readdir(-1)
	if err != nil {
		return nil, nil, err
	}
	return files, d, nil
}

// closeDirectory func - closes the directory.
func closeDirectory(directory *os.File) {
	err := directory.Close()
	if err != nil {
		log.Fatalf("Error while closing directory: %s\n", err)
	}
}

// isHidden func - checks if the file name starts with a dot (hidden file).
func isHidden(name string) bool {
	return len(name) > 0 && name[0] == '.'
}

// printWithIconAndPrefix func - prints files and folder with icons and prefix
func printWithIconAndPrefix(prefix string, file os.FileInfo) string {
	var format string
	// Directory icon logic
	if file.IsDir() {
		switch {
		case strings.HasSuffix(file.Name(), module.Github), strings.HasSuffix(file.Name(), module.GitFolder):
			format = printIconFolders(prefix, file, module.GitFolder)
		default:
			format = printIconFolders(prefix, file, module.Folder)
		}
		return format
	}

	// File icon logic
	switch {
	case strings.HasSuffix(file.Name(), module.Go), strings.HasSuffix(file.Name(), module.Mod), strings.HasSuffix(file.Name(), module.Sum):
		format = printIconFiles(prefix, file, module.Go)
	case strings.HasSuffix(file.Name(), module.Md):
		format = printIconFiles(prefix, file, module.Md)
	case strings.HasSuffix(file.Name(), module.Gitignore):
		format = printIconFiles(prefix, file, module.Gitignore)
	default:
		// Default file icon
		format = printIconFiles(prefix, file, module.File)
	}
	return format
}

func printIconFiles(prefix string, file os.FileInfo, extension string) string {
	var format string
	switch extension {
	case module.Go, module.Mod, module.Sum:
		format = fmt.Sprintf("%s%s%s %s%s", prefix, module.LightBlue, module.GoLang, module.Reset, file.Name())
	case module.Md:
		format = fmt.Sprintf("%s%s%s %s%s", prefix, module.Grey, module.Markdown, module.Reset, file.Name())
	case module.Gitignore:
		format = fmt.Sprintf("%s%s%s %s%s", prefix, module.Orange, module.Git, module.Reset, file.Name())
	default:
		format = fmt.Sprintf("%s%s %s", prefix, module.File, file.Name())
	}
	return format
}

func printIconFolders(prefix string, file os.FileInfo, extension string) string {
	var format string
	switch extension {
	case module.Github, module.GitFolder:
		format = fmt.Sprintf("%s%s%s%s %s%s", prefix, module.Bold, module.Orange, module.Git, file.Name(), module.Reset)
	default:
		format = fmt.Sprintf("%s%s%s%s %s%s", prefix, module.Bold, module.Pink, module.Folder, file.Name(), module.Reset)
	}
	return format
}
