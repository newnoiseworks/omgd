package utils

import "fmt"

// CheckProjectAndEnv doc
func CheckProjectAndEnv(args []string) bool {
	if len(args) < 2 {
		fmt.Println("You need both arguments dumb dumb")
		return false
	}

	var project = args[0]

	if project != "game" && project != "server" && project != "website" && project != "infra" {
		fmt.Println("Invalid project name")
		return false
	}

	// var environment = args[1]
	// TODO: check for existing profile yml against the given environment name, if none exists exit

	return true
}
