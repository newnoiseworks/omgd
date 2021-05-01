package cmd

import (
	"fmt"

	"github.com/newnoiseworks/tpl-fred/builder"
	"github.com/newnoiseworks/tpl-fred/utils"
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
	if utils.CheckProjectAndEnv(args) == false {
		return
	}

	var project = args[0]
	var environment = args[1]

	fmt.Println(fmt.Sprintf("build called with args %s %s", project, environment))
	switch project {
	case "game":
		builder.Game{
			Environment: environment,
			OutputDir:   OutputDir,
			CmdOnDir:    utils.CmdOnDir,
		}.Build()
		break
	case "server":
		builder.Server{
			Environment: environment,
			OutputDir:   OutputDir,
			CmdOnDir:    utils.CmdOnDir,
		}.Build()
		break
	}
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
