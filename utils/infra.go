package utils

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type InfraChange struct {
	OutputDir       string
	Profile         *ProfileConf
	CmdOnDir        func(string, string, string) string
	CmdOnDirWithEnv func(string, string, string, []string) string
}

func (infraChange *InfraChange) DeployClientAndServer() {
	BuildTemplatesFromPath(infraChange.Profile, infraChange.OutputDir, "tmpl", false)

	infraChange.CmdOnDir(
		fmt.Sprintf("terraform init -reconfigure -force-copy -backend-config bucket=%s -backend-config prefix=terraform/state/%s/%s", infraChange.Profile.Get("omgd.tfsettings.bucket"), infraChange.Profile.Get("omgd.name"), infraChange.Profile.Name),
		fmt.Sprintf("setting up terraform on profile %s", infraChange.Profile.Name),
		fmt.Sprintf("%s/server/infra/gcp/instance-setup/", infraChange.OutputDir),
	)

	ipAddress := infraChange.CmdOnDir(
		"terraform output -raw server_ip",
		"getting ip of newly created server...",
		fmt.Sprintf("%s/server/infra/gcp/instance-setup/", infraChange.OutputDir),
	)

	infraChange.Profile.UpdateProfile("omgd.nakama.host", ipAddress)

	infraChange.CmdOnDir(
		fmt.Sprintf("omgd build-clients --profile=%s", infraChange.Profile.path),
		"building game clients against profile",
		infraChange.OutputDir,
	)

	infraChange.CmdOnDir(
		fmt.Sprintf("cp -rf game/dist/web-%s/. server/nakama/website", infraChange.Profile.Name),
		"copy web build into server",
		infraChange.OutputDir,
	)

	infraChange.CmdOnDirWithEnv(
		"./deploy.sh",
		"deploying game server to gcp",
		fmt.Sprintf("%s/server/deploy/gcp", infraChange.OutputDir),
		[]string{
			fmt.Sprintf("GCP_PROJECT=%s", infraChange.Profile.Get("omgd.gcp.project")),
			fmt.Sprintf("GCP_ZONE=%s", infraChange.Profile.Get("omgd.gcp.zone")),
			fmt.Sprintf("OMGD_PROFILE=%s", infraChange.Profile.Name),
			fmt.Sprintf("OMGD_PROJECT=%s", infraChange.Profile.Get("omgd.name")),
		},
	)
}

func (infraChange *InfraChange) DeployInfra() {
	BuildTemplatesFromPath(infraChange.Profile, infraChange.OutputDir, "tmpl", false)

	infraChange.CmdOnDir(
		fmt.Sprintf("terraform init -reconfigure -force-copy -backend-config bucket=%s -backend-config prefix=terraform/state/%s/%s", infraChange.Profile.Get("omgd.tfsettings.bucket"), infraChange.Profile.Get("omgd.name"), infraChange.Profile.Name),
		"setting up terraform locally",
		fmt.Sprintf("%s/server/infra/gcp/instance-setup/", infraChange.OutputDir),
	)

	infraChange.CmdOnDir(
		"terraform apply -auto-approve",
		"updating cloud infra if needed",
		fmt.Sprintf("%s/server/infra/gcp/instance-setup/", infraChange.OutputDir),
	)

	ipAddress := infraChange.CmdOnDir(
		"terraform output -raw server_ip",
		"getting ip of newly created server...",
		fmt.Sprintf("%s/server/infra/gcp/instance-setup/", infraChange.OutputDir),
	)

	infraChange.Profile.UpdateProfile("omgd.nakama.host", ipAddress)
}

func (infraChange *InfraChange) DestroyInfra() {
	BuildTemplatesFromPath(infraChange.Profile, infraChange.OutputDir, "tmpl", false)

	infraChange.CmdOnDir(
		fmt.Sprintf("terraform init -reconfigure -force-copy -backend-config bucket=%s -backend-config prefix=terraform/state/%s/%s", infraChange.Profile.Get("omgd.tfsettings.bucket"), infraChange.Profile.Get("omgd.name"), infraChange.Profile.Name),
		fmt.Sprintf("setting up terraform on profile %s", infraChange.Profile.Name),
		fmt.Sprintf("%s/server/infra/gcp/instance-setup/", infraChange.OutputDir),
	)

	infraChange.CmdOnDir(
		"terraform destroy -auto-approve",
		"destroying infrastructure",
		fmt.Sprintf("%s/server/infra/gcp/instance-setup/", infraChange.OutputDir),
	)
}

