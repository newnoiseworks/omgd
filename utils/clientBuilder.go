package utils

import (
	"fmt"
	"strings"
)

type ClientBuilder struct {
	Profile             *ProfileConf
	CmdOnDirWithEnv     func(string, string, string, []string) string
	Targets             string
	CopyStaticDirectory func(string, string) error
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
		fmt.Sprintf("docker compose -p %s-%s up %s", cb.Profile.OMGD.Name, cb.Profile.Name, buildFor),
		fmt.Sprintf("Building %s game clients into game/dist folder", cb.Profile.Name),
		"game",
		[]string{
			fmt.Sprintf("BUILD_ENV=%s", cb.Profile.Name),
		},
	)

	for _, target := range cb.Profile.OMGD.Game.Targets {
		if target.Copy != "" && target.To != "" {
			err := cb.CopyStaticDirectory(target.Copy, target.To)

			if err != nil {
				LogFatal(fmt.Sprintf("Error on copying build from game to server folders %s", err))
			}
		}
	}
}
