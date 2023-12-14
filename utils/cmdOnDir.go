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

	LogDebug(fmt.Sprintf("providing env args %s", strings.Join(env, " ")))

	return runCmd(cmd, cmdStr, cmdDir)
}

func CmdOnDirToStdOut(cmdStr string, cmdDesc string, cmdDir string, env []string) {
	cmd := getCmd(cmdStr, cmdDesc, cmdDir)

	cmd.Env = os.Environ()
	for _, envVar := range env {
		cmd.Env = append(cmd.Env, envVar)
	}

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err := cmd.Run()

	if err != nil {
		LogError(fmt.Sprint(aurora.Red("Error!\n")))
		LogError(fmt.Sprint(err))
		LogFatal(fmt.Sprint(aurora.Yellow(fmt.Sprintf("Attempted to run: %s\n on dir: %s\n", cmdStr, cmdDir))))
	}
}

func getCmd(cmdStr string, cmdDesc string, cmdDir string) *exec.Cmd {
	quoted := false

	str := strings.FieldsFunc(cmdStr, func(r rune) bool {
		if r == '"' {
			quoted = !quoted
		}
		return !quoted && r == ' '
	})

	LogTrace(strings.Join(str, ", "))

	cmd := exec.Command(str[0], str[1:]...)

	if cmdDir == "" {
		cmd.Dir = "."
	} else {
		cmd.Dir = cmdDir
	}

	LogDebug(fmt.Sprint(aurora.Cyan(fmt.Sprintf("%s... ", cmdDesc))))
	LogDebug(fmt.Sprint(aurora.Cyan(fmt.Sprintf("running command %s... ", cmdStr))))

	if GetEnvLogLevel() >= DEBUG_LOG {
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
