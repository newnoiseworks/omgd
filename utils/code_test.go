package utils

import (
	"os"
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

	testFileShouldExist(t, "static/test/newProject/resources")
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
		"static/test/newProject/game/RootScenes/MatchChannelMUDController.gd",
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
}

// tests generation of omgd channels with event args
func TestCodeGenCmdOMGDChannelCreationWithEventArgs(t *testing.T) {
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
		Args:        "movement trade",
		SkipCleanup: true,
	}

	codePlan.Generate()

	// checks to make sure local.yml in tmp folder gets updated w/ events
	localProfile := GetProfile("static/test/newProject/.omgdtmp/profiles/local")
	expectedArr := [2]string{"movement", "trade"}
	receivedArr := localProfile.GetArray("omgd.channel_events")

	for i := 0; i < 2; i++ {
		if expectedArr[i] != receivedArr[i].(string) {
			testLogComparison(expectedArr, receivedArr)

			t.Fatalf("Profile didn't update with channel events in tmp folder")
		}
	}

	// check for match_channel_events.yml file
	testForFileAndRegexpMatch(
		t,
		"static/test/newProject/.omgdtmp/resources/match_channel_events.yml",
		`- movement\n- trade`,
	)

	codePlan.Cleanup()

	// TODO: move buildProfiles code to utils
	CmdOnDir("omgdd build-profiles", "", "static/test/newProject", false)

	BuildTemplatesFromPath(
		"static/test/newProject/.gg/local",
		"static/test/newProject",
		"tmpl",
		false,
		false,
	)

	// check for MatchChannelEvent.gd file
	testForFileAndRegexpMatch(
		t,
		"static/test/newProject/game/Autoloads/MatchChannelEvent.gd",
		`MOVEMENT = 0,\n  TRADE = 1`,
	)
}
