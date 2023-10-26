package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/logrusorgru/aurora"
)

// CmdOnDir d
func CmdOnDir(cmdStr string, cmdDesc string, cmdDir string) string {
	cmd := getCmd(cmdStr, cmdDesc, cmdDir)

	return runCmd(cmd, cmdStr, cmdDir)
}

func CmdOnDirWithEnv(cmdStr string, cmdDesc string, cmdDir string, env []string) string {
	cmd := getCmd(cmdStr, cmdDesc, cmdDir)

	cmd.Env = os.Environ()
	for _, envVar := range env {
		cmd.Env = append(cmd.Env, envVar)
	}

	return runCmd(cmd, cmdStr, cmdDir)
}

func getCmd(cmdStr string, cmdDesc string, cmdDir string) *exec.Cmd {
	str := strings.Split(cmdStr, " ")

	cmd := exec.Command(str[0], str[1:]...)

	if cmdDir == "" {
		cmd.Dir = "."
	} else {
		cmd.Dir = cmdDir
	}

	LogDebug(fmt.Sprint(aurora.Cyan(fmt.Sprintf("%s... ", cmdDesc))))

	if GetEnvLogLevel() == DEBUG_LOG {
		cmd.Stderr = os.Stderr
	}

	return cmd
}

func runCmd(cmd *exec.Cmd, cmdStr string, cmdDir string) string {
	output, err := cmd.Output()

	if err != nil {
		LogError(string(output))
		LogError(fmt.Sprint(aurora.Red("Error!\n")))
		LogError(fmt.Sprint(err))
		LogFatal(fmt.Sprint(aurora.Yellow(fmt.Sprintf("Attempted to run: %s\n on dir: %s\n", cmdStr, cmdDir))))
	}

	LogDebug(string(output))

	return string(output)
}
