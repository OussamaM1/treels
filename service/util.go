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
			format = printIconFolders(prefix, file, module.FolderIcon)
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
	case strings.HasSuffix(file.Name(), module.Js):
		format = printIconFiles(prefix, file, module.Js)
	case strings.HasSuffix(file.Name(), module.Jsx):
		format = printIconFiles(prefix, file, module.Jsx)
	case strings.HasSuffix(file.Name(), module.Ts), strings.HasSuffix(file.Name(), module.Tsx):
		format = printIconFiles(prefix, file, module.Ts)
	case strings.HasSuffix(file.Name(), module.HTML), strings.HasSuffix(file.Name(), module.Htm):
		format = printIconFiles(prefix, file, module.HTML)
	case strings.HasSuffix(file.Name(), module.CSS), strings.HasSuffix(file.Name(), module.Scss), strings.HasSuffix(file.Name(), module.Sass):
		format = printIconFiles(prefix, file, module.CSS)
	case strings.HasSuffix(file.Name(), module.Vue):
		format = printIconFiles(prefix, file, module.Vue)
	case strings.HasSuffix(file.Name(), module.Py):
		format = printIconFiles(prefix, file, module.Py)
	case strings.HasSuffix(file.Name(), module.Rs):
		format = printIconFiles(prefix, file, module.Rs)
	case strings.HasSuffix(file.Name(), module.Rb), strings.HasSuffix(file.Name(), module.Rake), file.Name() == module.Gemfile:
		format = printIconFiles(prefix, file, module.Rb)
	case strings.HasSuffix(file.Name(), module.Php):
		format = printIconFiles(prefix, file, module.Php)
	case strings.HasSuffix(file.Name(), module.Swift):
		format = printIconFiles(prefix, file, module.Swift)
	case strings.HasSuffix(file.Name(), module.Kt), strings.HasSuffix(file.Name(), module.Kts):
		format = printIconFiles(prefix, file, module.Kt)
	case strings.HasSuffix(file.Name(), module.Cs), strings.HasSuffix(file.Name(), module.Csx):
		format = printIconFiles(prefix, file, module.Cs)
	case strings.HasSuffix(file.Name(), module.XML):
		format = printIconFiles(prefix, file, module.XML)
	case strings.HasSuffix(file.Name(), module.Sh), strings.HasSuffix(file.Name(), module.Bash), strings.HasSuffix(file.Name(), module.Zsh):
		format = printIconFiles(prefix, file, module.Sh)
	case file.Name() == module.Dockerfile, strings.HasSuffix(file.Name(), module.Dockerignore):
		format = printIconFiles(prefix, file, module.Dockerfile)
	case strings.HasSuffix(file.Name(), module.Conf), strings.HasSuffix(file.Name(), module.Cfg), strings.HasSuffix(file.Name(), module.Ini), strings.HasSuffix(file.Name(), module.Env):
		format = printIconFiles(prefix, file, module.Conf)
	case file.Name() == module.Makefile, strings.HasSuffix(file.Name(), module.Make):
		format = printIconFiles(prefix, file, module.Makefile)
	case file.Name() == module.PackageJSON:
		format = printIconFiles(prefix, file, module.PackageJSON)
	case strings.HasSuffix(file.Name(), module.Tf), strings.HasSuffix(file.Name(), module.Tfvars):
		format = printIconFiles(prefix, file, module.Tf)
	case strings.HasSuffix(file.Name(), module.Png), strings.HasSuffix(file.Name(), module.Jpg), strings.HasSuffix(file.Name(), module.Jpeg), strings.HasSuffix(file.Name(), module.Gif), strings.HasSuffix(file.Name(), module.Svg), strings.HasSuffix(file.Name(), module.Ico):
		format = printIconFiles(prefix, file, module.Png)
	case strings.HasSuffix(file.Name(), module.Mp4), strings.HasSuffix(file.Name(), module.Avi), strings.HasSuffix(file.Name(), module.Mov), strings.HasSuffix(file.Name(), module.Mkv):
		format = printIconFiles(prefix, file, module.Mp4)
	case strings.HasSuffix(file.Name(), module.Mp3), strings.HasSuffix(file.Name(), module.Wav), strings.HasSuffix(file.Name(), module.Flac):
		format = printIconFiles(prefix, file, module.Mp3)
	case strings.HasSuffix(file.Name(), module.Zip), strings.HasSuffix(file.Name(), module.Tar), strings.HasSuffix(file.Name(), module.Gz), strings.HasSuffix(file.Name(), module.Rar), strings.HasSuffix(file.Name(), module.SevenZ):
		format = printIconFiles(prefix, file, module.Zip)
	case strings.HasSuffix(file.Name(), module.Pdf):
		format = printIconFiles(prefix, file, module.Pdf)
	case strings.HasSuffix(file.Name(), module.Lock):
		format = printIconFiles(prefix, file, module.Lock)
	case strings.HasSuffix(file.Name(), module.Key), strings.HasSuffix(file.Name(), module.Pem), strings.HasSuffix(file.Name(), module.Crt), strings.HasSuffix(file.Name(), module.Pub):
		format = printIconFiles(prefix, file, module.Key)
	case strings.HasSuffix(file.Name(), module.Log):
		format = printIconFiles(prefix, file, module.Log)
	case strings.HasSuffix(file.Name(), module.Txt):
		format = printIconFiles(prefix, file, module.Txt)
	default:
		// Default file icon
		format = printIconFiles(prefix, file, module.FileIcon)
	}
	return format
}

