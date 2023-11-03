/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/newnoiseworks/omgd/utils"
	"github.com/spf13/cobra"
)

// BuildTargets sets the targets to be built
var BuildTargets string

// buildClientsCmd represents the buildClients command
var buildClientsCmd = &cobra.Command{
	Use:   "build-clients",
	Short: "Builds game clients into the game/dist folder.",
	// Long:  `Builds local game clients into the game/dist folder based on your local profile.`,
	Run: func(cmd *cobra.Command, args []string) {
		profile := utils.GetProfile(ProfilePath)

		cb := utils.ClientBuilder{
			Profile:         profile,
			CmdOnDirWithEnv: utils.CmdOnDirWithEnv,
		}

		cb.Build()
	},
}

func init() {
	rootCmd.AddCommand(buildClientsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildClientsCmd.PersistentFlags().String("foo", "", "A help for foo")
	buildClientsCmd.PersistentFlags().StringVar(
		&BuildTargets,
		"targets",
		"",
		"specify build targets as listed in game/docker-compose.yml separated by a single space",
	)

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildClientsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
