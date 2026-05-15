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
	cmd.PersistentFlags().BoolVar(&flags.HideIcon, "no-icons", false, "Disable icons")
	cmd.PersistentFlags().BoolVarP(&flags.ShowReadableSize, "readable", "r", false, "Show human-readable size for each file and directory")
	cmd.PersistentFlags().BoolVarP(&flags.ShowVersion, "version", "v", false, "Show treels version")
}
