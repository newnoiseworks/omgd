package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"runtime/debug"
	"testing"
)

func testLogComparison(expected interface{}, received interface{}) {
	LogWarn(fmt.Sprintf("received %s", received))
	LogWarn(fmt.Sprintf("expected %s", expected))
}

// tests file exists and contains a string
func testForFileAndRegexpMatch(t *testing.T, filePath string, search string) {
	// checks to make sure templates were created and properly named
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		debug.PrintStack()
		log.Println(fmt.Sprintf("Cannot find file: %s\n", err))
		t.Fail()
	}

	// makes sure templates were adjusted with proper variables
	matches, err := regexp.Match(search, file)
	if err != nil {
		debug.PrintStack()
		t.Fatal(err)
	}
	if !matches {
		debug.PrintStack()
		// log.Println("regexp for %s didn't match in file %s with contents: \n %s \n", search, filePath, file)
		log.Println(fmt.Sprintf("regexp for %s didn't match in file %s", search, filePath))
		t.Fail()
	}
}

// tests for file not existing
func testFileShouldNotExist(t *testing.T, filePath string) {
	_, err := os.Stat(filePath)
	if !os.IsNotExist(err) {
		debug.PrintStack()
		log.Println(fmt.Sprintf("File exists but should have been cleaned up at %s\n %s", filePath, err))
		t.Fail()
	}
}

// tests for file existence
func testFileShouldExist(t *testing.T, filePath string) {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		debug.PrintStack()
		log.Println(fmt.Sprintf("File does not exist but should have been created up at %s\n %s", filePath, err))
		t.Fail()
	}
}
