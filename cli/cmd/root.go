package cmd

import (
	"github.com/devlup-labs/spok/cli/cmd/alias"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "spok",
	Short: "SPoK - Sans Password or Key",
	Long:  `This project uses the OpenPubKey project for passwordless SSH using email authentication`,
}

func Execute() {
	rootCmd.AddCommand(alias.AliasCmd)

	cobra.CheckErr(rootCmd.Execute())
}
