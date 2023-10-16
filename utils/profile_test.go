package utils

import (
	"testing"
)

func TestGetProfile(t *testing.T) {
	testDir := "static/test/infra_test_dir"

	profile := GetProfileFromDir("profiles/staging.yml", testDir)

	if profile.Name != "staging" {
		t.Fatalf("profile Name not properly formatted")
	}

	if profile.path != "profiles/staging.yml" {
		t.Fatalf("profile path not properly set")
	}

	if profile.Get("omgd.name") != "top-level-name" {
		t.Fatalf("Profile not inheriting from top level omgd.yml profile")
	}

	if profile.Get("omgd.override") != "overriden" {
		t.Fatalf("Profile not overriding properly")
	}
}
