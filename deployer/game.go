package deployer

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/newnoiseworks/tpl-fred/utils"
)

var gamePath string

// DeployGame d
func DeployGame(environment string, buildPath string) {
	fmt.Println("deploying game")

	_gamePath, err := filepath.Abs(fmt.Sprintf("%s/game", buildPath))
	if err != nil {
		log.Fatal(err)
		return
	}

	gamePath = _gamePath

	switch environment {
	case "local":
		fmt.Println("Need to make local deployment commands")
		break
	default:
		deployGameBasedOnProfile(environment, buildPath, "mac")
		deployGameBasedOnProfile(environment, buildPath, "windows")
		deployGameBasedOnProfile(environment, buildPath, "x11")
		break
	}
}

func deployGameBasedOnProfile(environment string, buildPath string, distro string) {
	itchGame := "the-promised-land"

	if environment != "production" {
		itchGame = itchGame + "-dev"
	}

	cmd := fmt.Sprintf("butler push ./dist/%s newnoiseworks/%s:%s", distro, itchGame, distro)

	if environment == "production" {
		config := utils.GetProfile(environment)
		cmd = fmt.Sprintf("%s --userversion %s", cmd, config.Game.Version)
	} else {
		cmd = fmt.Sprintf("%s-%s", cmd, environment)
	}

	runCmdOnGameDir(
		cmd,
		fmt.Sprintf("butler push on %s build", distro),
	)
}

func runCmdOnGameDir(cmdStr string, cmdDesc string) {
	utils.CmdOnDir(cmdStr, cmdDesc, gamePath)
}
