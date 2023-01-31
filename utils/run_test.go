package utils

import (
	"log"
	"strconv"
	"testing"
)

func validCompare(expected interface{}, received interface{}) {
	log.Printf("received %s", received)
	log.Println()
	log.Printf("expected %s", expected)
	log.Println()
}

func TestRunnerCmd(t *testing.T) {
	testCmdResponses = nil

	profile := GetProfile("../profiles/test")

	runner := Run{
		OutputDir:   ".",
		CmdDir:      testCmdOnDir,
		Profile:     profile,
		ProfilePath: "profiles/test",
		Verbosity:   false,
	}

	validResponseSet := []testCmdOnDirResponse{
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

	if len(validResponseSet) != len(testCmdResponses) {
		t.Errorf("Run main project doesn't have enough commands")
		validCompare(strconv.Itoa(len(validResponseSet)), strconv.Itoa(len(testCmdResponses)))
	}

	for i := range validResponseSet {
		if validResponseSet[i] != testCmdResponses[i] {
			t.Errorf("Run main project failed on step %s", strconv.Itoa(i))
			validCompare(validResponseSet[i], testCmdResponses[i])
		}
	}
}
