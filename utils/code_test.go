package utils

import (
	"io/ioutil"
	"os"
	"regexp"
	"testing"
)

// tests generation of new godot projects
func TestCodeGenCmdNewProjectWritesAndCleansUpFiles(t *testing.T) {
	t.Cleanup(func() {
		err := os.RemoveAll("static/test/newProject")

		if err != nil {
			t.Fatal(err)
		}
	})

	codePlan := CodeGenerationPlan{
		OutputDir: "static/test",
		Target:    "newProject",
		Plan:      "new",
	}

	codePlan.Generate()

	localProfile := GetProfile("static/test/newProject/profiles/local")
	expected := "newProject"
	received := localProfile.Get("game.name")

	if expected != received {
		testLogComparison(expected, received)

		t.Fatalf("Profile didn't update with game name in static folder")
	}

	_, err := os.Stat("static/test/newProject/game/project.godot.newomgdtpl")
	if !os.IsNotExist(err) {
		t.Fatal("Templates didn't cleanup afterwards")
	}
}

// tests generation of godot example 2d player movement project
func TestCodeGenCmdExample2DPlayerMovement(t *testing.T) {
	t.Cleanup(func() {
		err := os.RemoveAll("static/test/newProject")

		if err != nil {
			t.Fatal(err)
		}
	})

	// generates a new project to work in
	newProjectCodePlan := CodeGenerationPlan{
		OutputDir: "static/test",
		Target:    "newProject",
		Plan:      "new",
	}

	newProjectCodePlan.Generate()

	codePlan := CodeGenerationPlan{
		OutputDir:   "static/test/newProject",
		Plan:        "example-2d-player-movement",
		Target:      "movement",
		SkipCleanup: true,
	}

	codePlan.Generate()

	localProfile := GetProfile("static/test/newProject/.omgdtmp/profiles/local")
	expected := "movement"
	received := localProfile.Get("omgd.channel_name")

	if expected != received {
		testLogComparison(expected, received)

		t.Fatalf("Profile didn't update with channel name in tmp folder")
	}

	file, err := ioutil.ReadFile("static/test/newProject/.omgdtmp/game/Character/CharacterController.gd")
	if err != nil {
		t.Fatalf("Cannot find file: %s\n", err)
	}

	matches, err := regexp.Match(`MovementEvent`, file)
	if err != nil {
		t.Fatal(err)
	}
	if !matches {
		t.Fatalf("build-templates didn't adjust static/test/newProject/.omgdtmp/game/CharacterController.gd")
	}

	_, err = os.Stat("static/test/newProject/game/profiles/local.yml")
	if os.IsNotExist(err) == false {
		t.Fatal("Profiles folder didn't cleanup afterwards")
	}

	codePlan.Cleanup()

	_, err = os.Stat("static/test/newProject/game/Character/CharacterController.gd")
	if os.IsNotExist(err) {
		t.Fatal("Files were not moved into main project folder post cleanup")
	}

	_, err = os.Stat("static/test/newProject/.omgdtmp")
	if !os.IsNotExist(err) {
		t.Fatal("Temporary folder was not cleaned up")
	}
}

// tests generation of omgd channels
func TestCodeGenCmdOMGDChannelCreation(t *testing.T) {

}
