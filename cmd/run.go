package cmd

import (
	"log"
	"strconv"

	"github.com/newnoiseworks/omgd/utils"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs commands based on the chosen profile",
	Long: `Works in conjunction with the --profile you've set, or profiles/local.yml by default. Runs either the base project command or any tasks you specify if chosen.

Usage: 

Run entire project:
$ omgd run

Run part of project:
$ omgd run [name-of-project-step] [number-of-step (optional)]

Run task:
$ omgd run task [name-of-task] [number-of-step (optional)]
`,
	Run: func(cmd *cobra.Command, args []string) {
		profile := utils.GetProfile(ProfilePath)

		runner := utils.Run{
			Profile:   profile,
			OutputDir: OutputDir,
			CmdDir:    utils.CmdOnDir,
			Verbosity: Verbosity,
		}

		switch len(args) {
		case 0:
			runner.Run()
		case 1:
			runner.RunProjectStep(args[0])
		case 2:
			if args[0] == "task" {
				runner.RunTask(args[1])
			} else {
				idx, err := strconv.Atoi(args[1])

				if err != nil {
					log.Fatal(err)
				}

				runner.RunProjectSubStep(args[0], idx)
			}
		case 3:
			idx, err := strconv.Atoi(args[2])

			if err != nil {
				log.Fatal(err)
			}

			runner.RunTaskSubStep(args[1], idx)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
