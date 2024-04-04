// Package service - service/ls.go
package service

import (
	"fmt"
	"github.com/oussamaM1/treels/module"
	"log"
	"os"
)

// Options struct - Contains configuration options for directory listing.
type Options struct {
	Directory string
	Flags     module.Flags
}

// ListDirectory func - List content of the directory
func ListDirectory(options Options) {
	checkDefaultDirectory(&options.Directory)
	// Open the directory
	d, err := os.Open(options.Directory)
	if err != nil {
		log.Fatalln("Error:", err)
	}
	defer func(d *os.File) {
		err := d.Close()
		if err != nil {
			log.Fatalln("Error while closing directory: ", err)
		}
	}(d)
	// Read the directory contents
	files, err := d.Readdir(-1)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	// Print files and directories
	fmt.Println("Files and directories in", options.Directory)
	for _, file := range files {
		if !isHidden(file.Name()) || options.Flags.ShowHidden {
			fmt.Println(file.Name())
		}
	}
}

// isHidden func - checks if the file name starts with a dot (hidden file)
func isHidden(name string) bool {
	return len(name) > 0 && name[0] == '.'
}

// checkDefaultDirectory func - returns current working directory if no directory is specified
func checkDefaultDirectory(directory *string) {
	if *directory == "" {
		// Get the current working directory - default value
		var err error
		*directory, err = os.Getwd()
		if err != nil {
			log.Fatalln("Error:", err)
		}
	}
}
