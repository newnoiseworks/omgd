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

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Manages OMGD Docker containers",
	Long: `Manages OMGD Docker containers

$ omgd server start          | starts local docker server containers
$ omgd server stop           | stops local docker server containers
$ omgd server reset-data     | stops containers and resets the data volumes
$ omgd server logs           | prints logs from docker containers
$ omgd server logs --verbose | tails / follows logs continuously
$ omgd server status         | prints status of running docker containers
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
			utils.CmdOnDir(
				fmt.Sprintf("docker compose -p %s-%s-servers up %s -d", profile.OMGD.Name, profile.Name, services),
				fmt.Sprintf("spinning up docker containers"),
				"servers",
			)
		case "stop":
			utils.CmdOnDir(
				fmt.Sprintf("docker compose -p %s-%s-servers down", profile.OMGD.Name, profile.Name),
				fmt.Sprintf("stopping docker containers"),
				"servers",
			)
		case "reset-data":
			utils.CmdOnDir(
				fmt.Sprintf("docker compose -p %s-%s-servers down -v", profile.OMGD.Name, profile.Name),
				fmt.Sprintf("removing data volumes and stopping docker containers"),
				"servers",
			)
		case "logs":
			if utils.GetEnvLogLevel() < utils.DEBUG_LOG {
				utils.SetEnvLogLevel(utils.DEBUG_LOG)
			}

			utils.CmdOnDirToStdOut(
				fmt.Sprintf("docker compose -p %s-%s-servers logs --follow", profile.OMGD.Name, profile.Name),
				"printing server logs",
				"servers",
				[]string{},
			)
		case "status":
			if utils.GetEnvLogLevel() < utils.DEBUG_LOG {
				utils.SetEnvLogLevel(utils.DEBUG_LOG)
			}

			utils.CmdOnDir(
				fmt.Sprintf("docker compose -p %s-%s-servers ps", profile.OMGD.Name, profile.Name),
				fmt.Sprintf("printing server status"),
				"servers",
			)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
