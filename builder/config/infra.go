package config

import (
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/newnoiseworks/tpl-fred/utils"

	"github.com/logrusorgru/aurora"
)

// InfraConfig - builds config files needed for infrastructure
func InfraConfig(environment string, buildPath string) {
	var profile = utils.GetProfile(environment)

	config := map[string]string{
		"gcpProject": profile.Gcloud.Project,
		"gcpRegion":  profile.Gcloud.Region,
		"gcpZone":    profile.Gcloud.Zone,
	}

	fmt.Println(aurora.Green("building server config files"))

	buildTerraformVarsFile(buildPath, config)
}

func buildTerraformVarsFile(buildPath string, config map[string]string) {
	fmt.Println(aurora.Yellow(" >> building terraform.tfvars.tmpl >> server/infra/gcp/terraform.tfvars"))

	t, err := template.ParseFiles("builder/config/templates/terraform.tfvars.tmpl")
	if err != nil {
		log.Print(err)
		return
	}

	path := fmt.Sprintf("%s/server/infra/gcp/terraform.tfvars", buildPath)

	f, err := os.Create(path)
	if err != nil {
		log.Println("create file: ", err)
		return
	}

	err = t.Execute(f, config)
	if err != nil {
		log.Print("execute: ", err)
		return
	}
}
