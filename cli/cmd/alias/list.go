package alias

import (
	"fmt"
	"log"
	"os"

	"github.com/devlup-labs/spok/internal/pkg/selector"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// listCmd represents the alias list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all aliases",
	Long:  `List all aliases`,
	Run: func(cmd *cobra.Command, args []string) {
		data, err := os.ReadFile(AliasFilePath)
		if err != nil {
			return
		}

		aliases := new(Aliases)

		err = yaml.Unmarshal(data, aliases)
		if err != nil {
			log.Println("Error reading the aliases file")

			return
		}

		menu := selector.NewMenu("Choose your alias:")

		for _, alias := range aliases.Aliases {
			menu.AddItem(alias.Name, alias.Value)
		}

		choice := menu.Display()

		numLinesToClear := len(menu.MenuItems) + 1
		selector.ClearMenu(numLinesToClear)

		fmt.Println(choice)
	},
}

func init() {
	AliasCmd.AddCommand(listCmd)
}
