package utils

import (
	"fmt"
	"os"
	"testing"
)

func TestDeployInfra(t *testing.T) {
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

		profile.UpdateProfile("omgd.nakama.host", "???")
	})

	// profile := GetProfileFromDir("profiles/staging.yml", testDir)
	profile := GetProfile(fmt.Sprintf("%s/profiles/staging.yml", testDir))

	infraChange := InfraChange{
		OutputDir:   "static/test/infra_test_dir",
		Profile:     profile,
		CmdOnDir:    testCmdOnDir,
		SkipCleanup: true,
	}

	infraChange.DeployInfra()

	// 1. Should create or empty .omgdtmp directory to work in
	testFileShouldExist(t, fmt.Sprintf("%s/.omgd", testDir))

	// 2. Should clone repo at base of dir (? how to test w/o submodules? clone entire base repo maybe?)
	testFileShouldExist(t, fmt.Sprintf("%s/game", testDir))
	testFileShouldExist(t, fmt.Sprintf("%s/profiles", testDir))

	// 3. Copy profiles directory into new .omgdtmp dir (add staging.yml to static/test/infraDir)
	testFileShouldExist(t, fmt.Sprintf("%s/profiles/staging.yml", testDir))

	// 5. BuildTemplates runs
	testFileShouldExist(t, fmt.Sprintf("%s/.omgd/infra/gcp/instance-setup/terraform.tfvars", testDir))

	testForFileAndRegexpMatch(t, fmt.Sprintf("%s/.omgd/infra/gcp/instance-setup/terraform.tfvars", testDir), "gcp_project = \"test\"")

	cmdDirStrTf := fmt.Sprintf("%s/.omgd/infra/gcp/instance-setup/", testDir)

	testCmdOnDirValidResponseSet = []testCmdOnDirResponse{
		{
			cmdStr:  fmt.Sprintf("terraform init -reconfigure -backend-config bucket=%s -backend-config prefix=terraform/state/%s/%s", profile.Get("omgd.tfsettings.bucket"), profile.Get("omgd.name"), profile.Name),
			cmdDesc: "setting up terraform locally",
			cmdDir:  cmdDirStrTf,
		},
		{
			cmdStr:  "terraform apply -auto-approve",
			cmdDesc: "updating cloud infra if needed",
			cmdDir:  cmdDirStrTf,
		},
		{
			cmdStr:  "terraform output -raw server_ip",
			cmdDesc: "getting ip of newly created server...",
			cmdDir:  cmdDirStrTf,
		},
	}

	// 5. Run main task in new .omgdtmp dir profiles/profile.yml file
	testCmdOnDirValidCmdSet(t, "DeployInfra")

	testForFileAndRegexpMatch(t, fmt.Sprintf("%s/profiles/staging.yml", testDir), "127.6.6.6")

	infraChange.PerformCleanup()

	testFileShouldNotExist(t, fmt.Sprintf("%s/.omgd/infra/gcp/instance-setup/main.tf", testDir))
	testFileShouldNotExist(t, fmt.Sprintf("%s/.omgd", testDir))
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
				"%s/.omgd/infra/gcp/instance-setup/terraform.tfvars",
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
		OutputDir:   testDir,
		Profile:     profile,
		CmdOnDir:    testCmdOnDir,
		SkipCleanup: true,
	}

	infraChange.DestroyInfra()

	// 1. Should create or empty .omgdtmp directory to work in
	testFileShouldExist(t, fmt.Sprintf("%s/.omgd", testDir))

	testFileShouldExist(t, fmt.Sprintf("%s/game", testDir))
	testFileShouldExist(t, fmt.Sprintf("%s/profiles", testDir))

	testFileShouldExist(t, fmt.Sprintf("%s/profiles/staging.yml", testDir))

	cmdDirStrTf := fmt.Sprintf("%s/.omgd/infra/gcp/instance-setup/", testDir)

	testCmdOnDirValidResponseSet = []testCmdOnDirResponse{
		{
			cmdStr:  fmt.Sprintf("terraform init -reconfigure -backend-config bucket=%s -backend-config prefix=terraform/state/%s/%s", profile.Get("omgd.tfsettings.bucket"), profile.Get("omgd.name"), profile.Name),
			cmdDesc: fmt.Sprintf("setting up terraform on profile %s", profile.Name),
			cmdDir:  cmdDirStrTf,
		},
		{
			cmdStr:  "terraform destroy -auto-approve",
			cmdDesc: "destroying infrastructure",
			cmdDir:  cmdDirStrTf,
		},
	}

	// 5. Run destroy-infra task in new .omgdtmp dir profiles/profile.yml file
	testCmdOnDirValidCmdSet(t, "DestroyInfra")

	infraChange.PerformCleanup()

	testFileShouldNotExist(t, fmt.Sprintf("%s/.omgd/infra/gcp/instance-setup/main.tf", testDir))
	testFileShouldNotExist(t, fmt.Sprintf("%s/.omgd", testDir))
}

