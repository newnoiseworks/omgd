/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// buildClientsCmd represents the buildClients command
var buildClientsCmd = &cobra.Command{
	Use:   "build-clients",
	Short: "Builds local game clients into the game/dist folder based on your local profile.",
	// Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("buildClients called")
		Profile = ".gg/local"
		runCmd.Run(cmd, []string{})
	},
}

func init() {
	rootCmd.AddCommand(buildClientsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildClientsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildClientsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
