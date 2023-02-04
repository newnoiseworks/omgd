package utils

import (
	"fmt"
	"log"
	"os"
)

type CodeGenerationPlan struct {
	OutputDir   string
	Target      string
	Plan        string
	Verbosity   bool
	SkipCleanup bool
}

// generates code per plan
func (cp *CodeGenerationPlan) Generate() {
	switch cp.Plan {
	case "new":
		cp.generateNew()
	case "example-2d-player-movement":
		cp.generateExample2DPlayerMovement()
	}

	// optionally skip cleanup to observe files, mostly for testing
	if cp.SkipCleanup == false {
		cp.Cleanup()
	}
}

// generates code needed for new projects
func (cp *CodeGenerationPlan) generateNew() {
	outputPath := fmt.Sprintf("%s/%s", cp.OutputDir, cp.Target)
	err := CopyStaticDirectory("static/new", outputPath)

	if err != nil {
		log.Fatal(err)
	}

	newProfile := GetProfile(fmt.Sprintf("%s/profiles/local", outputPath))
	newProfile.UpdateProfile("game.name", cp.Target)

	if err != nil {
		log.Fatal(err)
	}

	BuildTemplatesFromPath(
		fmt.Sprintf("%s/profiles/local", outputPath),
		outputPath,
		"newomgdtpl",
		!cp.SkipCleanup,
		cp.Verbosity,
	)
}

// generates example 2d player movement code
func (cp *CodeGenerationPlan) generateExample2DPlayerMovement() {
	tmpDir := fmt.Sprintf("%s/.omgdtmp", cp.OutputDir)

	err := CopyStaticDirectory("static/example-2d-player-movement", tmpDir)
	if err != nil {
		log.Fatal(err)
	}

	newProfile := GetProfile(fmt.Sprintf("%s/profiles/local", tmpDir))
	newProfile.UpdateProfile("omgd.channel_name", cp.Target)

	BuildTemplatesFromPath(
		fmt.Sprintf("%s/profiles/local", tmpDir),
		tmpDir,
		"omgdtpl",
		true,
		cp.Verbosity,
	)
}

// cleanup code
func (cp *CodeGenerationPlan) Cleanup() {
	switch cp.Plan {
	// case "new":
	case "example-2d-player-movement":
		cp.cleanupExample2DPlayerMovement()
	}
}

// cleans up example 2d player movement code
func (cp *CodeGenerationPlan) cleanupExample2DPlayerMovement() {
	tmpDir := fmt.Sprintf("%s/.omgdtmp", cp.OutputDir)

	err := os.RemoveAll(fmt.Sprintf("%s/profiles", tmpDir))
	if err != nil {
		log.Fatal(err)
	}

	err = CopyStaticDirectory(tmpDir, cp.OutputDir)
	if err != nil {
		log.Fatal(err)
	}

	err = os.RemoveAll(tmpDir)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: use without this, then if something goes wrong, try it
	// 3. build templates across entire project (for some reason? not sure how to test / prove this / why it happens tbh)
}
