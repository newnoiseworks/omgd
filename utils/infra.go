package utils

import (
	"fmt"
	"log"
	"os"
)

type InfraChange struct {
	OutputDir   string
	ProfilePath string
	CmdOnDir    func(string, string, string, bool)
	Verbosity   bool
}

func (infraChange *InfraChange) DeployInfra() {
	infraChange.setup()

	// NOTE: Would like to discourage this in favor of using utils.Run but testing is easier this way
	infraChange.CmdOnDir(
		"omgd run --profile=.omgd/staging",
		"",
		fmt.Sprintf("%s/.omgdtmp", infraChange.OutputDir),
		infraChange.Verbosity,
	)
}

func (infraChange *InfraChange) DestroyInfra() {
	infraChange.setup()

	// NOTE: Would like to discourage this in favor of using utils.Run but testing is easier this way
	infraChange.CmdOnDir(
		"omgd run task destroy-infra --profile=.omgd/staging",
		"",
		fmt.Sprintf("%s/.omgdtmp", infraChange.OutputDir),
		infraChange.Verbosity,
	)
}

func (infraChange *InfraChange) setup() {
	// 1. Should create or empty .omgdtmp directory to work in
	tmpDir := fmt.Sprintf("%s/.omgdtmp", infraChange.OutputDir)

	_, err := os.Stat(tmpDir)
	if !os.IsNotExist(err) {
		err = os.RemoveAll(tmpDir)

		if err != nil {
			log.Fatal(err)
		}
	}

	err = os.Mkdir(tmpDir, 0755)

	if err != nil {
		log.Fatal(err)
	}

	// 2. Should clone repo at base of dir (? how to test w/o submodules? clone entire base repo maybe?)
	sccp := StaticCodeCopyPlan{
		skipPaths: []string{
			tmpDir,
		},
	}

	sccp.CopyStaticDirectory(infraChange.OutputDir, tmpDir)

	// 3. Build profiles directory in new .omgdtmp dir
	BuildProfiles(tmpDir, infraChange.Verbosity)
}
