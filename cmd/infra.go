/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/newnoiseworks/tpl-fred/builder/config"
	"github.com/newnoiseworks/tpl-fred/utils"
	"github.com/spf13/cobra"
)

// infraCmd represents the infra command
var infraCmd = &cobra.Command{
	Use:   "infra",
	Short: "Manages cloud based infrastructure for your project",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		environment := args[0]

		config.ServerConfig(environment, OutputDir)

		path, err := filepath.Abs(fmt.Sprintf("%s/server", OutputDir))
		if err != nil {
			log.Fatal(err)
			return
		}

		utils.CmdOnDir("terraform init", "Initing terraform...", path)
		utils.CmdOnDir("./gcp_tf_import.sh", "Importing existing resources into terraform...", path)

		exitCode := terraformPlan(path)
		if exitCode == 2 {
			utils.CmdOnDir("terraform apply -auto-approve", "Applying changes to infra", path)
		}

		getAndSetHostIPFromTerraform(path, environment)
	},
}

type serverIPData struct {
	Value string `json:"value"`
}

type infraResponse struct {
	ServerIP serverIPData `json:"server_ip"`
}

func getAndSetHostIPFromTerraform(path string, environment string) {
	cmd := exec.Command("bash", "-c", "terraform output -json")
	cmd.Dir = path

	fmt.Print(aurora.Cyan("Getting IP from terraform... "))

	out, err := cmd.Output()

	if err != nil {
		fmt.Print(aurora.Red("Error!\n"))
		fmt.Printf("%s", out)
		fmt.Println(err)
		log.Fatal("Error getting IP from terraform")
	}

	var response infraResponse
	json.Unmarshal(out, &response)

	conf := utils.GetProfile(environment)
	conf.Nakama.Host = response.ServerIP.Value
	utils.SaveProfile(conf, environment)

	fmt.Println(aurora.Green("Success!"))
	fmt.Println(fmt.Printf("Check your profile, Nakama.host should be set to %s", response.ServerIP.Value))
}

func terraformPlan(path string) int {
	cmd := exec.Command("bash", "-c", "terraform plan -detailed-exitcode")
	cmd.Dir = path

	fmt.Print(aurora.Cyan("Running terraform plan and apply as needed... "))

	out, err := cmd.Output()

	if err != nil {
		if strings.Contains(err.Error(), "exit status 2") {
			return 2
		} else if strings.Contains(err.Error(), "exit status 0") {
			return 0
		}

		fmt.Println(aurora.Red("error executing terraform plan -detailed-exitcode"))
		fmt.Println(string(out[:]))
		log.Fatal(err)
		return -1
	} else {
		return 0
	}
}

func init() {
	rootCmd.AddCommand(infraCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infraCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infraCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}