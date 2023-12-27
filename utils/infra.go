package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type InfraChange struct {
	OutputDir       string
	Profile         *ProfileConf
	CmdOnDir        func(string, string, string) string
	CmdOnDirWithEnv func(string, string, string, []string) string
	SkipCleanup     bool
}

func (infraChange *InfraChange) InstanceSetup() {
	infraChange.setupInstanceInfraFiles()

	BuildTemplatesFromPath(
		infraChange.Profile,
		filepath.Join(infraChange.OutputDir, ".omgd"),
		"tmpl",
		false,
	)

	infraChange.CmdOnDir(
		fmt.Sprintf("terraform init -reconfigure -backend-config bucket=%s -backend-config prefix=terraform/state/%s/%s", infraChange.Profile.OMGD.GCP.Bucket, infraChange.Profile.OMGD.Name, infraChange.Profile.Name),
		"setting up terraform locally",
		filepath.Join(infraChange.OutputDir, ".omgd", "infra", "gcp", "instance-setup"),
	)

	infraChange.CmdOnDir(
		"terraform apply -auto-approve",
		"updating cloud infra if needed",
		filepath.Join(infraChange.OutputDir, ".omgd", "infra", "gcp", "instance-setup"),
	)

	ipAddress := infraChange.CmdOnDir(
		"terraform output -raw server_ip",
		"getting ip of newly created server...",
		filepath.Join(infraChange.OutputDir, ".omgd", "infra", "gcp", "instance-setup"),
	)

	infraChange.Profile.UpdateProfile("omgd.servers.host", ipAddress)

	if !infraChange.SkipCleanup {
		infraChange.PerformCleanup()
	}
}

func (infraChange *InfraChange) InstanceDestroy() {
	infraChange.setupInstanceInfraFiles()

	BuildTemplatesFromPath(
		infraChange.Profile,
		filepath.Join(infraChange.OutputDir, ".omgd"),
		"tmpl",
		false,
	)

	infraChange.CmdOnDir(
		fmt.Sprintf("terraform init -reconfigure -backend-config bucket=%s -backend-config prefix=terraform/state/%s/%s", infraChange.Profile.OMGD.GCP.Bucket, infraChange.Profile.OMGD.Name, infraChange.Profile.Name),
		fmt.Sprintf("setting up terraform on profile %s", infraChange.Profile.Name),
		filepath.Join(infraChange.OutputDir, ".omgd", "infra", "gcp", "instance-setup"),
	)

	infraChange.CmdOnDir(
		"terraform destroy -auto-approve",
		"destroying infrastructure",
		filepath.Join(infraChange.OutputDir, ".omgd", "infra", "gcp", "instance-setup"),
	)

	if !infraChange.SkipCleanup {
		infraChange.PerformCleanup()
	}

	infraChange.Profile.UpdateProfile("omgd.servers.host", "???")
}

func (infraChange *InfraChange) ProjectSetup() {
	infraChange.setupProjectInfraFiles()

	BuildTemplatesFromPath(
		infraChange.Profile,
		filepath.Join(infraChange.OutputDir, ".omgd"),
		"tmpl",
		false,
	)

	infraChange.CmdOnDir(
		fmt.Sprintf("terraform init -reconfigure -backend-config path=%s", filepath.Join("..", "..", "..", "..", ".omgd", "terraform.tfstate")),
		"setting up terraform locally",
		filepath.Join(infraChange.OutputDir, ".omgd", "infra", "gcp", "project-setup"),
	)

	infraChange.CmdOnDir(
		"terraform apply -auto-approve",
		"setting up initial infra",
		filepath.Join(infraChange.OutputDir, ".omgd", "infra", "gcp", "project-setup"),
	)

	bucketName := infraChange.CmdOnDir(
		"terraform output -raw bucket_name",
		"getting newly created bucket name",
		filepath.Join(infraChange.OutputDir, ".omgd", "infra", "gcp", "project-setup"),
	)

	omgdProfile := GetProfileFromDir(strings.Replace(
		infraChange.Profile.path,
		fmt.Sprintf("%s.yml", infraChange.Profile.Name),
		"omgd.cloud.yml",
		1,
	), infraChange.Profile.rootDir)
	omgdProfile.UpdateProfile("omgd.gcp.bucket", bucketName)

	infraChange.Profile = infraChange.Profile.LoadProfile()

	tfFilePath := filepath.Join(infraChange.OutputDir, ".omgd", "infra", "gcp", "project-setup", "main.tf")
	infraChange.alterInfraBackendFile(tfFilePath, true)

	infraChange.CmdOnDir(
		fmt.Sprintf("terraform init -force-copy -backend-config bucket=%s -backend-config prefix=terraform/state/%s", infraChange.Profile.OMGD.GCP.Bucket, infraChange.Profile.OMGD.Name),
		"setting up terraform to use gcs backend",
		filepath.Join(infraChange.OutputDir, ".omgd", "infra", "gcp", "project-setup"),
	)

	if !infraChange.SkipCleanup {
		infraChange.PerformCleanup()
	}
}

