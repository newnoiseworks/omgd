/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/newnoiseworks/tpl-fred/builder"
	"github.com/newnoiseworks/tpl-fred/deployer"
	"github.com/newnoiseworks/tpl-fred/utils"

	"github.com/spf13/cobra"
)

// cbdAllCmd represents the cbdAll command
var cbdAllCmd = &cobra.Command{
	Use:   "cbd-all",
	Short: "Clones, builds, and deploys the app",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		var environment = args[0]

		CloneLibs(environment, false, false)

		// We "deploy" infra on the first step by checking if we need to make changes.
		// Then, we get the IP address of the server and put that into the yml profile
		deployer.Infra{
			Environment: environment,
			OutputDir:   OutputDir,
			CmdOnDir:    utils.CmdOnDir,
		}.Deploy()

		builder.Game{
			Environment: environment,
			OutputDir:   OutputDir,
			CmdOnDir:    utils.CmdOnDir,
		}.Build()

		builder.Server{
			Environment: environment,
			OutputDir:   OutputDir,
			CmdOnDir:    utils.CmdOnDir,
		}.Build()

		deployer.Game{
			Environment: environment,
			OutputDir:   OutputDir,
			CmdOnDir:    utils.CmdOnDir,
		}.Deploy()

		deployer.Server{
			Environment: environment,
			OutputDir:   OutputDir,
			CmdOnDir:    utils.CmdOnDir,
			VolumeReset: VolumeReset,
		}.Deploy()

	},
}

func init() {
	rootCmd.AddCommand(cbdAllCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cbdAllCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cbdAllCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
