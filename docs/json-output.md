# JSON output

## Table of contents

- [Flat JSON example](#flat-json-example)
- [Tree JSON example](#tree-json-example)
- [Schema](#schema)
- [Flag behavior](#flag-behavior)
- [Stability notes](#stability-notes)

Use `--json` when `treels` output needs to be consumed by scripts or other tools.

```bash
treels --json [path]
treels --tree --json --depth 2 [path]
```

JSON output is pretty-printed and always includes a `summary` object.

## Flat JSON example

```bash
treels --json service
```

```json
{
  "root": "service",
  "tree": false,
  "entries": [
    {
      "name": "gitignore.go",
      "path": "service/gitignore.go",
      "type": "file",
      "size": 3946
    },
    {
      "name": "json.go",
      "path": "service/json.go",
      "type": "file",
      "size": 3822
    }
  ],
  "summary": {
    "directories": 0,
    "files": 2
  }
}
```

## Tree JSON example

```bash
treels --tree --json --depth 2 cmd
```

```json
{
  "root": "cmd",
  "tree": true,
  "entries": [
    {
      "name": "flag.go",
      "path": "cmd/flag.go",
      "type": "file",
      "size": 1174
    },
    {
      "name": "root.go",
      "path": "cmd/root.go",
      "type": "file",
      "size": 1404
    }
  ],
  "summary": {
    "directories": 0,
    "files": 2
  }
}
```

Directories in tree mode may include `children`:

```json
{
  "name": "cmd",
  "path": "./cmd",
  "type": "directory",
  "size": 128,
  "children": [
    {
      "name": "root.go",
      "path": "cmd/root.go",
      "type": "file",
      "size": 1404
    }
  ]
}
```

## Schema

| Field | Type | Description |
| --- | --- | --- |
| `root` | string | Directory that was listed. |
| `tree` | boolean | Whether recursive tree mode was enabled. |
| `entries` | array | Visible entries after filtering. |
| `entries[].name` | string | Entry name. |
| `entries[].path` | string | Entry path joined with the listed parent path. |
| `entries[].type` | string | Either `file` or `directory`. |
| `entries[].size` | number | Size in bytes from the filesystem. |
| `entries[].children` | array | Child entries for directories in tree JSON mode. Omitted when empty. |
| `summary.directories` | number | Count of visible directories after filtering and depth limits. |
| `summary.files` | number | Count of visible files after filtering and depth limits. |

## Flag behavior

Most filtering flags affect JSON output:

| Flag | Effect on JSON |
| --- | --- |
| `--tree` | Enables recursive JSON and directory `children`. |
| `--depth N` | Limits recursive JSON depth when combined with `--tree`. |
| `--all` | Includes hidden entries. |
| `--dirs-only` | Omits file entries. |
| `--gitignore` | Omits entries matched by the target directory's `.gitignore`. |
| `--include PATTERN` | Includes only matching entries. In tree mode, parent directories stay visible when they contain included descendants. |
| `--exclude PATTERN` | Omits matching entries. Excluded directories are not traversed. |
| `--sort name|size|modified|type` | Sorts JSON entries using the selected field. |
| `--reverse` | Reverses JSON entry order for the selected sort. |
| `--dirs-first` | Groups directory entries before file entries. |

Text formatting flags do not affect JSON output:

| Flag | JSON effect |
| --- | --- |
| `--long` | No effect. |
| `--readable` | No effect; sizes remain raw bytes. |
| `--no-icons` | No effect. |
| `--no-summary` | No effect; JSON always includes `summary`. |

## Stability notes

The JSON format is intended for automation. Avoid relying on text output when writing scripts; use `--json` instead.
