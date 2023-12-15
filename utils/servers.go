package utils

import (
	"fmt"
	"os"
	"strings"
)

type ServersChange struct {
	OutputDir       string
	Profile         *ProfileConf
	CmdOnDir        func(string, string, string) string
	CmdOnDirWithEnv func(string, string, string, []string) string
	SkipCleanup     bool
}

func (serversChange *ServersChange) Deploy() {
	serversChange.setupInstanceInfraFiles()
	serversChange.setupDeployFiles()

	BuildTemplatesFromPath(serversChange.Profile, serversChange.OutputDir, "tmpl", false)

	serversChange.CmdOnDir(
		fmt.Sprintf("terraform init -reconfigure -backend-config bucket=%s -backend-config prefix=terraform/state/%s/%s", serversChange.Profile.OMGD.GCP.Bucket, serversChange.Profile.OMGD.Name, serversChange.Profile.Name),
		fmt.Sprintf("setting up terraform on profile %s", serversChange.Profile.Name),
		fmt.Sprintf("%s/.omgd/infra/gcp/instance-setup/", serversChange.OutputDir),
	)

	ipAddress := serversChange.CmdOnDir(
		"terraform output -raw server_ip",
		"getting ip of newly created server...",
		fmt.Sprintf("%s/.omgd/infra/gcp/instance-setup/", serversChange.OutputDir),
	)

	serversChange.Profile.UpdateProfile("omgd.servers.host", ipAddress)

	serversChange.CmdOnDir(
		fmt.Sprintf("omgd game build --profile=%s", serversChange.Profile.path),
		"building game clients against profile",
		serversChange.OutputDir,
	)

	serviceArray := []string{}

	for _, service := range serversChange.Profile.OMGD.Servers.Services {
		serviceArray = append(serviceArray, service.BuildService)
	}

	services := strings.Join(serviceArray, " ")

	serversChange.CmdOnDirWithEnv(
		"./deploy.sh",
		"deploying game server to gcp",
		fmt.Sprintf("%s/.omgd/deploy/gcp", serversChange.OutputDir),
		[]string{
			fmt.Sprintf("GCP_PROJECT=%s", serversChange.Profile.OMGD.GCP.Project),
			fmt.Sprintf("GCP_ZONE=%s", serversChange.Profile.OMGD.GCP.Zone),
			fmt.Sprintf("OMGD_PROFILE=%s", serversChange.Profile.Name),
			fmt.Sprintf("OMGD_PROJECT=%s", serversChange.Profile.OMGD.Name),
			fmt.Sprintf("OMGD_SERVER_SERVICES=%s", services),
			fmt.Sprintf("CLOUDSDK_AUTH_CREDENTIAL_FILE_OVERRIDE=%s", serversChange.Profile.OMGD.GCP.CredsFile),
		},
	)

	if !serversChange.SkipCleanup {
		serversChange.PerformCleanup()
	}
}

func (serversChange *ServersChange) setupInstanceInfraFiles() {
	sccp := StaticCodeCopyPlan{}

	sccp.CopyStaticDirectory(
		"static/infra/gcp/instance-setup",
		fmt.Sprintf("%s/.omgd/infra/gcp/instance-setup", serversChange.OutputDir),
	)
}

func (serversChange *ServersChange) setupDeployFiles() {
	sccp := StaticCodeCopyPlan{}

	sccp.CopyStaticDirectory(
		"static/deploy/gcp",
		fmt.Sprintf("%s/.omgd/deploy/gcp", serversChange.OutputDir),
	)
}

func (serversChange *ServersChange) PerformCleanup() {
	err := os.RemoveAll(
		fmt.Sprintf(
			"%s/.omgd",
			serversChange.OutputDir,
		),
	)

	if err != nil {
		LogFatal(fmt.Sprint(err))
	}
}
