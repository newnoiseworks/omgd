package utils

import (
	"fmt"
	"os"
	"testing"
)

func TestDeployInfra(t *testing.T) {
	testDir := "static/test/infra_test_dir"

	t.Cleanup(func() {
		testCmdOnDirResponses = []testCmdOnDirResponse{}

		profile := GetProfile(fmt.Sprintf("%s/profiles/staging.yml", testDir))

		profile.UpdateProfile("omgd.deploy.server.gcloud.host", "???")
	})

	// profile := GetProfileFromDir("profiles/staging.yml", testDir)
	profile := GetProfile(fmt.Sprintf("%s/profiles/staging.yml", testDir))

	infraChange := InfraChange{
		OutputDir: "static/test/infra_test_dir",
		Profile:   profile,
		CmdOnDir:  testCmdOnDir,
		Verbosity: false,
	}

	infraChange.DeployInfra()

	// 1. Should create or empty .omgdtmp directory to work in
	testFileShouldExist(t, fmt.Sprintf("%s", testDir))

	// 2. Should clone repo at base of dir (? how to test w/o submodules? clone entire base repo maybe?)
	testFileShouldExist(t, fmt.Sprintf("%s/game", testDir))
	testFileShouldExist(t, fmt.Sprintf("%s/server", testDir))
	testFileShouldExist(t, fmt.Sprintf("%s/profiles", testDir))

	// 3. Copy profiles directory into new .omgdtmp dir (add staging.yml to static/test/infraDir)
	testFileShouldExist(t, fmt.Sprintf("%s/profiles/staging.yml", testDir))

	// 5. BuildTemplates runs
	testFileShouldExist(t, fmt.Sprintf("%s/server/infra/gcp/terraform.tfvars", testDir))

	testForFileAndRegexpMatch(t, fmt.Sprintf("%s/server/infra/gcp/terraform.tfvars", testDir), "gcp_project = \"test\"")

	cmdDirStrTf := fmt.Sprintf("%s/server/infra/gcp/", testDir)

	testCmdOnDirValidResponseSet = []testCmdOnDirResponse{
		{
			cmdStr:    fmt.Sprintf("terraform init -reconfigure -force-copy --backend-config path=.omgd/%s/terraform.tfstate", profile.Name),
			cmdDesc:   "setting up terraform locally",
			cmdDir:    cmdDirStrTf,
			verbosity: false,
		},
		{
			cmdStr:    "terraform apply -auto-approve",
			cmdDesc:   "updating cloud infra if needed",
			cmdDir:    cmdDirStrTf,
			verbosity: false,
		},
		{
			cmdStr:    "terraform output -raw server_ip",
			cmdDesc:   "getting ip of newly created server...",
			cmdDir:    cmdDirStrTf,
			verbosity: false,
		},
	}

	// 5. Run main task in new .omgdtmp dir profiles/profile.yml file
	testCmdOnDirValidCmdSet(t, "DeployInfra")

	testForFileAndRegexpMatch(t, fmt.Sprintf("%s/profiles/staging.yml", testDir), "127.6.6.6")
}

