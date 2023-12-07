package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type InfraChange struct {
	OutputDir       string
	Profile         *ProfileConf
	CmdOnDir        func(string, string, string) string
	CmdOnDirWithEnv func(string, string, string, []string) string
	SkipCleanup     bool
}

func (infraChange *InfraChange) DeployClientAndServer() {
	infraChange.setupInstanceInfraFiles()
	infraChange.setupDeployFiles()

	BuildTemplatesFromPath(infraChange.Profile, infraChange.OutputDir, "tmpl", false)

	infraChange.CmdOnDir(
		fmt.Sprintf("terraform init -reconfigure -backend-config bucket=%s -backend-config prefix=terraform/state/%s/%s", infraChange.Profile.OMGD.GCP.Bucket, infraChange.Profile.OMGD.Name, infraChange.Profile.Name),
		fmt.Sprintf("setting up terraform on profile %s", infraChange.Profile.Name),
		fmt.Sprintf("%s/.omgd/infra/gcp/instance-setup/", infraChange.OutputDir),
	)

	ipAddress := infraChange.CmdOnDir(
		"terraform output -raw server_ip",
		"getting ip of newly created server...",
		fmt.Sprintf("%s/.omgd/infra/gcp/instance-setup/", infraChange.OutputDir),
	)

	infraChange.Profile.UpdateProfile("omgd.servers.host", ipAddress)

	infraChange.CmdOnDir(
		fmt.Sprintf("omgd build-clients --profile=%s", infraChange.Profile.path),
		"building game clients against profile",
		infraChange.OutputDir,
	)

	infraChange.CmdOnDirWithEnv(
		"./deploy.sh",
		"deploying game server to gcp",
		fmt.Sprintf("%s/.omgd/deploy/gcp", infraChange.OutputDir),
		[]string{
			fmt.Sprintf("GCP_PROJECT=%s", infraChange.Profile.OMGD.GCP.Project),
			fmt.Sprintf("GCP_ZONE=%s", infraChange.Profile.OMGD.GCP.Zone),
			fmt.Sprintf("OMGD_PROFILE=%s", infraChange.Profile.Name),
			fmt.Sprintf("OMGD_PROJECT=%s", infraChange.Profile.OMGD.Name),
			fmt.Sprintf("OMGD_SERVER_SERVICES=%s", strings.Join(infraChange.Profile.OMGD.Servers.Services, " ")),
		},
	)

	if !infraChange.SkipCleanup {
		infraChange.PerformCleanup()
	}
}

func (infraChange *InfraChange) DeployInfra() {
	infraChange.setupInstanceInfraFiles()

	BuildTemplatesFromPath(infraChange.Profile, infraChange.OutputDir, "tmpl", false)

	infraChange.CmdOnDir(
		fmt.Sprintf("terraform init -reconfigure -backend-config bucket=%s -backend-config prefix=terraform/state/%s/%s", infraChange.Profile.OMGD.GCP.Bucket, infraChange.Profile.OMGD.Name, infraChange.Profile.Name),
		"setting up terraform locally",
		fmt.Sprintf("%s/.omgd/infra/gcp/instance-setup/", infraChange.OutputDir),
	)

	infraChange.CmdOnDir(
		"terraform apply -auto-approve",
		"updating cloud infra if needed",
		fmt.Sprintf("%s/.omgd/infra/gcp/instance-setup/", infraChange.OutputDir),
	)

	ipAddress := infraChange.CmdOnDir(
		"terraform output -raw server_ip",
		"getting ip of newly created server...",
		fmt.Sprintf("%s/.omgd/infra/gcp/instance-setup/", infraChange.OutputDir),
	)

	infraChange.Profile.UpdateProfile("omgd.servers.host", ipAddress)

	if !infraChange.SkipCleanup {
		infraChange.PerformCleanup()
	}
}

func (infraChange *InfraChange) DestroyInfra() {
	infraChange.setupInstanceInfraFiles()

	BuildTemplatesFromPath(infraChange.Profile, infraChange.OutputDir, "tmpl", false)

	infraChange.CmdOnDir(
		fmt.Sprintf("terraform init -reconfigure -backend-config bucket=%s -backend-config prefix=terraform/state/%s/%s", infraChange.Profile.OMGD.GCP.Bucket, infraChange.Profile.OMGD.Name, infraChange.Profile.Name),
		fmt.Sprintf("setting up terraform on profile %s", infraChange.Profile.Name),
		fmt.Sprintf("%s/.omgd/infra/gcp/instance-setup/", infraChange.OutputDir),
	)

	infraChange.CmdOnDir(
		"terraform destroy -auto-approve",
		"destroying infrastructure",
		fmt.Sprintf("%s/.omgd/infra/gcp/instance-setup/", infraChange.OutputDir),
	)

	if !infraChange.SkipCleanup {
		infraChange.PerformCleanup()
	}
}

