package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "maven-resource",
	Short: "Implementation of a concourse resource that queries maven dependencies",
	Long:  `Implementation of a concourse resource that queries maven dependencies`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.AddCommand(NewInCmd().Command)
	rootCmd.AddCommand(NewCheckCmd().Command)

	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf("Some kind of error: %s", err)
		os.Exit(1)
	}
}
