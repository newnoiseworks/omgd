package utils

import "fmt"

// CheckProjectAndEnv doc
func CheckProjectAndEnv(args []string) bool {
	if len(args) < 2 {
		fmt.Println("You need both arguments dumb dumb")
		return false
	}

	var project = args[0]
	var environment = args[1]

	if project != "game" && project != "server" && project != "website" && project != "config-files" {
		fmt.Println("Invalid project name")
		return false
	}

	if environment != "production" && environment != "staging" && environment != "local" {
		fmt.Println("Invalid environment name")
		return false
	}

	return true
}
