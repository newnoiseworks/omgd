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
	ipAddress := infraChange.CmdOnDir(
		"terraform output -raw server_ip",
		"getting ip of newly created server...",
		fmt.Sprintf("%s/server/infra/gcp/", infraChange.OutputDir),
	)

	infraChange.Profile.UpdateProfile("omgd.gcp.host", ipAddress)

	infraChange.CmdOnDir(
		fmt.Sprintf("omgd build-templates --profile=%s", infraChange.Profile.path),
		"",
		infraChange.OutputDir,
	)

	infraChange.CmdOnDir(
		fmt.Sprintf("omgd build-clients --profile=%s", infraChange.Profile.path),
		"",
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
		},
	)
}

func (infraChange *InfraChange) DeployInfra() {
	BuildTemplatesFromPath(infraChange.Profile, infraChange.OutputDir, "tmpl", false)

	infraChange.CmdOnDir(
		fmt.Sprintf("terraform init -reconfigure -force-copy --backend-config path=.omgd/%s/terraform.tfstate", infraChange.Profile.Name),
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
		fmt.Sprintf("terraform init -reconfigure -force-copy --backend-config path=.omgd/%s/terraform.tfstate", infraChange.Profile.Name),
		fmt.Sprintf("setting up terraform on profile %s", infraChange.Profile.Name),
		fmt.Sprintf("%s/server/infra/gcp/", infraChange.OutputDir),
	)

	infraChange.CmdOnDir(
		"terraform destroy -auto-approve",
		"destroying infrastructure",
		fmt.Sprintf("%s/server/infra/gcp/", infraChange.OutputDir),
	)
}
