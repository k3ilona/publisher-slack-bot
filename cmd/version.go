package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var appVersion = "Version"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print IlonaBot version",
	Long:  `print IlonaBot version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(appVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
