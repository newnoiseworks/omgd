package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestInstanceSetup(t *testing.T) {
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

	profile := GetProfile(filepath.Join(testDir, "profiles", "staging.yml"))

	infraChange := InfraChange{
		OutputDir:   testDir,
		Profile:     profile,
		CmdOnDir:    testCmdOnDir,
		SkipCleanup: true,
	}

	infraChange.InstanceSetup()

	testFileShouldExist(t, filepath.Join(testDir, ".omgd"))

	testFileShouldExist(t, filepath.Join(testDir, "game"))
	testFileShouldExist(t, filepath.Join(testDir, "profiles"))

	testFileShouldExist(t, filepath.Join(testDir, "profiles", "staging.yml"))

	// testFileShouldExist(t, filepath.Join(testDir, ".omgd", "infra", "gcp", "instance-setup", "terraform.tfvars"))

	// testForFileAndRegexpMatch(t, filepath.Join(testDir, ".omgd", "infra", "gcp", "instance-setup", "terraform.tfvars"), "gcp_project = \"test\"")

	cmdDirStrTf := filepath.Join(testDir, ".omgd", "infra", "gcp", "instance-setup")

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

	testForFileAndRegexpMatch(t, filepath.Join(testDir, "profiles", "staging.yml"), "127.6.6.6")

	infraChange.PerformCleanup()

	testFileShouldNotExist(t, filepath.Join(cmdDirStrTf, "main.tf"))
	testFileShouldNotExist(t, filepath.Join(testDir, ".omgd"))
}

func TestInstanceDestroy(t *testing.T) {
	testDir := filepath.Join("static", "test", "infra_test_dir")

	t.Cleanup(func() {
		err := os.RemoveAll(filepath.Join(testDir, ".omgd"))

		if err != nil {
			LogError(fmt.Sprint(err))
			t.Fail()
		}

		testCmdOnDirResponses = []testCmdOnDirResponse{}
	})

	// profile := GetProfile(fmt.Sprintf("%s/profiles/staging.yml", testDir))
	profile := GetProfileFromDir(filepath.Join("profiles", "staging.yml"), testDir)

	infraChange := InfraChange{
		OutputDir:   testDir,
		Profile:     profile,
		CmdOnDir:    testCmdOnDir,
		SkipCleanup: true,
	}

	infraChange.InstanceDestroy()

	testFileShouldExist(t, filepath.Join(testDir, ".omgd"))

	testFileShouldExist(t, filepath.Join(testDir, "game"))
	testFileShouldExist(t, filepath.Join(testDir, "profiles"))

	testFileShouldExist(t, filepath.Join(testDir, "profiles", "staging.yml"))

	cmdDirStrTf := filepath.Join(testDir, ".omgd", "infra", "gcp", "instance-setup")

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

	testFileShouldNotExist(t, filepath.Join(testDir, ".omgd", "infra", "gcp", "instance-setup", "main.tf"))
	testFileShouldNotExist(t, filepath.Join(testDir, ".omgd"))
}

func TestProjectSetup(t *testing.T) {
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

		GetProfileFromDir(filepath.Join("profiles", "omgd.cloud.yml"), testDir).UpdateProfile("omgd.gcp.bucket", "???")
	})

	profile := GetProfileFromDir(filepath.Join("profiles", "staging.yml"), testDir)

	infraChange := InfraChange{
		OutputDir:       filepath.Join("static", "test", "infra_test_dir"),
		Profile:         profile,
		CmdOnDir:        testCmdOnDir,
		CmdOnDirWithEnv: testCmdOnDirWithEnv,
		SkipCleanup:     true,
	}

	infraChange.ProjectSetup()

	// 1. Should create or empty .omgdtmp directory to work in
	testFileShouldExist(t, filepath.Join(testDir, ".omgd"))

	cmdDirStrTf := filepath.Join(testDir, ".omgd", "infra", "gcp", "project-setup")

	testForFileAndRegexpMatch(t, filepath.Join(testDir, ".omgd", "infra", "gcp", "project-setup", "main.tf"), "gcs")

	if GetProfileFromDir(filepath.Join("profiles", "staging.yml"), testDir).OMGD.GCP.Bucket != "omgd.gcp.bucket" {
		LogError("Bucket name not being set in profile")
		t.Fail()
	}

	testCmdOnDirValidResponseSet = []testCmdOnDirResponse{
		{
			cmdStr:  fmt.Sprintf("terraform init -reconfigure -backend-config path=%s", filepath.Join("..", "..", "..", "..", ".omgd", "terraform.tfstate")),
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

	testFileShouldNotExist(t, filepath.Join(testDir, ".omgd", "infra", "gcp", "project-setup", "main.tf"))
	testFileShouldNotExist(t, filepath.Join(testDir, ".omgd"))
}

func TestProjectDestroy(t *testing.T) {
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

	profile := GetProfileFromDir(filepath.Join("profiles", "staging.yml"), testDir)

	infraChange := InfraChange{
		OutputDir:       testDir,
		Profile:         profile,
		CmdOnDir:        testCmdOnDir,
		CmdOnDirWithEnv: testCmdOnDirWithEnv,
		SkipCleanup:     true,
	}

	infraChange.ProjectDestroy()

	// 1. Should create or empty .omgdtmp directory to work in
	testFileShouldExist(t, filepath.Join(testDir, ".omgd"))

	testForFileAndRegexpMatch(t, filepath.Join(testDir, ".omgd", "infra", "gcp", "project-setup", "main.tf"), "local")

	cmdDirStrTf := filepath.Join(testDir, ".omgd", "infra", "gcp", "project-setup")

	testCmdOnDirValidResponseSet = []testCmdOnDirResponse{
		{
			cmdStr:  fmt.Sprintf("terraform init -reconfigure -backend-config bucket=%s -backend-config prefix=terraform/state/%s", "???", "top-level-name"),
			cmdDesc: "setting up terraform to local backend",
			cmdDir:  cmdDirStrTf,
		},
		{
			cmdStr:  fmt.Sprintf("terraform init -force-copy -backend-config path=%s", filepath.Join("..", "..", "..", "..", ".omgd", "terraform.tfstate")),
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

	testFileShouldNotExist(t, filepath.Join(testDir, ".omgd", "infra", "gcp", "project-setup", "main.tf"))
	testFileShouldNotExist(t, filepath.Join(testDir, ".omgd"))
}
