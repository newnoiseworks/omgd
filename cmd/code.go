/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/newnoiseworks/tpl-fred/utils"
	"github.com/spf13/cobra"
)

// codeCmd represents the code command
var codeCmd = &cobra.Command{
	Use:   "code [plan] [folder name (optional, defaults to plan name)]",
	Short: "Generates code samples and necessary files for development.",
	Long: `Generates code samples and necessary files for development.

See hopeful non existent at the moment documentation sometime in the future. Similiar to rails generate commands if you're familiar.
`,
	Run: func(cmd *cobra.Command, args []string) {
		plan := args[0]
		target := args[0]

		if len(args) > 1 {
			target = args[1]
		}

		cp := utils.CodeGenerationPlan{
			OutputDir: OutputDir,
			Target:    target,
			Plan:      plan,
			Verbosity: Verbosity,
		}

		cp.Generate()

		fmt.Printf("Code generated")
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
