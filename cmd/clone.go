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
		// cloneTPLFred()
		cloneGame()
		cloneServer()
		cloneWebsite()
	},
}

func cloneTPLFred() {
	_, err := git.PlainClone(fmt.Sprintf("%s/tpl-fred", OutputDir), false, &git.CloneOptions{
		URL:      "git@github.com:newnoiseworks/tpl-fred.git",
		Progress: os.Stdout,
	})

	if err != nil {
		log.Fatal(err)
	}
}

func cloneGame() {
	_, err := git.PlainClone(fmt.Sprintf("%s/game", OutputDir), false, &git.CloneOptions{
		URL:           "git@github.com:newnoiseworks/not-stardew.git",
		Progress:      os.Stdout,
		ReferenceName: "refs/heads/golang-refactor",
	})

	if err != nil {
		log.Fatal(err)
	}
}

func cloneServer() {
	_, err := git.PlainClone(fmt.Sprintf("%s/server", OutputDir), false, &git.CloneOptions{
		URL:           "git@github.com:newnoiseworks/not-stardew-backend.git",
		Progress:      os.Stdout,
		ReferenceName: "refs/heads/golang-build-refactor",
	})

	if err != nil {
		log.Fatal(err)
	}
}

func cloneWebsite() {
	_, err := git.PlainClone(fmt.Sprintf("%s/website", OutputDir), false, &git.CloneOptions{
		URL:      "git@github.com:newnoiseworks/tpl-website.git",
		Progress: os.Stdout,
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
