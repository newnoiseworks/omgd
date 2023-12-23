package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type ServersChange struct {
	OutputDir       string
	Profile         *ProfileConf
	CmdOnDir        func(string, string, string) string
	CmdOnDirWithEnv func(string, string, string, []string) string
	SkipCleanup     bool
}

func RemoteGCPCommand(cmd string, dir string, profile *ProfileConf) {
	CmdOnDirToStdOut(
		fmt.Sprintf("gcloud compute ssh omgd-sa@%s-omgd-dev-instance-%s --project=%s --zone=%s -- %s",
			profile.OMGD.Name,
			profile.Name,
			profile.OMGD.GCP.Project,
			profile.OMGD.GCP.Zone,
			cmd,
		),
		"printing server logs from GCP compute instance",
		dir,
		[]string{
			fmt.Sprintf("CLOUDSDK_AUTH_CREDENTIAL_FILE_OVERRIDE=%s", profile.OMGD.GCP.CredsFile),
		},
	)
}

func (serversChange *ServersChange) Deploy() {
	serversChange.setupInstanceInfraFiles()
	serversChange.setupDeployFiles()

	BuildTemplatesFromPath(serversChange.Profile, serversChange.OutputDir, "tmpl", false)

	serversChange.CmdOnDir(
		fmt.Sprintf("terraform init -reconfigure -backend-config bucket=%s -backend-config prefix=terraform/state/%s/%s", serversChange.Profile.OMGD.GCP.Bucket, serversChange.Profile.OMGD.Name, serversChange.Profile.Name),
		fmt.Sprintf("setting up terraform on profile %s", serversChange.Profile.Name),
		filepath.Join(serversChange.OutputDir, ".omgd", "infra", "gcp", "instance-setup"),
	)

	ipAddress := serversChange.CmdOnDir(
		"terraform output -raw server_ip",
		"getting ip of newly created server...",
		filepath.Join(serversChange.OutputDir, ".omgd", "infra", "gcp", "instance-setup"),
	)

	serversChange.Profile.UpdateProfile("omgd.servers.host", ipAddress)

	serviceArray := []string{}

	for _, service := range serversChange.Profile.OMGD.Servers.Services {
		serviceArray = append(serviceArray, service.BuildService)
	}

	services := strings.Join(serviceArray, " ")

	deployCmd := "./deploy.sh"

	if runtime.GOOS == "windows" {
		deployCmd = "cmd.exe /C deploy.bat"
	}

	serversChange.CmdOnDirWithEnv(
		deployCmd,
		"deploying game server to gcp",
		filepath.Join(serversChange.OutputDir, ".omgd", "deploy", "gcp"),
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
		filepath.Join("static", "infra", "gcp", "instance-setup"),
		filepath.Join(serversChange.OutputDir, ".omgd", "infra", "gcp", "instance-setup"),
	)
}

func (serversChange *ServersChange) setupDeployFiles() {
	sccp := StaticCodeCopyPlan{}

	sccp.CopyStaticDirectory(
		filepath.Join("static", "deploy", "gcp"),
		filepath.Join(serversChange.OutputDir, ".omgd", "deploy", "gcp"),
	)
}

func (serversChange *ServersChange) PerformCleanup() {
	err := os.RemoveAll(filepath.Join(serversChange.OutputDir, ".omgd"))

	if err != nil {
		LogFatal(fmt.Sprint(err))
	}
}
