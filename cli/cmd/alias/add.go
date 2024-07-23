package alias

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// addCmd represents the alias add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new alias",
	Long:  "Add a new alias to the aliases file",
	Run: func(cmd *cobra.Command, args []string) {
		var alias string
		var value string
		var description string

		if len(os.Args) < 4 {
			fmt.Println("Invalid number of arguments")

			os.Exit(1)
		}

		value = os.Args[3]

		fmt.Print("Please provide a short, memorable name for the alias: ")
		_, err := fmt.Scanln(&alias)
		cobra.CheckErr(err)

		fmt.Print("Provide a short description for the alias (Press ENTER to leave blank): ")
		reader := bufio.NewReader(os.Stdin)
		description, err = reader.ReadString('\n')
		cobra.CheckErr(err)
		description = description[:len(description)-1]

		aliases := new(Aliases)
		cobra.CheckErr(aliases.ReadFromFile())

		aliases.Add(alias, value, description)
		cobra.CheckErr(aliases.WriteToFile())

		fmt.Printf("Successfully added the alias to the file: %s\n", AliasFilePath)
	},
}

func init() {
	AliasCmd.AddCommand(addCmd)
}
