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
var removeCmd = &cobra.Command{
	Use:   "remove",
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

		menu := selector.NewMenu("Choose Alias for removal:")

		for aliasKey, aliasVal := range aliases.Aliases {
			menu.AddItem(aliasKey, aliasVal.Name)
		}

		choice := menu.Display()

		numLinesToClear := len(menu.MenuItems) + 1
		selector.ClearMenu(numLinesToClear)
		
		delete(aliases.Aliases, choice)
	
		f, err := os.Create(AliasFilePath)
	
		if err != nil {
			fmt.Println("File not Found aliases.yml")

			return
		}
		defer f.Close()

		yamlData, err := yaml.Marshal(&aliases)
		if err != nil {
			fmt.Println("Error while Marshalling. ", err)

			return
		}

		_, err = f.WriteString(string(yamlData))
		if err != nil {
			fmt.Println(err)
			
			return
		}
},
}

func init() {
	AliasCmd.AddCommand(removeCmd)
}
