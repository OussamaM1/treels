// Package service - service/util.go
package service

import (
	"fmt"
	"github.com/oussamaM1/treels/module"
	"log"
	"os"
	"sort"
	"strings"
)

// Define constant for formatting
const (
	IconFileFormat   = "%s%s%s %s%s"
	IconFolderFormat = "%s%s%s%s %s%s"
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
		case strings.HasSuffix(file.Name(), module.IntellijFolder):
			format = printIconFolders(prefix, file, module.IntellijFolder)
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
	case strings.HasSuffix(file.Name(), module.Yaml), strings.HasSuffix(file.Name(), module.Yml):
		format = printIconFiles(prefix, file, module.Yml)
	case strings.HasSuffix(file.Name(), module.JSON):
		format = printIconFiles(prefix, file, module.JSON)
	case strings.HasSuffix(file.Name(), module.SQL), strings.HasSuffix(file.Name(), module.Pls), strings.HasSuffix(file.Name(), module.Plb):
		format = printIconFiles(prefix, file, module.SQL)
	case strings.HasSuffix(file.Name(), module.Java), strings.HasSuffix(file.Name(), module.Class):
		format = printIconFiles(prefix, file, module.Java)
	case strings.HasSuffix(file.Name(), module.Cpp):
		format = printIconFiles(prefix, file, module.Cpp)
	case strings.HasSuffix(file.Name(), module.C):
		format = printIconFiles(prefix, file, module.C)
	default:
		// Default file icon
		format = printIconFiles(prefix, file, module.File)
	}
	return format
}

// printIconFiles func - prints files with icons
func printIconFiles(prefix string, file os.FileInfo, extension string) string {
	var format string
	switch extension {
	case module.Go, module.Mod, module.Sum:
		format = fmt.Sprintf(IconFileFormat, prefix, module.LightBlue, module.GoLang, module.Reset, file.Name())
	case module.Md:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Grey, module.Markdown, module.Reset, file.Name())
	case module.Gitignore:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Orange, module.Git, module.Reset, file.Name())
	case module.JSON:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Yellow, module.JSONIcon, module.Reset, file.Name())
	case module.Yml, module.Yaml:
		format = fmt.Sprintf(IconFileFormat, prefix, module.LightGreen, module.YamlIcon, module.Reset, file.Name())
	case module.Pls, module.Plb, module.SQL:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Red, module.Database, module.Reset, file.Name())
	case module.Java, module.Class:
		format = fmt.Sprintf(IconFileFormat, prefix, module.LightRed, module.JavaLang, module.Reset, file.Name())
	case module.Cpp:
		format = fmt.Sprintf(IconFileFormat, prefix, module.LightBlue, module.CppLang, module.Reset, file.Name())
	case module.C:
		format = fmt.Sprintf(IconFileFormat, prefix, module.LightBlue, module.CLang, module.Reset, file.Name())
	default:
		format = fmt.Sprintf("%s%s %s", prefix, module.File, file.Name())
	}
	return format
}

// printIconFolders func - prints folders with icons
func printIconFolders(prefix string, file os.FileInfo, extension string) string {
	var format string
	switch extension {
	case module.Github, module.GitFolder:
		format = fmt.Sprintf(IconFolderFormat, prefix, module.Bold, module.Orange, module.Git, file.Name(), module.Reset)
	case module.IntellijFolder:
		format = fmt.Sprintf(IconFolderFormat, prefix, module.Bold, module.LightBlue, module.Intellij, file.Name(), module.Reset)
	default:
		format = fmt.Sprintf(IconFolderFormat, prefix, module.Bold, module.Pink, module.Folder, file.Name(), module.Reset)
	}
	return format
}

// printFilesAndFolderWithoutIcons func - prints Files/Folder without icons
func printFilesAndFolderWithoutIcons(prefix string, file os.FileInfo) string {
	var format string
	if file.IsDir() {
		format = fmt.Sprintf("%s%s%s%s%s", prefix, module.Bold, module.Pink, file.Name(), module.Reset)
	} else {
		format = fmt.Sprintf("%s%s", prefix, file.Name())
	}
	return format
}

// sortSlice func - sorts a slice of os.FileInfo objects alphabetically by file name.
// It modifies the original slice in place.
func sortSlice(files []os.FileInfo) {
	// Sort files by name
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})
}