func TestDestroyInfra(t *testing.T) {
	testDir := "static/test/infra_test_dir"

	t.Cleanup(func() {
		err := os.RemoveAll(fmt.Sprintf("%s/.omgd", testDir))

		if err != nil {
			t.Fatal(err)
		}

		err = os.RemoveAll(
			fmt.Sprintf(
				"%s/server/infra/gcp/terraform.tfvars",
				testDir,
			),
		)

		if err != nil {
			t.Fatal(err)
		}

		testCmdOnDirResponses = []testCmdOnDirResponse{}
	})

	// profile := GetProfile(fmt.Sprintf("%s/profiles/staging.yml", testDir))
	profile := GetProfileFromDir("profiles/staging.yml", testDir)

	infraChange := InfraChange{
		OutputDir: testDir,
		Profile:   profile,
		CmdOnDir:  testCmdOnDir,
		Verbosity: false,
	}

	infraChange.DestroyInfra()

	testFileShouldExist(t, fmt.Sprintf("%s/game", testDir))
	testFileShouldExist(t, fmt.Sprintf("%s/server", testDir))
	testFileShouldExist(t, fmt.Sprintf("%s/profiles", testDir))

	testFileShouldExist(t, fmt.Sprintf("%s/profiles/staging.yml", testDir))

	cmdDirStrTf := fmt.Sprintf("%s/server/infra/gcp/", testDir)

	// - cmd: terraform init -reconfigure -force-copy --backend-config path=.omgd/{{ .profile.Name }}/terraform.tfstate
	//   desc: setting up terraform on profile {{ .profile.Name }}
	// - cmd: terraform destroy -auto-approve
	//   desc: destroying infrastructure

	testCmdOnDirValidResponseSet = []testCmdOnDirResponse{
		{
			cmdStr:    fmt.Sprintf("terraform init -reconfigure -force-copy --backend-config path=.omgd/%s/terraform.tfstate", profile.Name),
			cmdDesc:   fmt.Sprintf("setting up terraform on profile %s", profile.Name),
			cmdDir:    cmdDirStrTf,
			verbosity: false,
		},
		{
			cmdStr:    "terraform destroy -auto-approve",
			cmdDesc:   "destroying infrastructure",
			cmdDir:    cmdDirStrTf,
			verbosity: false,
		},
	}

	// 5. Run destroy-infra task in new .omgdtmp dir profiles/profile.yml file
	testCmdOnDirValidCmdSet(t, "DestroyInfra")
}

func TestDeployClientAndServer(t *testing.T) {
	testDir := "static/test/infra_test_dir"

	t.Cleanup(func() {
		testCmdOnDirResponses = []testCmdOnDirResponse{}

		profile := GetProfile(fmt.Sprintf("%s/profiles/staging.yml", testDir))

		profile.UpdateProfile("omgd.deploy.server.gcloud.host", "???")
	})

	profile := GetProfileFromDir("profiles/staging.yml", testDir)

	infraChange := InfraChange{
		OutputDir:       "static/test/infra_test_dir",
		Profile:         profile,
		CmdOnDir:        testCmdOnDir,
		CmdOnDirWithEnv: testCmdOnDirWithEnv,
		Verbosity:       false,
	}

	infraChange.DeployClientAndServer()

	testFileShouldExist(t, fmt.Sprintf("%s/game", testDir))
	testFileShouldExist(t, fmt.Sprintf("%s/server", testDir))
	testFileShouldExist(t, fmt.Sprintf("%s/profiles", testDir))

	testFileShouldExist(t, fmt.Sprintf("%s/profiles/staging.yml", testDir))

	cmdDirStrTf := fmt.Sprintf("%s/server/infra/gcp/", testDir)

	testCmdOnDirValidResponseSet = []testCmdOnDirResponse{
		{
			cmdStr:    "terraform output -raw server_ip",
			cmdDesc:   "getting ip of newly created server...",
			cmdDir:    cmdDirStrTf,
			verbosity: false,
		},
		{
			cmdStr:    "omgd build-templates --profile=profiles/staging.yml",
			cmdDesc:   "",
			cmdDir:    testDir,
			verbosity: false,
		},
		{
			cmdStr:    "omgd build-clients --profile=profiles/staging.yml",
			cmdDesc:   "",
			cmdDir:    testDir,
			verbosity: false,
		},
		{
			cmdStr:    "cp -rf ../game/dist/web-staging/* nakama/website",
			cmdDesc:   "copy web build into server",
			cmdDir:    testDir,
			verbosity: false,
		},
		{
			cmdStr:    "./deploy.sh",
			env:       []string{"GCP_PROJECT=test", "GCP_ZONE=us-east4c"},
			cmdDesc:   "deploying game server to gcp",
			cmdDir:    testDir,
			verbosity: false,
		},
	}

	testForFileAndRegexpMatch(t, fmt.Sprintf("%s/profiles/staging.yml", testDir), "127.6.6.6")

	// 5. Run main task in new .omgdtmp dir profiles/profile.yml file
	testCmdOnDirValidCmdSet(t, "DeployInfra")
}
