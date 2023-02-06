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

	testFileShouldNotExist(t, "static/test/newProject/game/project.godot.newomgdtpl")
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

	// generates code plan for 2d movement, skips cleanup for testing
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

	// check to see if profile was edited
	if expected != received {
		testLogComparison(expected, received)

		t.Fatalf("Profile didn't update with channel name in tmp folder")
	}

	// check to see if buildTemplates created new file
	testForFileAndRegexpMatch(
		t,
		"static/test/newProject/.omgdtmp/game/Character/CharacterController.gd",
		`MovementEvent`,
	)

	codePlan.Cleanup()

	// check to make sure profiles dir in tmp folder cleaned up
	testFileShouldNotExist(t, "static/test/newProject/.omgdtmp/profiles/local.yml")

	// check to make sure templates were moved into main folder
	testFileShouldExist(t, "static/test/newProject/game/Character/CharacterController.gd")

	// make sure .omgdtmp folder is cleaned up
	testFileShouldNotExist(t, "static/test/newProject/.omgdtmp")
}

// tests generation of omgd channels
func TestCodeGenCmdOMGDChannelCreation(t *testing.T) {
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

	// generates channel code within new project
	codePlan := CodeGenerationPlan{
		OutputDir:   "static/test/newProject",
		Plan:        "channel",
		Target:      "match_channel",
		SkipCleanup: true,
	}

	codePlan.Generate()

	// checks to make sure local.yml in tmp folder gets updated
	localProfile := GetProfile("static/test/newProject/.omgdtmp/profiles/local")
	expected := "match_channel"
	received := localProfile.Get("omgd.channel_name")

	if expected != received {
		testLogComparison(expected, received)

		t.Fatalf("Profile didn't update with channel name in tmp folder")
	}

	// checks to make sure templates were created and properly renamed
	testForFileAndRegexpMatch(
		t,
		"static/test/newProject/.omgdtmp/game/Autoloads/MatchChannelManager.gd",
		`match_channel`,
	)

	// check for channel event change files
	testForFileAndRegexpMatch(
		t,
		"static/test/newProject/.omgdtmp/game/Autoloads/MatchChannelEvent.gd.tmpl",
		`match_channel`,
	)

	testForFileAndRegexpMatch(
		t,
		"static/test/newProject/.omgdtmp/game/Autoloads/MatchChannelEvent.gd.tmpl",
		`MatchChannel`,
	)

	// check for match_channel_events.yml file
	testForFileAndRegexpMatch(
		t,
		"static/test/newProject/.omgdtmp/resources/match_channel_events.yml",
		`match_channel_events`,
	)

	codePlan.Cleanup()

	// check to make sure profiles dir in tmp folder cleaned up
	testFileShouldNotExist(t, "static/test/newProject/.omgdtmp/profiles/local.yml")

	// check for MatchChannelMUD.tscn file
	testForFileAndRegexpMatch(
		t,
		"static/test/newProject/game/RootScenes/MatchChannelMUD.tscn",
		`MatchChannel`,
	)

	// check for MatchChannelMUDController.tscn file
	testForFileAndRegexpMatch(
		t,
		"static/test/newProject/game/RootScenes/MatchChannelMUDController.tscn",
		`MatchChannel`,
	)

	// check for match_channel.lua file
	testForFileAndRegexpMatch(
		t,
		"static/test/newProject/server/nakama/data/modules/match_channel.lua",
		`match_channel_size`,
	)

	// check for match_channel_manager.lua file
	testForFileAndRegexpMatch(
		t,
		"static/test/newProject/server/nakama/data/modules/match_channel_manager.lua",
		`max_match_channel_size`,
	)

	//
}

// tests file exists and contains a string
func testForFileAndRegexpMatch(t *testing.T, filePath string, search string) {
	// checks to make sure templates were created and properly named
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Cannot find file: %s\n", err)
	}

	// makes sure templates were adjusted with proper variables
	matches, err := regexp.Match(search, file)
	if err != nil {
		t.Fatal(err)
	}
	if !matches {
		t.Fatalf("build-templates didn't adjust %s", filePath)
	}
}

// tests for file not existing
func testFileShouldNotExist(t *testing.T, filePath string) {
	_, err := os.Stat(filePath)
	if !os.IsNotExist(err) {
		t.Fatalf("File exists but should have been cleaned up at %s\n %s", filePath, err)
	}
}

// tests for file existence
func testFileShouldExist(t *testing.T, filePath string) {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		t.Fatalf("File does not exist but should have been created up at %s\n %s", filePath, err)
	}
}
