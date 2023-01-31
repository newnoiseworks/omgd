package utils

// Run doc
type CodePlan struct {
	Profile   *ProfileConf
	OutputDir string
	CmdDir    func(string, string, string, bool)
	Verbosity bool
}

func CodeGenerate(plan string) {

}
