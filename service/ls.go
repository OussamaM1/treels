// Package service - service/ls.go
package service

import (
	"fmt"
	"log"
	"os"
)

// Options struct - Contains configuration options for directory listing.
type Options struct {
	Directory string
}

// ListDirectory func - List content of the directory
func ListDirectory(options Options) {
	if options.Directory == "" {
		// Get the current working directory - default value
		var err error
		options.Directory, err = os.Getwd()
		if err != nil {
			log.Fatalln("Error:", err)
		}
	}

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
	// todo: Add Option to show hidden files/directories
	fmt.Println("Files and directories in", options.Directory)
	for _, file := range files {
		if !isHidden(file.Name()) {
			fmt.Println(file.Name())
		}
	}
}

// isHidden func - checks if the file name starts with a dot (hidden file)
func isHidden(name string) bool {
	return len(name) > 0 && name[0] == '.'
}
