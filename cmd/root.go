// Package cmd - cmd/root.go
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/oussamaM1/treels/module"
	"github.com/oussamaM1/treels/service"
	"github.com/spf13/cobra"
)

const version = "v1.3.1"

// Execute func - runs the root command.
func Execute() {
	execute(newRootCmd(), os.Stderr, os.Exit)
}

func execute(cmd *cobra.Command, errorOutput io.Writer, exit func(int)) {
	if err := cmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(errorOutput, err)
		exit(1)
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
			if cmd.Flags().Changed("depth") && flag.TreeDepth < 0 {
				return fmt.Errorf("--depth must be greater than or equal to 0")
			}
			flag.LimitTreeDepth = cmd.Flags().Changed("depth")

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
