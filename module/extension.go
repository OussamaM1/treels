// Package module - module/extension .go
package module

const (
	// Go represents the Golang file extension
	Go string = ".go"
	// Mod represents the Go module file extension
	Mod string = ".mod"
	// Sum represents the Go sum file extension
	Sum string = ".sum"

	// Documentation

	// Md represents the Markdown file extension
	Md string = ".md"
	// Txt represents the text file extension
	Txt string = ".txt"

	// Git

	// Gitignore represents the gitignore file
	Gitignore string = ".gitignore"
	// GitFolder represents the git folder
	GitFolder string = ".git"
	// Github represents the github folder
	Github string = ".github"

	// IDEs

	// IntellijFolder represents the IntelliJ IDEA folder
	IntellijFolder string = ".idea"

	// Data formats

	// JSON represents the JSON file extension
	JSON string = ".json"
	// Yml represents the YAML file extension
	Yml string = ".yml"
	// Yaml represents the YAML file extension (alternative)
	Yaml string = ".yaml"
	// XML represents the XML file extension
	XML string = ".xml"

	// Database

	// SQL represents the SQL file extension
	SQL string = ".sql"
	// Pls represents the PL/SQL file extension
	Pls string = ".pls"
	// Plb represents the PL/SQL binary file extension
	Plb string = ".plb"

	// Java

	// Java represents the Java source file extension
	Java string = ".java"
	// Class represents the Java compiled class file extension
	Class string = ".class"

	// C/C++

	// Cpp represents the C++ source file extension
	Cpp string = ".cpp"
	// C represents the C source file extension
	C string = ".c"

	// JavaScript/TypeScript

	// Js represents the JavaScript file extension
	Js string = ".js"
	// Jsx represents the React JSX file extension
	Jsx string = ".jsx"
	// Ts represents the TypeScript file extension
	Ts string = ".ts"
	// Tsx represents the TypeScript React file extension
	Tsx string = ".tsx"

	// Web

	// HTML represents the HTML file extension
	HTML string = ".html"
	// Htm represents the HTM file extension
	Htm string = ".htm"
	// CSS represents the CSS file extension
	CSS string = ".css"
	// Scss represents the SCSS file extension
	Scss string = ".scss"
	// Sass represents the Sass file extension
	Sass string = ".sass"
	// Vue represents the Vue.js file extension
	Vue string = ".vue"

	// Other languages

	// Rs represents the Rust source file extension
	Rs string = ".rs"
	// Py represents the Python file extension
	Py string = ".py"
	// Rb represents the Ruby file extension
	Rb string = ".rb"
	// Rake represents the Rake file extension
	Rake string = ".rake"
	// Gemfile represents the Ruby Gemfile
	Gemfile string = "Gemfile"
	// Php represents the PHP file extension
	Php string = ".php"
	// Swift represents the Swift file extension
	Swift string = ".swift"
	// Kt represents the Kotlin file extension
	Kt string = ".kt"
	// Kts represents the Kotlin script file extension
	Kts string = ".kts"
	// Cs represents the C# file extension
	Cs string = ".cs"
	// Csx represents the C# script file extension
	Csx string = ".csx"

	// Shell

	// Sh represents the shell script file extension
	Sh string = ".sh"
	// Bash represents the Bash script file extension
	Bash string = ".bash"
	// Zsh represents the Zsh script file extension
	Zsh string = ".zsh"

	// Docker

	// Dockerfile represents the Dockerfile
	Dockerfile string = "Dockerfile"
	// Dockerignore represents the dockerignore file
	Dockerignore string = ".dockerignore"

	// Config files

	// Conf represents the config file extension
	Conf string = ".conf"
	// Cfg represents the config file extension
	Cfg string = ".cfg"
	// Ini represents the INI file extension
	Ini string = ".ini"
	// Env represents the environment file
	Env string = ".env"

	// Build files

	// Makefile represents the Makefile
	Makefile string = "Makefile"
	// Make represents the make file extension
	Make string = ".make"

	// Package managers

	// PackageJSON represents the npm package.json file
	PackageJSON string = "package.json"

	// Infrastructure as Code

	// Tf represents the Terraform file extension
	Tf string = ".tf"
	// Tfvars represents the Terraform variables file extension
	Tfvars string = ".tfvars"

	// Images

	// Png represents the PNG image file extension
	Png string = ".png"
	// Jpg represents the JPG image file extension
	Jpg string = ".jpg"
	// Jpeg represents the JPEG image file extension
	Jpeg string = ".jpeg"
	// Gif represents the GIF image file extension
	Gif string = ".gif"
	// Svg represents the SVG image file extension
	Svg string = ".svg"
	// Ico represents the icon file extension
	Ico string = ".ico"

	// Video

	// Mp4 represents the MP4 video file extension
	Mp4 string = ".mp4"
	// Avi represents the AVI video file extension
	Avi string = ".avi"
	// Mov represents the MOV video file extension
	Mov string = ".mov"
	// Mkv represents the MKV video file extension
	Mkv string = ".mkv"

	// Audio

	// Mp3 represents the MP3 audio file extension
	Mp3 string = ".mp3"
	// Wav represents the WAV audio file extension
	Wav string = ".wav"
	// Flac represents the FLAC audio file extension
	Flac string = ".flac"

	// Archives

	// Zip represents the ZIP archive file extension
	Zip string = ".zip"
	// Tar represents the TAR archive file extension
	Tar string = ".tar"
	// Gz represents the GZIP archive file extension
	Gz string = ".gz"
	// Rar represents the RAR archive file extension
	Rar string = ".rar"
	// SevenZ represents the 7z archive file extension
	SevenZ string = ".7z"

	// Documents

	// Pdf represents the PDF file extension
	Pdf string = ".pdf"

	// Security/Keys

	// Lock represents the lock file
	Lock string = ".lock"
	// Key represents the key file extension
	Key string = ".key"
	// Pem represents the PEM file extension
	Pem string = ".pem"
	// Crt represents the certificate file extension
	Crt string = ".crt"
	// Pub represents the public key file extension
	Pub string = ".pub"

	// Logs

	// Log represents the log file extension
	Log string = ".log"

	// Scala represents the Scala file extension
	Scala string = ".scala"

	// Ex represents the Elixir file extension
	Ex string = ".ex"
	// Exs represents the Elixir script file extension
	Exs string = ".exs"

	// Hs represents the Haskell file extension
	Hs string = ".hs"

	// Clj represents the Clojure file extension
	Clj string = ".clj"

	// Dart represents the Dart file extension
	Dart string = ".dart"

	// R represents the R programming language file extension
	R string = ".r"
	// Rmd represents the R Markdown file extension
	Rmd string = ".rmd"

	// Lua represents the Lua file extension
	Lua string = ".lua"

	// Pl represents the Perl file extension
	Pl string = ".pl"
	// Pm represents the Perl module file extension
	Pm string = ".pm"

	// Psql represents the PostgreSQL file extension
	Psql string = ".psql"

	// Wasm represents the WebAssembly file extension
	Wasm string = ".wasm"

	// Less represents the LESS file extension
	Less string = ".less"
	// Styl represents the Stylus file extension
	Styl string = ".styl"

	// NgComponent represents Angular component files
	NgComponent string = ".component.ts"
	// NgModule represents Angular module files
	NgModule string = ".module.ts"
	// NgService represents Angular service files
	NgService string = ".ngservice.ts"

	// Svelte represents the Svelte file extension
	Svelte string = ".svelte"

	// Prisma represents the Prisma schema file extension
	Prisma string = ".prisma"

	// Graphql represents the GraphQL file extension
	Graphql string = ".graphql"
	// Gql represents the GraphQL query file extension
	Gql string = ".gql"

	// Proto represents the Protocol Buffers file extension
	Proto string = ".proto"

	// License represents the LICENSE file
	License string = "LICENSE"
	// LicenseMd represents the LICENSE.md file
	LicenseMd string = "LICENSE.md"

	// Readme represents the README file
	Readme string = "README"
	// ReadmeMd represents the README.md file
	ReadmeMd string = "README.md"

	// YmlCI represents CI/CD yaml files
	YmlCI string = ".gitlab-ci.yml"
	// Jenkinsfile JenkinsFile represents Jenkins pipeline file
	Jenkinsfile string = "Jenkinsfile"

	// Vagrantfile represents the Vagrantfile
	Vagrantfile string = "Vagrantfile"

	// Procfile represents the Procfile (Heroku)
	Procfile string = "Procfile"

	// Editorconfig represents the editorconfig file
	Editorconfig string = ".editorconfig"

	// Eslintrc represents the eslint config file
	Eslintrc string = ".eslintrc"
	// EslintrcJSON represents the eslint JSON config
	EslintrcJSON string = ".eslintrc.json"

	// Prettierrc represents the prettier config file
	Prettierrc string = ".prettierrc"
	// Prettierignore represents the prettier ignore file
	Prettierignore string = ".prettierignore"

	// Babelrc represents the babel config file
	Babelrc string = ".babelrc"

	// WebpackConfig represents the webpack config file
	WebpackConfig string = "webpack.config.js"

	// ViteConfig represents the vite config file
	ViteConfig string = "vite.config.js"

	// TsConfig represents the TypeScript config file
	TsConfig string = "tsconfig.json"

	// CMakeLists represents the CMake file
	CMakeLists string = "CMakeLists.txt"
	// Cmake represents the CMake file extension
	Cmake string = ".cmake"

	// Gradle represents the Gradle file extension
	Gradle string = ".gradle"
	// GradleKts represents the Gradle Kotlin DSL file extension
	GradleKts string = ".gradle.kts"

	// Pom represents the Maven pom.xml file
	Pom string = "pom.xml"

	// Requirements represents Python requirements file
	Requirements string = "requirements.txt"

	// PoetryLock represents Poetry lock file
	PoetryLock string = "poetry.lock"
	// Pyproject represents Python project config
	Pyproject string = "pyproject.toml"

	// CargoToml represents Rust Cargo.toml file
	CargoToml string = "Cargo.toml"
	// CargoLock represents Rust Cargo.lock file
	CargoLock string = "Cargo.lock"

	// Sqlite represents SQLite database file
	Sqlite string = ".sqlite"
	// Db represents generic database file
	Db string = ".db"

	// Exe represents Windows executable
	Exe string = ".exe"
	// Dll represents Windows DLL
	Dll string = ".dll"
	// So represents Linux shared object
	So string = ".so"
	// Dylib represents macOS dynamic library
	Dylib string = ".dylib"

	// Doc represents Word document
	Doc string = ".doc"
	// Docx represents Word document (new format)
	Docx string = ".docx"
	// Xls represents Excel spreadsheet
	Xls string = ".xls"
	// Xlsx represents Excel spreadsheet (new format)
	Xlsx string = ".xlsx"
	// Ppt represents PowerPoint presentation
	Ppt string = ".ppt"
	// Pptx represents PowerPoint presentation (new format)
	Pptx string = ".pptx"

	// Ttf represents TrueType font
	Ttf string = ".ttf"
	// Otf represents OpenType font
	Otf string = ".otf"
	// Woff represents Web Open Font Format
	Woff string = ".woff"
	// Woff2 represents Web Open Font Format 2
	Woff2 string = ".woff2"

	// Cer represents certificate file
	Cer string = ".cer"
	// P12 represents PKCS#12 certificate
	P12 string = ".p12"

	// Toml represents TOML config file
	Toml string = ".toml"

	// Nix represents Nix file extension
	Nix string = ".nix"
)
