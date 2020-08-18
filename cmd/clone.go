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
	"fmt"
	"log"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/newnoiseworks/tpl-fred/utils"
	"github.com/spf13/cobra"
)

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var environment = args[0]

		conf := utils.GetProfile(environment)
		// cloneTPLFred()
		cloneGame(conf)
		cloneServer(conf)
		// cloneWebsite()
	},
}

func cloneGame(conf utils.ProfileConf) {
	gitClone("git@github.com:newnoiseworks/not-stardew.git", "game", conf.Git.GameBranch)
}

func cloneServer(conf utils.ProfileConf) {
	gitClone("git@github.com:newnoiseworks/not-stardew-backend.git", "server", conf.Git.ServerBranch)
}

func cloneWebsite() {
	gitClone("git@github.com:newnoiseworks/tpl-website.git", "website", "")
}

func cloneTPLFred() {
	gitClone("git@github.com:newnoiseworks/tpl-fred.git", "tpl-fred", "")
}

func gitClone(repo string, dir string, confBranch string) {
	branchName := "master"

	if confBranch != "" {
		branchName = confBranch
	}

	refVal := plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branchName))

	_, err := git.PlainClone(fmt.Sprintf("%s/%s", OutputDir, dir), false, &git.CloneOptions{
		URL:           repo,
		Progress:      os.Stdout,
		ReferenceName: refVal,
	})

	if err != nil {
		log.Fatal(err)
	}

}

func init() {
	rootCmd.AddCommand(cloneCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cloneCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cloneCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
