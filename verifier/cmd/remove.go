package cmd

import (
	"fmt"
	"os"

	"github.com/devlup-labs/sos/internal/pkg/policy"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Removes user and principles from the policy.yaml file",
	Long: `Removes the user and principles from the policy.yaml file maintained by SOS.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(os.Args) != 4 {
			fmt.Println("Invalid number of arguments for add, should be `verifier add <Email> <User (TOKEN u)>`")

			os.Exit(1)
		}

		emailArgs := os.Args[2]
		userArgs := os.Args[3]

		policy.RemovePolicy(emailArgs, userArgs)
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
