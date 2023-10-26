package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/logrusorgru/aurora"
)

// CmdOnDir d
func CmdOnDir(cmdStr string, cmdDesc string, cmdDir string, verbosity bool) string {
	cmd := getCmd(cmdStr, cmdDesc, cmdDir, verbosity)

	return runCmd(cmd, cmdStr, cmdDir, verbosity)
}

func CmdOnDirWithEnv(cmdStr string, cmdDesc string, cmdDir string, env []string, verbosity bool) string {
	cmd := getCmd(cmdStr, cmdDesc, cmdDir, verbosity)

	cmd.Env = os.Environ()
	for _, envVar := range env {
		cmd.Env = append(cmd.Env, envVar)
	}

	return runCmd(cmd, cmdStr, cmdDir, verbosity)
}

func getCmd(cmdStr string, cmdDesc string, cmdDir string, verbosity bool) *exec.Cmd {
	str := strings.Split(cmdStr, " ")

	cmd := exec.Command(str[0], str[1:]...)

	if cmdDir == "" {
		cmd.Dir = "."
	} else {
		cmd.Dir = cmdDir
	}

	LogDebug(fmt.Sprint(aurora.Cyan(fmt.Sprintf("%s... ", cmdDesc))))

	if verbosity {
		cmd.Stderr = os.Stderr
	}

	return cmd
}

func runCmd(cmd *exec.Cmd, cmdStr string, cmdDir string, verbosity bool) string {
	output, err := cmd.Output()

	if err != nil {
		LogError(string(output))
		LogError(fmt.Sprint(aurora.Red("Error!\n")))
		LogError(fmt.Sprint(err))
		LogFatal(fmt.Sprint(aurora.Yellow(fmt.Sprintf("Attempted to run: %s\n on dir: %s\n", cmdStr, cmdDir))))
	} else if verbosity {
		LogError(string(output))
	}

	return string(output)
}
