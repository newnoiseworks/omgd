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
		"if it is something with the static mod, you may need to adjust these tests and see what go has wrt mocks / stubs / and last but not least, spies(!) when it comes to testing",
		"./server/infra",
		cp.Verbosity,
	)
}
