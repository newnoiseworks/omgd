/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/newnoiseworks/omgd/utils"
	"github.com/spf13/cobra"
)

// codeCmd represents the code command
var codeCmd = &cobra.Command{
	Use:   "code [plan] [plan args]*",
	Short: "Generates code samples and necessary files for development.",
	Long: `Generates code samples and necessary files for development. Similiar to rails generate commands if you're familiar.

Current plans available:
new [project_name (will name folder)] - Starts a new OMGD project using Godot as the game engine
channel [channel_name (must be snake and lowercase!)] [events (must be snake and lowercase, separated by blank spaces)]* - Creates a new OMGD multiplayer channel to communicate with

Example code plans:
example-2d-player-movement - Complete example, demonstrates 2d player movement.
`,
	Run: func(cmd *cobra.Command, args []string) {
		plan := args[0]
		target := args[0]
		codePlanArgs := ""

		if len(args) == 1 {
			if target == plan && strings.Contains(target, "/") {
				splits := strings.Split(target, "/")
				target = splits[len(splits)-1]
			}
		}

		if len(args) > 1 {
			target = args[1]
		}

		if len(args) > 2 {
			codePlanArgs = args[2]
		}

		cp := utils.CodeGenerationPlan{
			OutputDir: OutputDir,
			Target:    target,
			Plan:      plan,
			Args:      codePlanArgs,
		}

		cp.Generate()

		utils.LogWarn(fmt.Sprintf("Code generated based on plan to %s dir %s", cp.Plan, cp.OutputDir))
	},
	ValidArgs: []string{
		"godot3-nakama3/new",
		"godot3-nakama3/channel",
		"godot3-nakama3/example-2d-player-movement",
		"bare",
	},
}

func init() {
	rootCmd.AddCommand(codeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// codeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// codeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
