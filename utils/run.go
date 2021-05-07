package utils

import (
	"fmt"
	"strings"
)

type Run struct {
	Profile   *ProfileConf
	OutputDir string
	CmdDir    func(string, string, string)
}

func (r *Run) runCmdOnDir(cmd string, cmdDesc string, cmdDir string) {
	baseCmd := strings.Split(cmd, " ")[0]

	if strings.HasSuffix(baseCmd, "gg") {
		if strings.HasPrefix(cmdDir, "/") {
			cmdDir = cmdDir[1:]
		}

		if strings.HasSuffix(cmd, "/") {
			cmdDir = cmdDir[:len(cmdDir)-1]
		}

		path_prepend := "../"

		for i := 0; i < len(strings.Split(cmdDir, "/"))-1; i++ {
			path_prepend += "../"
		}

		cmd = cmd + " --profile=" + path_prepend + "profiles/" + r.Profile.Name
	}

	r.CmdDir(cmd, cmdDesc, cmdDir)
}

func (r *Run) Run() {
	for _, project := range r.Profile.Main {
		dir := r.OutputDir

		if project.Dir != "" {
			dir = fmt.Sprintf("%s/%s", dir, project.Dir)
		}

		for _, step := range project.Steps {
			stepDir := dir

			if step.Dir != "" {
				stepDir = fmt.Sprintf("%s/%s", stepDir, step.Dir)
			}

			r.runCmdOnDir(step.Cmd, "", stepDir)
		}
	}
}

func (r *Run) RunProjectStep(projectStep string) {
	for _, project := range r.Profile.Main {
		dir := r.OutputDir

		if project.Name != projectStep {
			continue
		}

		if project.Dir != "" {
			dir = fmt.Sprintf("%s/%s", dir, project.Dir)
		}

		for _, step := range project.Steps {
			stepDir := dir

			if step.Dir != "" {
				stepDir = fmt.Sprintf("%s/%s", stepDir, step.Dir)
			}

			r.runCmdOnDir(step.Cmd, "", stepDir)
		}
	}
}

func (r *Run) RunProjectSubStep(projectStep string, index int) {
	for _, project := range r.Profile.Main {
		dir := r.OutputDir

		if project.Name != projectStep {
			continue
		}

		if project.Dir != "" {
			dir = fmt.Sprintf("%s/%s", dir, project.Dir)
		}

		for i, step := range project.Steps {
			if i != index {
				continue
			}

			if step.Dir != "" {
				dir = fmt.Sprintf("%s/%s", dir, step.Dir)
			}

			r.runCmdOnDir(step.Cmd, "", dir)
		}
	}
}

func (r *Run) RunTask(task string) {
	for _, project := range r.Profile.Tasks {
		dir := r.OutputDir

		if project.Name != task {
			continue
		}

		if project.Dir != "" {
			dir = fmt.Sprintf("%s/%s", dir, project.Dir)
		}

		for _, step := range project.Steps {
			stepDir := dir

			if step.Dir != "" {
				stepDir = fmt.Sprintf("%s/%s", stepDir, step.Dir)
			}

			r.runCmdOnDir(step.Cmd, "", stepDir)
		}
	}
}

func (r *Run) RunTaskSubStep(task string, index int) {
	for _, project := range r.Profile.Tasks {
		dir := r.OutputDir

		if project.Name != task {
			continue
		}

		if project.Dir != "" {
			dir = fmt.Sprintf("%s/%s", dir, project.Dir)
		}

		for i, step := range project.Steps {
			if i != index {
				continue
			}

			if step.Dir != "" {
				dir = fmt.Sprintf("%s/%s", dir, step.Dir)
			}

			r.runCmdOnDir(step.Cmd, "", dir)
		}
	}
}
