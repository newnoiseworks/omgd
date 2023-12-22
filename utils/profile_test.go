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
		t.Fatalf("profile Name not properly formatted")
	}

	if profile.path != filepath.Join("profiles", "staging.yml") {
		t.Fatalf("profile path not properly set")
	}

	if profile.Get("omgd.name") != "top-level-name" {
		t.Fatalf("Profile not inheriting from top level omgd.yml profile")
	}

	if profile.Get("omgd.override") != "overriden" {
		t.Fatalf("Profile not overriding properly")
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
