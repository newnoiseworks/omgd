package builder

import (
	"fmt"

	"github.com/newnoiseworks/tpl-fred/builder/config"
)

type Game struct {
	Environment string
	OutputDir   string
	CmdOnDir    func(string, string, string)
}

func (g Game) Build() {
	config.GameConfig(g.Environment, g.OutputDir)

	path := fmt.Sprintf("%s/game", g.OutputDir)

	g.buildDistro("mac", path)
	g.buildDistro("windows", path)
	g.buildDistro("x11", path)
	g.buildDistro("web", path)
}

func (g Game) buildDistro(target string, gamePath string) {
	cmdStr := fmt.Sprintf("BUILD_ENV=%s docker-compose run build-%s", g.Environment, target)

	g.CmdOnDir(cmdStr, fmt.Sprintf("building %s distro", target), gamePath)

	g.destroyDockerImage(target, gamePath)
}

func (g Game) destroyDockerImage(target string, gamePath string) {
	dockerImgName := fmt.Sprintf("newnoiseworks/game-build-%s-%s", target, g.Environment)
	cmdStr := fmt.Sprintf("docker rmi -f %s", dockerImgName)

	g.CmdOnDir(cmdStr, fmt.Sprintf("docker rmi destroys image post build for %s", target), gamePath)
}
