package builder

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"

	"github.com/logrusorgru/aurora"
	"github.com/newnoiseworks/tpl-fred/builder/config"
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
	destroyDockerImages(buildPath, path)
}

func buildDistro(target string, environment string, buildPath string, gamePath string) {
	cmdStr := fmt.Sprintf("BUILD_ENV=%s docker-compose run build-%s", environment, target)
	cmd := exec.Command("/bin/sh", "-c", cmdStr)
	cmd.Dir = gamePath

	fmt.Print(aurora.Cyan(fmt.Sprintf("Building %s distro...\n", target)))

	out, err := cmd.Output()

	if err != nil {
		fmt.Printf("%s", out)
		log.Fatal(aurora.Red(fmt.Sprintf("Error building %s distribution using command %s \n", target, cmdStr)))
		return
	}

	fmt.Print(aurora.Green(fmt.Sprintf("%s distro built!\n", target)))
}

func destroyDockerImages(buildPath string, gamePath string) {
	cmdStr := "docker-compose down --rmi=all"
	cmd := exec.Command("/bin/sh", "-c", cmdStr)
	cmd.Dir = gamePath

	path, err := filepath.Abs(fmt.Sprintf("%s/game", buildPath))
	if err != nil {
		log.Fatal(err)
		return
	}

	cmd.Dir = path

	fmt.Print(aurora.Cyan("Destroying docker images...\n"))

	out, err := cmd.Output()

	if err != nil {
		fmt.Printf("%s", out)
		log.Fatal(aurora.Red("Error destroying images"))
		return
	}

	fmt.Print(aurora.Magenta("Docker images destroyed\n"))
}
