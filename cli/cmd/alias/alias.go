package alias

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type Alias struct {
	Description string `yaml:"description"`
	Value       string `yaml:"value"`
}

type Aliases struct {
	Aliases map[string]Alias `yaml:"aliases"`
}

func (aliases *Aliases) ReadFromFile() error {
	data, err := os.ReadFile(AliasFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		fmt.Println("Error reading the aliases file")

		return err
	}

	err = yaml.Unmarshal(data, aliases)
	if err != nil {
		fmt.Println("Invalid yaml file")

		return err
	}

	return nil
}

func (aliases *Aliases) Add(name, value, description string) {
	if _, ok := aliases.Aliases[name]; ok {
		fmt.Printf(
			"An alias with the name \"%s\" already exists. Please try with a different name",
			name,
		)

		return
	}

	aliases.Aliases[name] = Alias{
		Value:       value,
		Description: description,
	}
}

func (aliases *Aliases) Update(name, fieldToUpdate, fieldNewValue string) bool {
	if fieldNewValue == "" {
		return false
	}

	switch fieldToUpdate {
	case "name":
		aliases.Aliases[fieldNewValue] = aliases.Aliases[name]

		delete(aliases.Aliases, name)
	case "value":
		description := aliases.Aliases[name].Description

		aliases.Aliases[name] = Alias{
			Value:       fieldNewValue,
			Description: description,
		}
	case "description":
		value := aliases.Aliases[name].Value

		aliases.Aliases[name] = Alias{
			Value:       value,
			Description: fieldNewValue,
		}
	}

	return true
}

func (aliases *Aliases) Remove(name string) {
	delete(aliases.Aliases, name)
}

func (aliases *Aliases) RemoveAll() {
	aliases.Aliases = map[string]Alias{}
}

func (aliases *Aliases) WriteToFile() error {
	data, err := yaml.Marshal(aliases)
	if err != nil {
		fmt.Println("Error marshalling the aliases")

		return err
	}

	err = os.WriteFile(AliasFilePath, data, 0644)
	if err != nil {
		fmt.Println("Error writing the aliases file")

		return err
	}

	return nil
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
