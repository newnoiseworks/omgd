package utils

import (
	"testing"
)

func TestCodeGenCmd(t *testing.T) {
	testCmdResponses = nil

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
	testValidResponseSet = []testCmdOnDirResponse{
		{
			cmdStr:  "something with git or static most likely",
			cmdDesc: "if it is something with the static mod, you may need to adjust these tests and see what go has wrt mocks / stubs / and last but not least, spies(!) when it comes to testing",
			cmdDir:  "./server/infra",
		},
	}

	codePlan.Generate("plan")

	testValidCmdSet(t, "CodeGenerationPlan#Generate")
}
