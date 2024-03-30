package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "verifier",
	Short: "A verifier for verifying the key for SSH",
	Long:  `A verifier for verifying the key for SSH and adding users.`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
