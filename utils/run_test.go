package utils

import (
	"testing"
)

func TestRunnerCmd(t *testing.T) {
	testCmdOnDirResponses = nil

	profile := GetProfile("../profiles/test")

	runner := Run{
		OutputDir:   ".",
		CmdDir:      testCmdOnDir,
		Profile:     profile,
		ProfilePath: "profiles/test",
		Verbosity:   false,
	}

	testCmdOnDirValidResponseSet = []testCmdOnDirResponse{
		{
			cmdStr:  "gg build-templates . --profile=../../profiles/test",
			cmdDesc: "builds infra templates",
			cmdDir:  "./server/infra",
		},
		{
			cmdStr:  "./infra_deploy.sh",
			cmdDesc: "",
			cmdDir:  "./server/infra/gcp",
		},
		{
			cmdStr:  "gg build-templates . --profile=../profiles/test",
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