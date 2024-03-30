/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"github.com/devlup-labs/sos/internal/pkg/policy"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add called")
			if len(os.Args) != 4 {
				fmt.Println("Invalid number of arguments for add, should be `verifier add <Email> <User (TOKEN u)>`")

				os.Exit(1)
			}

			emailArgs := os.Args[2]
			userArgs := os.Args[3]

			policy.AddPolicy(emailArgs, userArgs)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
