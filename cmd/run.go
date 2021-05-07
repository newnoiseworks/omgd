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
	"fmt"

	"github.com/newnoiseworks/tpl-fred/utils"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs commands based on the chosen profile",
	Long: `Works in conjunction with the --profile you've set, or profiles/local.yml by default. Runs either the base project command or any tasks you specify if chosen.

Usage: 

Run entire project:
$ gg run

Run part of project:
$ gg run [name-of-project-step] [number-of-step (optional)]

Run task:
$ gg run task [name-of-task] [number-of-step]
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run called")

		profile := utils.GetProfile(Profile)

		dir := OutputDir

		for _, project := range profile.Project {
			if project.Dir != "" {
				dir = fmt.Sprintf("%s/%s", dir, project.Dir)
			}

			for _, step := range project.Steps {
				if step.Dir != "" {
					dir = fmt.Sprintf("%s/%s", dir, step.Dir)
				}

				utils.CmdOnDir(step.Cmd, "", dir)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
