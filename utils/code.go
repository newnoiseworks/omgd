package utils

import (
	"fmt"
	"os"
	"strings"
)

type CodeGenerationPlan struct {
	OutputDir   string
	Target      string
	Plan        string
	Args        string
	SkipCleanup bool
}

// generates code per plan
func (cp *CodeGenerationPlan) Generate() {
	cp.resetOMGDTmpDir()

	switch cp.Plan {
	case "new":
		cp.generateNew()
	case "example-partial-2d-player-movement":
		cp.generateExample2DPlayerMovement()
	case "example-2d-player-movement":
		cp.generateExampleComplete2DPlayerMovement()
	case "channel":
		cp.generateChannel()
	default:
		LogFatal(fmt.Sprintf("Found no code plan matching %s", cp.Plan))
	}

	// optionally skip cleanup to observe files, mostly for testing
	if cp.SkipCleanup == false {
		cp.Cleanup()
	}
}

// generates code needed for new projects
func (cp *CodeGenerationPlan) generateNew() {
	outputPath := fmt.Sprintf("%s/%s", cp.OutputDir, cp.Target)

	sccp := StaticCodeCopyPlan{}

	err := sccp.CopyStaticDirectory("static/new", outputPath)

	if err != nil {
		LogFatal(fmt.Sprint(err))
	}

	newProfile := GetProfile(fmt.Sprintf("%s/profiles/omgd.yml", outputPath))
	newProfile.UpdateProfile("omgd.name", cp.Target)

	if err != nil {
		LogFatal(fmt.Sprint(err))
	}

	BuildTemplatesFromPath(
		newProfile,
		outputPath,
		"omgdtpl",
		!cp.SkipCleanup,
	)

	err = os.Mkdir(fmt.Sprintf("%s/resources", outputPath), 0755)
	if err != nil {
		LogFatal(fmt.Sprint(err))
	}
}

// generates example 2d player movement code
// TODO: Write tests
func (cp *CodeGenerationPlan) generateExampleComplete2DPlayerMovement() {
	outputPath := fmt.Sprintf("%s/%s", cp.OutputDir, cp.Target)

	sccp := StaticCodeCopyPlan{}

	err := sccp.CopyStaticDirectory("static/example-2d-player-movement-complete", outputPath)
	if err != nil {
		LogFatal(fmt.Sprint(err))
	}
}

// generates example 2d player movement code
func (cp *CodeGenerationPlan) generateExample2DPlayerMovement() {
	tmpDir := fmt.Sprintf("%s/.omgdtmp", cp.OutputDir)

	sccp := StaticCodeCopyPlan{}

	err := sccp.CopyStaticDirectory("static/example-2d-player-movement", tmpDir)
	if err != nil {
		LogFatal(fmt.Sprint(err))
	}

	newProfile := GetProfile(fmt.Sprintf("%s/profiles/local.yml", tmpDir))
	newProfile.UpdateProfile("omgd.channel_name", cp.Target)

	BuildTemplatesFromPath(
		newProfile,
		tmpDir,
		"omgdtpl",
		true,
	)
}

