package deployer

import (
	"fmt"
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

func TestDeployGame(t *testing.T) {
	cmdResponses = nil

	Game{
		Environment: "test",
		OutputDir:   ".tmp",
		CmdOnDir:    cmdOnDir,
	}.Deploy()

	validResponse := cmdOnDirResponse{
		cmdStr:  "butler push ./dist/mac newnoiseworks/the-promised-land-dev:mac-test",
		cmdDesc: "butler push on mac build",
		cmdDir:  ".tmp/game",
	}

	if cmdResponses[0] != validResponse {
		t.Error("Proper commands not being structured")
		fmt.Printf("received %s", cmdResponses[0])
		fmt.Println()
		fmt.Printf("expected %s", validResponse)
		fmt.Println()
	}
}
