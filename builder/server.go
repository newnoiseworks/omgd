package builder

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/newnoiseworks/tpl-fred/builder/config"
	"github.com/newnoiseworks/tpl-fred/utils"
)

// BuildServer doinit
func BuildServer(environment string, buildPath string) {
	config.ServerConfig(environment, buildPath)

	path, err := filepath.Abs(fmt.Sprintf("%s/server", buildPath))
	if err != nil {
		log.Fatal(err)
		return
	}

	gamePath, err := filepath.Abs(fmt.Sprintf("%s/game", buildPath))
	if err != nil {
		log.Fatal(err)
		return
	}

	cmdStr := fmt.Sprintf("cp -rf %s/dist/web-%s/* %s/nakama/website", gamePath, environment, path)

	utils.CmdOnDir(cmdStr, fmt.Sprintf("Packing website into nakama folder..."), path)
	// no need to do anything more, just static files
}