// generates channel management code
func (cp *CodeGenerationPlan) generateChannel() {
	tmpDir := fmt.Sprintf("%s/.omgdtmp", cp.OutputDir)
	staticDir := "static/channel"

	// TODO: There should probably be a check to make sure the
	// user passes in a snake_case_channel_name meaning no caps
	// and no spaces just underscores, shouldn't be too hard
	snakeChannel := cp.Target
	camelChannel := StrToCamel(cp.Target)

	sccp := StaticCodeCopyPlan{
		filePathAlterations: []StaticCodeFilePathAlteration{
			{
				filePathToRead:  fmt.Sprintf("%s/game/Autoloads/ChannelManager.gd.omgdtpl", staticDir),
				filePathToWrite: fmt.Sprintf("%s/game/Autoloads/%sManager.gd.omgdtpl", tmpDir, camelChannel),
			},
			{
				filePathToRead:  fmt.Sprintf("%s/game/Autoloads/ChannelEvent.gd.tmpl.omgdtpl", staticDir),
				filePathToWrite: fmt.Sprintf("%s/game/Autoloads/%sEvent.gd.tmpl.omgdtpl", tmpDir, camelChannel),
			},
			{
				filePathToRead:  fmt.Sprintf("%s/resources/channel_events.yml.omgdtpl", staticDir),
				filePathToWrite: fmt.Sprintf("%s/resources/%s_events.yml.omgdtpl", tmpDir, snakeChannel),
			},
			{
				filePathToRead:  fmt.Sprintf("%s/game/RootScenes/ChannelMUX.tscn.omgdtpl", staticDir),
				filePathToWrite: fmt.Sprintf("%s/game/RootScenes/%sMUX.tscn.omgdtpl", tmpDir, camelChannel),
			},
			{
				filePathToRead:  fmt.Sprintf("%s/game/RootScenes/ChannelMUXController.gd.omgdtpl", staticDir),
				filePathToWrite: fmt.Sprintf("%s/game/RootScenes/%sMUXController.gd.omgdtpl", tmpDir, camelChannel),
			},
			{
				filePathToRead:  fmt.Sprintf("%s/server/nakama/data/modules/channel.lua.omgdtpl", staticDir),
				filePathToWrite: fmt.Sprintf("%s/server/nakama/data/modules/%s.lua.omgdtpl", tmpDir, snakeChannel),
			},
			{
				filePathToRead:  fmt.Sprintf("%s/server/nakama/data/modules/channel_manager.lua.omgdtpl", staticDir),
				filePathToWrite: fmt.Sprintf("%s/server/nakama/data/modules/%s_manager.lua.omgdtpl", tmpDir, snakeChannel),
			},
		},
	}

	err := sccp.CopyStaticDirectory("static/channel", tmpDir)
	if err != nil {
		LogFatal(fmt.Sprint(err))
	}

	newProfile := GetProfile(fmt.Sprintf("%s/profiles/local.yml", tmpDir))
	newProfile.UpdateProfile("omgd.channel_name", cp.Target)

	if cp.Args != "" {
		newProfile.UpdateProfile("omgd.channel_events", strings.Split(cp.Args, " "))
	}

	BuildTemplatesFromPath(
		newProfile,
		tmpDir,
		"omgdtpl",
		true,
	)
}

// cleanup code
func (cp *CodeGenerationPlan) Cleanup() {
	switch cp.Plan {
	// case "new":
	case "example-partial-2d-player-movement":
		cp.cleanupExample2DPlayerMovement()
	case "channel":
		cp.cleanupChannel()
	}
}

// cleans up example 2d player movement code
func (cp *CodeGenerationPlan) cleanupExample2DPlayerMovement() {
	tmpDir := fmt.Sprintf("%s/.omgdtmp", cp.OutputDir)

	err := os.RemoveAll(fmt.Sprintf("%s/profiles", tmpDir))
	if err != nil {
		LogFatal(fmt.Sprint(err))
	}

	sccp := StaticCodeCopyPlan{}

	err = sccp.CopyStaticDirectory(tmpDir, cp.OutputDir)
	if err != nil {
		LogFatal(fmt.Sprint(err))
	}

	err = os.RemoveAll(tmpDir)
	if err != nil {
		LogFatal(fmt.Sprint(err))
	}
}

// cleans up channel code
func (cp *CodeGenerationPlan) cleanupChannel() {
	tmpDir := fmt.Sprintf("%s/.omgdtmp", cp.OutputDir)

	err := os.RemoveAll(fmt.Sprintf("%s/profiles", tmpDir))
	if err != nil {
		LogFatal(fmt.Sprint(err))
	}

	sccp := StaticCodeCopyPlan{}

	err = sccp.CopyStaticDirectory(tmpDir, cp.OutputDir)

	if err != nil {
		LogFatal(fmt.Sprint(err))
	}

	err = os.RemoveAll(tmpDir)
	if err != nil {
		LogFatal(fmt.Sprint(err))
	}
}

// cleans up and resets omgdtmp dir
func (cp *CodeGenerationPlan) resetOMGDTmpDir() {
	tmpDir := fmt.Sprintf("%s/.omgdtmp", cp.OutputDir)

	err := os.RemoveAll(tmpDir)
	if err != nil {
		LogFatal(fmt.Sprint(err))
	}

	err = os.Mkdir(tmpDir, 0755)
	if err != nil {
		LogFatal(fmt.Sprint(err))
	}
}
