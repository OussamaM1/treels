// Package cmd - cmd/root.go
package cmd

import (
	"fmt"
	"github.com/oussamaM1/treels/module"
	"github.com/oussamaM1/treels/service"
	"github.com/oussamaM1/treels/utils"
	"github.com/spf13/cobra"
	"log"
)

var flag module.Flags

// RootCommand Command - is the root command for the application.
var RootCommand = &cobra.Command{
	Use:   "treels",
	Short: "⚡ Treels is a CLI tool crafted in Go, merging tree and ls commands with intuitive merging and beautification features.",
	Run:   run,
}

// Execute func - runs the root command.
func Execute() {
	FlagDefinition(&flag)
	if err := RootCommand.Execute(); err != nil {
		log.Fatalf(err.Error())
	}
}

func run(_ *cobra.Command, args []string) {
	if utils.ValidateDirectoryArgs(args) {
		service.Dispatcher(module.Options{Directory: args[len(args)-1], Flags: flag})
	} else if len(args) == 0 {
		service.Dispatcher(module.Options{Flags: flag})
	} else {
		fmt.Println("❌ Usage: treels [options] <directory_path>")
	}
}
