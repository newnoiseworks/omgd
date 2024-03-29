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

var volumeDrop = false

// serversCmd represents the servers command
var serversCmd = &cobra.Command{
	Use:   "servers",
	Short: "Manages OMGD Docker containers",
	Long: `Manages OMGD Docker containers

$ omgd servers start          | starts local docker servers containers
$ omgd servers stop           | stops local docker servers containers
$ omgd servers logs           | prints logs from docker containers
$ omgd servers status         | prints status of running docker containers
$ omgd servers deploy         | deploys servers folder to profile target, not usable with local profile
	`,
	Run: func(cmd *cobra.Command, args []string) {
		profile := utils.GetProfile(ProfilePath)

		serviceArray := []string{}

		for _, service := range profile.OMGD.Servers.Services {
			serviceArray = append(serviceArray, service.BuildService)
		}

		services := strings.Join(serviceArray, " ")

		switch args[0] {
		case "start":
			utils.LogInfo("Starting OMGD servers containers...")

			if profile.Name == "local" {
				utils.CmdOnDir(
					fmt.Sprintf("docker compose -p %s-%s-servers up %s -d", profile.OMGD.Name, profile.Name, services),
					fmt.Sprintf("spinning up docker containers"),
					"servers",
				)
			} else if profile.OMGD.Servers.Host != "" {
				cmd := "docker compose up -d"

				if profile.OMGD.GCP.Bucket != "" {
					utils.RemoteGCPCommand(cmd, OutputDir, profile)
				}
			}
		case "stop":
			dropVolumeArg := ""

			if volumeDrop {
				dropVolumeArg = " -v"
				utils.LogInfo("Stopping OMGD servers containers and dropping data volumes...")
			} else {
				utils.LogInfo("Stopping OMGD servers containers...")
			}

			if profile.Name == "local" {
				utils.CmdOnDir(
					fmt.Sprintf("docker compose -p %s-%s-servers down%s", profile.OMGD.Name, profile.Name, dropVolumeArg),
					fmt.Sprintf("stopping docker containers"),
					"servers",
				)
			} else if profile.OMGD.Servers.Host != "" {
				cmd := fmt.Sprintf("docker compose down%s", dropVolumeArg)

				if profile.OMGD.GCP.Bucket != "" {
					utils.RemoteGCPCommand(cmd, OutputDir, profile)
				}
			}
		case "logs":
			utils.LogInfo("Getting OMGD servers logs...")

			if utils.GetEnvLogLevel() < utils.DEBUG_LOG {
				utils.SetEnvLogLevel(utils.DEBUG_LOG)
			}

			if profile.Name == "local" {
				utils.CmdOnDirToStdOut(
					fmt.Sprintf("docker compose -p %s-%s-servers logs --follow", profile.OMGD.Name, profile.Name),
					"printing server logs",
					"servers",
					[]string{},
				)
			} else if profile.OMGD.Servers.Host != "" {
				cmd := "docker compose logs"

				if profile.OMGD.GCP.Bucket != "" {
					utils.RemoteGCPCommand(cmd, OutputDir, profile)
				}
			}
		case "status":
			utils.LogInfo("Getting status of servers containers...")

			if utils.GetEnvLogLevel() < utils.DEBUG_LOG {
				utils.SetEnvLogLevel(utils.DEBUG_LOG)
			}

			if profile.Name == "local" {
				utils.CmdOnDir(
					fmt.Sprintf("docker compose -p %s-%s-servers ps", profile.OMGD.Name, profile.Name),
					fmt.Sprintf("printing servers status"),
					"servers",
				)
			} else if profile.OMGD.Servers.Host != "" {
				cmd := "docker compose ps"

				if profile.OMGD.GCP.Bucket != "" {
					utils.RemoteGCPCommand(cmd, OutputDir, profile)
				}
			}
		case "deploy":
			if profile.Name == "local" {
				utils.LogWarn("omgd servers deploy requires a non local yml profile to be passed in with -p")
				return
			}

			utils.LogInfo("Deploying servers to cloud... WARNING: OMGD is not responsible for managing your server costs.")

			sc := utils.ServersChange{
				OutputDir:       OutputDir,
				Profile:         profile,
				CmdOnDir:        utils.CmdOnDir,
				CmdOnDirWithEnv: utils.CmdOnDirWithEnv,
			}

			sc.Deploy()
		}
	},
}

func init() {
	rootCmd.AddCommand(serversCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serversCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serversCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	serversCmd.Flags().BoolVarP(&volumeDrop, "volume-drop", "v", false, "Used with omgd servers stop -v, when passed will drop data volumes with containers, resetting data on next server start")
}
