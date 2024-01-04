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

$ omgd infra project-setup | Initial one time project level infra setup
$ omgd infra project-destroy | Destroy initial project infra setup
$ omgd infra instance-deploy | Sets up a cloud VM instance of the game servers against a non local profile
$ omgd infra instance-destroy | Destroys cloud VM created against supplied profile`,
	Run: func(cmd *cobra.Command, args []string) {
		profile := utils.GetProfile(ProfilePath)
		command := args[0]

		if (profile.Name == "local" || profile.Name == "omgd") && (command != "project-setup" && command != "project-destroy") {
			utils.LogFatal("Cannot run infra commands against local or top level omgd profile, please supply a profile with -p")
		}

		infraChange := utils.InfraChange{
			OutputDir:       OutputDir,
			Profile:         profile,
			CmdOnDir:        utils.CmdOnDir,
			CmdOnDirWithEnv: utils.CmdOnDirWithEnv,
		}

		switch command {
		case "instance-setup":
			if profile.OMGD.GCP.Bucket == "" || profile.OMGD.GCP.Bucket == "???" {
				utils.LogWarn("No bucket setup in omgd.cloud.yml file -- you need to run omgd infra project-setup first.")
				return
			}

			utils.LogInfo("Setting up instance on cloud servers... WARNING: OMGD is not responsible for managing your server costs.")
			infraChange.InstanceSetup()
		case "instance-destroy":
			if profile.OMGD.GCP.Bucket == "" || profile.OMGD.GCP.Bucket == "???" {
				utils.LogWarn("No bucket setup in omgd.cloud.yml file -- you need to run omgd infra project-setup first.")
				return
			}

			utils.LogInfo("Destroying instance on cloud servers...")
			infraChange.InstanceDestroy()
		case "project-setup":
			if profile.OMGD.GCP.Project == "" || profile.OMGD.GCP.Project == "???" || profile.OMGD.GCP.Project == "your-project-name" {
				utils.LogWarn("No GCP project setup in omgd.cloud.yml file -- you need to create a GCP project and insert it's ID in that file")
				return
			}

			utils.LogInfo("Setting up OMGD project to work with cloud servers... WARNING: OMGD is not responsible for managing your server costs.")
			infraChange.ProjectSetup()
		case "project-destroy":
			if profile.OMGD.GCP.Project == "" || profile.OMGD.GCP.Project == "???" || profile.OMGD.GCP.Project == "your-project-name" {
				utils.LogWarn("No GCP project setup in omgd.cloud.yml file -- you need to create a GCP project and insert it's ID in that file")
				return
			}

			if profile.OMGD.GCP.Bucket == "" || profile.OMGD.GCP.Bucket == "???" {
				utils.LogWarn("No bucket setup in omgd.cloud.yml file -- you need to run omgd infra project-setup first.")
				return
			}

			utils.LogInfo("Destroying OMGD project setup on cloud servers...")
			infraChange.ProjectDestroy()
		default:
			utils.LogWarn(fmt.Sprintf("Found no infra command for %s", args[0]))
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
