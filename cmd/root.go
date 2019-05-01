package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{}
)

// RootCommand the entry point for the cmd interface
func RootCommand() *cobra.Command {
	rootCmd.AddCommand(&tickCmd)

	return rootCmd
}
