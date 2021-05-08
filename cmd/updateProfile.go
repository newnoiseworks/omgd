/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"strings"

	"github.com/newnoiseworks/tpl-fred/utils"
	"github.com/spf13/cobra"
)

// updateProfileCmd represents the updateProfile command
var updateProfileCmd = &cobra.Command{
	Use:   "update-profile",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		profile := utils.GetProfile(Profile)
		profileMap := profile.GetProfileAsMap()

		keys := strings.Split(args[0], ".")

		// TODO: The below should work recursively
		for k, v := range profileMap {
			if key, ok := k.(string); ok {
				if key == keys[0] {
					for kA := range v.(map[interface{}]interface{}) {
						if keyTwo, ok := kA.(string); ok {
							if keyTwo == keys[1] {
								v.(map[interface{}]interface{})[keyTwo] = args[1]
							}
						}
					}
				}
			}
		}

		profile.SaveProfileFromMap(&profileMap)
	},
}

func init() {
	rootCmd.AddCommand(updateProfileCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateProfileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateProfileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
