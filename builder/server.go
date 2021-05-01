package builder

import (
	"fmt"

	"github.com/newnoiseworks/tpl-fred/builder/config"
	"github.com/newnoiseworks/tpl-fred/utils"
)

// BuildServer doinit
func BuildServer(environment string, buildPath string) {
	config.ServerConfig(environment, buildPath)

	cmdStr := fmt.Sprintf("cp -rf game/dist/web-%s/* server/nakama/website", environment)

	utils.CmdOnDir(cmdStr, fmt.Sprintf("copy web build into nakama folder..."), buildPath)
	// no need to do anything more, just static files
}
