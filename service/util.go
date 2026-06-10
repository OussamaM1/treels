// Package service - service/util.go
package service

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/oussamaM1/treels/module"
	"golang.org/x/term"
)

// Define constant for formatting
const (
	IconFileFormat   = "%s%s%s %s%s"
	IconFolderFormat = "%s%s%s%s %s%s"
	longDateFormat   = "2006-01-02"
)

var (
	getWorkingDirectory = os.Getwd
	terminalSize        = term.GetSize
)

type iconStyle struct {
	color string
	icon  string
}

var defaultFileIconStyle = iconStyle{icon: module.FileIcon}

var defaultFolderIconStyle = iconStyle{
	color: module.Pink,
	icon:  module.FolderIcon,
}

var folderIconStyles = map[string]iconStyle{
	module.Github: {
		color: module.Orange,
		icon:  module.GitIcon,
	},
	module.GitFolder: {
		color: module.Orange,
		icon:  module.GitIcon,
	},
	module.IntellijFolder: {
		color: module.LightBlue,
		icon:  module.IntellijIcon,
	},
}

var filenameIconStyles = map[string]iconStyle{
	module.Readme: {
		color: module.Cyan,
		icon:  module.ReadmeIcon,
	},
	module.ReadmeMd: {
		color: module.Cyan,
		icon:  module.ReadmeIcon,
	},
	module.Gitignore: {
		color: module.Orange,
		icon:  module.GitIcon,
	},
	module.TsConfig: {
		color: module.Blue,
		icon:  module.TypeScriptIcon,
	},
	module.WebpackConfig: {
		color: module.Blue,
		icon:  module.WebpackIcon,
	},
	module.ViteConfig: {
		color: module.Purple,
		icon:  module.ViteIcon,
	},
	module.Requirements: {
		color: module.Blue,
		icon:  module.PipIcon,
	},
	module.CargoToml: {
		color: module.Orange,
		icon:  module.CargoIcon,
	},
	module.CargoLock: {
		color: module.Orange,
		icon:  module.CargoIcon,
	},
	module.Gemfile: {
		color: module.Red,
		icon:  module.RubyIcon,
	},
	module.Dockerfile: {
		color: module.Cyan,
		icon:  module.DockerIcon,
	},
	module.Pom: {
		color: module.Red,
		icon:  module.MavenIcon,
	},
	module.Makefile: {
		color: module.Orange,
		icon:  module.MakefileIcon,
	},
	module.CMakeLists: {
		color: module.Red,
		icon:  module.CMakeIcon,
	},
	module.PackageJSON: {
		color: module.Red,
		icon:  module.NPMIcon,
	},
	module.Jenkinsfile: {
		color: module.Red,
		icon:  module.CIIcon,
	},
	module.Vagrantfile: {
		color: module.Blue,
		icon:  module.VagrantIcon,
	},
	module.Procfile: {
		color: module.Purple,
		icon:  module.CIIcon,
	},
	module.License: {
		color: module.Yellow,
		icon:  module.LicenseIcon,
	},
	module.LicenseMd: {
		color: module.Yellow,
		icon:  module.LicenseIcon,
	},
}

