package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

// Profile this is the yml profile you're using
var Profile string

// OutputDir this is where all builds and build artifacts will be written to
var OutputDir string

// Verbosity how verbose should output be
var Verbosity bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tpl-fred",
	Short: "Build tool for The Promised Land",
	Long: `This tool builds the game The Promised Land, it's server, the
website, and helps deploy each to various targets. Also, it
should setup the project for you if you're just starting. Maybe.
	
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&Profile, "profile", "profiles/local", "yml profile representing this build in the build/profiles folder")

	rootCmd.PersistentFlags().BoolVar(&Verbosity, "verbose", false, "include for verbose output of commands")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	OutputDir = "."
}
