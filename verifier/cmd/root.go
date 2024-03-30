package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "verifier",
	Short: "A verifier for comparing the key for SSH",
	Long: `A verifier for comparing the key for SSH and adding users.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
