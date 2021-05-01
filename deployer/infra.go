package deployer

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
)

type Infra struct {
	Environment string
	OutputDir   string
	CmdOnDir    func(string, string, string)
}

// DeployInfra Deploys infrastructure using terraform
func (di Infra) Deploy() {
	fmt.Println("deploying server")

	serverPath, err := filepath.Abs(fmt.Sprintf("%s/server/infra/gcp", di.OutputDir))
	if err != nil {
		log.Fatal(err)
		return
	}

	switch di.Environment {
	case "local":
		fmt.Println("Need to make local infra deployment commands")
		break
	case "production":
		di.deployInfraBasedOnProfile(serverPath)
		break
	default:
		di.deployInfraBasedOnProfile(serverPath)
		break
	}
}

func (di Infra) deployInfraBasedOnProfile(serverPath string) {
	config.InfraConfig(di.Environment, di.OutputDir)

	di.CmdOnDir("terraform init", "Initing terraform...", serverPath)
	di.CmdOnDir("./tf_import.sh", "Importing existing resources into terraform...", serverPath)

	exitCode := terraformPlan(serverPath)
	if exitCode == 2 {
		di.CmdOnDir("terraform apply -auto-approve", "Applying changes to infra", serverPath)
	}

	di.getAndSetHostIPFromTerraform(serverPath)
}

type serverIPData struct {
	Value string `json:"value"`
}

type infraResponse struct {
	ServerIP serverIPData `json:"server_ip"`
}

func (di Infra) getAndSetHostIPFromTerraform(path string) {
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

	conf := utils.GetProfile(di.Environment)
	conf.Nakama.Host = response.ServerIP.Value
	utils.SaveProfile(conf, di.Environment)

	fmt.Println(aurora.Green("Success!"))
	fmt.Println(fmt.Printf("Check your profile, Nakama.host should be set to %s --", response.ServerIP.Value))
	fmt.Println("---")
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