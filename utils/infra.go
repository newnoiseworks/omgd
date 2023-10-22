package utils

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type InfraChange struct {
	OutputDir    string
	Profile      *ProfileConf
	CmdOnDir     func(string, string, string, bool) string
	Verbosity    bool
	CopyToTmpDir bool
	tmpDir       string
}

func (infraChange *InfraChange) DeployClientAndServer() {
	infraChange.setup()

	tmpProfilePath := strings.ReplaceAll(infraChange.Profile.path, "profiles/", ".omgd/")

	infraChange.CmdOnDir(
		fmt.Sprintf("omgd run task set-ip-to-profile --profile=%s", tmpProfilePath),
		"",
		infraChange.OutputDir,
		infraChange.Verbosity,
	)

	infraChange.CmdOnDir(
		fmt.Sprintf("omgd build-templates --profile=%s", tmpProfilePath),
		"",
		infraChange.OutputDir,
		infraChange.Verbosity,
	)

	infraChange.CmdOnDir(
		fmt.Sprintf("omgd build-clients --profile=%s", tmpProfilePath),
		"",
		infraChange.OutputDir,
		infraChange.Verbosity,
	)

	infraChange.CmdOnDir(
		fmt.Sprintf("omgd run nakama-server --profile=%s", tmpProfilePath),
		"",
		infraChange.OutputDir,
		infraChange.Verbosity,
	)
}

func (infraChange *InfraChange) DeployInfra() {
	infraChange.setup()

	tmpProfile := GetProfile(strings.ReplaceAll(infraChange.Profile.path, "profiles/", ".omgd/"))

	BuildTemplatesFromPath(tmpProfile, infraChange.OutputDir, "tmpl", false, infraChange.Verbosity)

	infraChange.CmdOnDir(
		"terraform init",
		"setting up terraform locally",
		fmt.Sprintf("%s/server/infra/gcp/", infraChange.OutputDir),
		infraChange.Verbosity,
	)

	infraChange.CmdOnDir(
		"terraform apply -auto-approve",
		"updating cloud infra if needed",
		fmt.Sprintf("%s/server/infra/gcp/", infraChange.OutputDir),
		infraChange.Verbosity,
	)

	ipAddress := infraChange.CmdOnDir(
		"terraform output -raw server_ip",
		"getting ip of newly created server...",
		fmt.Sprintf("%s/server/infra/gcp/", infraChange.OutputDir),
		infraChange.Verbosity,
	)

	tmpProfile.UpdateProfile("nakama.host", ipAddress)
	infraChange.Profile.UpdateProfile("omgd.deploy.server.gcloud.host", ipAddress)
}

func (infraChange *InfraChange) DestroyInfra() {
	infraChange.setup()
	tmpProfilePath := strings.ReplaceAll(infraChange.Profile.path, "profiles/", ".omgd/")

	// NOTE: Would like to discourage this in favor of using utils.Run but testing is easier this way
	infraChange.CmdOnDir(
		fmt.Sprintf("omgd run task destroy-infra --profile=%s", tmpProfilePath),
		"",
		infraChange.OutputDir,
		infraChange.Verbosity,
	)
}

func (infraChange *InfraChange) setup() {
	// 1. Should create or empty .omgdtmp directory to work in
	if infraChange.CopyToTmpDir {
		infraChange.tmpDir = fmt.Sprintf("%s/.omgdtmp", infraChange.OutputDir)

		if infraChange.OutputDir == "." {
			infraChange.tmpDir = ".omgdtmp"
		}

		_, err := os.Stat(infraChange.tmpDir)
		if !os.IsNotExist(err) {
			err = os.RemoveAll(infraChange.tmpDir)

			if err != nil {
				log.Fatal(err)
			}
		}

		err = os.Mkdir(infraChange.tmpDir, 0755)

		if err != nil {
			log.Fatal(err)
		}

		// 2. Should clone repo at base of dir (? how to test w/o submodules? clone entire base repo maybe?)
		sccp := StaticCodeCopyPlan{
			skipPaths: []string{
				infraChange.tmpDir,
				".git",
			},
		}

		cwd, err := os.Getwd()

		if err != nil {
			log.Fatal(err)
		}

		if infraChange.OutputDir == "." {
			sccp.CopyStaticDirectory(cwd, infraChange.tmpDir)
		} else {
			sccp.CopyStaticDirectory(infraChange.OutputDir, infraChange.tmpDir)
		}

		infraChange.OutputDir = infraChange.tmpDir
	}

	// 2. Build profiles directory
	BuildProfiles(infraChange.OutputDir, infraChange.Verbosity)
}
