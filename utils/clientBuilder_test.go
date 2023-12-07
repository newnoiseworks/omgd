package utils

import (
	"fmt"
	"testing"
)

func TestClientBuilderBuildFromProfile(t *testing.T) {
	t.Cleanup(func() {
		testCmdOnDirResponses = []testCmdOnDirResponse{}
		testCopyStaticDirectoryResponses = []testCopyStaticDirectoryResponse{}
	})

	testDir := "static/test/client_builder_dir"

	profile := GetProfileFromDir("profiles/local.yml", testDir)

	cb := ClientBuilder{
		Profile:             profile,
		CmdOnDirWithEnv:     testCmdOnDirWithEnv,
		CopyStaticDirectory: testCopyStaticDirectory,
	}

	cb.Build()

	testCmdOnDirValidResponseSet = []testCmdOnDirResponse{
		{
			cmdStr:  "docker compose up build-windows build-mac build-web build-x11",
			cmdDesc: fmt.Sprintf("Building %s game clients into game/dist folder", profile.Name),
			cmdDir:  "game",
			env:     []string{fmt.Sprintf("BUILD_ENV=%s", profile.Name)},
		},
	}

	testCmdOnDirValidCmdSet(t, "ClientBuilder#Build")
}

func TestClientBuilderBuildFromArgs(t *testing.T) {
	t.Cleanup(func() {
		testCmdOnDirResponses = []testCmdOnDirResponse{}
		testCopyStaticDirectoryResponses = []testCopyStaticDirectoryResponse{}
	})

	testDir := "static/test/client_builder_dir"

	profile := GetProfileFromDir("profiles/local.yml", testDir)

	cb := ClientBuilder{
		Profile:             profile,
		CmdOnDirWithEnv:     testCmdOnDirWithEnv,
		Targets:             "build-mac build-x11",
		CopyStaticDirectory: testCopyStaticDirectory,
	}

	cb.Build()

	testCmdOnDirValidResponseSet = []testCmdOnDirResponse{
		{
			cmdStr:  "docker compose up build-mac build-x11",
			cmdDesc: fmt.Sprintf("Building %s game clients into game/dist folder", profile.Name),
			cmdDir:  "game",
			env:     []string{fmt.Sprintf("BUILD_ENV=%s", profile.Name)},
		},
	}

	testCmdOnDirValidCmdSet(t, "ClientBuilder#Build")
}

func TestClientBuilderBuildFromProfileWithOverrides(t *testing.T) {
	t.Cleanup(func() {
		testCmdOnDirResponses = []testCmdOnDirResponse{}
		testCopyStaticDirectoryResponses = []testCopyStaticDirectoryResponse{}
	})

	testDir := "static/test/client_builder_dir"

	profile := GetProfileFromDir("profiles/override.yml", testDir)

	cb := ClientBuilder{
		Profile:             profile,
		CmdOnDirWithEnv:     testCmdOnDirWithEnv,
		CopyStaticDirectory: testCopyStaticDirectory,
	}

	cb.Build()

	testCmdOnDirValidResponseSet = []testCmdOnDirResponse{
		{
			cmdStr:  "docker compose up build-windows build-web",
			cmdDesc: fmt.Sprintf("Building %s game clients into game/dist folder", profile.Name),
			cmdDir:  "game",
			env:     []string{fmt.Sprintf("BUILD_ENV=%s", profile.Name)},
		},
	}

	testCmdOnDirValidCmdSet(t, "ClientBuilder#Build")
}

func TestClientBuilderBuildCopiesFiles(t *testing.T) {
	t.Cleanup(func() {
		testCmdOnDirResponses = []testCmdOnDirResponse{}
		testCopyStaticDirectoryResponses = []testCopyStaticDirectoryResponse{}
	})

	testDir := "static/test/client_builder_dir"

	profile := GetProfileFromDir("profiles/local.yml", testDir)

	cb := ClientBuilder{
		Profile:             profile,
		CmdOnDirWithEnv:     testCmdOnDirWithEnv,
		CopyStaticDirectory: testCopyStaticDirectory,
	}

	cb.Build()

	testCmdOnDirValidResponseSet = []testCmdOnDirResponse{
		{
			cmdStr:  "docker compose up build-windows build-mac build-web build-x11",
			cmdDesc: fmt.Sprintf("Building %s game clients into game/dist folder", profile.Name),
			cmdDir:  "game",
			env:     []string{fmt.Sprintf("BUILD_ENV=%s", profile.Name)},
		},
	}

	testCmdOnDirValidCmdSet(t, "ClientBuilder#Build")

	testCopyStaticDirectoryValidResponseSet = []testCopyStaticDirectoryResponse{
		{
			pathToCopy:   "test/dir",
			pathToCopyTo: "test/dir_to",
		},
	}

	testCopyStaticDirectoryValidCmdSet(t, "ClientBuilder#Build on copy static directory")
}
