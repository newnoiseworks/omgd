package utils

type testCmdOnDirResponse struct {
	cmdStr  string
	cmdDesc string
	cmdDir  string
}

var testCmdResponses = []testCmdOnDirResponse{}

var testCmdOnDir = func(cmdStr string, cmdDesc string, cmdDir string, verbosity bool) {
	testCmdResponses = append(testCmdResponses, testCmdOnDirResponse{
		cmdStr:  cmdStr,
		cmdDesc: cmdDesc,
		cmdDir:  cmdDir,
	})
}
