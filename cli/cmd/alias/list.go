package alias

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

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

		if run, _ := cmd.Flags().GetBool("run"); run {
			sshCommand := strings.Split(choice, " ")

			cmd := exec.Command(sshCommand[0], sshCommand[1:]...)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cobra.CheckErr(cmd.Run())
		} else {
			fmt.Println(choice)
		}
	},
}

func init() {
	listCmd.Flags().BoolP("run", "r", false, "Run the selected alias")
	AliasCmd.AddCommand(listCmd)
}
