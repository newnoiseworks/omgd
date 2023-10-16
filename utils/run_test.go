package utils

import (
	"testing"
)

func TestRunnerCmd(t *testing.T) {
	testCmdOnDirResponses = nil

	profile := GetProfileFromDir("profiles/test.yml", "..")

	runner := Run{
		OutputDir: ".",
		CmdDir:    testCmdOnDir,
		Profile:   profile,
		Verbosity: false,
	}

	testCmdOnDirValidResponseSet = []testCmdOnDirResponse{
		{
			cmdStr:  "omgd build-templates --profile=../../profiles/test.yml",
			cmdDesc: "builds infra templates",
			cmdDir:  "./server/infra",
		},
		{
			cmdStr:  "./infra_deploy.sh",
			cmdDesc: "",
			cmdDir:  "./server/infra/gcp",
		},
		{
			cmdStr:  "omgd build-templates --profile=../profiles/test.yml",
			cmdDesc: "",
			cmdDir:  "./game",
		},
		{
			cmdStr:  "./deploy/gcp/deploy.sh",
			cmdDesc: "",
			cmdDir:  "./server",
		},
	}

	runner.Run()

	testCmdOnDirValidCmdSet(t, "Run#Run")
}
