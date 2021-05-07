package utils

import "fmt"

type Run struct {
	Profile   string
	OutputDir string
	CmdDir    func(string, string, string)
}

func (r *Run) Run() {
	profile := GetProfile(r.Profile)

	for _, project := range profile.Main {
		dir := r.OutputDir

		if project.Dir != "" {
			dir = fmt.Sprintf("%s/%s", dir, project.Dir)
		}

		for _, step := range project.Steps {
			if step.Dir != "" {
				dir = fmt.Sprintf("%s/%s", dir, step.Dir)
			}

			r.CmdDir(step.Cmd, "", dir)
		}
	}
}

func (r *Run) RunProjectStep(projectStep string) {
	profile := GetProfile(r.Profile)

	for _, project := range profile.Main {
		dir := r.OutputDir

		if project.Name != projectStep {
			continue
		}

		if project.Dir != "" {
			dir = fmt.Sprintf("%s/%s", dir, project.Dir)
		}

		for _, step := range project.Steps {
			if step.Dir != "" {
				dir = fmt.Sprintf("%s/%s", dir, step.Dir)
			}

			r.CmdDir(step.Cmd, "", dir)
		}
	}
}

func (r *Run) RunProjectSubStep(projectStep string, index int) {
	profile := GetProfile(r.Profile)

	for _, project := range profile.Main {
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

			r.CmdDir(step.Cmd, "", dir)
		}
	}
}

func (r *Run) RunTask(task string) {
	profile := GetProfile(r.Profile)

	for _, project := range profile.Tasks {
		dir := r.OutputDir

		if project.Name != task {
			continue
		}

		if project.Dir != "" {
			dir = fmt.Sprintf("%s/%s", dir, project.Dir)
		}

		for _, step := range project.Steps {
			if step.Dir != "" {
				dir = fmt.Sprintf("%s/%s", dir, step.Dir)
			}

			r.CmdDir(step.Cmd, "", dir)
		}
	}
}

func (r *Run) RunTaskSubStep(task string, index int) {
	profile := GetProfile(r.Profile)

	for _, project := range profile.Tasks {
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

			r.CmdDir(step.Cmd, "", dir)
		}
	}
}
