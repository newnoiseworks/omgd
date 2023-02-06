package utils

import (
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"testing"
)

func testLogComparison(expected interface{}, received interface{}) {
	log.Printf("received %s", received)
	log.Printf("expected %s", expected)
}

// tests file exists and contains a string
func testForFileAndRegexpMatch(t *testing.T, filePath string, search string) {
	// checks to make sure templates were created and properly named
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Cannot find file: %s\n", err)
	}

	// makes sure templates were adjusted with proper variables
	matches, err := regexp.Match(search, file)
	if err != nil {
		t.Fatal(err)
	}
	if !matches {
		t.Fatalf("regexp for %s didn't match in file %s", search, filePath)
	}
}

// tests for file not existing
func testFileShouldNotExist(t *testing.T, filePath string) {
	_, err := os.Stat(filePath)
	if !os.IsNotExist(err) {
		t.Fatalf("File exists but should have been cleaned up at %s\n %s", filePath, err)
	}
}

// tests for file existence
func testFileShouldExist(t *testing.T, filePath string) {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		t.Fatalf("File does not exist but should have been created up at %s\n %s", filePath, err)
	}
}