func TestDeployClientAndServer(t *testing.T) {
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

		profile.UpdateProfile("omgd.nakama.host", "???")
	})

	profile := GetProfileFromDir("profiles/staging.yml", testDir)

	infraChange := InfraChange{
		OutputDir:       "static/test/infra_test_dir",
		Profile:         profile,
		CmdOnDir:        testCmdOnDir,
		CmdOnDirWithEnv: testCmdOnDirWithEnv,
		SkipCleanup:     true,
	}

	infraChange.DeployClientAndServer()

	// 1. Should create or empty .omgdtmp directory to work in
	testFileShouldExist(t, fmt.Sprintf("%s/.omgd", testDir))

	testFileShouldExist(t, fmt.Sprintf("%s/game", testDir))
	testFileShouldExist(t, fmt.Sprintf("%s/profiles", testDir))

	testFileShouldExist(t, fmt.Sprintf("%s/profiles/staging.yml", testDir))

	cmdDirStrTf := fmt.Sprintf("%s/.omgd/infra/gcp/instance-setup/", testDir)

	gcloudEnvVars := []string{
		"CLOUDSDK_AUTH_CREDENTIAL_FILE_OVERRIDE=~/.config/gcloud/application_default_credentials.json"}

	testCmdOnDirValidResponseSet = []testCmdOnDirResponse{
		{
			cmdStr:  fmt.Sprintf("terraform init -reconfigure -backend-config bucket=%s -backend-config prefix=terraform/state/%s/%s", profile.Get("omgd.tfsettings.bucket"), profile.Get("omgd.name"), profile.Name),
			cmdDesc: fmt.Sprintf("setting up terraform on profile %s", profile.Name),
			cmdDir:  cmdDirStrTf,
		},
		{
			cmdStr:  "terraform output -raw server_ip",
			cmdDesc: "getting ip of newly created server...",
			cmdDir:  cmdDirStrTf,
		},
		{
			cmdStr:  "omgd build-clients --profile=profiles/staging.yml",
			cmdDesc: "building game clients against profile",
			cmdDir:  testDir,
		},
		{
			cmdStr:  "cp -rf game/dist/web-staging/. server/nakama/website",
			cmdDesc: "copy web build into server",
			cmdDir:  testDir,
		},
		{
			cmdStr: fmt.Sprintf(
				"gcloud compute ssh --project %s --zone %s --command \"truncate -s 0 /var/log/docker.log\" %s-omgd-dev-instance-%s",
				profile.Get("omgd.gcp.project"),
				profile.Get("omgd.gcp.zone"),
				profile.Get("omgd.name"),
				profile.Name,
			),
			env:     gcloudEnvVars,
			cmdDesc: "truncating docker log on server",
			cmdDir:  fmt.Sprintf("%s/server", testDir),
		},
		{
			cmdStr: fmt.Sprintf(
				"gcloud compute scp --project %s --zone %s --force-key-file-overwrite docker-compose.yml %s-omgd-dev-instance-%s:",
				profile.Get("omgd.gcp.project"),
				profile.Get("omgd.gcp.zone"),
				profile.Get("omgd.name"),
				profile.Name,
			),
			env:     gcloudEnvVars,
			cmdDesc: "uploading docker-compose.yml file",
			cmdDir:  fmt.Sprintf("%s/server", testDir),
		},
		{
			cmdStr: fmt.Sprintf(
				"gcloud compute ssh --project %s --zone %s --command \"docker run --rm -v /var/run/docker.sock:/var/run/docker.sock -v \"$PWD:$PWD\" -w=\"$PWD\" docker/compose:1.27.0 down\" %s-omgd-dev-instance-%s",
				profile.Get("omgd.gcp.project"),
				profile.Get("omgd.gcp.zone"),
				profile.Get("omgd.name"),
				profile.Name,
			),
			env:     gcloudEnvVars,
			cmdDesc: "downing running containers on server",
			cmdDir:  fmt.Sprintf("%s/server", testDir),
		},
		{
			cmdStr: fmt.Sprintf(
				"gcloud compute scp --project %s --zone %s --recurse --force-key-file-overwrite nakama %s-omgd-dev-instance-%s:",
				profile.Get("omgd.gcp.project"),
				profile.Get("omgd.gcp.zone"),
				profile.Get("omgd.name"),
				profile.Name,
			),
			env:     gcloudEnvVars,
			cmdDesc: "uploading nakama modules",
			cmdDir:  fmt.Sprintf("%s/server", testDir),
		},
		{
			cmdStr: fmt.Sprintf(
				"gcloud compute ssh --project %s --zone %s --command \"docker run --rm -v /var/run/docker.sock:/var/run/docker.sock -v \"$PWD:$PWD\" -w=\"$PWD\" docker/compose:1.27.0 up -d\" %s-omgd-dev-instance-%s",
				profile.Get("omgd.gcp.project"),
				profile.Get("omgd.gcp.zone"),
				profile.Get("omgd.name"),
				profile.Name,
			),
			env:     gcloudEnvVars,
			cmdDesc: "spinning up docker containers on server",
			cmdDir:  fmt.Sprintf("%s/server", testDir),
		},
	}

	testForFileAndRegexpMatch(t, fmt.Sprintf("%s/profiles/staging.yml", testDir), "127.6.6.6")

	// 5. Run main task in new .omgdtmp dir profiles/profile.yml file
	testCmdOnDirValidCmdSet(t, "DeployInfra")

	infraChange.PerformCleanup()

	testFileShouldNotExist(t, fmt.Sprintf("%s/.omgd/infra/gcp/instance-setup/main.tf", testDir))
	testFileShouldNotExist(t, fmt.Sprintf("%s/.omgd", testDir))
}

