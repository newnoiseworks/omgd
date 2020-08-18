package builder

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/newnoiseworks/tpl-fred/builder/config"
	"github.com/newnoiseworks/tpl-fred/utils"
)

// BuildGame doinit
func BuildGame(environment string, buildPath string) {
	config.GameConfig(environment, buildPath)

	path, err := filepath.Abs(fmt.Sprintf("%s/game", buildPath))
	if err != nil {
		log.Fatal(err)
		return
	}

	buildDistro("mac", environment, buildPath, path)
	buildDistro("windows", environment, buildPath, path)
	buildDistro("x11", environment, buildPath, path)
}

func buildDistro(target string, environment string, buildPath string, gamePath string) {
	cmdStr := fmt.Sprintf("BUILD_ENV=%s docker-compose run build-%s", environment, target)

	utils.CmdOnDir(cmdStr, fmt.Sprintf("building %s distro", target), gamePath)

	destroyDockerImage(environment, buildPath, gamePath, target)
}

func destroyDockerImage(environment string, buildPath string, gamePath string, distro string) {
	dockerImgName := fmt.Sprintf("newnoiseworks/game-build-%s-%s", distro, environment)
	cmdStr := fmt.Sprintf("docker rmi -f %s", dockerImgName)

	utils.CmdOnDir(cmdStr, fmt.Sprintf("docker rmi destroys image post build for %s", distro), gamePath)
}
