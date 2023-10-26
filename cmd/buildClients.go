/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

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

		buildFor := strings.ReplaceAll(BuildTargets, ",", " ")

		utils.CmdOnDirWithEnv(
			// TODO: break below into optional builds per OS based on... profile probably?
			fmt.Sprintf("docker compose up %s", buildFor),
			fmt.Sprintf("Building %s game clients into game/dist folder", profile.Name),
			"game",
			[]string{
				fmt.Sprintf("BUILD_ENV=%s", profile.Name),
			},
		)
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
		"build-web,build-windows,build-mac,build-x11",
		"specify build targets comma separated no spaces",
	)

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildClientsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