var extensionIconStyles = map[string]iconStyle{
	module.Go: {
		color: module.LightBlue,
		icon:  module.GoLangIcon,
	},
	module.Mod: {
		color: module.LightBlue,
		icon:  module.GoLangIcon,
	},
	module.Sum: {
		color: module.LightBlue,
		icon:  module.GoLangIcon,
	},
	module.Md: {
		color: module.Grey,
		icon:  module.MarkdownIcon,
	},
	module.Yaml: {
		color: module.LightGreen,
		icon:  module.YamlIcon,
	},
	module.Yml: {
		color: module.LightGreen,
		icon:  module.YamlIcon,
	},
	module.YmlCI: {
		color: module.LightGreen,
		icon:  module.YamlIcon,
	},
	module.JSON: {
		color: module.Yellow,
		icon:  module.JSONIcon,
	},
	module.SQL: {
		color: module.Red,
		icon:  module.DatabaseIcon,
	},
	module.Pls: {
		color: module.Red,
		icon:  module.DatabaseIcon,
	},
	module.Plb: {
		color: module.Red,
		icon:  module.DatabaseIcon,
	},
	module.Psql: {
		color: module.Red,
		icon:  module.DatabaseIcon,
	},
	module.Sqlite: {
		color: module.LightBlue,
		icon:  module.SQLiteIcon,
	},
	module.Db: {
		color: module.LightBlue,
		icon:  module.SQLiteIcon,
	},
	module.Java: {
		color: module.LightRed,
		icon:  module.JavaLangIcon,
	},
	module.Class: {
		color: module.LightRed,
		icon:  module.JavaLangIcon,
	},
	module.Scala: {
		color: module.Red,
		icon:  module.ScalaIcon,
	},
	module.Cpp: {
		color: module.LightBlue,
		icon:  module.CppLangIcon,
	},
	module.C: {
		color: module.LightBlue,
		icon:  module.CLangIcon,
	},
	module.Js: {
		color: module.LightYellow,
		icon:  module.JavascriptLangIcon,
	},
	module.Jsx: {
		color: module.Cyan,
		icon:  module.ReactIcon,
	},
	module.Ts: {
		color: module.Blue,
		icon:  module.TypeScriptIcon,
	},
	module.Tsx: {
		color: module.Blue,
		icon:  module.TypeScriptIcon,
	},
	module.HTML: {
		color: module.Orange,
		icon:  module.HTMLIcon,
	},
	module.Htm: {
		color: module.Orange,
		icon:  module.HTMLIcon,
	},
	module.CSS: {
		color: module.LightBlue,
		icon:  module.CSSIcon,
	},
	module.Scss: {
		color: module.LightBlue,
		icon:  module.CSSIcon,
	},
	module.Sass: {
		color: module.LightBlue,
		icon:  module.CSSIcon,
	},
	module.Less: {
		color: module.Blue,
		icon:  module.LessIcon,
	},
	module.Vue: {
		color: module.LightGreen,
		icon:  module.VueIcon,
	},
	module.Svelte: {
		color: module.Orange,
		icon:  module.SvelteIcon,
	},
	module.Py: {
		color: module.Blue,
		icon:  module.PythonLangIcon,
	},
	module.Pyproject: {
		color: module.Blue,
		icon:  module.PoetryIcon,
	},
	module.PoetryLock: {
		color: module.Blue,
		icon:  module.PoetryIcon,
	},
	module.Rs: {
		color: module.Orange,
		icon:  module.RustLangIcon,
	},
	module.Rb: {
		color: module.Red,
		icon:  module.RubyIcon,
	},
	module.Rake: {
		color: module.Red,
		icon:  module.RubyIcon,
	},
	module.Php: {
		color: module.Purple,
		icon:  module.PHPIcon,
	},
	module.Swift: {
		color: module.Orange,
		icon:  module.SwiftIcon,
	},
	module.Kt: {
		color: module.Purple,
		icon:  module.KotlinIcon,
	},
	module.Kts: {
		color: module.Purple,
		icon:  module.KotlinIcon,
	},
	module.Cs: {
		color: module.Green,
		icon:  module.CSharpIcon,
	},
	module.Csx: {
		color: module.Green,
		icon:  module.CSharpIcon,
	},
	module.Dart: {
		color: module.Cyan,
		icon:  module.DartIcon,
	},
	module.Ex: {
		color: module.Purple,
		icon:  module.ElixirIcon,
	},
	module.Exs: {
		color: module.Purple,
		icon:  module.ElixirIcon,
	},
	module.Hs: {
		color: module.Purple,
		icon:  module.HaskellIcon,
	},
	module.Clj: {
		color: module.Green,
		icon:  module.ClojureIcon,
	},
	module.R: {
		color: module.Blue,
		icon:  module.RIcon,
	},
	module.Rmd: {
		color: module.Blue,
		icon:  module.RIcon,
	},
	module.Lua: {
		color: module.Blue,
		icon:  module.LuaIcon,
	},
	module.Pl: {
		color: module.Blue,
		icon:  module.PerlIcon,
	},
	module.Pm: {
		color: module.Blue,
		icon:  module.PerlIcon,
	},
	module.XML: {
		color: module.Orange,
		icon:  module.XMLIcon,
	},
	module.Graphql: {
		color: module.Magenta,
		icon:  module.GraphQLIcon,
	},
	module.Gql: {
		color: module.Magenta,
		icon:  module.GraphQLIcon,
	},
	module.Prisma: {
		color: module.Blue,
		icon:  module.PrismaIcon,
	},
	module.Proto: {
		color: module.Blue,
		icon:  module.ProtoIcon,
	},
	module.Wasm: {
		color: module.Purple,
		icon:  module.WasmIcon,
	},
	module.Sh: {
		color: module.Green,
		icon:  module.ShellIcon,
	},
	module.Bash: {
		color: module.Green,
		icon:  module.ShellIcon,
	},
	module.Zsh: {
		color: module.Green,
		icon:  module.ShellIcon,
	},
	module.Dockerignore: {
		color: module.Cyan,
		icon:  module.DockerIcon,
	},
	module.Conf: {
		color: module.Grey,
		icon:  module.ConfigIcon,
	},
	module.Cfg: {
		color: module.Grey,
		icon:  module.ConfigIcon,
	},
	module.Ini: {
		color: module.Grey,
		icon:  module.ConfigIcon,
	},
	module.Env: {
		color: module.Grey,
		icon:  module.ConfigIcon,
	},
	module.Toml: {
		color: module.Grey,
		icon:  module.TomlIcon,
	},
	module.Editorconfig: {
		color: module.Grey,
		icon:  module.EditorConfigIcon,
	},
	module.Eslintrc: {
		color: module.Purple,
		icon:  module.ESLintIcon,
	},
	module.EslintrcJSON: {
		color: module.Purple,
		icon:  module.ESLintIcon,
	},
	module.Prettierrc: {
		color: module.Grey,
		icon:  module.PrettierIcon,
	},
	module.Prettierignore: {
		color: module.Grey,
		icon:  module.PrettierIcon,
	},
	module.Babelrc: {
		color: module.Yellow,
		icon:  module.BabelIcon,
	},
	module.Make: {
		color: module.Orange,
		icon:  module.MakefileIcon,
	},
	module.Cmake: {
		color: module.Red,
		icon:  module.CMakeIcon,
	},
	module.Gradle: {
		color: module.Green,
		icon:  module.GradleIcon,
	},
	module.GradleKts: {
		color: module.Green,
		icon:  module.GradleIcon,
	},
	module.Tf: {
		color: module.Purple,
		icon:  module.TerraformIcon,
	},
	module.Tfvars: {
		color: module.Purple,
		icon:  module.TerraformIcon,
	},
	module.Nix: {
		color: module.Blue,
		icon:  module.NixIcon,
	},
	module.Png: {
		color: module.Magenta,
		icon:  module.ImageIcon,
	},
	module.Jpg: {
		color: module.Magenta,
		icon:  module.ImageIcon,
	},
	module.Jpeg: {
		color: module.Magenta,
		icon:  module.ImageIcon,
	},
	module.Gif: {
		color: module.Magenta,
		icon:  module.ImageIcon,
	},
	module.Svg: {
		color: module.Magenta,
		icon:  module.ImageIcon,
	},
	module.Ico: {
		color: module.Magenta,
		icon:  module.ImageIcon,
	},
	module.Mp4: {
		color: module.Purple,
		icon:  module.VideoIcon,
	},
	module.Avi: {
		color: module.Purple,
		icon:  module.VideoIcon,
	},
	module.Mov: {
		color: module.Purple,
		icon:  module.VideoIcon,
	},
	module.Mkv: {
		color: module.Purple,
		icon:  module.VideoIcon,
	},
	module.Mp3: {
		color: module.Cyan,
		icon:  module.AudioIcon,
	},
	module.Wav: {
		color: module.Cyan,
		icon:  module.AudioIcon,
	},
	module.Flac: {
		color: module.Cyan,
		icon:  module.AudioIcon,
	},
	module.Zip: {
		color: module.Yellow,
		icon:  module.ArchiveIcon,
	},
	module.Tar: {
		color: module.Yellow,
		icon:  module.ArchiveIcon,
	},
	module.Gz: {
		color: module.Yellow,
		icon:  module.ArchiveIcon,
	},
	module.Rar: {
		color: module.Yellow,
		icon:  module.ArchiveIcon,
	},
	module.SevenZ: {
		color: module.Yellow,
		icon:  module.ArchiveIcon,
	},
	module.Pdf: {
		color: module.Red,
		icon:  module.PDFIcon,
	},
	module.Doc: {
		color: module.Blue,
		icon:  module.WordIcon,
	},
	module.Docx: {
		color: module.Blue,
		icon:  module.WordIcon,
	},
	module.Xls: {
		color: module.Green,
		icon:  module.ExcelIcon,
	},
	module.Xlsx: {
		color: module.Green,
		icon:  module.ExcelIcon,
	},
	module.Ppt: {
		color: module.Orange,
		icon:  module.PowerPointIcon,
	},
	module.Pptx: {
		color: module.Orange,
		icon:  module.PowerPointIcon,
	},
	module.Ttf: {
		color: module.Grey,
		icon:  module.FontIcon,
	},
	module.Otf: {
		color: module.Grey,
		icon:  module.FontIcon,
	},
	module.Woff: {
		color: module.Grey,
		icon:  module.FontIcon,
	},
	module.Woff2: {
		color: module.Grey,
		icon:  module.FontIcon,
	},
	module.Exe: {
		color: module.Red,
		icon:  module.BinaryIcon,
	},
	module.Dll: {
		color: module.Red,
		icon:  module.BinaryIcon,
	},
	module.So: {
		color: module.Red,
		icon:  module.BinaryIcon,
	},
	module.Dylib: {
		color: module.Red,
		icon:  module.BinaryIcon,
	},
	module.Lock: {
		color: module.Yellow,
		icon:  module.LockIcon,
	},
	module.Key: {
		color: module.Yellow,
		icon:  module.CertificateIcon,
	},
	module.Pem: {
		color: module.Yellow,
		icon:  module.CertificateIcon,
	},
	module.Crt: {
		color: module.Yellow,
		icon:  module.CertificateIcon,
	},
	module.Pub: {
		color: module.Yellow,
		icon:  module.CertificateIcon,
	},
	module.Cer: {
		color: module.Yellow,
		icon:  module.CertificateIcon,
	},
	module.P12: {
		color: module.Yellow,
		icon:  module.CertificateIcon,
	},
	module.Log: {
		color: module.Grey,
		icon:  module.LogIcon,
	},
	module.Txt: {
		color: module.Grey,
		icon:  module.TextIcon,
	},
}

