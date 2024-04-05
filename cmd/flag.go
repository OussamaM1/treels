// Package cmd - cmd/flag.go
package cmd

import "github.com/oussamaM1/treels/module"

// FlagDefinition func - defines flags for the application.
func FlagDefinition(flags *module.Flags) {
	RootCommand.PersistentFlags().BoolVarP(&flags.ShowHidden, "all", "a", false, "List all files and directories")
	RootCommand.PersistentFlags().BoolVarP(&flags.ShowTreeView, "tree", "t", false, "Tree view of the directory")
}
