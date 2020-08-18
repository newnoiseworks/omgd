package utils

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/logrusorgru/aurora"
)

// CmdOnDir d
func CmdOnDir(cmdStr string, cmdDesc string, cmdDir string) {
	cmd := exec.Command("/bin/sh", "-c", cmdStr)
	cmd.Dir = cmdDir

	fmt.Print(aurora.Cyan(fmt.Sprintf("Running %s...\n", cmdDesc)))

	out, err := cmd.Output()

	if err != nil {
		fmt.Printf("%s", out)
		fmt.Println(err)
		log.Print(aurora.Red(fmt.Sprintf("Error running %s\n", cmdDesc)))
		log.Fatal(aurora.Yellow(fmt.Sprintf("Attempted to run: %s\n", cmdStr)))
	}

	fmt.Print(aurora.Magenta(fmt.Sprintf("Success on %s\n", cmdDesc)))
}
