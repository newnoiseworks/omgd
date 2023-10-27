package utils

import (
	"fmt"
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
		fmt.Sprintf("terraform init -reconfigure -force-copy -backend-config bucket=%s-bucket-tfstate -backend-config prefix=terraform/state/%s", infraChange.Profile.Get("omgd.name"), infraChange.Profile.Name),
		fmt.Sprintf("setting up terraform on profile %s", infraChange.Profile.Name),
		fmt.Sprintf("%s/server/infra/gcp/", infraChange.OutputDir),
	)

	ipAddress := infraChange.CmdOnDir(
		"terraform output -raw server_ip",
		"getting ip of newly created server...",
		fmt.Sprintf("%s/server/infra/gcp/", infraChange.OutputDir),
	)

	infraChange.Profile.UpdateProfile("omgd.gcp.host", ipAddress)

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
			fmt.Sprintf("PROFILE=%s", infraChange.Profile.Name),
		},
	)
}

func (infraChange *InfraChange) DeployInfra() {
	BuildTemplatesFromPath(infraChange.Profile, infraChange.OutputDir, "tmpl", false)

	infraChange.CmdOnDir(
		fmt.Sprintf("terraform init -reconfigure -force-copy -backend-config bucket=%s-bucket-tfstate -backend-config prefix=terraform/state/%s", infraChange.Profile.Get("omgd.name"), infraChange.Profile.Name),
		"setting up terraform locally",
		fmt.Sprintf("%s/server/infra/gcp/", infraChange.OutputDir),
	)

	infraChange.CmdOnDir(
		"terraform apply -auto-approve",
		"updating cloud infra if needed",
		fmt.Sprintf("%s/server/infra/gcp/", infraChange.OutputDir),
	)

	ipAddress := infraChange.CmdOnDir(
		"terraform output -raw server_ip",
		"getting ip of newly created server...",
		fmt.Sprintf("%s/server/infra/gcp/", infraChange.OutputDir),
	)

	infraChange.Profile.UpdateProfile("omgd.gcp.host", ipAddress)
}

func (infraChange *InfraChange) DestroyInfra() {
	BuildTemplatesFromPath(infraChange.Profile, infraChange.OutputDir, "tmpl", false)

	infraChange.CmdOnDir(
		fmt.Sprintf("terraform init -reconfigure -force-copy -backend-config bucket=%s-bucket-tfstate -backend-config prefix=terraform/state/%s", infraChange.Profile.Get("omgd.name"), infraChange.Profile.Name),
		fmt.Sprintf("setting up terraform on profile %s", infraChange.Profile.Name),
		fmt.Sprintf("%s/server/infra/gcp/", infraChange.OutputDir),
	)

	infraChange.CmdOnDir(
		"terraform destroy -auto-approve",
		"destroying infrastructure",
		fmt.Sprintf("%s/server/infra/gcp/", infraChange.OutputDir),
	)
}

func (infraChange *InfraChange) ProjectSetup() {
	BuildTemplatesFromPath(infraChange.Profile, infraChange.OutputDir, "tmpl", false)

	infraChange.CmdOnDir(
		"terraform init -reconfigure -force-copy -backend-config path=../../../../.omgd/terraform.tfstate",
		"setting up terraform locally",
		fmt.Sprintf("%s/server/infra/project-setup/gcp/", infraChange.OutputDir),
	)

	infraChange.CmdOnDir(
		"terraform apply -auto-approve",
		"setting up initial infra",
		fmt.Sprintf("%s/server/infra/project-setup/gcp/", infraChange.OutputDir),
	)
}

func (infraChange *InfraChange) ProjectDestroy() {
	BuildTemplatesFromPath(infraChange.Profile, infraChange.OutputDir, "tmpl", false)

	infraChange.CmdOnDir(
		"terraform init -reconfigure -force-copy -backend-config path=../../../../.omgd/terraform.tfstate",
		"setting up terraform locally",
		fmt.Sprintf("%s/server/infra/project-setup/gcp/", infraChange.OutputDir),
	)

	infraChange.CmdOnDir(
		"terraform destroy -auto-approve",
		"destroying initial infra",
		fmt.Sprintf("%s/server/infra/project-setup/gcp/", infraChange.OutputDir),
	)
}
