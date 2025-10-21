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
	case file.Name() == module.Readme, file.Name() == module.ReadmeMd:
		format = printIconFiles(prefix, file, module.Readme)
	case strings.HasSuffix(file.Name(), module.Gitignore):
		format = printIconFiles(prefix, file, module.Gitignore)
	case strings.HasSuffix(file.Name(), module.Yaml), strings.HasSuffix(file.Name(), module.Yml), strings.HasSuffix(file.Name(), module.YmlCI):
		format = printIconFiles(prefix, file, module.Yml)
	case strings.HasSuffix(file.Name(), module.JSON):
		format = printIconFiles(prefix, file, module.JSON)
	case file.Name() == module.TsConfig:
		format = printIconFiles(prefix, file, module.TsConfig)
	case strings.HasSuffix(file.Name(), module.SQL), strings.HasSuffix(file.Name(), module.Pls), strings.HasSuffix(file.Name(), module.Plb), strings.HasSuffix(file.Name(), module.Psql):
		format = printIconFiles(prefix, file, module.SQL)
	case strings.HasSuffix(file.Name(), module.Sqlite), strings.HasSuffix(file.Name(), module.Db):
		format = printIconFiles(prefix, file, module.Sqlite)
	case strings.HasSuffix(file.Name(), module.Java), strings.HasSuffix(file.Name(), module.Class):
		format = printIconFiles(prefix, file, module.Java)
	case strings.HasSuffix(file.Name(), module.Scala):
		format = printIconFiles(prefix, file, module.Scala)
	case strings.HasSuffix(file.Name(), module.Cpp):
		format = printIconFiles(prefix, file, module.Cpp)
	case strings.HasSuffix(file.Name(), module.C):
		format = printIconFiles(prefix, file, module.C)
	case strings.HasSuffix(file.Name(), module.Js):
		format = printIconFiles(prefix, file, module.Js)
	case file.Name() == module.WebpackConfig:
		format = printIconFiles(prefix, file, module.WebpackConfig)
	case file.Name() == module.ViteConfig:
		format = printIconFiles(prefix, file, module.ViteConfig)
	case strings.HasSuffix(file.Name(), module.Jsx):
		format = printIconFiles(prefix, file, module.Jsx)
	case strings.HasSuffix(file.Name(), module.Ts), strings.HasSuffix(file.Name(), module.Tsx):
		format = printIconFiles(prefix, file, module.Ts)
	case strings.HasSuffix(file.Name(), module.HTML), strings.HasSuffix(file.Name(), module.Htm):
		format = printIconFiles(prefix, file, module.HTML)
	case strings.HasSuffix(file.Name(), module.CSS), strings.HasSuffix(file.Name(), module.Scss), strings.HasSuffix(file.Name(), module.Sass):
		format = printIconFiles(prefix, file, module.CSS)
	case strings.HasSuffix(file.Name(), module.Less):
		format = printIconFiles(prefix, file, module.Less)
	case strings.HasSuffix(file.Name(), module.Vue):
		format = printIconFiles(prefix, file, module.Vue)
	case strings.HasSuffix(file.Name(), module.Svelte):
		format = printIconFiles(prefix, file, module.Svelte)
	case strings.HasSuffix(file.Name(), module.Py):
		format = printIconFiles(prefix, file, module.Py)
	case file.Name() == module.Requirements:
		format = printIconFiles(prefix, file, module.Requirements)
	case strings.HasSuffix(file.Name(), module.Pyproject), strings.HasSuffix(file.Name(), module.PoetryLock):
		format = printIconFiles(prefix, file, module.Pyproject)
	case strings.HasSuffix(file.Name(), module.Rs):
		format = printIconFiles(prefix, file, module.Rs)
	case file.Name() == module.CargoToml, file.Name() == module.CargoLock:
		format = printIconFiles(prefix, file, module.CargoToml)
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
	case strings.HasSuffix(file.Name(), module.Dart):
		format = printIconFiles(prefix, file, module.Dart)
	case strings.HasSuffix(file.Name(), module.Ex), strings.HasSuffix(file.Name(), module.Exs):
		format = printIconFiles(prefix, file, module.Ex)
	case strings.HasSuffix(file.Name(), module.Hs):
		format = printIconFiles(prefix, file, module.Hs)
	case strings.HasSuffix(file.Name(), module.Clj):
		format = printIconFiles(prefix, file, module.Clj)
	case strings.HasSuffix(file.Name(), module.R), strings.HasSuffix(file.Name(), module.Rmd):
		format = printIconFiles(prefix, file, module.R)
	case strings.HasSuffix(file.Name(), module.Lua):
		format = printIconFiles(prefix, file, module.Lua)
	case strings.HasSuffix(file.Name(), module.Pl), strings.HasSuffix(file.Name(), module.Pm):
		format = printIconFiles(prefix, file, module.Pl)
	case strings.HasSuffix(file.Name(), module.XML):
		format = printIconFiles(prefix, file, module.XML)
	case file.Name() == module.Pom:
		format = printIconFiles(prefix, file, module.Pom)
	case strings.HasSuffix(file.Name(), module.Graphql), strings.HasSuffix(file.Name(), module.Gql):
		format = printIconFiles(prefix, file, module.Graphql)
	case strings.HasSuffix(file.Name(), module.Prisma):
		format = printIconFiles(prefix, file, module.Prisma)
	case strings.HasSuffix(file.Name(), module.Proto):
		format = printIconFiles(prefix, file, module.Proto)
	case strings.HasSuffix(file.Name(), module.Wasm):
		format = printIconFiles(prefix, file, module.Wasm)
	case strings.HasSuffix(file.Name(), module.Sh), strings.HasSuffix(file.Name(), module.Bash), strings.HasSuffix(file.Name(), module.Zsh):
		format = printIconFiles(prefix, file, module.Sh)
	case file.Name() == module.Dockerfile, strings.HasSuffix(file.Name(), module.Dockerignore):
		format = printIconFiles(prefix, file, module.Dockerfile)
	case strings.HasSuffix(file.Name(), module.Conf), strings.HasSuffix(file.Name(), module.Cfg), strings.HasSuffix(file.Name(), module.Ini), strings.HasSuffix(file.Name(), module.Env):
		format = printIconFiles(prefix, file, module.Conf)
	case strings.HasSuffix(file.Name(), module.Toml):
		format = printIconFiles(prefix, file, module.Toml)
	case strings.HasSuffix(file.Name(), module.Editorconfig):
		format = printIconFiles(prefix, file, module.Editorconfig)
	case strings.HasSuffix(file.Name(), module.Eslintrc), strings.HasSuffix(file.Name(), module.EslintrcJSON):
		format = printIconFiles(prefix, file, module.Eslintrc)
	case strings.HasSuffix(file.Name(), module.Prettierrc), strings.HasSuffix(file.Name(), module.Prettierignore):
		format = printIconFiles(prefix, file, module.Prettierrc)
	case strings.HasSuffix(file.Name(), module.Babelrc):
		format = printIconFiles(prefix, file, module.Babelrc)
	case file.Name() == module.Makefile, strings.HasSuffix(file.Name(), module.Make):
		format = printIconFiles(prefix, file, module.Makefile)
	case file.Name() == module.CMakeLists, strings.HasSuffix(file.Name(), module.Cmake):
		format = printIconFiles(prefix, file, module.CMakeLists)
	case file.Name() == module.PackageJSON:
		format = printIconFiles(prefix, file, module.PackageJSON)
	case strings.HasSuffix(file.Name(), module.Gradle), strings.HasSuffix(file.Name(), module.GradleKts):
		format = printIconFiles(prefix, file, module.Gradle)
	case file.Name() == module.Jenkinsfile:
		format = printIconFiles(prefix, file, module.Jenkinsfile)
	case file.Name() == module.Vagrantfile:
		format = printIconFiles(prefix, file, module.Vagrantfile)
	case file.Name() == module.Procfile:
		format = printIconFiles(prefix, file, module.Procfile)
	case strings.HasSuffix(file.Name(), module.Tf), strings.HasSuffix(file.Name(), module.Tfvars):
		format = printIconFiles(prefix, file, module.Tf)
	case strings.HasSuffix(file.Name(), module.Nix):
		format = printIconFiles(prefix, file, module.Nix)
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
	case strings.HasSuffix(file.Name(), module.Doc), strings.HasSuffix(file.Name(), module.Docx):
		format = printIconFiles(prefix, file, module.Doc)
	case strings.HasSuffix(file.Name(), module.Xls), strings.HasSuffix(file.Name(), module.Xlsx):
		format = printIconFiles(prefix, file, module.Xls)
	case strings.HasSuffix(file.Name(), module.Ppt), strings.HasSuffix(file.Name(), module.Pptx):
		format = printIconFiles(prefix, file, module.Ppt)
	case strings.HasSuffix(file.Name(), module.Ttf), strings.HasSuffix(file.Name(), module.Otf), strings.HasSuffix(file.Name(), module.Woff), strings.HasSuffix(file.Name(), module.Woff2):
		format = printIconFiles(prefix, file, module.Ttf)
	case strings.HasSuffix(file.Name(), module.Exe), strings.HasSuffix(file.Name(), module.Dll), strings.HasSuffix(file.Name(), module.So), strings.HasSuffix(file.Name(), module.Dylib):
		format = printIconFiles(prefix, file, module.Exe)
	case file.Name() == module.License, file.Name() == module.LicenseMd:
		format = printIconFiles(prefix, file, module.License)
	case strings.HasSuffix(file.Name(), module.Lock):
		format = printIconFiles(prefix, file, module.Lock)
	case strings.HasSuffix(file.Name(), module.Key), strings.HasSuffix(file.Name(), module.Pem), strings.HasSuffix(file.Name(), module.Crt), strings.HasSuffix(file.Name(), module.Pub), strings.HasSuffix(file.Name(), module.Cer), strings.HasSuffix(file.Name(), module.P12):
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
	case module.Readme, module.ReadmeMd:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Cyan, module.ReadmeIcon, module.Reset, file.Name())
	case module.Gitignore:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Orange, module.GitIcon, module.Reset, file.Name())
	case module.JSON:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Yellow, module.JSONIcon, module.Reset, file.Name())
	case module.TsConfig:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Blue, module.TypeScriptIcon, module.Reset, file.Name())
	case module.Yml, module.Yaml, module.YmlCI:
		format = fmt.Sprintf(IconFileFormat, prefix, module.LightGreen, module.YamlIcon, module.Reset, file.Name())
	case module.Pls, module.Plb, module.SQL, module.Psql:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Red, module.DatabaseIcon, module.Reset, file.Name())
	case module.Sqlite, module.Db:
		format = fmt.Sprintf(IconFileFormat, prefix, module.LightBlue, module.SQLiteIcon, module.Reset, file.Name())
	case module.Java, module.Class:
		format = fmt.Sprintf(IconFileFormat, prefix, module.LightRed, module.JavaLangIcon, module.Reset, file.Name())
	case module.Scala:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Red, module.ScalaIcon, module.Reset, file.Name())
	case module.Cpp:
		format = fmt.Sprintf(IconFileFormat, prefix, module.LightBlue, module.CppLangIcon, module.Reset, file.Name())
	case module.C:
		format = fmt.Sprintf(IconFileFormat, prefix, module.LightBlue, module.CLangIcon, module.Reset, file.Name())
	case module.Js:
		format = fmt.Sprintf(IconFileFormat, prefix, module.LightYellow, module.JavascriptLangIcon, module.Reset, file.Name())
	case module.WebpackConfig:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Blue, module.WebpackIcon, module.Reset, file.Name())
	case module.ViteConfig:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Purple, module.ViteIcon, module.Reset, file.Name())
	case module.Rs:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Orange, module.RustLangIcon, module.Reset, file.Name())
	case module.CargoToml, module.CargoLock:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Orange, module.CargoIcon, module.Reset, file.Name())
	case module.Py:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Blue, module.PythonLangIcon, module.Reset, file.Name())
	case module.Requirements:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Blue, module.PipIcon, module.Reset, file.Name())
	case module.Pyproject, module.PoetryLock:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Blue, module.PoetryIcon, module.Reset, file.Name())
	case module.Ts, module.Tsx:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Blue, module.TypeScriptIcon, module.Reset, file.Name())
	case module.HTML, module.Htm:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Orange, module.HTMLIcon, module.Reset, file.Name())
	case module.CSS, module.Scss, module.Sass:
		format = fmt.Sprintf(IconFileFormat, prefix, module.LightBlue, module.CSSIcon, module.Reset, file.Name())
	case module.Less:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Blue, module.LessIcon, module.Reset, file.Name())
	case module.Jsx:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Cyan, module.ReactIcon, module.Reset, file.Name())
	case module.Vue:
		format = fmt.Sprintf(IconFileFormat, prefix, module.LightGreen, module.VueIcon, module.Reset, file.Name())
	case module.Svelte:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Orange, module.SvelteIcon, module.Reset, file.Name())
	case module.Dart:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Cyan, module.DartIcon, module.Reset, file.Name())
	case module.Ex, module.Exs:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Purple, module.ElixirIcon, module.Reset, file.Name())
	case module.Hs:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Purple, module.HaskellIcon, module.Reset, file.Name())
	case module.Clj:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Green, module.ClojureIcon, module.Reset, file.Name())
	case module.R, module.Rmd:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Blue, module.RIcon, module.Reset, file.Name())
	case module.Lua:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Blue, module.LuaIcon, module.Reset, file.Name())
	case module.Pl, module.Pm:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Blue, module.PerlIcon, module.Reset, file.Name())
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
	case module.Pom:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Red, module.MavenIcon, module.Reset, file.Name())
	case module.Graphql, module.Gql:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Magenta, module.GraphQLIcon, module.Reset, file.Name())
	case module.Prisma:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Blue, module.PrismaIcon, module.Reset, file.Name())
	case module.Proto:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Blue, module.ProtoIcon, module.Reset, file.Name())
	case module.Wasm:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Purple, module.WasmIcon, module.Reset, file.Name())
	case module.Pdf:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Red, module.PDFIcon, module.Reset, file.Name())
	case module.Doc, module.Docx:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Blue, module.WordIcon, module.Reset, file.Name())
	case module.Xls, module.Xlsx:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Green, module.ExcelIcon, module.Reset, file.Name())
	case module.Ppt, module.Pptx:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Orange, module.PowerPointIcon, module.Reset, file.Name())
	case module.Ttf, module.Otf, module.Woff, module.Woff2:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Grey, module.FontIcon, module.Reset, file.Name())
	case module.Exe, module.Dll, module.So, module.Dylib:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Red, module.BinaryIcon, module.Reset, file.Name())
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
	case module.Toml:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Grey, module.TomlIcon, module.Reset, file.Name())
	case module.Editorconfig:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Grey, module.EditorConfigIcon, module.Reset, file.Name())
	case module.Eslintrc, module.EslintrcJSON:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Purple, module.ESLintIcon, module.Reset, file.Name())
	case module.Prettierrc, module.Prettierignore:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Grey, module.PrettierIcon, module.Reset, file.Name())
	case module.Babelrc:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Yellow, module.BabelIcon, module.Reset, file.Name())
	case module.License, module.LicenseMd:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Yellow, module.LicenseIcon, module.Reset, file.Name())
	case module.Lock:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Yellow, module.LockIcon, module.Reset, file.Name())
	case module.Key, module.Pem, module.Crt, module.Pub, module.Cer, module.P12:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Yellow, module.CertificateIcon, module.Reset, file.Name())
	case module.Log:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Grey, module.LogIcon, module.Reset, file.Name())
	case module.Makefile, module.Make:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Orange, module.MakefileIcon, module.Reset, file.Name())
	case module.CMakeLists, module.Cmake:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Red, module.CMakeIcon, module.Reset, file.Name())
	case module.PackageJSON:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Red, module.NPMIcon, module.Reset, file.Name())
	case module.Gradle, module.GradleKts:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Green, module.GradleIcon, module.Reset, file.Name())
	case module.Jenkinsfile:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Red, module.CIIcon, module.Reset, file.Name())
	case module.Vagrantfile:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Blue, module.VagrantIcon, module.Reset, file.Name())
	case module.Procfile:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Purple, module.CIIcon, module.Reset, file.Name())
	case module.Tf, module.Tfvars:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Purple, module.TerraformIcon, module.Reset, file.Name())
	case module.Nix:
		format = fmt.Sprintf(IconFileFormat, prefix, module.Blue, module.NixIcon, module.Reset, file.Name())
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
