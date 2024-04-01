package cmd

import (
	"fmt"
	"os"

	"github.com/devlup-labs/spok/internal/pkg/policy"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds user to the policy.yaml file",
	Long:  `Adds user to the policy structure maintained in the server by our SOS.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println("Invalid number of arguments for add, should be `verifier add <Email> <User>`")

			os.Exit(1)
		}

		emailArgs := args[0]
		userArgs := args[1]

		policy.AddPolicy(emailArgs, userArgs)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