func TestProjectSetup(t *testing.T) {
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

		profile.UpdateProfile("omgd.nakama.host", "???")

		GetProfileFromDir("profiles/omgd.yml", testDir).UpdateProfile("omgd.tfsettings.bucket", "???")
	})

	profile := GetProfileFromDir("profiles/staging.yml", testDir)

	infraChange := InfraChange{
		OutputDir:       "static/test/infra_test_dir",
		Profile:         profile,
		CmdOnDir:        testCmdOnDir,
		CmdOnDirWithEnv: testCmdOnDirWithEnv,
		SkipCleanup:     true,
	}

	infraChange.ProjectSetup()

	// 1. Should create or empty .omgdtmp directory to work in
	testFileShouldExist(t, fmt.Sprintf("%s/.omgd", testDir))

	cmdDirStrTf := fmt.Sprintf("%s/.omgd/infra/gcp/project-setup/", testDir)

	testForFileAndRegexpMatch(t, fmt.Sprintf("%s/.omgd/infra/gcp/project-setup/main.tf", testDir), "gcs")

	if GetProfileFromDir("profiles/omgd.yml", testDir).Get("omgd.tfsettings.bucket") != "omgd.tfsettings.bucket" {
		LogError("Bucket name not being set in profile")
		t.Fail()
	}

	testCmdOnDirValidResponseSet = []testCmdOnDirResponse{
		{
			cmdStr:  "terraform init -reconfigure -backend-config path=../../../../.omgd/terraform.tfstate",
			cmdDesc: "setting up terraform locally",
			cmdDir:  cmdDirStrTf,
		},
		{
			cmdStr:  "terraform apply -auto-approve",
			cmdDesc: "setting up initial infra",
			cmdDir:  cmdDirStrTf,
		},
		{
			cmdStr:  "terraform output -raw bucket_name",
			cmdDesc: "getting newly created bucket name",
			cmdDir:  cmdDirStrTf,
		},
		{
			cmdStr:  fmt.Sprintf("terraform init -force-copy -backend-config bucket=%s -backend-config prefix=terraform/state/%s", profile.Get("omgd.tfsettings.bucket"), profile.Get("omgd.name")),
			cmdDesc: "setting up terraform to use gcs backend",
			cmdDir:  cmdDirStrTf,
		},
	}

	testCmdOnDirValidCmdSet(t, "ProjectSetup")

	infraChange.PerformCleanup()

	testFileShouldNotExist(t, fmt.Sprintf("%s/.omgd/infra/gcp/project-setup/main.tf", testDir))
	testFileShouldNotExist(t, fmt.Sprintf("%s/.omgd", testDir))
}

