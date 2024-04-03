// Package cmd - cmd/root.go
package cmd

import (
	"fmt"
	"github.com/oussamaM1/treels/service"
	"github.com/oussamaM1/treels/utils"
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

func run(_ *cobra.Command, args []string) {
	if utils.ValidateDirectoryArgs(args) {
		service.ListDirectory(service.Options{Directory: args[0]})
	} else if len(args) == 0 {
		service.ListDirectory(service.Options{})
	} else {
		fmt.Println("❌ Usage: treels <directory_path>")
	}
}
