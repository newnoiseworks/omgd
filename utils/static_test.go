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

// 3. Test for copying a directory
func TestStaticCopyStaticDirectoryCmd(t *testing.T) {
	t.Cleanup(func() {
		err := os.RemoveAll("static/test/test_dir_post_copying")

		if err != nil {
			t.Fatal(err)
		}
	})

	// 1. copy static/test/test_dir_to_copy to static/test/test_dir_post_copying
	err := CopyStaticDirectory("static/test/test_dir_to_copy", "static/test/test_dir_post_copying")
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

// 4. Test for copying a file w/ a replaced string or two

// 5. Test for combining the above into one direction or command?
