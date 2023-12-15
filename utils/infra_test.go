package utils

import (
	"fmt"
	"os"
	"testing"
)

func TestInstanceSetup(t *testing.T) {
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

	// profile := GetProfileFromDir("profiles/staging.yml", testDir)
	profile := GetProfile(fmt.Sprintf("%s/profiles/staging.yml", testDir))

	infraChange := InfraChange{
		OutputDir:   "static/test/infra_test_dir",
		Profile:     profile,
		CmdOnDir:    testCmdOnDir,
		SkipCleanup: true,
	}

	infraChange.InstanceSetup()

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
			cmdStr:  fmt.Sprintf("terraform init -reconfigure -backend-config bucket=%s -backend-config prefix=terraform/state/%s/%s", "???", "top-level-name", profile.Name),
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
	testCmdOnDirValidCmdSet(t, "InstanceSetup")

	testForFileAndRegexpMatch(t, fmt.Sprintf("%s/profiles/staging.yml", testDir), "127.6.6.6")

	infraChange.PerformCleanup()

	testFileShouldNotExist(t, fmt.Sprintf("%s/.omgd/infra/gcp/instance-setup/main.tf", testDir))
	testFileShouldNotExist(t, fmt.Sprintf("%s/.omgd", testDir))
}

func TestInstanceDestroy(t *testing.T) {
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

	infraChange.InstanceDestroy()

	// 1. Should create or empty .omgdtmp directory to work in
	testFileShouldExist(t, fmt.Sprintf("%s/.omgd", testDir))

	testFileShouldExist(t, fmt.Sprintf("%s/game", testDir))
	testFileShouldExist(t, fmt.Sprintf("%s/profiles", testDir))

	testFileShouldExist(t, fmt.Sprintf("%s/profiles/staging.yml", testDir))

	cmdDirStrTf := fmt.Sprintf("%s/.omgd/infra/gcp/instance-setup/", testDir)

	testCmdOnDirValidResponseSet = []testCmdOnDirResponse{
		{
			cmdStr:  fmt.Sprintf("terraform init -reconfigure -backend-config bucket=%s -backend-config prefix=terraform/state/%s/%s", "???", "top-level-name", profile.Name),
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
	testCmdOnDirValidCmdSet(t, "InstanceDestroy")

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

		profile.UpdateProfile("omgd.servers.host", "???")

		GetProfileFromDir("profiles/omgd.cloud.yml", testDir).UpdateProfile("omgd.gcp.bucket", "???")
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

	if GetProfileFromDir("profiles/staging.yml", testDir).OMGD.GCP.Bucket != "omgd.gcp.bucket" {
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
			cmdStr:  fmt.Sprintf("terraform init -force-copy -backend-config bucket=%s -backend-config prefix=terraform/state/%s", "omgd.gcp.bucket", "top-level-name"),
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

		profile.UpdateProfile("omgd.servers.host", "???")
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
			cmdStr:  fmt.Sprintf("terraform init -reconfigure -backend-config bucket=%s -backend-config prefix=terraform/state/%s", "???", "top-level-name"),
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
