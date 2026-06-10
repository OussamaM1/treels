// Package cmd - cmd/flag.go
package cmd

import (
	"github.com/oussamaM1/treels/module"
	"github.com/spf13/cobra"
)

// FlagDefinition func - defines flags for the application.
func FlagDefinition(cmd *cobra.Command, flags *module.Flags) {
	cmd.PersistentFlags().BoolVarP(&flags.ShowHidden, "all", "a", false, "List all files and directories")
	cmd.PersistentFlags().BoolVarP(&flags.ShowTreeView, "tree", "t", false, "Tree view of the directory")
	cmd.PersistentFlags().BoolVar(&flags.ShowDirsOnly, "dirs-only", false, "List directories only")
	cmd.PersistentFlags().BoolVar(&flags.ShowJSON, "json", false, "Output machine-readable JSON")
	cmd.PersistentFlags().BoolVarP(&flags.ShowLongFormat, "long", "l", false, "Show detailed file metadata")
	cmd.PersistentFlags().BoolVar(&flags.HideIcon, "no-icons", false, "Disable icons")
	cmd.PersistentFlags().BoolVar(&flags.HideSummary, "no-summary", false, "Hide the final file and directory count")
	cmd.PersistentFlags().BoolVar(&flags.RespectGitIgnore, "gitignore", false, "Respect .gitignore rules")
	cmd.PersistentFlags().StringArrayVar(&flags.IncludePatterns, "include", nil, "Show only entries matching a glob pattern; can be used multiple times")
	cmd.PersistentFlags().StringArrayVar(&flags.ExcludePatterns, "exclude", nil, "Hide entries matching a glob pattern; can be used multiple times")
	cmd.PersistentFlags().StringVar(&flags.SortBy, "sort", "name", "Sort entries by name, size, modified, or type")
	cmd.PersistentFlags().BoolVar(&flags.ReverseSort, "reverse", false, "Reverse sort order")
	cmd.PersistentFlags().BoolVar(&flags.DirsFirst, "dirs-first", false, "Show directories before files")
	cmd.PersistentFlags().IntVar(&flags.TreeDepth, "depth", -1, "Limit tree view recursion depth")
	cmd.PersistentFlags().Lookup("depth").DefValue = "unlimited"
	cmd.PersistentFlags().BoolVarP(&flags.ShowReadableSize, "readable", "r", false, "Show human-readable size for each file and directory")
	cmd.PersistentFlags().BoolVarP(&flags.ShowVersion, "version", "v", false, "Show treels version")
}
