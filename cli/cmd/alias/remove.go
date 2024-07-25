package alias

import (
	"fmt"
	"os"
	"strings"

	"github.com/devlup-labs/spok/internal/pkg/selector"
	"github.com/spf13/cobra"
)

// removeCmd represents the alias remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove an/all alias(es)",
	Long:  "Remove an/all alias(es)",
	Run: func(cmd *cobra.Command, args []string) {
		aliases := new(Aliases)
		cobra.CheckErr(aliases.ReadFromFile())

		all, _ := cmd.Flags().GetBool("all")

		if all {
			userInput := "n"

			fmt.Print("Are you sure you want to remove all the aliases? [y/N]: ")
			fmt.Scanln(&userInput)

			if strings.ToLower(userInput) == "y" {
				aliases.RemoveAll()
				cobra.CheckErr(aliases.WriteToFile())

				fmt.Printf("Removed all the aliases from the file: %s\n", AliasFilePath)

				os.Exit(0)
			} else {
				os.Exit(0)
			}
		}

		menu := selector.NewMenu("Choose alias to remove:")

		for alias, _ := range aliases.Aliases {
			menu.AddItem(alias, alias)
		}

		choice := menu.Display()

		if choice == "" {
			os.Exit(0)
		}

		menu.Clear()

		userInput := "n"

		fmt.Printf("Are you sure you want to remove the alias \"%s\"? [y/N]: ", choice)
		fmt.Scanln(&userInput)

		if strings.ToLower(userInput) == "y" {
			aliases.Remove(choice)
			cobra.CheckErr(aliases.WriteToFile())

			fmt.Printf("Removed the alias \"%s\" from the file: %s\n", choice, AliasFilePath)

			os.Exit(0)
		} else {
			os.Exit(0)
		}
		cobra.CheckErr(aliases.WriteToFile())

		fmt.Printf(
			"Successfully removed the alias \"%s\" from the file: %s\n",
			choice,
			AliasFilePath,
		)
	},
}

func init() {
	removeCmd.Flags().BoolP("all", "a", false, "To remove all the aliases")
	AliasCmd.AddCommand(removeCmd)
}
