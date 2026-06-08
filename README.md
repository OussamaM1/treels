# ⚡ treels
Treels, a CLI tool built in Go, merges the tree and ls commands while introducing intuitive merging and beautification features, 
simplifying directory navigation and enhancing the command-line experience.

> [!IMPORTANT]  
> To ensure that icons are displayed correctly in the terminal, it's recommended to use Nerd fonts. For example, you can download the FiraCode Nerd Font from [here](https://github.com/ryanoasis/nerd-fonts/releases/).

## 🚀 Installation

To install the `treels` command-line tool, ensure you have Go 1.25.0 or newer installed on your system. If not, you can download and install it from the [official Go website](https://golang.org/dl/).

Once you have Go installed, open a terminal or command prompt and run the following command:

```bash
go install github.com/oussamaM1/treels@latest
```

This command will download the repository, build the `treels` executable, and place it in your Go binary directory. Make sure your Go binary directory is in your system's PATH so that you can execute `treels` from any directory.

## ⚡ Usage

```bash
treels [Flags] [Path]
```

## 🏷️ Flags

- `-a, --all`: List all files and directories
- `-h, --help`: Help for treels
- `-t, --tree`: Tree view of the directory
- `--no-icons`: Disable icons
- `--gitignore`: Respect .gitignore rules
- `--depth N`: Limit tree view recursion depth
- `-r, --readable`: Show human-readable size for each file and directory
- `-v, --version`: Show treels version

## 📋 Example

```bash
treels -a
```

This command will list all files and directories in the `current` directory.

![](example/example-treels-a.png)

```bash
treels -t /Project/treels
```

This command will display the tree view of the `/Project/treels/` directory.

![](example/example-treels-t.png)

```bash
treels -t --no-icons /Project/treels
```

This command will display the tree view without icons of the `/Project/treels/` directory.

![](example/example-treels-it.png)
