package utils

import (
	"os"
	"testing"
)

// tests generation of new projects
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
