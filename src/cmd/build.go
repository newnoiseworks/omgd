package cmd

import (
	"fmt"

	"github.com/newnoiseworks/tpl-fred/builder"
	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds various parts of the stack",
	Long: `tpl-fred build command

The main command. Builds all components of the stack.

Usage: $ tpl-fred build [project] [environment] [target]
`,
	Run: run,
}

func run(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		fmt.Println("You need both arguments dumb dumb")
		return
	}

	var project = args[0]
	var environment = args[1]

	if project != "game" && project != "server" && project != "website" && project != "config-files" {
		fmt.Println("Invalid project name")
		return
	}

	if environment != "production" && environment != "staging" && environment != "local" {
		fmt.Println("Invalid environment name")
		return
	}

	fmt.Println(fmt.Sprintf("build called with args %s %s", project, environment))
	builder.Builder(project, environment)
}

func init() {
	rootCmd.AddCommand(buildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
