package cmd

import (
	"fmt"
	"os"

	"github.com/devlup-labs/spok/internal/pkg/policy"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Removes user and principles from the policy.yaml file",
	Long:  `Removes the user and principles from the policy.yaml file maintained by SOS.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println("Invalid number of arguments for remove, should be `verifier remove <Email> <User>`")

			os.Exit(1)
		}

		emailArgs := args[0]
		userArgs := args[1]

		policy.RemovePolicy(emailArgs, userArgs)
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
