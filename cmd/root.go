package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "lokstra",
	Short: "Lokstra CLI: Build scalable, structured Go backend apps",
	Long: `Lokstra CLI is the official tool to scaffold, lint, and manage
Lokstra-based backend applications.`,
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
