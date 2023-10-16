package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

// ProfilePath this is the yml profile you're using
var ProfilePath string

// OutputDir this is where all builds and build artifacts will be written to
var OutputDir string

// Verbosity how verbose should output be
var Verbosity bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "omgd",
	Short: "Open Multiplayer Game Development",
	Long: `Open Multiplayer Game Development aims to make multiplayer game development easier, as well as multiverse development, if you're into that sort of thing.

omgd uses open source game and cloud frameworks to help you start development quickly and scale as needed before production, with a focus on development across teams.

Godot is the current game engine of focus with a likelihood of expansion to Unity next, then Unreal.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&ProfilePath, "profile", "profiles/local.yml", "yml profile representing this build in the build/profiles folder")

	rootCmd.PersistentFlags().StringVar(&OutputDir, "output-dir", ".", "output dir of files that are generated etc")

	rootCmd.PersistentFlags().BoolVar(&Verbosity, "verbose", false, "include for verbose output of commands")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
}
