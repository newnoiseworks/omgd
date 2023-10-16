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
}
