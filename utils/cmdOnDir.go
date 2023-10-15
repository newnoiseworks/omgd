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
func CmdOnDir(cmdStr string, cmdDesc string, cmdDir string, verbosity bool) {
	cmd := getCmd(cmdStr, cmdDesc, cmdDir, verbosity)

	err := cmd.Run()

	if err != nil {
		log.Print(aurora.Red("Error!\n"))
		log.Println(err)
		log.Fatal(aurora.Yellow(fmt.Sprintf("Attempted to run: %s\n on dir: %s\n", cmdStr, cmdDir)))
	}

	// log.Print(aurora.Green("Success!\n"))
}

func CmdOnDirWithEnv(cmdStr string, cmdDesc string, cmdDir string, env []string, verbosity bool) {
	cmd := getCmd(cmdStr, cmdDesc, cmdDir, verbosity)

	cmd.Env = os.Environ()
	for _, envVar := range env {
		cmd.Env = append(cmd.Env, envVar)
	}

	err := cmd.Run()

	if err != nil {
		log.Print(aurora.Red("Error!\n"))
		log.Println(err)
		log.Fatal(aurora.Yellow(fmt.Sprintf("Attempted to run: %s\n on dir: %s\n", cmdStr, cmdDir)))
	}
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
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	return cmd
}
