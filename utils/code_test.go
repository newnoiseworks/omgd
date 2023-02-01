package utils

import (
	"testing"
)

func TestCodeGenCmdNewProject(t *testing.T) {
	testCmdOnDirResponses = nil

	profile := GetProfile("../profiles/test")

	codePlan := CodeGenerationPlan{
		Profile:     profile,
		ProfilePath: "profiles/test",
		OutputDir:   ".",
		CmdDir:      testCmdOnDir,
		Verbosity:   false,
	}

	// TODO: old rust code post making new project dir, run commands
	// let update_profile = format!("gg update-profile game.name {}", &name);
	// utils::run_cmd_on_dir(&update_profile, "updates profile w/ game name", &name);

	// utils::run_cmd_on_dir("gg build-templates . --ext=newomgdtpl", "builds templates", &name);
	// utils::run_cmd_on_dir("rm -rf **/*.newomgdtpl", "cleaning...", &name);

	// TODO: Write the code that makes this work after you finish
	// static.go code most likely
	testCmdOnDirValidResponseSet = []testCmdOnDirResponse{
		{
			cmdStr:  "something with git or static most likely",
			cmdDesc: "if it is something with the static mod, you'll probably need to adjust the struct to take in a method similar to the CmdDir approach. Consider organizing those \"stub\" methods when you get a chance",
			cmdDir:  "./server/infra",
		},
	}

	codePlan.Generate("plan")

	testCmdOnDirValidCmdSet(t, "CodeGenerationPlan#Generate")
}
