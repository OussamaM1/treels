// Package service - service/util.go
package service

import (
	"fmt"
	"github.com/oussamaM1/treels/module"
	"log"
	"os"
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
	if file.IsDir() {
		format = fmt.Sprintf("%s%s%s%s %s%s", prefix, module.LightPurple, module.Bold, module.Folder, file.Name(), module.Reset)
	} else {
		format = fmt.Sprintf("%s%s %s", prefix, module.File, file.Name())
	}
	return format
}
