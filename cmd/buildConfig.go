package cmd

import (
	"log"

	"github.com/newnoiseworks/tpl-fred/builder/config"
	"github.com/spf13/cobra"
)

// buildConfigCmd represents the buildConfig command
var buildConfigCmd = &cobra.Command{
	Use:   "build-config",
	Short: "Builds necessary config files for apps within The Promised Land",
	Long: `tpl-fred build-config command

Builds all config files and artifacts needed to build 
and run full applications

Usage: $ tpl-fred build-config [project]
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatal("Specify project in argument, see omgd build-config --help")
		}

		var project = args[0]

		switch project {
		case "game":
			config.GameConfig(Profile, OutputDir)
			break
		case "server":
			config.ServerConfig(Profile, OutputDir)
			break
		case "infra":
			config.InfraConfig(Profile, OutputDir)
			break
		}
	},
}

func init() {
	rootCmd.AddCommand(buildConfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildConfigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildConfigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
