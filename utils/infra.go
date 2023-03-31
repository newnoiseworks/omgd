package utils

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type InfraChange struct {
	OutputDir    string
	ProfilePath  string
	CmdOnDir     func(string, string, string, bool)
	Verbosity    bool
	tmpDir       string
	CopyToTmpDir bool
}

func (infraChange *InfraChange) DeployClientAndServer() {
	infraChange.setup()

	infraChange.CmdOnDir(
		fmt.Sprintf("omgd run task set-ip-to-profile --profile=%s", infraChange.ProfilePath),
		"",
		fmt.Sprintf("%s/.omgdtmp", infraChange.OutputDir),
		infraChange.Verbosity,
	)

	infraChange.CmdOnDir(
		fmt.Sprintf("omgd build-templates --profile=%s", infraChange.ProfilePath),
		"",
		fmt.Sprintf("%s/.omgdtmp", infraChange.OutputDir),
		infraChange.Verbosity,
	)

	infraChange.CmdOnDir(
		fmt.Sprintf("omgd build-clients --profile=%s", infraChange.ProfilePath),
		"",
		fmt.Sprintf("%s/.omgdtmp", infraChange.OutputDir),
		infraChange.Verbosity,
	)

	infraChange.CmdOnDir(
		fmt.Sprintf("omgd run nakama-server --profile=%s", infraChange.ProfilePath),
		"",
		fmt.Sprintf("%s/.omgdtmp", infraChange.OutputDir),
		infraChange.Verbosity,
	)
}

func (infraChange *InfraChange) DeployInfra() {
	infraChange.setup()

	// NOTE: Would like to discourage this in favor of using utils.Run but testing is easier this way
	infraChange.CmdOnDir(
		fmt.Sprintf("omgd run task deploy-infra --profile=%s", infraChange.ProfilePath),
		"",
		infraChange.OutputDir,
		infraChange.Verbosity,
	)
}

func (infraChange *InfraChange) DestroyInfra() {
	infraChange.setup()

	// NOTE: Would like to discourage this in favor of using utils.Run but testing is easier this way
	infraChange.CmdOnDir(
		fmt.Sprintf("omgd run task destroy-infra --profile=%s", infraChange.ProfilePath),
		"",
		infraChange.OutputDir,
		infraChange.Verbosity,
	)
}

func (infraChange *InfraChange) setup() {
	infraChange.ProfilePath = strings.ReplaceAll(infraChange.ProfilePath, "profiles/", ".omgd/")

	// 1. Should create or empty .omgdtmp directory to work in
	if infraChange.CopyToTmpDir {
		infraChange.tmpDir = infraChange.OutputDir
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

	// 3. Build profiles directory
	BuildProfiles(infraChange.OutputDir, infraChange.Verbosity)
}
