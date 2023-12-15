package utils

import (
	"fmt"
	"os"
	"testing"
)

func TestServersDeploy(t *testing.T) {
	testDir := "static/test/infra_test_dir"

	t.Cleanup(func() {
		err := os.RemoveAll(
			fmt.Sprintf(
				"%s/.omgd",
				testDir,
			),
		)

		if err != nil {
			t.Fatal(err)
		}

		testCmdOnDirResponses = []testCmdOnDirResponse{}

		profile := GetProfile(fmt.Sprintf("%s/profiles/staging.yml", testDir))

		profile.UpdateProfile("omgd.servers.host", "???")
	})

	profile := GetProfileFromDir("profiles/staging.yml", testDir)

	serversChange := ServersChange{
		OutputDir:       "static/test/infra_test_dir",
		Profile:         profile,
		CmdOnDir:        testCmdOnDir,
		CmdOnDirWithEnv: testCmdOnDirWithEnv,
		SkipCleanup:     true,
	}

	serversChange.Deploy()

	// 1. Should create or empty .omgdtmp directory to work in
	testFileShouldExist(t, fmt.Sprintf("%s/.omgd", testDir))

	testFileShouldExist(t, fmt.Sprintf("%s/.omgd/infra", testDir))

	testFileShouldExist(t, fmt.Sprintf("%s/.omgd/deploy", testDir))
	testFileShouldExist(t, fmt.Sprintf("%s/.omgd/deploy/gcp/deploy.sh", testDir))

	testFileShouldExist(t, fmt.Sprintf("%s/game", testDir))
	testFileShouldExist(t, fmt.Sprintf("%s/profiles", testDir))

	testFileShouldExist(t, fmt.Sprintf("%s/profiles/staging.yml", testDir))

	cmdDirStrTf := fmt.Sprintf("%s/.omgd/infra/gcp/instance-setup/", testDir)

	homeDir, err := os.UserHomeDir()

	if err != nil {
		LogFatal(fmt.Sprintf("Error finding user's home directory %s", err))
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
			cmdStr:  "omgd game build --profile=profiles/staging.yml",
			cmdDesc: "building game clients against profile",
			cmdDir:  testDir,
		},
		{
			cmdStr: "./deploy.sh",
			env: []string{
				"GCP_PROJECT=test",
				"GCP_ZONE=us-east4c",
				"OMGD_PROFILE=staging",
				"OMGD_PROJECT=top-level-name",
				"OMGD_SERVER_SERVICES=central web",
				fmt.Sprintf("CLOUDSDK_AUTH_CREDENTIAL_FILE_OVERRIDE=%s/.config/gcloud/application_default_credentials.json", homeDir),
			},
			cmdDesc: "deploying game server to gcp",
			cmdDir:  fmt.Sprintf("%s/.omgd/deploy/gcp", testDir),
		},
	}

	testForFileAndRegexpMatch(t, fmt.Sprintf("%s/profiles/staging.yml", testDir), "127.6.6.6")

	// 5. Run main task in new .omgdtmp dir profiles/profile.yml file
	testCmdOnDirValidCmdSet(t, "Servers#Deploy")

	serversChange.PerformCleanup()

	testFileShouldNotExist(t, fmt.Sprintf("%s/.omgd/infra/gcp/instance-setup/main.tf", testDir))
	testFileShouldNotExist(t, fmt.Sprintf("%s/.omgd", testDir))
}
