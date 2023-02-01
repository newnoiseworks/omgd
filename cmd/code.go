/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/newnoiseworks/tpl-fred/utils"
	"github.com/spf13/cobra"
)

// codeCmd represents the code command
var codeCmd = &cobra.Command{
	Use:   "code [plan]",
	Short: "Generates code samples and necessary files for development.",
	Long: `Generates code samples and necessary files for development.

See hopeful non existent at the moment documentation sometime in the future. Similiar to rails generate commands if you're familiar.
`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: needs to take in a plan as a string with possible arguments hopefully in a splat?
		fmt.Println("code called")

		err := utils.CopyStaticDirectory("static/test/test_dir_to_copy", "utils/static/test/test_dir_post_copying")

		if err != nil {
			log.Fatal(err)
		}

		// Utils can have a repo / git repository / codegen and this can be a wrapper
		// seems to match previous setup more or less
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
