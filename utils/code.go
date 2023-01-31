package utils

// TODO: Rename file and test file to code_generate_plan.go

// Run doc
type CodeGenerationPlan struct {
	Profile     *ProfileConf
	ProfilePath string
	OutputDir   string
	CmdDir      func(string, string, string, bool)
	Verbosity   bool
}

func (cp *CodeGenerationPlan) Generate(plan string) {
	cp.CmdDir(
		"something with git or static most likely",
		"if it is something with the static mod, you'll probably need to adjust the struct to take in a method similar to the CmdDir approach. Consider organizing those \"stub\" methods when you get a chance",
		"./server/infra",
		cp.Verbosity,
	)
}
