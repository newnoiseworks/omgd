/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

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
		switch args[0] {
		case "start":
			utils.CmdOnDir(
				"docker-compose up -d",
				fmt.Sprintf("spinning up docker containers"),
				"server",
				true,
			)
		case "stop":
			utils.CmdOnDir(
				"docker-compose down",
				fmt.Sprintf("stopping docker containers"),
				"server",
				true,
			)
		case "reset-data":
			utils.CmdOnDir(
				"docker-compose down -v",
				fmt.Sprintf("removing data volumes and stopping docker containers"),
				"server",
				true,
			)
		case "logs":
			cmd := "docker-compose logs"

			if Verbosity {
				cmd = "docker-compose logs --follow"
			}

			utils.CmdOnDir(
				cmd,
				fmt.Sprintf("printing server logs via $ %s", cmd),
				"server",
				true,
			)
		case "status":
			utils.CmdOnDir(
				"docker-compose ps",
				fmt.Sprintf("printing server status via $ docker-compose ps"),
				"server",
				true,
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
