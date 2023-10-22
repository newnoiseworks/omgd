package utils

import (
	"strconv"
	"testing"
)

type testCmdOnDirResponse struct {
	cmdStr    string
	cmdDesc   string
	cmdDir    string
	verbosity bool
}

var testCmdOnDirResponses = []testCmdOnDirResponse{}
var testCmdOnDirValidResponseSet = []testCmdOnDirResponse{}

var testCmdOnDir = func(cmdStr string, cmdDesc string, cmdDir string, verbosity bool) string {
	testCmdOnDirResponses = append(testCmdOnDirResponses, testCmdOnDirResponse{
		cmdStr:    cmdStr,
		cmdDesc:   cmdDesc,
		cmdDir:    cmdDir,
		verbosity: verbosity,
	})

	if cmdStr == "terraform output -raw server_ip" {
		return "127.6.6.6"
	}

	return ""
}

func testCmdOnDirValidCmdSet(t *testing.T, method string) {
	if len(testCmdOnDirValidResponseSet) != len(testCmdOnDirResponses) {
		t.Errorf("%s didn't create enough commands", method)
		testLogComparison(strconv.Itoa(len(testCmdOnDirValidResponseSet)), strconv.Itoa(len(testCmdOnDirResponses)))
	}

	for i := range testCmdOnDirValidResponseSet {
		if testCmdOnDirValidResponseSet[i] != testCmdOnDirResponses[i] {
			t.Errorf("%s failed on step %s", method, strconv.Itoa(i))
			testLogComparison(testCmdOnDirValidResponseSet[i], testCmdOnDirResponses[i])
		}
	}

	testCmdOnDirValidResponseSet = nil
}
