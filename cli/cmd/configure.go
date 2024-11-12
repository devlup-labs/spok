package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/devlup-labs/spok/cli/cmd/alias"
	"github.com/devlup-labs/spok/internal/pkg"
	"github.com/spf13/cobra"
)

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configures the target server for OpenPubkey authentication",
	Long:  `This command will configure the target server for OpenPubkey authentication. It will copy the necessary files to the server and run the configuration script.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("---Setup Initiated----")

		serverFlag, _ := cmd.Flags().GetString("server")
		port, _ := cmd.Flags().GetString("port")
		keyFlag, _ := cmd.Flags().GetString("key")
		emailFlag, _ := cmd.Flags().GetString("email")

		emailArgs := emailFlag
		userArgs := serverFlag

		principal := strings.Split(userArgs, "@")[0]

		privateKeyAuth := false
		privateKeyPath := ""

		platform := runtime.GOOS

		if keyFlag != "" {
			fmt.Println("Private Key Mode selected.")

			privateKeyAuth = true
			privateKeyPath = keyFlag
		}

		var configDirs []string
		var configPath string

		switch platform {
		case "windows":
			homeDir, err := os.UserHomeDir()
			cobra.CheckErr(err)

			configDirs = []string{
				homeDir + "/scoop/apps/spok/" + pkg.Version,
				"C:/ProgramData/SPoK",
			}
		default:
			configDirs = []string{
				"/etc/spok", "/opt/homebrew/etc/spok", "/usr/local/etc/spok",
			}
		}

		for i, dir := range configDirs {
			configPath = dir + "/scripts/configure-spok-server.sh"

			_, err := os.Stat(configPath)
			if err == nil {
				break
			} else if i == len(configDirs)-1 {
				log.Fatal("Configuration script not found.")
			}
		}

		fmt.Println("Configuration script found at:", configPath)

		var serverConfigPath string

		var osCommand *exec.Cmd

		if privateKeyAuth && privateKeyPath != "" {
			osCommand = exec.Command(
				"ssh",
				"-p",
				port,
				"-i",
				privateKeyPath,
				userArgs,
				"uname",
				"-o",
			)
		} else {
			osCommand = exec.Command(
				"ssh",
				"-p",
				port,
				userArgs,
				"uname",
				"-o",
			)
		}

		output, err := osCommand.Output()
		cobra.CheckErr(err)

		if strings.Contains(string(output), "Linux") {
			fmt.Println("Remote platform detected: Linux")

			serverConfigPath = "/root/configure-spok-server.sh"
		} else if strings.Contains(string(output), "Darwin") {
			fmt.Println("Remote platform detected: Mac")

			serverConfigPath = "/var/root/configure-spok-server.sh"
		} else {
			log.Fatal("Unsupported remote platform.")
		}

		var scpCommandScript []string
		var sshCommandChmod []string
		var sshCommandConfigure []string

		if privateKeyAuth && privateKeyPath != "" {
			scpCommandScript = []string{
				"scp",
				"-P",
				port,
				"-i",
				privateKeyPath,
				configPath,
				userArgs + ":" + serverConfigPath,
			}
			sshCommandChmod = []string{
				"ssh",
				"-p",
				port,
				"-i",
				privateKeyPath,
				userArgs,
				"chmod",
				"+x",
				serverConfigPath,
			}
			sshCommandConfigure = []string{
				"ssh",
				"-p",
				port,
				"-i",
				privateKeyPath,
				userArgs,
				serverConfigPath,
				emailArgs,
				principal,
				pkg.Version,
			}
		} else {
			scpCommandScript = []string{
				"scp",
				"-P",
				port,
				configPath,
				userArgs + ":" + serverConfigPath,
			}
			sshCommandChmod = []string{
				"ssh",
				"-p",
				port,
				userArgs,
				"chmod",
				"+x",
				serverConfigPath,
			}
			sshCommandConfigure = []string{
				"ssh",
				"-p",
				port,
				userArgs,
				serverConfigPath,
				emailArgs,
				principal,
				pkg.Version,
			}
		}

		scpCmdScript := exec.Command(
			scpCommandScript[0], scpCommandScript[1:]...,
		)

		fmt.Printf(
			"Copying configuration script to server at %s ...\n",
			serverConfigPath,
		)

		scpCmdScript.Stdout = os.Stdout
		scpCmdScript.Stderr = os.Stderr

		err = scpCmdScript.Run()
		if err != nil {
			log.Fatal(err)
		}

		sshCmdChmod := exec.Command(
			sshCommandChmod[0], sshCommandChmod[1:]...,
		)

		fmt.Println("Making configuration script executable...")

		sshCmdChmod.Stdout = os.Stdout
		sshCmdChmod.Stderr = os.Stderr

		err = sshCmdChmod.Run()
		if err != nil {
			log.Fatal(err)
		}

		sshCmdConfigure := exec.Command(
			sshCommandConfigure[0], sshCommandConfigure[1:]...,
		)

		fmt.Println("Configuring SPoK server...")

		sshCmdConfigure.Stdout = os.Stdout
		sshCmdConfigure.Stderr = os.Stderr

		err = sshCmdConfigure.Run()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Configured SPoK server for:", emailArgs)

		userInput := "y"

		fmt.Print("Would you like to add an alias for this server now? [Y/n]: ")
		fmt.Scanln(&userInput)

		if strings.ToLower(userInput) == "n" {
			os.Exit(0)
		} else {
			var aliasName string
			var value string
			var description string

			value = fmt.Sprintf("ssh -p %s %s", port, userArgs)

			fmt.Print("Enter an alias for this server: ")
			fmt.Scanln(&aliasName)

			fmt.Print("Provide a short description for the alias (Press ENTER to leave blank): ")
			reader := bufio.NewReader(os.Stdin)
			description, err = reader.ReadString('\n')
			cobra.CheckErr(err)
			description = description[:len(description)-1]

			aliases := new(alias.Aliases)
			cobra.CheckErr(aliases.ReadFromFile())

			aliases.Add(aliasName, value, description)
			cobra.CheckErr(aliases.WriteToFile())

			fmt.Printf("Successfully added the alias to the file: %s\n", alias.AliasFilePath)
		}
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)

	configureCmd.Flags().StringP(
		"server", "s", "", "The server you want to configure in",
	)
	configureCmd.Flags().StringP(
		"email", "e", "", "The email you want to configure with",
	)
	configureCmd.Flags().StringP(
		"key", "i", "", "To choose if you want to use private key for setup",
	)
	configureCmd.Flags().StringP(
		"port", "p", "22", "The port to connect to the server on",
	)
}
