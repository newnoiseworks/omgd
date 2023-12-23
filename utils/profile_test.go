package utils

import (
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
