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
	Short: "Clones repositories for build and deploy cycles, or local setup",
	Long: `Clones repos for the build, or for local steup

	Arguments: 

	Example:
	tpl clone --profile=[profile]
	`,
	Run: func(cmd *cobra.Command, args []string) {
		CloneLibs()
	},
}

// CloneLibs clones dem libs
func CloneLibs() {
	conf := utils.GetProfile(Profile)

	repo := conf.Git.Repo

	if repo == "" {
		repo = "."
	}

	err := os.RemoveAll(OutputDir)

	if err != nil {
		log.Fatalf("Failure on removing directory at %s \r\n %s", OutputDir, err)
	}

	gitClone(repo, conf.Git.GameBranch)
}

func gitClone(repo string, confBranch string) {
	branchName := "master"

	if confBranch != "" {
		branchName = confBranch
	}

	refVal := plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branchName))

	_, err := git.PlainClone(OutputDir, false, &git.CloneOptions{
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
