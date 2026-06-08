// Package cmd - cmd/root.go
package cmd

import (
	"fmt"
	"os"

	"github.com/oussamaM1/treels/module"
	"github.com/oussamaM1/treels/service"
	"github.com/spf13/cobra"
)

const version = "v1.3.1"

// Execute func - runs the root command.
func Execute() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	var flag module.Flags

	cmd := &cobra.Command{
		Use:   "treels [path]",
		Short: "⚡ Treels is a CLI tool crafted in Go, merging tree and ls commands with intuitive merging and beautification features.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if flag.ShowVersion {
				_, err := fmt.Fprintf(cmd.OutOrStdout(), "treels %s\n", version)
				return err
			}

			options := module.Options{Flags: flag}
			if len(args) == 1 {
				options.Directory = args[0]
			}

			return service.Dispatcher(options)
		},
	}

	FlagDefinition(cmd, &flag)
	cmd.SilenceUsage = true

	return cmd
}
