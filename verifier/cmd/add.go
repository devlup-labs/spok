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
	Short: "Adds user to the policy.yaml file",
	Long: `Adds user to the policy structure maintained in the server by our SOS.`,
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