// CheckDefaultDirectory func - returns current working directory if no directory is specified
func CheckDefaultDirectory(directory *string) error {
	if *directory == "" {
		// Get the current working directory - default value
		var err error
		*directory, err = getWorkingDirectory()
		if err != nil {
			return fmt.Errorf("%w: %w", errGetwd, err)
		}
	}

	return nil
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
		if closeErr := d.Close(); closeErr != nil {
			return nil, nil, fmt.Errorf("read directory %q: %w; close directory: %v", directory, err, closeErr)
		}
		return nil, nil, err
	}
	return files, d, nil
}

// closeDirectory func - closes the directory.
func closeDirectory(directory *os.File) error {
	if directory == nil {
		return nil
	}

	err := directory.Close()
	if err != nil {
		return fmt.Errorf("close directory: %w", err)
	}

	return nil
}

// isHidden func - checks if the file name starts with a dot (hidden file).
func isHidden(name string) bool {
	return len(name) > 0 && name[0] == '.'
}

// printWithIconAndPrefix func - prints files and folder with icons and prefix
func printWithIconAndPrefix(prefix string, file os.FileInfo) string {
	if file.IsDir() {
		return formatFolderWithIcon(prefix, file.Name())
	}

	return formatFileWithIcon(prefix, file.Name())
}

