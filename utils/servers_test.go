package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestServersDeploy(t *testing.T) {
	testDir := filepath.Join("static", "test", "infra_test_dir")

	t.Cleanup(func() {
		err := os.RemoveAll(filepath.Join(testDir, ".omgd"))

		if err != nil {
			LogError(fmt.Sprint(err))
			t.Fail()
		}

		testCmdOnDirResponses = []testCmdOnDirResponse{}

		profile := GetProfile(filepath.Join(testDir, "profiles", "staging.yml"))

		profile.UpdateProfile("omgd.servers.host", "???")
	})

	profile := GetProfileFromDir("profiles/staging.yml", testDir)

	serversChange := ServersChange{
		OutputDir:       testDir,
		Profile:         profile,
		CmdOnDir:        testCmdOnDir,
		CmdOnDirWithEnv: testCmdOnDirWithEnv,
		SkipCleanup:     true,
	}

	serversChange.Deploy()

	testFileShouldExist(t, filepath.Join(testDir, ".omgd"))

	testFileShouldExist(t, filepath.Join(testDir, ".omgd", "infra"))

	testForFileAndRegexpMatch(t, filepath.Join(testDir, ".omgd", "infra", "gcp", "instance-setup", "terraform.tfvars"), "us-east4c")

	testFileShouldExist(t, filepath.Join(testDir, ".omgd", "deploy"))
	testFileShouldExist(t, filepath.Join(testDir, ".omgd", "deploy", "gcp", "deploy.sh"))

	testFileShouldExist(t, filepath.Join(testDir, "game"))
	testFileShouldExist(t, filepath.Join(testDir, "profiles"))

	testFileShouldExist(t, filepath.Join(testDir, "profiles", "staging.yml"))

	cmdDirStrTf := filepath.Join(testDir, ".omgd", "infra", "gcp", "instance-setup")

	homeDir, err := os.UserHomeDir()

	if err != nil {
		LogError(fmt.Sprintf("Error finding user's home directory %s", err))
		t.Fail()
	}

	deployCmd := "./deploy.sh"

	if runtime.GOOS == "Windows" {
		deployCmd = "deploy.bat"
	}

	testCmdOnDirValidResponseSet = []testCmdOnDirResponse{
		{
			cmdStr:  fmt.Sprintf("terraform init -reconfigure -backend-config bucket=%s -backend-config prefix=terraform/state/%s/%s", "???", "top-level-name", profile.Name),
			cmdDesc: fmt.Sprintf("setting up terraform on profile %s", profile.Name),
			cmdDir:  cmdDirStrTf,
		},
		{
			cmdStr:  "terraform output -raw server_ip",
			cmdDesc: "getting ip of newly created server...",
			cmdDir:  cmdDirStrTf,
		},
		{
			cmdStr:  fmt.Sprintf("omgd game build --profile=%s", filepath.Join("profiles", "staging.yml")),
			cmdDesc: "building game clients against profile",
			cmdDir:  testDir,
		},
		{
			cmdStr: deployCmd,
			env: []string{
				"GCP_PROJECT=test",
				"GCP_ZONE=us-east4c",
				"OMGD_PROFILE=staging",
				"OMGD_PROJECT=top-level-name",
				"OMGD_SERVER_SERVICES=central web",
				fmt.Sprintf("CLOUDSDK_AUTH_CREDENTIAL_FILE_OVERRIDE=%s/.config/gcloud/application_default_credentials.json", homeDir),
			},
			cmdDesc: "deploying game server to gcp",
			cmdDir:  filepath.Join(testDir, ".omgd", "deploy", "gcp"),
		},
	}

	testForFileAndRegexpMatch(t, fmt.Sprintf("%s/profiles/staging.yml", testDir), "127.6.6.6")
	testForFileAndRegexpMatch(t, filepath.Join(testDir, "profiles", "staging.yml"), "127.6.6.6")

	testCmdOnDirValidCmdSet(t, "Servers#Deploy")

	serversChange.PerformCleanup()

	testFileShouldNotExist(t, filepath.Join(testDir, ".omgd", "infra", "gcp", "instance-setup", "main.tf"))
	testFileShouldNotExist(t, filepath.Join(testDir, ".omgd"))
}
