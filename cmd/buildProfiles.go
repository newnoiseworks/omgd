/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/newnoiseworks/tpl-fred/utils"
	"github.com/spf13/cobra"
)

// buildProfilesCmd represents the buildProfiles command
var buildProfilesCmd = &cobra.Command{
	Use:   "build-profiles",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		files, err := ioutil.ReadDir("profiles")
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			splits := strings.Split(file.Name(), ".")
			ext := splits[len(splits)-1]

			if ext == "yml" && splits[0] != "example" {
				utils.BuildTemplateFromPath(
					".gg/profile.yml.omgdtpl",
					fmt.Sprintf("profiles/%s", strings.Replace(file.Name(), ".yml", "", 1)),
					OutputDir,
					"omgdtpl",
				)

				wd, err := os.Getwd()

				if err != nil {
					log.Fatal(err)
				}

				err = os.Rename(
					fmt.Sprintf("/%s/.gg/profile.yml", wd),
					fmt.Sprintf("/%s/.gg/%s", wd, file.Name()),
				)

				if err != nil {
					log.Fatal(err)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(buildProfilesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildProfilesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildProfilesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
