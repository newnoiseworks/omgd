package utils

import (
	"fmt"
	"os"
	"testing"
)

func TestDeployInfra(t *testing.T) {
	testDir := "static/test/infra_test_dir"

	t.Cleanup(func() {
		err := os.RemoveAll(fmt.Sprintf("%s/.omgdtmp", testDir))

		if err != nil {
			t.Fatal(err)
		}

		testCmdOnDirResponses = []testCmdOnDirResponse{}
	})

	infraChange := InfraChange{
		OutputDir: "static/test/infra_test_dir",
		Profile:   "staging",
		CmdOnDir:  testCmdOnDir,
		Verbosity: true,
	}

	infraChange.DeployInfra()

	// 1. Should create or empty .omgdtmp directory to work in
	testFileShouldExist(t, fmt.Sprintf("%s/.omgdtmp", testDir))

	// 2. Should clone repo at base of dir (? how to test w/o submodules? clone entire base repo maybe?)
	testFileShouldExist(t, fmt.Sprintf("%s/.omgdtmp/game", testDir))
	testFileShouldExist(t, fmt.Sprintf("%s/.omgdtmp/server", testDir))
	testFileShouldExist(t, fmt.Sprintf("%s/.omgdtmp/profiles", testDir))

	// 3. Copy profiles directory into new .omgdtmp dir (add staging.yml to static/test/infraDir)
	testFileShouldExist(t, fmt.Sprintf("%s/.omgdtmp/profiles/staging.yml", testDir))

	// 4. Build profiles directory in new .omgdtmp dir
	testFileShouldExist(t, fmt.Sprintf("%s/.omgdtmp/.omgd/staging.yml", testDir))

	testCmdOnDirValidResponseSet = []testCmdOnDirResponse{
		{
			cmdStr:    "omgd run --profile=.omgd/staging",
			cmdDesc:   "",
			cmdDir:    fmt.Sprintf("%s/.omgdtmp", testDir),
			verbosity: true,
		},
	}

	// 5. Run main task in new .omgdtmp dir profiles/profile.yml file
	testCmdOnDirValidCmdSet(t, "DeployInfra")
}

func TestDestroyInfra(t *testing.T) {
	testDir := "static/test/infra_test_dir"

	t.Cleanup(func() {
		err := os.RemoveAll(fmt.Sprintf("%s/.omgdtmp", testDir))

		if err != nil {
			t.Fatal(err)
		}

		testCmdOnDirResponses = []testCmdOnDirResponse{}
	})

	infraChange := InfraChange{
		OutputDir: "static/test/infra_test_dir",
		Profile:   "staging",
		CmdOnDir:  testCmdOnDir,
		Verbosity: true,
	}

	infraChange.DestroyInfra()

	// 1. Should create or empty .omgdtmp directory to work in
	testFileShouldExist(t, fmt.Sprintf("%s/.omgdtmp", testDir))

	// 2. Should clone repo at base of dir (? how to test w/o submodules? clone entire base repo maybe?)
	testFileShouldExist(t, fmt.Sprintf("%s/.omgdtmp/game", testDir))
	testFileShouldExist(t, fmt.Sprintf("%s/.omgdtmp/server", testDir))
	testFileShouldExist(t, fmt.Sprintf("%s/.omgdtmp/profiles", testDir))

	// 3. Copy profiles directory into new .omgdtmp dir (add staging.yml to static/test/infraDir)
	testFileShouldExist(t, fmt.Sprintf("%s/.omgdtmp/profiles/staging.yml", testDir))

	// 4. Build profiles directory in new .omgdtmp dir
	testFileShouldExist(t, fmt.Sprintf("%s/.omgdtmp/.omgd/staging.yml", testDir))

	testCmdOnDirValidResponseSet = []testCmdOnDirResponse{
		{
			cmdStr:    "omgd run task destroy-infra --profile=.omgd/staging",
			cmdDesc:   "",
			cmdDir:    fmt.Sprintf("%s/.omgdtmp", testDir),
			verbosity: true,
		},
	}

	// 5. Run destroy-infra task in new .omgdtmp dir profiles/profile.yml file
	testCmdOnDirValidCmdSet(t, "DestroyInfra")
}
