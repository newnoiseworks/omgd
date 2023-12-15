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

// gameCmd represents the game command
var gameCmd = &cobra.Command{
	Use:   "game",
	Short: "Manages game clients.",
	Long: `Manages game clients.

$ omgd game build | builds game clients into the game/dist folder based on provided profile
	`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.LogInfo("Building game clients and copying to servers directory as specified in profile...")

		profile := utils.GetProfile(ProfilePath)

		switch args[0] {
		case "build":
			buildTemplatesCmd.Run(cmd, args)

			sccp := utils.StaticCodeCopyPlan{}

			cb := utils.ClientBuilder{
				Profile:             profile,
				CmdOnDirWithEnv:     utils.CmdOnDirWithEnv,
				CopyStaticDirectory: sccp.CopyStaticDirectory,
				Targets:             BuildTargets,
			}

			cb.Build()
		}
	},
}

func init() {
	rootCmd.AddCommand(gameCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gameCmd.PersistentFlags().String("foo", "", "A help for foo")
	gameCmd.PersistentFlags().StringVar(
		&BuildTargets,
		"targets",
		"",
		"specify build targets as listed in game/docker-compose.yml separated by a single space",
	)

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gameCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
