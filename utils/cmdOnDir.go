package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/logrusorgru/aurora"
)

// CmdOnDir d
func CmdOnDir(cmdStr string, cmdDesc string, cmdDir string, verbosity bool) {
	cmd := exec.Command("bash", "-c", cmdStr)

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

	err := cmd.Run()

	if err != nil {
		log.Print(aurora.Red("Error!\n"))
		log.Println(err)
		log.Fatal(aurora.Yellow(fmt.Sprintf("Attempted to run: %s\n on dir: %s\n", cmdStr, cmdDir)))
	}

	log.Print(aurora.Green("Success!\n"))
}
