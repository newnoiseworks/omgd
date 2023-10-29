package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestDeployInfra(t *testing.T) {
	testDir := "static/test/infra_test_dir"

	t.Cleanup(func() {
		err := os.RemoveAll(
			fmt.Sprintf(
				"%s/server/infra/gcp/terraform.tfvars",
				testDir,
			),
		)

		if err != nil {
			t.Fatal(err)
		}

		testCmdOnDirResponses = []testCmdOnDirResponse{}

		profile := GetProfile(fmt.Sprintf("%s/profiles/staging.yml", testDir))

		profile.UpdateProfile("omgd.gcp.host", "???")
	})

	// profile := GetProfileFromDir("profiles/staging.yml", testDir)
	profile := GetProfile(fmt.Sprintf("%s/profiles/staging.yml", testDir))

	infraChange := InfraChange{
		OutputDir: "static/test/infra_test_dir",
		Profile:   profile,
		CmdOnDir:  testCmdOnDir,
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
			cmdStr:  fmt.Sprintf("terraform init -reconfigure -force-copy -backend-config bucket=%s -backend-config prefix=terraform/state/%s", profile.Get("omgd.tfsettings.bucket"), profile.Name),
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
	}

	infraChange.DestroyInfra()

	testFileShouldExist(t, fmt.Sprintf("%s/game", testDir))
	testFileShouldExist(t, fmt.Sprintf("%s/server", testDir))
	testFileShouldExist(t, fmt.Sprintf("%s/profiles", testDir))

	testFileShouldExist(t, fmt.Sprintf("%s/profiles/staging.yml", testDir))

	cmdDirStrTf := fmt.Sprintf("%s/server/infra/gcp/", testDir)

	testCmdOnDirValidResponseSet = []testCmdOnDirResponse{
		{
			cmdStr:  fmt.Sprintf("terraform init -reconfigure -force-copy -backend-config bucket=%s -backend-config prefix=terraform/state/%s", profile.Get("omgd.tfsettings.bucket"), profile.Name),
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
}

func TestDeployClientAndServer(t *testing.T) {
	testDir := "static/test/infra_test_dir"

	t.Cleanup(func() {
		err := os.RemoveAll(
			fmt.Sprintf(
				"%s/server/infra/gcp/terraform.tfvars",
				testDir,
			),
		)

		if err != nil {
			t.Fatal(err)
		}

		testCmdOnDirResponses = []testCmdOnDirResponse{}

		profile := GetProfile(fmt.Sprintf("%s/profiles/staging.yml", testDir))

		profile.UpdateProfile("omgd.gcp.host", "???")
	})

	profile := GetProfileFromDir("profiles/staging.yml", testDir)

	infraChange := InfraChange{
		OutputDir:       "static/test/infra_test_dir",
		Profile:         profile,
		CmdOnDir:        testCmdOnDir,
		CmdOnDirWithEnv: testCmdOnDirWithEnv,
	}

	infraChange.DeployClientAndServer()

	testFileShouldExist(t, fmt.Sprintf("%s/game", testDir))
	testFileShouldExist(t, fmt.Sprintf("%s/server", testDir))
	testFileShouldExist(t, fmt.Sprintf("%s/profiles", testDir))

	testFileShouldExist(t, fmt.Sprintf("%s/profiles/staging.yml", testDir))

	cmdDirStrTf := fmt.Sprintf("%s/server/infra/gcp/", testDir)

	testCmdOnDirValidResponseSet = []testCmdOnDirResponse{
		{
			cmdStr:  fmt.Sprintf("terraform init -reconfigure -force-copy -backend-config bucket=%s -backend-config prefix=terraform/state/%s", profile.Get("omgd.tfsettings.bucket"), profile.Name),
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
			cmdStr:  "./deploy.sh",
			env:     []string{"GCP_PROJECT=test", "GCP_ZONE=us-east4c", "OMGD_PROFILE=staging", "OMGD_PROJECT=top-level-name"},
			cmdDesc: "deploying game server to gcp",
			cmdDir:  fmt.Sprintf("%s/server/deploy/gcp", testDir),
		},
	}

	testForFileAndRegexpMatch(t, fmt.Sprintf("%s/profiles/staging.yml", testDir), "127.6.6.6")

	// 5. Run main task in new .omgdtmp dir profiles/profile.yml file
	testCmdOnDirValidCmdSet(t, "DeployInfra")
}

func TestProjectSetup(t *testing.T) {
	testDir := "static/test/infra_test_dir"

	t.Cleanup(func() {
		err := os.RemoveAll(
			fmt.Sprintf(
				"%s/server/infra/gcp/terraform.tfvars",
				testDir,
			),
		)

		if err != nil {
			t.Fatal(err)
		}

		tfFilePath := fmt.Sprintf("%s/server/infra/project-setup/gcp/main.tf", testDir)
		input, err := ioutil.ReadFile(tfFilePath)
		if err != nil {
			LogFatal(fmt.Sprint(err))
		}

		lines := strings.Split(string(input), "\n")

		for i, line := range lines {
			if strings.Contains(line, "backend \"gcs\"") {
				lines[i] = strings.Replace(lines[i], "backend \"gcs\"", "backend \"local\"", 1)
			}
		}
		output := strings.Join(lines, "\n")
		err = ioutil.WriteFile(tfFilePath, []byte(output), 0644)
		if err != nil {
			LogFatal(fmt.Sprint(err))
		}

		testCmdOnDirResponses = []testCmdOnDirResponse{}

		profile := GetProfile(fmt.Sprintf("%s/profiles/staging.yml", testDir))

		profile.UpdateProfile("omgd.gcp.host", "???")

		GetProfileFromDir("profiles/omgd.yml", testDir).UpdateProfile("omgd.tfsettings.bucket", "???")
	})

	profile := GetProfileFromDir("profiles/staging.yml", testDir)

	infraChange := InfraChange{
		OutputDir:       "static/test/infra_test_dir",
		Profile:         profile,
		CmdOnDir:        testCmdOnDir,
		CmdOnDirWithEnv: testCmdOnDirWithEnv,
	}

	infraChange.ProjectSetup()

	cmdDirStrTf := fmt.Sprintf("%s/server/infra/project-setup/gcp/", testDir)

	testForFileAndRegexpMatch(t, fmt.Sprintf("%s/server/infra/project-setup/gcp/main.tf", testDir), "gcs")

	if GetProfileFromDir("profiles/omgd.yml", testDir).Get("omgd.tfsettings.bucket") != "omgd.tfsettings.bucket" {
		LogError("Bucket name not being set in profile")
		t.Fail()
	}

	testCmdOnDirValidResponseSet = []testCmdOnDirResponse{
		{
			cmdStr:  "terraform init -reconfigure -force-copy -backend-config path=../../../../.omgd/terraform.tfstate",
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
}

func TestProjectDestroy(t *testing.T) {
	testDir := "static/test/infra_test_dir"

	t.Cleanup(func() {
		err := os.RemoveAll(
			fmt.Sprintf(
				"%s/server/infra/gcp/terraform.tfvars",
				testDir,
			),
		)

		if err != nil {
			t.Fatal(err)
		}

		tfFilePath := fmt.Sprintf("%s/server/infra/project-setup/gcp/main.tf", testDir)
		input, err := ioutil.ReadFile(tfFilePath)
		if err != nil {
			LogFatal(fmt.Sprint(err))
		}

		lines := strings.Split(string(input), "\n")

		for i, line := range lines {
			if strings.Contains(line, "backend \"gcs\"") {
				lines[i] = strings.Replace(lines[i], "backend \"gcs\"", "backend \"local\"", 1)
			}
		}
		output := strings.Join(lines, "\n")
		err = ioutil.WriteFile(tfFilePath, []byte(output), 0644)
		if err != nil {
			LogFatal(fmt.Sprint(err))
		}

		testCmdOnDirResponses = []testCmdOnDirResponse{}

		profile := GetProfile(fmt.Sprintf("%s/profiles/staging.yml", testDir))

		profile.UpdateProfile("omgd.gcp.host", "???")
	})

	tfFilePath := fmt.Sprintf("%s/server/infra/project-setup/gcp/main.tf", testDir)
	input, err := ioutil.ReadFile(tfFilePath)
	if err != nil {
		LogFatal(fmt.Sprint(err))
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "backend \"local\"") {
			lines[i] = strings.Replace(lines[i], "backend \"local\"", "backend \"gcs\"", 1)
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(tfFilePath, []byte(output), 0644)
	if err != nil {
		LogFatal(fmt.Sprint(err))
	}

	profile := GetProfileFromDir("profiles/staging.yml", testDir)

	infraChange := InfraChange{
		OutputDir:       "static/test/infra_test_dir",
		Profile:         profile,
		CmdOnDir:        testCmdOnDir,
		CmdOnDirWithEnv: testCmdOnDirWithEnv,
	}

	infraChange.ProjectDestroy()

	testForFileAndRegexpMatch(t, fmt.Sprintf("%s/server/infra/project-setup/gcp/main.tf", testDir), "local")

	cmdDirStrTf := fmt.Sprintf("%s/server/infra/project-setup/gcp/", testDir)

	testCmdOnDirValidResponseSet = []testCmdOnDirResponse{
		{
			cmdStr:  fmt.Sprintf("terraform init -reconfigure -force-copy -backend-config bucket=%s -backend-config prefix=terraform/state/%s", profile.Get("omgd.tfsettings.bucket"), profile.Get("omgd.name")),
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
}
