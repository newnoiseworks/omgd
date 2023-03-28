/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/newnoiseworks/omgd/utils"
	"github.com/spf13/cobra"
)

// buildClientsCmd represents the buildClients command
var buildClientsCmd = &cobra.Command{
	Use:   "build-clients",
	Short: "Builds game clients into the game/dist folder.",
	// Long:  `Builds local game clients into the game/dist folder based on your local profile.`,
	Run: func(cmd *cobra.Command, args []string) {
		profile := utils.GetProfile(ProfilePath)

		utils.CmdOnDirWithEnv(
			// TODO: break below into optional builds per OS based on... profile probably?
			"docker compose up build-web build-windows build-mac build-x11",
			fmt.Sprintf("Building local game clients into game/dist folder"),
			"game",
			[]string{
				fmt.Sprintf("BUILD_ENV=%s", profile.Name),
			},
			Verbosity,
		)
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