func formatFileWithIcon(prefix, name string) string {
	style := resolveFileIconStyle(name)
	if style.color == "" {
		return fmt.Sprintf("%s%s %s", prefix, style.icon, name)
	}
	return fmt.Sprintf(IconFileFormat, prefix, style.color, style.icon, module.Reset, name)
}

func formatFolderWithIcon(prefix, name string) string {
	style := resolveFolderIconStyle(name)
	return fmt.Sprintf(IconFolderFormat, prefix, module.Bold, style.color, style.icon, name, module.Reset)
}

func resolveFileIconStyle(name string) iconStyle {
	if style, ok := filenameIconStyles[name]; ok {
		return style
	}

	lowerName := strings.ToLower(name)
	if style, ok := filenameIconStyles[lowerName]; ok {
		return style
	}

	if style, ok := extensionIconStyles[lowerName]; ok {
		return style
	}

	extension := strings.ToLower(filepath.Ext(name))
	if style, ok := extensionIconStyles[extension]; ok {
		return style
	}

	return defaultFileIconStyle
}

func resolveFolderIconStyle(name string) iconStyle {
	if style, ok := folderIconStyles[name]; ok {
		return style
	}

	if style, ok := folderIconStyles[strings.ToLower(name)]; ok {
		return style
	}

	return defaultFolderIconStyle
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

func appendReadableSize(formatted string, size int64) string {
	return fmt.Sprintf("%s (%s)", formatted, humanReadableSize(size))
}

func humanReadableSize(size int64) string {
	if size < 0 {
		size = 0
	}
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	}

	units := []string{"KB", "MB", "GB", "TB", "PB"}
	value := float64(size)
	unitIndex := -1
	for value >= 1024 && unitIndex < len(units)-1 {
		value /= 1024
		unitIndex++
	}

	return fmt.Sprintf("%.1f %s", value, units[unitIndex])
}

