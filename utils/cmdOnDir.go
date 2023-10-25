package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/logrusorgru/aurora"
)

// CmdOnDir d
func CmdOnDir(cmdStr string, cmdDesc string, cmdDir string, verbosity bool) string {
	cmd := getCmd(cmdStr, cmdDesc, cmdDir, verbosity)

	output, err := cmd.Output()

	if err != nil {
		log.Println(string(output))
		log.Println(aurora.Red("Error!\n"))
		log.Println(err)
		log.Fatal(aurora.Yellow(fmt.Sprintf("Attempted to run: %s\n on dir: %s\n", cmdStr, cmdDir)))
	} else if verbosity {
		log.Println(string(output))
	}

	return string(output)
}

func CmdOnDirWithEnv(cmdStr string, cmdDesc string, cmdDir string, env []string, verbosity bool) string {
	cmd := getCmd(cmdStr, cmdDesc, cmdDir, verbosity)

	cmd.Env = os.Environ()
	for _, envVar := range env {
		cmd.Env = append(cmd.Env, envVar)
	}

	output, err := cmd.Output()

	if err != nil {
		log.Println(string(output))
		log.Println(aurora.Red("Error!\n"))
		log.Println(err)
		log.Fatal(aurora.Yellow(fmt.Sprintf("Attempted to run: %s\n on dir: %s\n", cmdStr, cmdDir)))
	} else if verbosity {
		log.Println(string(output))
	}

	return string(output)
}

func getCmd(cmdStr string, cmdDesc string, cmdDir string, verbosity bool) *exec.Cmd {
	str := strings.Split(cmdStr, " ")

	cmd := exec.Command(str[0], str[1:]...)

	if cmdDir == "" {
		cmd.Dir = "."
	} else {
		cmd.Dir = cmdDir
	}

	log.Print(aurora.Cyan(fmt.Sprintf("%s... ", cmdDesc)))

	if verbosity {
		cmd.Stderr = os.Stderr
	}

	return cmd
}
