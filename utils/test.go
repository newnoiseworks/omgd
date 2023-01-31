package utils

import (
	"log"
	"strconv"
	"testing"
)

type testCmdOnDirResponse struct {
	cmdStr  string
	cmdDesc string
	cmdDir  string
}

var testCmdResponses = []testCmdOnDirResponse{}
var testValidResponseSet = []testCmdOnDirResponse{}

var testCmdOnDir = func(cmdStr string, cmdDesc string, cmdDir string, verbosity bool) {
	testCmdResponses = append(testCmdResponses, testCmdOnDirResponse{
		cmdStr:  cmdStr,
		cmdDesc: cmdDesc,
		cmdDir:  cmdDir,
	})
}

func testValidCompare(expected interface{}, received interface{}) {
	log.Printf("received %s", received)
	log.Println()
	log.Printf("expected %s", expected)
	log.Println()
}

func testValidCmdSet(t *testing.T, method string) {
	if len(testValidResponseSet) != len(testCmdResponses) {
		t.Errorf("CodeGenerationPlan#Generate didn't create enough commands")
		testValidCompare(strconv.Itoa(len(testValidResponseSet)), strconv.Itoa(len(testCmdResponses)))
	}

	for i := range testValidResponseSet {
		if testValidResponseSet[i] != testCmdResponses[i] {
			t.Errorf("CodeGenerationPlan#Generate failed on step %s", strconv.Itoa(i))
			testValidCompare(testValidResponseSet[i], testCmdResponses[i])
		}
	}

	testValidResponseSet = nil
}
