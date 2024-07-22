package alias

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// listCmd represents the alias list command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "List all aliases",
	Long:  `List all aliases`,
	Run: func(cmd *cobra.Command, args []string) {	
		serverAddress := os.Args[3]
		serverAlias := os.Args[4]

		data, err := os.ReadFile(AliasFilePath)

		aliases := new(Aliases)

		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("File not found, creating...")
				err = os.WriteFile(AliasFilePath, []byte("File created!"), 0644)

				if err != nil {
					fmt.Println("Error creating file:", err)
					return
				}

				fmt.Println("File created successfully!")
				
				aliases.Aliases = map[string]Alias{}

			} else {
				fmt.Println("Error reading file:", err)
				return
			}
		} else {
			fmt.Println("File exists!")
			err = yaml.Unmarshal(data, aliases)
			if err != nil {
				log.Println("Error reading the aliases file")

				return
			}
		}		

		_, exist := aliases.Aliases[serverAlias]

		if !exist {
			fmt.Println("Creating a New Server Alias", serverAlias)
			aliases.Aliases[serverAlias] = Alias{
				Name:        serverAlias,
				Value:       serverAddress,
				Description: "Testing12345",
			}
		} else {
			fmt.Println("Server Alias Already Exists!")
		}

		
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
	AliasCmd.AddCommand(addCmd)
}
