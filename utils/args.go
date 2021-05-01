package utils

import "fmt"

type BasicArgs struct {
	Environment string
	OutputDir   string
}

// CheckProjectAndEnv doc
func CheckProjectAndEnv(args []string) bool {
	var project = args[0]

	if project != "game" && project != "server" && project != "website" && project != "infra" {
		fmt.Println("Invalid project name")
		return false
	}

	return true
}
