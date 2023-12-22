package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuildTemplateFromPath(t *testing.T) {
	testDir := filepath.Join("static", "test", "build_templates_dir")

	t.Cleanup(func() {
		os.Remove(filepath.Join(testDir, "template_example.gd"))
	})

	profile := GetProfileFromDir("profile.yml", testDir)

	BuildTemplateFromPath(
		filepath.Join(testDir, "template_example.gd.tmpl"),
		profile,
		testDir,
		"tmpl",
		false,
	)

	testFileShouldExist(t, filepath.Join(testDir, "template_example.gd"))

	testForFileAndRegexpMatch(t, filepath.Join(testDir, "template_example.gd"), "127.6.6.6")
}
