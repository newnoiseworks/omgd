package utils

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
)

type testCmdOnDirResponse struct {
	cmdStr  string
	cmdDesc string
	cmdDir  string
	env     []string
}

var testCmdOnDirResponses = []testCmdOnDirResponse{}
var testCmdOnDirValidResponseSet = []testCmdOnDirResponse{}

var testCmdOnDir = func(cmdStr string, cmdDesc string, cmdDir string) string {
	testCmdOnDirResponses = append(testCmdOnDirResponses, testCmdOnDirResponse{
		cmdStr:  cmdStr,
		cmdDesc: cmdDesc,
		cmdDir:  cmdDir,
	})

	if cmdStr == "terraform output -raw server_ip" {
		return "127.6.6.6"
	} else if cmdStr == "terraform output -raw bucket_name" {
		return "omgd.gcp.bucket"
	}

	return ""
}

var testCmdOnDirWithEnv = func(cmdStr string, cmdDesc string, cmdDir string, env []string) string {
	testCmdOnDirResponses = append(testCmdOnDirResponses, testCmdOnDirResponse{
		cmdStr:  cmdStr,
		cmdDesc: cmdDesc,
		cmdDir:  cmdDir,
		env:     env,
	})

	if strings.Contains(cmdStr, "gcloud compute instances list") {
		return "[]"
	}

	return ""
}

func testCmdOnDirValidCmdSet(t *testing.T, method string) {
	if len(testCmdOnDirValidResponseSet) != len(testCmdOnDirResponses) {
		t.Errorf("%s didn't create enough commands", method)
		testLogComparison(strconv.Itoa(len(testCmdOnDirValidResponseSet)), strconv.Itoa(len(testCmdOnDirResponses)))
	}

	for i := range testCmdOnDirValidResponseSet {
		if !reflect.DeepEqual(testCmdOnDirValidResponseSet[i], testCmdOnDirResponses[i]) {
			t.Errorf("%s failed on step %s", method, strconv.Itoa(i))
			testLogComparison(testCmdOnDirValidResponseSet[i], testCmdOnDirResponses[i])
		}
	}

	testCmdOnDirValidResponseSet = nil
}
