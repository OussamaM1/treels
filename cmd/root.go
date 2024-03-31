// Package cmd provides functionality related to the command line.
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

// RootCommand Command - is the root command for the application.
var RootCommand = &cobra.Command{
	Use:   "treels",
	Short: "🌳 Treels is a CLI tool crafted in Go, merging tree and ls commands with intuitive merging and beautification features.",
	Run:   run,
}

// Execute func - runs the root command.
func Execute() {
	if err := RootCommand.Execute(); err != nil {
		log.Fatalf(err.Error())
	}
}

func run(_ *cobra.Command, _ []string) {
	fmt.Println("run function")
}