// sortSlice func - sorts a slice of os.FileInfo objects based on CLI flags.
// It modifies the original slice in place.
func sortSlice(files []os.FileInfo, flags module.Flags) {
	sort.Slice(files, func(i, j int) bool {
		return lessFileInfo(files[i], files[j], flags)
	})
}

func lessFileInfo(left, right os.FileInfo, flags module.Flags) bool {
	if flags.DirsFirst && left.IsDir() != right.IsDir() {
		return left.IsDir()
	}

	result := compareFileInfo(left, right, sortField(flags.SortBy))
	if result == 0 {
		result = compareName(left, right)
	}
	if flags.ReverseSort {
		return result > 0
	}

	return result < 0
}

func sortField(value string) string {
	if value == "" {
		return "name"
	}
	return value
}

func compareFileInfo(left, right os.FileInfo, sortBy string) int {
	switch sortBy {
	case "size":
		return compareInt64(left.Size(), right.Size())
	case "modified":
		return compareInt64(left.ModTime().UnixNano(), right.ModTime().UnixNano())
	case "type":
		return compareString(fileType(left), fileType(right))
	default:
		return compareName(left, right)
	}
}

func compareName(left, right os.FileInfo) int {
	return compareString(strings.ToLower(left.Name()), strings.ToLower(right.Name()))
}

func fileType(file os.FileInfo) string {
	if file.IsDir() {
		return ""
	}
	return strings.ToLower(filepath.Ext(file.Name()))
}

func compareString(left, right string) int {
	if left < right {
		return -1
	}
	if left > right {
		return 1
	}
	return 0
}

func compareInt64(left, right int64) int {
	if left < right {
		return -1
	}
	if left > right {
		return 1
	}
	return 0
}

// getTerminalWidth returns the width of the terminal
func getTerminalWidth() int {
	width, _, err := terminalSize(int(os.Stdout.Fd()))
	if err != nil || width == 0 {
		return 80 // default
	}
	return width
}

// stripANSI removes ANSI color codes to calculate actual string length
func stripANSI(str string) string {
	// Simple ANSI stripper - removes escape sequences
	inEscape := false
	var result strings.Builder

	for i := 0; i < len(str); i++ {
		if str[i] == '\x1b' && i+1 < len(str) && str[i+1] == '[' {
			inEscape = true
			continue
		}
		if inEscape {
			if (str[i] >= 'A' && str[i] <= 'Z') || (str[i] >= 'a' && str[i] <= 'z') {
				inEscape = false
			}
			continue
		}
		result.WriteByte(str[i])
	}
	return result.String()
}

// getVisibleLength returns the visible length of a string (excluding ANSI codes)
func getVisibleLength(str string) int {
	return len(stripANSI(str))
}

// printGrid prints entries in a grid layout.
func printGrid(output io.Writer, entries []string, maxLen int) error {
	if len(entries) == 0 {
		return nil
	}

	termWidth := getTerminalWidth()
	columnWidth := maxLen + 4 // Add padding between columns

	// Calculate number of columns that fit
	numColumns := termWidth / columnWidth
	if numColumns < 1 {
		numColumns = 1
	}

	// Print entries in grid
	for i := 0; i < len(entries); i++ {
		entry := entries[i]
		visibleLen := getVisibleLength(entry)

		// Print entry
		if _, err := fmt.Fprint(output, entry); err != nil {
			return err
		}

		// Add padding to align columns (except for last column in row)
		if (i+1)%numColumns != 0 && i < len(entries)-1 {
			padding := columnWidth - visibleLen
			if _, err := fmt.Fprint(output, strings.Repeat(" ", padding)); err != nil {
				return err
			}
		} else {
			// New line at end of row
			if _, err := fmt.Fprintln(output); err != nil {
				return err
			}
		}
	}

	// Add final newline if needed
	if len(entries)%numColumns != 0 {
		if _, err := fmt.Fprintln(output); err != nil {
			return err
		}
	}

	return nil
}