func (infraChange *InfraChange) ProjectDestroy() {
	if infraChange.hasRunningInstances() {
		LogWarn("It appears this project has running compute instances. Run omgd infra instance-destroy against all profiles that have created instances within this project, or delete the entire project on GCP and start again with a new one. OMGD is not responsible for managing your cloud billing.")

		return
	}

	infraChange.setupProjectInfraFiles()

	BuildTemplatesFromPath(
		infraChange.Profile,
		filepath.Join(infraChange.OutputDir, ".omgd"),
		"tmpl",
		false,
	)

	tfFilePath := filepath.Join(infraChange.OutputDir, ".omgd", "infra", "gcp", "project-setup", "main.tf")
	infraChange.alterInfraBackendFile(tfFilePath, true)

	infraChange.CmdOnDir(
		fmt.Sprintf("terraform init -reconfigure -backend-config bucket=%s -backend-config prefix=terraform/state/%s", infraChange.Profile.OMGD.GCP.Bucket, infraChange.Profile.OMGD.Name),
		"setting up terraform to local backend",
		filepath.Join(infraChange.OutputDir, ".omgd", "infra", "gcp", "project-setup"),
	)

	infraChange.alterInfraBackendFile(tfFilePath, false)

	infraChange.CmdOnDir(
		fmt.Sprintf("terraform init -force-copy -backend-config path=%s", filepath.Join("..", "..", "..", "..", ".omgd", "terraform.tfstate")),
		"setting up terraform to destroy project level infra",
		filepath.Join(infraChange.OutputDir, ".omgd", "infra", "gcp", "project-setup"),
	)

	infraChange.CmdOnDir(
		"terraform destroy -auto-approve",
		"destroying initial infra",
		filepath.Join(infraChange.OutputDir, ".omgd", "infra", "gcp", "project-setup"),
	)

	if !infraChange.SkipCleanup {
		infraChange.PerformCleanup()
	}

	omgdProfile := GetProfileFromDir(strings.Replace(
		infraChange.Profile.path,
		fmt.Sprintf("%s.yml", infraChange.Profile.Name),
		"omgd.cloud.yml",
		1,
	), infraChange.Profile.rootDir)

	omgdProfile.UpdateProfile("omgd.gcp.bucket", "???")
}

func (infraChange *InfraChange) PerformCleanup() {
	err := os.RemoveAll(
		filepath.Join(infraChange.OutputDir, ".omgd"),
	)

	if err != nil {
		LogFatal(fmt.Sprintf("Error on infraChange#PerformCleanup %s", err))
	}
}

type gcloudInstanceResponse struct {
	Name string `json:"name"`
}

func (infraChange *InfraChange) hasRunningInstances() bool {
	instanceList := infraChange.CmdOnDirWithEnv(
		fmt.Sprintf(
			"gcloud compute instances list --format=json --project=%s",
			infraChange.Profile.OMGD.GCP.Project,
		),
		"checking for running compute instances",
		infraChange.OutputDir,
		[]string{
			fmt.Sprintf("CLOUDSDK_AUTH_CREDENTIAL_FILE_OVERRIDE=%s", infraChange.Profile.OMGD.GCP.CredsFile),
		},
	)

	instances := []gcloudInstanceResponse{}

	err := json.Unmarshal([]byte(instanceList), &instances)

	if err != nil {
		LogFatal(fmt.Sprintf("Error unmarshalling json from gcloud instance list check %s", err))
	}

	return len(instances) > 0
}

func (infraChange *InfraChange) setupDeployFiles() {
	sccp := StaticCodeCopyPlan{}

	sccp.CopyStaticDirectory(
		filepath.Join("static", "deploy", "gcp"),
		filepath.Join(infraChange.OutputDir, ".omgd", "deploy", "gcp"),
	)
}

func (infraChange *InfraChange) setupInstanceInfraFiles() {
	sccp := StaticCodeCopyPlan{}

	sccp.CopyStaticDirectory(
		filepath.Join("static", "infra", "gcp", "instance-setup"),
		filepath.Join(infraChange.OutputDir, ".omgd", "infra", "gcp", "instance-setup"),
	)
}

func (infraChange *InfraChange) setupProjectInfraFiles() {
	sccp := StaticCodeCopyPlan{}

	sccp.CopyStaticDirectory(
		filepath.Join("static", "infra", "gcp", "project-setup"),
		filepath.Join(infraChange.OutputDir, ".omgd", "infra", "gcp", "project-setup"),
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
