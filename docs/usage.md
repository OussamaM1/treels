# Usage guide

## Table of contents

- [Output modes](#output-modes)
- [Common workflows](#common-workflows)
- [Flags](#flags)
- [Flag interactions](#flag-interactions)
- [Exit codes](#exit-codes)

`treels` lists one directory at a time. If no path is provided, it lists the current working directory.

```bash
treels [flags] [path]
```

Examples in this guide use `--no-icons` when showing text output so the examples render correctly without a Nerd Font.

## Output modes

### Compact mode

Compact mode is the default. It lists the direct children of the target directory in a terminal-width-aware grid.

```bash
treels --no-icons
```

```text
.
LICENSE      README.md    cmd          docs         go.mod
main.go      module       service      utils

6 directories, 4 files
```

### Tree mode

Use `--tree` / `-t` for recursive tree output.

```bash
treels --tree --no-icons
```

```text
.
├── LICENSE
├── README.md
├── cmd
│   ├── flag.go
│   ├── root.go
│   └── root_test.go
└── main.go

1 directories, 4 files
```

### Long mode

Use `--long` / `-l` to show file metadata.

```bash
treels --long --no-icons
```

```text
.
-rw-r--r--        1.1 KB  2026-06-10  README.md
drwxr-xr-x        4.0 KB  2026-06-10  cmd
-rw-r--r--          123 B  2026-06-10  main.go

1 directories, 2 files
```

Long mode columns are:

| Column | Description |
| --- | --- |
| mode | File mode / permissions, such as `-rw-r--r--` or `drwxr-xr-x`. |
| size | Size in bytes by default, or human-readable with `--readable`. |
| date | Modification date in `YYYY-MM-DD` format. |
| name | Entry name, including icons unless `--no-icons` is used. |

Long mode also works with tree output:

```bash
treels --tree --long --readable --no-icons
```

### JSON mode

Use `--json` for machine-readable output.

```bash
treels --json
```

See [json-output.md](json-output.md) for the JSON schema and examples.

## Common workflows

### Quickly scan the current directory

```bash
treels
```

### Include hidden files

```bash
treels --all
```

### Inspect a project tree without generated files

```bash
treels --tree --depth 2 --gitignore
```

### Show only folder structure

```bash
treels --tree --dirs-only
```

### Show detailed metadata with readable sizes

```bash
treels --long --readable
```

### Focus on specific file types

```bash
treels --include "*.go"
treels --include "*.go" --include "*.md"
```

### Hide noisy files or directories

```bash
treels --exclude "*.log"
treels --exclude "vendor/**"
```

### Combine include and exclude filters

```bash
treels --tree --include "*.go" --exclude "vendor/**"
```

### Show Git status decorations

```bash
treels --tree --git-status --no-icons
```

Example symbols:

| Symbol | Meaning | Color |
| --- | --- | --- |
| `M` | Modified | Yellow |
| `A` | Added | Green |
| `D` | Deleted | Red |
| `?` | Untracked | Cyan |
| `!` | Ignored | Grey |
| space | Clean or no Git status | Uncolored |

### Sort by largest files first

```bash
treels --sort size --reverse --long --readable
```

### Show directories before files

```bash
treels --dirs-first
```

### Disable icons for plain terminals or logs

```bash
treels --no-icons
```

### Use from scripts

```bash
treels --json
```

## Flags

| Flag | Description |
| --- | --- |
| `-a`, `--all` | List hidden files and directories. |
| `-t`, `--tree` | Show recursive tree view. |
| `--dirs-only` | Show only directories. |
| `--depth N` | Limit tree recursion depth. |
| `--gitignore` | Respect `.gitignore` rules from the target directory. |
| `--include PATTERN` | Show only entries matching a glob pattern. Can be used multiple times. |
| `--exclude PATTERN` | Hide entries matching a glob pattern. Can be used multiple times. |
| `--git-status` | Show Git status symbols next to entries in text output. |
| `--json` | Output machine-readable JSON. |
| `-l`, `--long` | Show detailed file metadata. |
| `--sort name|size|modified|type` | Sort entries by name, size, modification time, or file type. Defaults to `name`. |
| `--reverse` | Reverse the selected sort order. |
| `--dirs-first` | Group directories before files. |
| `--no-icons` | Disable file and folder icons. |
| `--no-summary` | Hide the final text summary. |
| `-r`, `--readable` | Show human-readable file and directory sizes. |
| `-v`, `--version` | Print the current version. |
| `-h`, `--help` | Show help. |

## Flag interactions

| Combination | Behavior |
| --- | --- |
| `--tree --depth N` | Recurses up to `N` levels. `--depth 0` shows only the root line and summary. |
| `--tree --dirs-only` | Recursively shows directories while omitting files. |
| `--long --readable` | Shows human-readable sizes in the long metadata column. |
| `--tree --long` | Shows tree branches plus metadata for each entry. |
| `--include "*.go" --include "*.md"` | Shows entries matching either include pattern. |
| `--include "*.go" --exclude "vendor/**"` | Shows Go files except entries under `vendor`. |
| `--tree --include "*.go"` | Keeps parent directories visible when they contain included files. |
| `--tree --git-status` | Shows tree branches plus Git status symbols. |
| `--sort size --reverse` | Shows largest entries first. |
| `--sort modified --reverse` | Shows newest entries first. |
| `--dirs-first --reverse` | Keeps directories grouped first, then reverses the selected sort within each group. |
| `--json --tree` | Emits recursive JSON with `children` arrays for directories. |
| `--json --long` | JSON output is unchanged; `--long` only affects text output. |
| `--gitignore --all` | Hidden files are included only if they are not ignored by `.gitignore`. |
| `--json --no-summary` | No effect; JSON always includes the `summary` object. |

## Exit codes

| Code | Meaning |
| --- | --- |
| `0` | Command completed successfully. |
| `1` | Invalid arguments, invalid directory, or runtime error. |
