# treels

`treels` is a small Go CLI that blends the quick scan of `ls` with the structure of `tree`.
Use it to inspect a directory as a compact grid, expand it as a tree, hide project noise with
`.gitignore`, and keep large repositories readable with depth limits.

```bash
treels                 # compact ls-style view
treels -t              # tree view
treels -t --depth 2    # tree view, limited to two levels
treels -t --gitignore  # tree view, excluding .gitignore matches
```

> [!NOTE]
> File and folder icons look best with a Nerd Font installed. If your terminal does not support
> Nerd Font glyphs, use `--no-icons`.

## Preview

Compact directory listing:

```text
$ treels --no-icons file-icons-example
.
Dockerfile               Main.kt                  Program.cs
app.conf                 app.lock                 app.log
app.rb                   app.swift                backup.zip
c-file.c                 class-java-file.class    component.vue
config.xml               cpp-file.cpp             document.pdf
index.html               index.php                java-file.java
javascript-file.js       json-file.json           logo.png
main.tf                  package.json             plb-file.plb
pls-file.pls             python-file.py           react-component.jsx
react-typescript.tsx     rust-file.rs             script.sh
song.mp3                 sql-file.sql             styles.css
typescript-file.ts       video.mp4                yaml-file.yml

0 directories, 36 files
```

Tree view with a depth limit:

```text
$ treels -t --depth 1 --gitignore --no-icons
.
├── LICENSE
├── README.md
├── cmd
├── example
├── file-icons-example
├── go.mod
├── go.sum
├── main.go
├── module
├── service
└── utils

6 directories, 5 files
```

Focused tree view for one directory:

```text
$ treels -t --depth 2 --no-icons service
.
├── gitignore.go
├── service.go
├── service_test.go
└── util.go

0 directories, 4 files
```

## Features

- Compact grid output for fast directory scans.
- Recursive tree output with familiar branch characters.
- Optional Nerd Font icons and file-type colors.
- `--gitignore` support to skip generated files, dependencies, logs, and build output.
- `--depth N` to keep tree output readable in large repositories.
- `--readable` file sizes.
- `--all` support for hidden files.
- `--no-icons` fallback for terminals without icon fonts.

## Installation

Install with Go:

```bash
go install github.com/oussamaM1/treels@latest
```

Make sure your Go binary directory is in your `PATH`. You can check where Go installs binaries with:

```bash
go env GOPATH
```

The binary is usually placed in:

```text
$(go env GOPATH)/bin
```

## Usage

```bash
treels [flags] [path]
```

If no path is provided, `treels` lists the current directory.

## Examples

List the current directory:

```bash
treels
```

Include hidden files:

```bash
treels --all
```

Show a tree:

```bash
treels --tree
```

Limit tree depth:

```bash
treels --tree --depth 2
```

Respect `.gitignore` rules:

```bash
treels --tree --gitignore
```

Show readable sizes:

```bash
treels --readable
```

Disable icons:

```bash
treels --no-icons
```

## Flags

| Flag | Description |
| --- | --- |
| `-a`, `--all` | List hidden files and directories. |
| `-t`, `--tree` | Show recursive tree view. |
| `--depth N` | Limit tree recursion depth. |
| `--gitignore` | Respect `.gitignore` rules from the target directory. |
| `--no-icons` | Disable file and folder icons. |
| `-r`, `--readable` | Show human-readable file and directory sizes. |
| `-v`, `--version` | Print the current version. |
| `-h`, `--help` | Show help. |

## Fonts

Icons require a Nerd Font-compatible terminal. Recommended setup:

1. Install a Nerd Font, such as FiraCode Nerd Font.
2. Select that font in your terminal profile.
3. Run `treels` normally.

If icons render as empty boxes or odd symbols, use:

```bash
treels --no-icons
```

## Development

Run tests:

```bash
go test ./...
```

Run lint:

```bash
golangci-lint run
```

Build locally:

```bash
go build .
```
