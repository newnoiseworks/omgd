package builder

import (
	"fmt"

	"github.com/newnoiseworks/tpl-fred/builder/config"
)

type Server struct {
	Environment string
	OutputDir   string
	CmdOnDir    func(string, string, string)
}

func (s Server) Build() {
	config.ServerConfig(s.Environment, s.OutputDir)

	cmdStr := fmt.Sprintf("cp -rf game/dist/web-%s/* server/nakama/website", s.Environment)

	s.CmdOnDir(cmdStr, fmt.Sprintf("copy web build into nakama folder..."), s.OutputDir)
}
