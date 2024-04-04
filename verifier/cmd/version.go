package cmd

import (
	"fmt"

	"github.com/devlup-labs/spok/internal/pkg"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "SPoK Verifier version",
	Long:  `SPoK Verifier version; must be the same as the version of the spok tool`,
	Run: func (cmd *cobra.Command, args []string)  {
		fmt.Printf("SPoK Verifier Version - v%s\n", pkg.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
