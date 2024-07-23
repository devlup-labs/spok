package alias

import (
	"bufio"
	"fmt"
	"os"

	"github.com/devlup-labs/spok/internal/pkg/selector"
	"github.com/spf13/cobra"
)

// updateCmd represents the alias update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an alias",
	Long:  "Update an alias",
	Run: func(cmd *cobra.Command, args []string) {
		aliases := new(Aliases)
		cobra.CheckErr(aliases.ReadFromFile())

		menu := selector.NewMenu("Choose your alias:")

		for alias, _ := range aliases.Aliases {
			menu.AddItem(alias, alias)
		}

		aliasToUpdate := menu.Display()
		menu.Clear()

		menu = selector.NewMenu("What do you want to update?")

		menu.AddItem("Name", "name")
		menu.AddItem("Value", "value")
		menu.AddItem("Description", "description")

		fieldToUpdate := menu.Display()
		menu.Clear()

		var fieldNewValue string

		fmt.Printf("Enter a new %s for the alias (Leave blank to keep unchanged): ", fieldToUpdate)
		reader := bufio.NewReader(os.Stdin)
		fieldNewValue, err := reader.ReadString('\n')
		cobra.CheckErr(err)
		fieldNewValue = fieldNewValue[:len(fieldNewValue)-1]

		if ok := aliases.Update(aliasToUpdate, fieldToUpdate, fieldNewValue); !ok {
			os.Exit(0)
		}
		cobra.CheckErr(aliases.WriteToFile())

		fmt.Printf("Successfully updated the alias in the file: %s\n", AliasFilePath)
	},
}

func init() {
	AliasCmd.AddCommand(updateCmd)
}