// printIconFiles func - prints files with icons
func printIconFiles(prefix string, file os.FileInfo, extension string) string {
	var format string
	switch extension {
	case module.Go, module.Mod, module.Sum:
		format = fmt.Sprintf(IconFileFormat, prefix, module.LightBlue, module.GoLangIcon, module.Reset, file.Name())
	case module.Md:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Grey, module.MarkdownIcon, module.Reset, file.Name())
	case module.Gitignore:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Orange, module.GitIcon, module.Reset, file.Name())
	case module.JSON:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Yellow, module.JSONIcon, module.Reset, file.Name())
	case module.Yml, module.Yaml:
		format = fmt.Sprintf(IconFileFormat, prefix, module.LightGreen, module.YamlIcon, module.Reset, file.Name())
	case module.Pls, module.Plb, module.SQL:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Red, module.DatabaseIcon, module.Reset, file.Name())
	case module.Java, module.Class:
		format = fmt.Sprintf(IconFileFormat, prefix, module.LightRed, module.JavaLangIcon, module.Reset, file.Name())
	case module.Cpp:
		format = fmt.Sprintf(IconFileFormat, prefix, module.LightBlue, module.CppLangIcon, module.Reset, file.Name())
	case module.C:
		format = fmt.Sprintf(IconFileFormat, prefix, module.LightBlue, module.CLangIcon, module.Reset, file.Name())
	case module.Js:
		format = fmt.Sprintf(IconFileFormat, prefix, module.LightYellow, module.JavascriptLangIcon, module.Reset, file.Name())
	case module.Rs:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Orange, module.RustLangIcon, module.Reset, file.Name())
	case module.Py:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Blue, module.PythonLangIcon, module.Reset, file.Name())
	case module.Ts, module.Tsx:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Blue, module.TypeScriptIcon, module.Reset, file.Name())
	case module.HTML, module.Htm:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Orange, module.HTMLIcon, module.Reset, file.Name())
	case module.CSS, module.Scss, module.Sass:
		format = fmt.Sprintf(IconFileFormat, prefix, module.LightBlue, module.CSSIcon, module.Reset, file.Name())
	case module.Jsx:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Cyan, module.ReactIcon, module.Reset, file.Name())
	case module.Vue:
		format = fmt.Sprintf(IconFileFormat, prefix, module.LightGreen, module.VueIcon, module.Reset, file.Name())
	case module.Dockerfile, module.Dockerignore:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Cyan, module.DockerIcon, module.Reset, file.Name())
	case module.Sh, module.Bash, module.Zsh:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Green, module.ShellIcon, module.Reset, file.Name())
	case module.Rb, module.Rake, module.Gemfile:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Red, module.RubyIcon, module.Reset, file.Name())
	case module.Php:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Purple, module.PHPIcon, module.Reset, file.Name())
	case module.Swift:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Orange, module.SwiftIcon, module.Reset, file.Name())
	case module.Kt, module.Kts:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Purple, module.KotlinIcon, module.Reset, file.Name())
	case module.Cs, module.Csx:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Green, module.CSharpIcon, module.Reset, file.Name())
	case module.XML:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Orange, module.XMLIcon, module.Reset, file.Name())
	case module.Pdf:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Red, module.PDFIcon, module.Reset, file.Name())
	case module.Png, module.Jpg, module.Jpeg, module.Gif, module.Svg, module.Ico:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Magenta, module.ImageIcon, module.Reset, file.Name())
	case module.Mp4, module.Avi, module.Mov, module.Mkv:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Purple, module.VideoIcon, module.Reset, file.Name())
	case module.Mp3, module.Wav, module.Flac:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Cyan, module.AudioIcon, module.Reset, file.Name())
	case module.Zip, module.Tar, module.Gz, module.Rar, module.SevenZ:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Yellow, module.ArchiveIcon, module.Reset, file.Name())
	case module.Txt:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Grey, module.TextIcon, module.Reset, file.Name())
	case module.Conf, module.Cfg, module.Ini, module.Env:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Grey, module.ConfigIcon, module.Reset, file.Name())
	case module.Lock:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Yellow, module.LockIcon, module.Reset, file.Name())
	case module.Key, module.Pem, module.Crt, module.Pub:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Yellow, module.KeyIcon, module.Reset, file.Name())
	case module.Log:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Grey, module.LogIcon, module.Reset, file.Name())
	case module.Makefile, module.Make:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Orange, module.MakefileIcon, module.Reset, file.Name())
	case module.PackageJSON:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Red, module.NPMIcon, module.Reset, file.Name())
	case module.Tf, module.Tfvars:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Purple, module.TerraformIcon, module.Reset, file.Name())
	default:
		format = fmt.Sprintf("%s%s %s", prefix, module.FileIcon, file.Name())
	}
	return format
}

// printIconFolders func - prints folders with icons
func printIconFolders(prefix string, file os.FileInfo, extension string) string {
	var format string
	switch extension {
	case module.Github, module.GitFolder:
		format = fmt.Sprintf(IconFolderFormat, prefix, module.Bold, module.Orange, module.GitIcon, file.Name(), module.Reset)
	case module.IntellijFolder:
		format = fmt.Sprintf(IconFolderFormat, prefix, module.Bold, module.LightBlue, module.IntellijIcon, file.Name(), module.Reset)
	default:
		format = fmt.Sprintf(IconFolderFormat, prefix, module.Bold, module.Pink, module.FolderIcon, file.Name(), module.Reset)
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
