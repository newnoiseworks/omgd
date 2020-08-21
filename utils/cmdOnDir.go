package utils

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/logrusorgru/aurora"
)

// CmdOnDir d
func CmdOnDir(cmdStr string, cmdDesc string, cmdDir string) {
	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.Dir = cmdDir

	fmt.Print(aurora.Cyan(fmt.Sprintf("Running %s... ", cmdDesc)))

	out, err := cmd.Output()

	if err != nil {
		fmt.Print(aurora.Red("Error!\n"))
		fmt.Printf("%s", out)
		fmt.Println(err)
		log.Fatal(aurora.Yellow(fmt.Sprintf("Attempted to run: %s\n on dir: %s\n", cmdStr, cmdDir)))
	}

	fmt.Print(aurora.Green("Success!\n"))
}
