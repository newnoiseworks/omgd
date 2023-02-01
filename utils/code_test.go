package utils

import (
	"testing"
)

func TestCodeGenCmd(t *testing.T) {
	testCmdOnDirResponses = nil

	profile := GetProfile("../profiles/test")

	codePlan := CodeGenerationPlan{
		Profile:     profile,
		ProfilePath: "profiles/test",
		OutputDir:   ".",
		CmdDir:      testCmdOnDir,
		Verbosity:   false,
	}

	// TODO: Write the code that makes this work after you finish
	// static.go code most likely
	testCmdOnDirValidResponseSet = []testCmdOnDirResponse{
		{
			cmdStr:  "something with git or static most likely",
			cmdDesc: "if it is something with the static mod, you'll probably need to adjust the struct to take in a method similar to the CmdDir approach. Consider organizing those \"stub\" methods when you get a chance",
			cmdDir:  "./server/infra",
		},
	}

	codePlan.Generate("plan")

	testCmdOnDirValidCmdSet(t, "CodeGenerationPlan#Generate")
}
