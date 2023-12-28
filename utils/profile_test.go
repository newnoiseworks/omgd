package utils

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetProfile(t *testing.T) {
	testDir := filepath.Join("static", "test", "infra_test_dir")

	profile := GetProfileFromDir(filepath.Join("profiles", "staging.yml"), testDir)

	if profile.Name != "staging" {
		LogError("profile Name not properly formatted")
		t.Fail()
	}

	if profile.path != filepath.Join("profiles", "staging.yml") {
		LogError("profile path not properly set")
		t.Fail()
	}

	if profile.Get("omgd.name") != "top-level-name" {
		LogError("Profile not inheriting from top level omgd.yml profile")
		t.Fail()
	}

	if profile.Get("omgd.override") != "overriden" {
		LogError("Profile not overriding properly")
		t.Fail()
	}
}

func Test_setValueToKeyWithArray(t *testing.T) {
	testObj := map[interface{}]interface{}{}

	keys := strings.Split("this.is.a.test", ".")

	setValueToKeyWithArray(keys, 0, testObj, "success")

	if getValueToKeyWithArray(keys, 0, testObj) != "success" {
		LogError("Could not set a nested key onto a map that doesn't have those keys in the first place")
		t.Fail()
	}
}

func Test_validateProfileProjectNameLimit(t *testing.T) {
	setupLoggingTests()
	t.Cleanup(cleanupLoggingTests)

	testDir := filepath.Join("static", "test", "invalid_profiles")

	GetProfileFromDir("staging.yml", testDir)

	if testPrintFnOutput[0] != "OMGD project name exceeds 30 character limit" {
		fmt.Println("OMGD project name not limited to 30 character limit")
		t.Fail()
	}
}

func Test_validateProfileNameLimit(t *testing.T) {
	setupLoggingTests()
	t.Cleanup(cleanupLoggingTests)

	testDir := filepath.Join("static", "test", "invalid_profiles")

	GetProfileFromDir("twentyTwoCharacterLimit.yml", testDir)

	if testPrintFnOutput[0] != "OMGD profile filename exceeds 14 character limit (not counting .yml)" {
		fmt.Println("OMGD profile name not limited to 14 character limit")
		t.Fail()
	}
}
