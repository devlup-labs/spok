/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"log"
	"bufio"
	"gopkg.in/yaml.v2"
	"github.com/spf13/cobra"
)

type Alias struct {
	Details map[string]string `yaml:"alias"`
}

func getFileSize(filename string) (int64, error) {
	fileinfo, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}

	return fileinfo.Size(), nil
}

func (p *Alias) Unmarshal(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fileSize, _ := getFileSize(filename)
	data := make([]byte, fileSize)

	reader := bufio.NewReader(file)

	_, err = reader.Read(data)
	if err != nil {
		fmt.Println(err)

		return err
	}

	return yaml.Unmarshal(data, &p)
}

// aliasCmd represents the alias command
var aliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "SPoK functionality to introduce server aliases",
	Long: `It is for setting up server aliases for easier SSH experience.
		It has the following functionalities:
		- add : To add new server aliases
		- remove : To remove server aliases
		- list: To list server aliases
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Alias Called")
		
		checkArg := os.Args[2]

		if checkArg == "add" {
			fmt.Println(checkArg)
			serverAddress := os.Args[3]
			serverAlias := os.Args[4]
			
			policy := new(Alias)

			err := policy.Unmarshal("./aliases.yml")
			if err != nil {
				log.Fatal(err)
			}
			
			_, exist := policy.Details[serverAlias]

			if !exist {
				policy.Details[serverAlias] = serverAddress
			}
			
			f, err := os.Create("./aliases.yml")
			
			if err != nil {
				fmt.Println("File not Found aliases.yml")

				return
			}
			defer f.Close()

			yamlData, err := yaml.Marshal(&policy)
			if err != nil {
				fmt.Println("Error while Marshalling. ", err)

				return
			}

			_, err = f.WriteString(string(yamlData))
			if err != nil {
				fmt.Println(err)
				
				return
			}
		} else if checkArg == "list" {
			policy := new(Alias)

			err := policy.Unmarshal("./aliases.yml")
			if err != nil {
				log.Fatalln(err)
			}
			
			for k, v := range policy.Details {
				fmt.Println("Alias:",k,"Server:",v)
			}
		} else if checkArg == "remove" {
			serverAlias := os.Args[3]
			policy := new(Alias)

			err := policy.Unmarshal("./aliases.yml")
			if err != nil {
				log.Fatal(err)
			}
			
			_, exist := policy.Details[serverAlias]
			
			if exist {
				fmt.Println("Deleting Entry: ", serverAlias)
				delete(policy.Details, serverAlias)
				f, err := os.Create("./aliases.yml")
			
				if err != nil {
					fmt.Println("File not Found aliases.yml")

					return
				}
				defer f.Close()

				yamlData, err := yaml.Marshal(&policy)
				if err != nil {
					fmt.Println("Error while Marshalling. ", err)

					return
				}

				_, err = f.WriteString(string(yamlData))
				if err != nil {
					fmt.Println(err)
					
					return
				}
			} else {
				fmt.Println("Entry Does not Exist: ", serverAlias)
			}

		}



	},
}

func init() {
	rootCmd.AddCommand(aliasCmd)
}
