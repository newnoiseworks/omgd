package utils

import (
	"log"
	"os"
	"testing"
)

func TestStaticGetStaticFileCmd(t *testing.T) {
	// 1. Test for reading a simple one line file
	received, err := GetStaticFile("static/test/test.md")

	if err != nil {
		log.Fatal(err)
	}

	expected := "This is a test test test\n"

	if expected != received {
		t.Errorf("File read from static lib doesn't match")

		testLogComparison(expected, received)
	}
}

func TestStaticCopyStaticDirectoryCmd(t *testing.T) {
	t.Cleanup(func() {
		err := os.RemoveAll("static/test/test_dir_post_copying")

		if err != nil {
			t.Fatal(err)
		}
	})

	// 1. copy static/test/test_dir_to_copy to static/test/test_dir_post_copying
	sccPlan := StaticCodeCopyPlan{}

	err := sccPlan.CopyStaticDirectory("static/test/test_dir_to_copy", "static/test/test_dir_post_copying")
	if err != nil {
		t.Fatal(err)
	}

	// 2. validate files match
	file, err := os.ReadFile("static/test/test_dir_post_copying/test_one.md")
	if err != nil {
		t.Fatal(err)
	}

	fileTwo, err := os.ReadFile("static/test/test_dir_post_copying/test_two.md")
	if err != nil {
		t.Fatal(err)
	}

	fileThree, err := os.ReadFile("static/test/test_dir_post_copying/folder/test_one.md")
	if err != nil {
		t.Fatal(err)
	}

	expected := "test_one\n"
	received := string(file)

	if expected != received {
		t.Fatal("File static/test/test_dir_post_copying/test_one.md doesn't match expected contents")

		testLogComparison(expected, received)
	}

	expected = "test_two\n"
	received = string(fileTwo)

	if expected != received {
		t.Fatal("File static/test/test_dir_post_copying/test_two.md doesn't match expected contents")

		testLogComparison(expected, received)
	}

	expected = "test_one\n"
	received = string(fileThree)

	if expected != received {
		t.Fatal("File static/test/test_dir_post_copying/folder/test_one.md doesn't match expected contents")

		testLogComparison(expected, received)
	}
}

// func TestStaticCopyStaticFileWithChangedString(t *testing.T) {
// 	t.Cleanup(func() {
// 		err := os.RemoveAll("static/test/.omgdtmp")

// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 	})

// 	err := os.Mkdir("static/test/.omgdtmp", 0755)
// 	if err != nil && !os.IsExist(err) {
// 		t.Fatal(err)
// 	}

// 	scpp := StaticCodeCopyPlan{
// 		filePathAlterations: []StaticCodeFilePathAlteration{{
// 			filePathToRead:          "static/test/test.md",
// 			stringToReadForReplace:  "test",
// 			stringToWriteForReplace: "nothing",
// 		}},
// 	}

// 	scpp.CopyStaticFile("static/test/test.md", "static/test/.omgdtmp/test.md")

// 	testForFileAndRegexpMatch(t, "static/test/.omgdtmp/test.md", `This is a nothing nothing nothing`)
// }

func TestStaticCopyStaticFileWithChangedPath(t *testing.T) {
	t.Cleanup(func() {
		err := os.RemoveAll("static/test/.omgdtmp")

		if err != nil {
			t.Fatal(err)
		}
	})

	err := os.Mkdir("static/test/.omgdtmp", 0755)
	if err != nil && !os.IsExist(err) {
		t.Fatal(err)
	}

	CopyStaticFile("static/test/test.md", "static/test/.omgdtmp/test22.md")

	testForFileAndRegexpMatch(t, "static/test/.omgdtmp/test22.md", `This is a test test test`)
}

func TestStaticCopyStaticDirectoryWithEdits(t *testing.T) {
	t.Cleanup(func() {
		err := os.RemoveAll("static/test/.omgdtmp")

		if err != nil {
			t.Fatal(err)
		}
	})

	err := os.Mkdir("static/test/.omgdtmp", 0755)
	if err != nil && !os.IsExist(err) {
		t.Fatal(err)
	}

	scpp := StaticCodeCopyPlan{
		filePathAlterations: []StaticCodeFilePathAlteration{{
			filePathToRead:  "static/test/test_dir_to_copy/test_three.md",
			filePathToWrite: "static/test/.omgdtmp/test_dir_post_copying/test_trifecta.md",
			// stringToReadForReplace:  "test_three",
			// stringToWriteForReplace: "nothing",
		}},
	}

	scpp.CopyStaticDirectory("static/test/test_dir_to_copy", "static/test/.omgdtmp/test_dir_post_copying")

	testForFileAndRegexpMatch(t, "static/test/.omgdtmp/test_dir_post_copying/test_trifecta.md", `is everything`)

	// 2. validate other files moved over
	file, err := os.ReadFile("static/test/.omgdtmp/test_dir_post_copying/test_one.md")
	if err != nil {
		t.Fatal(err)
	}

	fileThree, err := os.ReadFile("static/test/.omgdtmp/test_dir_post_copying/folder/test_one.md")
	if err != nil {
		t.Fatal(err)
	}

	expected := "test_one\n"
	received := string(file)

	if expected != received {
		t.Fatal("File static/test/.omgdtmp/test_dir_post_copying/test_one.md doesn't match expected contents")

		testLogComparison(expected, received)
	}

	expected = "test_one\n"
	received = string(fileThree)

	if expected != received {
		t.Fatal("File static/test/.omgdtmp/test_dir_post_copying/folder/test_one.md doesn't match expected contents")

		testLogComparison(expected, received)
	}
}
