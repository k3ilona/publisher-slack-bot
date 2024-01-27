/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// promoteCmd represents the promote command
var promoteCmd = &cobra.Command{
	Use:   "promote",
	Short: "Go to the new version and confirm the successful execution of the operation.",
	Long: `This command allows the user to upgrade the application version to the next environment. 
	On command, the bot will perform the necessary actions to switch to the new version 
	and confirm the successful completion of the operation.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("promote called")
	},
}

func init() {
	rootCmd.AddCommand(promoteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// promoteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// promoteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
