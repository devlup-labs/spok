package alias

import (
	"fmt"
	"os"

	"github.com/devlup-labs/spok/internal/pkg/selector"
	"github.com/spf13/cobra"
)

// listCmd represents the alias list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all aliases",
	Long:  `List all aliases`,
	Run: func(cmd *cobra.Command, args []string) {
		aliases := new(Aliases)
		cobra.CheckErr(aliases.ReadFromFile())

		if len(aliases.Aliases) == 0 {
			fmt.Printf("No aliases found in the file: %s\n", AliasFilePath)

			os.Exit(0)
		}

		menu := selector.NewMenu("List of aliases:")

		for alias, value := range aliases.Aliases {
			menu.AddItem(alias, value.Value)
		}

		choice := menu.Display()
		menu.Clear()

		fmt.Println(choice)
	},
}

func init() {
	AliasCmd.AddCommand(listCmd)
}
