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
$ omgd infra destroy | Destroys cloud infrastructure via terraform`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("infra called")

		profile := utils.GetProfile(ProfilePath)

		infraChange := utils.InfraChange{
			OutputDir:       OutputDir,
			Profile:         profile,
			CmdOnDir:        utils.CmdOnDir,
			CmdOnDirWithEnv: utils.CmdOnDirWithEnv,
			Verbosity:       Verbosity,
			CopyToTmpDir:    CopyToTmpDir,
		}

		switch args[0] {
		case "deploy":
			infraChange.DeployInfra()
		case "game-deploy":
			infraChange.DeployClientAndServer()
		case "destroy":
			infraChange.DestroyInfra()
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
	infraCmd.Flags().BoolVar(&CopyToTmpDir, "copy", false, "Copies the project to an .omgdtmp folder and executes all commands there")
}
