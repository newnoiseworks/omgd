package utils

import (
	"fmt"
	"strings"
)

// Run doc
type Run struct {
	Profile     *ProfileConf
	ProfilePath string
	OutputDir   string
	CmdDir      func(string, string, string, bool)
	Verbosity   bool
}

func (r *Run) runCmdOnDir(cmd string, cmdDesc string, cmdDir string) {
	baseCmd := strings.Split(cmd, " ")[0]

	if strings.HasSuffix(baseCmd, "omgd") {
		dir := r.OutputDir

		dir = fmt.Sprintf("%s/%s", dir, cmdDir)

		if strings.HasPrefix(dir, "/") {
			dir = dir[1:]
		}

		if strings.HasPrefix(dir, "./") {
			dir = dir[2:]
		}

		if strings.HasSuffix(dir, "/") {
			dir = dir[:len(dir)-1]
		}

		pathPrepend := ""

		for i := 0; i < len(strings.Split(dir, "/"))-1; i++ {
			pathPrepend += "../"
		}

		cmd = cmd + " --profile=" + pathPrepend + r.ProfilePath
	}

	r.CmdDir(cmd, cmdDesc, cmdDir, r.Verbosity)

	if strings.HasSuffix(baseCmd, "omgd") && strings.Contains(cmd, "update-profile") {
		r.Profile = GetProfile(r.Profile.env)
	}
}

// Run doc
func (r *Run) Run() {
	for _, project := range r.Profile.Main {
		dir := "."

		if r.OutputDir != "" {
			dir = r.OutputDir
		}

		if project.Dir != "" {
			dir = fmt.Sprintf("%s/%s", dir, project.Dir)
		}

		for _, step := range project.Steps {
			stepDir := dir

			if step.Dir != "" {
				stepDir = fmt.Sprintf("%s/%s", stepDir, step.Dir)
			}

			r.runCmdOnDir(step.Cmd, step.Desc, stepDir)
		}
	}
}

// RunProjectStep doc
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

			r.runCmdOnDir(step.Cmd, step.Desc, stepDir)
		}
	}
}

// RunProjectSubStep doc
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

			r.runCmdOnDir(step.Cmd, step.Desc, dir)
		}
	}
}

// RunTask doc
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

			r.runCmdOnDir(step.Cmd, step.Desc, stepDir)
		}
	}
}

// RunTaskSubStep doc
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

			r.runCmdOnDir(step.Cmd, step.Desc, dir)
		}
	}
}