func (infraChange *InfraChange) ProjectSetup() {
	infraChange.setupProjectInfraFiles()

	BuildTemplatesFromPath(infraChange.Profile, infraChange.OutputDir, "tmpl", false)

	infraChange.CmdOnDir(
		"terraform init -reconfigure -backend-config path=../../../../.omgd/terraform.tfstate",
		"setting up terraform locally",
		fmt.Sprintf("%s/.omgd/infra/gcp/project-setup/", infraChange.OutputDir),
	)

	infraChange.CmdOnDir(
		"terraform apply -auto-approve",
		"setting up initial infra",
		fmt.Sprintf("%s/.omgd/infra/gcp/project-setup/", infraChange.OutputDir),
	)

	bucketName := infraChange.CmdOnDir(
		"terraform output -raw bucket_name",
		"getting newly created bucket name",
		fmt.Sprintf("%s/.omgd/infra/gcp/project-setup/", infraChange.OutputDir),
	)

	omgdProfile := GetProfileFromDir(strings.Replace(
		infraChange.Profile.path,
		fmt.Sprintf("%s.yml", infraChange.Profile.Name),
		"omgd.yml",
		1,
	), infraChange.Profile.rootDir)
	omgdProfile.UpdateProfile("omgd.gcp.bucket", bucketName)

	infraChange.Profile = infraChange.Profile.LoadProfile()

	tfFilePath := fmt.Sprintf("%s/.omgd/infra/gcp/project-setup/main.tf", infraChange.OutputDir)
	infraChange.alterInfraBackendFile(tfFilePath, true)

	infraChange.CmdOnDir(
		fmt.Sprintf("terraform init -force-copy -backend-config bucket=%s -backend-config prefix=terraform/state/%s", infraChange.Profile.OMGD.GCP.Bucket, infraChange.Profile.OMGD.Name),
		"setting up terraform to use gcs backend",
		fmt.Sprintf("%s/.omgd/infra/gcp/project-setup/", infraChange.OutputDir),
	)

	if !infraChange.SkipCleanup {
		infraChange.PerformCleanup()
	}
}

func (infraChange *InfraChange) ProjectDestroy() {
	infraChange.setupProjectInfraFiles()

	BuildTemplatesFromPath(infraChange.Profile, infraChange.OutputDir, "tmpl", false)

	tfFilePath := fmt.Sprintf("%s/.omgd/infra/gcp/project-setup/main.tf", infraChange.OutputDir)
	infraChange.alterInfraBackendFile(tfFilePath, true)

	infraChange.CmdOnDir(
		fmt.Sprintf("terraform init -reconfigure -backend-config bucket=%s -backend-config prefix=terraform/state/%s", infraChange.Profile.OMGD.GCP.Bucket, infraChange.Profile.OMGD.Name),
		"setting up terraform to local backend",
		fmt.Sprintf("%s/.omgd/infra/gcp/project-setup/", infraChange.OutputDir),
	)

	infraChange.alterInfraBackendFile(tfFilePath, false)

	infraChange.CmdOnDir(
		"terraform init -force-copy -backend-config path=../../../../.omgd/terraform.tfstate",
		"setting up terraform to destroy project level infra",
		fmt.Sprintf("%s/.omgd/infra/gcp/project-setup/", infraChange.OutputDir),
	)

	infraChange.CmdOnDir(
		"terraform destroy -auto-approve",
		"destroying initial infra",
		fmt.Sprintf("%s/.omgd/infra/gcp/project-setup/", infraChange.OutputDir),
	)

	if !infraChange.SkipCleanup {
		infraChange.PerformCleanup()
	}
}

func (infraChange *InfraChange) PerformCleanup() {
	err := os.RemoveAll(
		fmt.Sprintf(
			"%s/.omgd",
			infraChange.OutputDir,
		),
	)

	if err != nil {
		LogFatal(fmt.Sprint(err))
	}
}

func (infraChange *InfraChange) setupDeployFiles() {
	sccp := StaticCodeCopyPlan{}

	sccp.CopyStaticDirectory(
		"static/deploy/gcp",
		fmt.Sprintf("%s/.omgd/deploy/gcp", infraChange.OutputDir),
	)
}

func (infraChange *InfraChange) setupInstanceInfraFiles() {
	sccp := StaticCodeCopyPlan{}

	sccp.CopyStaticDirectory(
		"static/infra/gcp/instance-setup",
		fmt.Sprintf("%s/.omgd/infra/gcp/instance-setup", infraChange.OutputDir),
	)
}

func (infraChange *InfraChange) setupProjectInfraFiles() {
	sccp := StaticCodeCopyPlan{}

	sccp.CopyStaticDirectory(
		"static/infra/gcp/project-setup",
		fmt.Sprintf("%s/.omgd/infra/gcp/project-setup", infraChange.OutputDir),
	)
}

func (infraChange *InfraChange) alterInfraBackendFile(tfFilePath string, toProvider bool) {
	input, err := ioutil.ReadFile(tfFilePath)
	if err != nil {
		LogFatal(fmt.Sprint(err))
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if toProvider {
			if strings.Contains(line, "backend \"local\"") {
				lines[i] = strings.Replace(lines[i], "backend \"local\"", "backend \"gcs\"", 1)
			}
		} else {
			if strings.Contains(line, "backend \"gcs\"") {
				lines[i] = strings.Replace(lines[i], "backend \"gcs\"", "backend \"local\"", 1)
			}
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(tfFilePath, []byte(output), 0644)
	if err != nil {
		LogFatal(fmt.Sprint(err))
	}
}
