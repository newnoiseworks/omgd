package utils

import (
	"fmt"
	"strings"
)

type ClientBuilder struct {
	Profile         *ProfileConf
	CmdOnDirWithEnv func(string, string, string, []string) string
	Targets         string
}

func (cb *ClientBuilder) Build() {
	buildFor := cb.Targets

	if strings.TrimSpace(buildFor) == "" {
		for x, target := range cb.Profile.OMGD.Game.Targets {
			if x == 0 {
				buildFor = target.BuildService
			} else {
				buildFor = fmt.Sprintf("%s %s", buildFor, target.BuildService)
			}
		}
	}

	cb.CmdOnDirWithEnv(
		// TODO: break below into optional builds per OS based on... profile probably?
		fmt.Sprintf("docker compose up %s", buildFor),
		fmt.Sprintf("Building %s game clients into game/dist folder", cb.Profile.Name),
		"game",
		[]string{
			fmt.Sprintf("BUILD_ENV=%s", cb.Profile.Name),
		},
	)

	// TODO: Need to either internally store docker-compose.yml files ala terraform
	// files or find a way to make the below configuratble, either way, currently this
	// relies on a namespace to be set in those files to work
	if strings.Contains(buildFor, "build-web") {
		sccp := StaticCodeCopyPlan{}

		sccp.CopyStaticDirectory(fmt.Sprintf("game/dist/web-%s", cb.Profile.Name), "servers/web-build/src")
	}

	if strings.Contains(buildFor, "build-x11-server") {
		sccp := StaticCodeCopyPlan{}

		sccp.CopyStaticDirectory(fmt.Sprintf("game/dist/x11-server-%s", cb.Profile.Name), "servers/dedicated-build/src")
	}
}