func TestProjectDestroy(t *testing.T) {
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

		profile.UpdateProfile("omgd.nakama.host", "???")
	})

	profile := GetProfileFromDir("profiles/staging.yml", testDir)

	infraChange := InfraChange{
		OutputDir:       "static/test/infra_test_dir",
		Profile:         profile,
		CmdOnDir:        testCmdOnDir,
		CmdOnDirWithEnv: testCmdOnDirWithEnv,
		SkipCleanup:     true,
	}

	infraChange.ProjectDestroy()

	// 1. Should create or empty .omgdtmp directory to work in
	testFileShouldExist(t, fmt.Sprintf("%s/.omgd", testDir))

	testForFileAndRegexpMatch(t, fmt.Sprintf("%s/.omgd/infra/gcp/project-setup/main.tf", testDir), "local")

	cmdDirStrTf := fmt.Sprintf("%s/.omgd/infra/gcp/project-setup/", testDir)

	testCmdOnDirValidResponseSet = []testCmdOnDirResponse{
		{
			cmdStr:  fmt.Sprintf("terraform init -reconfigure -backend-config bucket=%s -backend-config prefix=terraform/state/%s", profile.Get("omgd.tfsettings.bucket"), profile.Get("omgd.name")),
			cmdDesc: "setting up terraform to local backend",
			cmdDir:  cmdDirStrTf,
		},
		{
			cmdStr:  "terraform init -force-copy -backend-config path=../../../../.omgd/terraform.tfstate",
			cmdDesc: "setting up terraform to destroy project level infra",
			cmdDir:  cmdDirStrTf,
		},
		{
			cmdStr:  "terraform destroy -auto-approve",
			cmdDesc: "destroying initial infra",
			cmdDir:  cmdDirStrTf,
		},
	}

	testCmdOnDirValidCmdSet(t, "ProjectDestroy")

	infraChange.PerformCleanup()

	testFileShouldNotExist(t, fmt.Sprintf("%s/.omgd/infra/gcp/project-setup/main.tf", testDir))
	testFileShouldNotExist(t, fmt.Sprintf("%s/.omgd", testDir))
}