func (infraChange *InfraChange) ProjectSetup() {
	BuildTemplatesFromPath(infraChange.Profile, infraChange.OutputDir, "tmpl", false)

	infraChange.CmdOnDir(
		"terraform init -reconfigure -force-copy -backend-config path=../../../../.omgd/terraform.tfstate",
		"setting up terraform locally",
		fmt.Sprintf("%s/server/infra/gcp/project-setup/", infraChange.OutputDir),
	)

	infraChange.CmdOnDir(
		"terraform apply -auto-approve",
		"setting up initial infra",
		fmt.Sprintf("%s/server/infra/gcp/project-setup/", infraChange.OutputDir),
	)

	bucketName := infraChange.CmdOnDir(
		"terraform output -raw bucket_name",
		"getting newly created bucket name",
		fmt.Sprintf("%s/server/infra/gcp/project-setup/", infraChange.OutputDir),
	)

	omgdProfile := GetProfileFromDir(strings.Replace(
		infraChange.Profile.path,
		fmt.Sprintf("%s.yml", infraChange.Profile.Name),
		"omgd.yml",
		1,
	), infraChange.Profile.rootDir)
	omgdProfile.UpdateProfile("omgd.tfsettings.bucket", bucketName)

	tfFilePath := fmt.Sprintf("%s/server/infra/gcp/project-setup/main.tf", infraChange.OutputDir)
	input, err := ioutil.ReadFile(tfFilePath)
	if err != nil {
		LogFatal(fmt.Sprint(err))
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "backend \"local\"") {
			lines[i] = strings.Replace(lines[i], "backend \"local\"", "backend \"gcs\"", 1)
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(tfFilePath, []byte(output), 0644)
	if err != nil {
		LogFatal(fmt.Sprint(err))
	}

	infraChange.CmdOnDir(
		fmt.Sprintf("terraform init -force-copy -backend-config bucket=%s -backend-config prefix=terraform/state/%s", infraChange.Profile.Get("omgd.tfsettings.bucket"), infraChange.Profile.Get("omgd.name")),
		"setting up terraform to use gcs backend",
		fmt.Sprintf("%s/server/infra/gcp/project-setup/", infraChange.OutputDir),
	)
}

func (infraChange *InfraChange) ProjectDestroy() {
	BuildTemplatesFromPath(infraChange.Profile, infraChange.OutputDir, "tmpl", false)

	infraChange.CmdOnDir(
		fmt.Sprintf("terraform init -reconfigure -force-copy -backend-config bucket=%s -backend-config prefix=terraform/state/%s", infraChange.Profile.Get("omgd.tfsettings.bucket"), infraChange.Profile.Get("omgd.name")),
		"setting up terraform to local backend",
		fmt.Sprintf("%s/server/infra/gcp/project-setup/", infraChange.OutputDir),
	)

	tfFilePath := fmt.Sprintf("%s/server/infra/gcp/project-setup/main.tf", infraChange.OutputDir)
	input, err := ioutil.ReadFile(tfFilePath)
	if err != nil {
		LogFatal(fmt.Sprint(err))
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "backend \"gcs\"") {
			lines[i] = strings.Replace(lines[i], "backend \"gcs\"", "backend \"local\"", 1)
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(tfFilePath, []byte(output), 0644)
	if err != nil {
		LogFatal(fmt.Sprint(err))
	}

	infraChange.CmdOnDir(
		"terraform init -force-copy -backend-config path=../../../../.omgd/terraform.tfstate",
		"setting up terraform to destroy project level infra",
		fmt.Sprintf("%s/server/infra/gcp/project-setup/", infraChange.OutputDir),
	)

	infraChange.CmdOnDir(
		"terraform destroy -auto-approve",
		"destroying initial infra",
		fmt.Sprintf("%s/server/infra/gcp/project-setup/", infraChange.OutputDir),
	)
}
