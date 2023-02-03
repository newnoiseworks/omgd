package utils

import (
	"fmt"
	"log"
	"os"
)

type CodeGenerationPlan struct {
	OutputDir string
	Target    string
	Plan      string
	Verbosity bool
}

func (cp *CodeGenerationPlan) Generate() {
	outputPath := fmt.Sprintf("%s/%s", cp.OutputDir, cp.Target)

	switch cp.Plan {
	case "new":
		cp.generateNew(outputPath)
	}
}

func (cp *CodeGenerationPlan) generateNew(outputPath string) {
	err := CopyStaticDirectory("static/new", outputPath)

	if err != nil {
		log.Fatal(err)
	}

	newProfile := GetProfile(fmt.Sprintf("%s/profiles/local", outputPath))
	newProfile.UpdateProfile("game.name", cp.Target)

	err = os.Mkdir(fmt.Sprintf("%s/resources", outputPath), 0755)

	if err != nil {
		log.Fatal(err)
	}

	BuildTemplatesFromPath(
		fmt.Sprintf("%s/profiles/local", outputPath),
		outputPath,
		"newomgdtpl",
		true,
		cp.Verbosity,
	)
}
