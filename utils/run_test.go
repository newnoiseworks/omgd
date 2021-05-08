package utils

import (
	"fmt"
	"strconv"
	"testing"
)

type cmdOnDirResponse struct {
	cmdStr  string
	cmdDesc string
	cmdDir  string
}

var cmdResponses = []cmdOnDirResponse{}

var cmdOnDir = func(cmdStr string, cmdDesc string, cmdDir string) {
	cmdResponses = append(cmdResponses, cmdOnDirResponse{
		cmdStr:  cmdStr,
		cmdDesc: cmdDesc,
		cmdDir:  cmdDir,
	})
}

func validCompare(expected interface{}, received interface{}) {
	fmt.Printf("received %s", received)
	fmt.Println()
	fmt.Printf("expected %s", expected)
	fmt.Println()
}

func TestRunnerCmd(t *testing.T) {
	cmdResponses = nil

	profile := GetProfile("../profiles/test")

	runner := Run{
		OutputDir: ".",
		CmdDir:    cmdOnDir,
		Profile:   profile,
	}

	validResponseSet := []cmdOnDirResponse{
		{
			cmdStr:  "gg build-templates . --profile=../../profiles/test",
			cmdDesc: "",
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

	if len(validResponseSet) != len(cmdResponses) {
		t.Errorf("Run main project doesn't have enough commands")
		validCompare(strconv.Itoa(len(validResponseSet)), strconv.Itoa(len(cmdResponses)))
	}

	for i := range validResponseSet {
		if validResponseSet[i] != cmdResponses[i] {
			t.Errorf("Run main project failed on step %s", strconv.Itoa(i))
			validCompare(validResponseSet[i], cmdResponses[i])
		}
	}
}
