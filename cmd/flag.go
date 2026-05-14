// Package cmd - cmd/flag.go
package cmd

import (
	"github.com/oussamaM1/treels/module"
	"github.com/spf13/cobra"
)

// FlagDefinition func - defines flags for the application.
func FlagDefinition(cmd *cobra.Command, flags *module.Flags) {
	flags.HideIcon = true
	cmd.PersistentFlags().BoolVarP(&flags.ShowHidden, "all", "a", false, "List all files and directories")
	cmd.PersistentFlags().BoolVarP(&flags.ShowTreeView, "tree", "t", false, "Tree view of the directory")
	cmd.PersistentFlags().BoolVarP(&flags.HideIcon, "icon", "i", false, "Disable icons (Enabled by default)")
}
