/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/newnoiseworks/omgd/utils"
	"github.com/spf13/cobra"
)

var CopyToTmpDir bool

// infraCmd represents the infra command
var infraCmd = &cobra.Command{
	Use:   "infra",
	Short: "Deploys and destroys cloud infrastructure",
	Long: `Deploys and destroys cloud infrastructure

$ omgd infra deploy | Deploys cloud infrastructure via terraform
$ omgd infra game-deploy | Builds and deploys clients and server to infra
$ omgd infra destroy | Destroys cloud infrastructure via terraform
$ omgd infra project-setup | Initial one time project level infra setup`,
	Run: func(cmd *cobra.Command, args []string) {
		profile := utils.GetProfile(ProfilePath)
		command := args[0]

		if (profile.Name == "local" || profile.Name == "omgd") && command != "project-setup" {
			utils.LogFatal("Cannot run infra commands against local or top level omgd profile, please supply a profile with -p")
		}

		infraChange := utils.InfraChange{
			OutputDir:       OutputDir,
			Profile:         profile,
			CmdOnDir:        utils.CmdOnDir,
			CmdOnDirWithEnv: utils.CmdOnDirWithEnv,
		}

		switch command {
		case "deploy":
			infraChange.DeployInfra()
		case "game-deploy":
			infraChange.DeployClientAndServer()
		case "destroy":
			infraChange.DestroyInfra()
		case "project-setup":
			infraChange.ProjectSetup()
		default:
			utils.LogFatal(fmt.Sprintf("Found no infra command for %s", args[0]))
			utils.LogWarn("hello")
		}
	},
}

func init() {
	rootCmd.AddCommand(infraCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infraCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infraCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
