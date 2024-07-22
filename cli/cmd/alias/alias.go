package alias

import (
	"os"

	"github.com/spf13/cobra"
)

type Alias struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Value       string `yaml:"value"`
}

type Aliases struct {
	Aliases map[string]Alias `yaml:"aliases"`
}

func (alias *Alias) Update(newAddress string) {
	alias.Value = newAddress
}

var AliasFilePath string

// AliasCmd represents the base alias command
var AliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "Use server aliases",
	Long:  "Use server aliases to connect to servers",
}

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}

	AliasFilePath = homeDir + "/.spok_aliases"
}
