package builder

import "github.com/newnoiseworks/tpl-fred/builder/config"

// BuildServer doinit
func BuildServer(environment string, buildPath string) {
	config.ServerConfig(environment, buildPath)

	// no need to do anything more, just static files
}
