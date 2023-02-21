/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Manages OMGD Docker containers",
	Long: `Manages OMGD Docker containers

$ omgd server start -- starts local docker server containers
$ omgd server stop -- stops local docker server containers
$ omgd server reset-data -- resets the data volumes
	`,
	Run: func(cmd *cobra.Command, args []string) {
		Profile = strings.ReplaceAll(Profile, "profiles/", ".omgd/")

		switch args[0] {
		case "start":
			runCmd.Run(cmd, []string{
				"task", "start-server",
			})
		case "stop":
			runCmd.Run(cmd, []string{
				"task", "stop-server",
			})
		case "reset-data":
			runCmd.Run(cmd, []string{
				"task", "reset-server-data",
			})
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
