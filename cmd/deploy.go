package cmd

import (
	"github.com/newnoiseworks/tpl-fred/deployer"
	"github.com/newnoiseworks/tpl-fred/utils"
	"github.com/spf13/cobra"
)

// VolumeReset whether or not to reset docker volumes on deploy
var VolumeReset bool

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Direct deploy command for the projects and infrastructure",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var project = args[0]

		switch project {
		case "server":
			deployer.Server{
				Environment: Profile,
				OutputDir:   OutputDir,
				CmdOnDir:    utils.CmdOnDir,
				VolumeReset: VolumeReset,
			}.Deploy()
			break
		case "game":
			deployer.Game{
				Environment: Profile,
				OutputDir:   OutputDir,
				CmdOnDir:    utils.CmdOnDir,
			}.Deploy()
			break
		case "infra":
			deployer.Infra{
				Environment: Profile,
				OutputDir:   OutputDir,
				CmdOnDir:    utils.CmdOnDir,
			}.Deploy()
			break
		}
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")
	deployCmd.PersistentFlags().BoolVar(&VolumeReset, "volume-reset", false, "Resets docker volumes on deploy if present")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
