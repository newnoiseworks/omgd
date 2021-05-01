package deployer

import (
	"fmt"

	"github.com/newnoiseworks/tpl-fred/utils"
)

type Game struct {
	Environment string
	OutputDir   string
	CmdOnDir    func(string, string, string)
}

func (dg Game) Deploy() {
	fmt.Println("deploying game")

	switch dg.Environment {
	case "local":
		fmt.Println("Need to make local deployment commands")
		break
	default:
		dg.deployGameBasedOnProfile("mac")
		dg.deployGameBasedOnProfile("windows")
		dg.deployGameBasedOnProfile("x11")
		break
	}
}

func (dg Game) deployGameBasedOnProfile(distro string) {
	itchGame := "the-promised-land"

	if dg.Environment != "production" {
		itchGame = itchGame + "-dev"
	}

	itchDistro := distro

	if distro == "x11" {
		itchDistro = "linux"
	}

	cmd := fmt.Sprintf("butler push ./dist/%s newnoiseworks/%s:%s", distro, itchGame, itchDistro)

	if dg.Environment == "production" {
		config := utils.GetProfile(dg.Environment)
		cmd = fmt.Sprintf("%s --userversion %s", cmd, config.Game.Version)
	} else {
		cmd = fmt.Sprintf("%s-%s", cmd, dg.Environment)
	}

	dg.CmdOnDir(
		cmd,
		fmt.Sprintf("butler push on %s build", distro),
		fmt.Sprintf("%s/game", dg.OutputDir),
	)
}
