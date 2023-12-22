package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"testing"
)

func getNewLine() string {
	newline := "\n"

	if runtime.GOOS == "windows" {
		newline = "\r\n"
	}

	return newline
}

func TestStaticGetStaticFileCmd(t *testing.T) {
	// 1. Test for reading a simple one line file
	received, err := GetStaticFile(filepath.Join("static", "test", "test.md"))

	if err != nil {
		LogFatal(fmt.Sprint(err))
	}

	expected := "This is a test test test" + getNewLine()

	if expected != received {
		t.Errorf("File read from static lib doesn't match")

		testLogComparison(expected, received)
	}
}

func TestStaticCopyStaticDirectoryCmd(t *testing.T) {
	t.Cleanup(func() {
		err := os.RemoveAll(filepath.Join("static", "test", "test_dir_post_copying"))

		if err != nil {
			LogDebug(fmt.Sprint(err))
			t.Fail()
		}
	})

	// 1. copy static/test/test_dir_to_copy to static/test/test_dir_post_copying
	sccPlan := StaticCodeCopyPlan{}

	err := sccPlan.CopyStaticDirectory(
		filepath.Join("static", "test", "test_dir_to_copy"),
		filepath.Join("static", "test", "test_dir_post_copying"),
	)
	if err != nil {
		LogDebug(fmt.Sprint(err))
		t.Fail()
	}

	// 2. validate files match
	file, err := os.ReadFile(filepath.Join("static", "test", "test_dir_post_copying", "test_one.md"))
	if err != nil {
		LogDebug(fmt.Sprint(err))
		t.Fail()
	}

	fileTwo, err := os.ReadFile(filepath.Join("static", "test", "test_dir_post_copying", "test_two.md"))
	if err != nil {
		LogDebug(fmt.Sprint(err))
		t.Fail()
	}

	fileThree, err := os.ReadFile(filepath.Join("static", "test", "test_dir_post_copying", "folder", "test_one.md"))
	if err != nil {
		LogDebug(fmt.Sprint(err))
		t.Fail()
	}

	expected := "test_one" + getNewLine()
	received := string(file)

	if expected != received {
		LogDebug(fmt.Sprintf("File %s doesn't match expected contents", filepath.Join("static", "test", "test_dir_post_copying", "test_one.md")))
		t.Fail()

		testLogComparison(expected, received)
	}

	expected = "test_two" + getNewLine()
	received = string(fileTwo)

	if expected != received {
		LogDebug(fmt.Sprintf("File %s doesn't match expected contents", filepath.Join("static", "test", "test_dir_post_copying", "test_two.md")))
		t.Fail()

		testLogComparison(expected, received)
	}

	expected = "test_one" + getNewLine()
	received = string(fileThree)

	if expected != received {
		LogDebug(fmt.Sprintf("File %s doesn't match expected contents", filepath.Join("static", "test", "test_dir_post_copying", "folder", "test_one.md")))
		t.Fail()

		testLogComparison(expected, received)
	}
}

func TestStaticCopyStaticFileWithChangedPath(t *testing.T) {
	t.Cleanup(func() {
		err := os.RemoveAll(filepath.Join("static", "test", ".omgdtmp"))

		if err != nil {
			LogDebug(fmt.Sprint(err))
			t.Fail()
		}
	})

	err := os.Mkdir(filepath.Join("static", "test", ".omgdtmp"), 0755)
	if err != nil && !os.IsExist(err) {
		LogDebug(fmt.Sprint(err))
		t.Fail()
	}

	CopyStaticFile(filepath.Join("static", "test", "test.md"), filepath.Join("static", "test", ".omgdtmp", "test22.md"))

	testForFileAndRegexpMatch(t, filepath.Join("static", "test", ".omgdtmp", "test22.md"), `This is a test test test`)
}

func TestStaticCopyStaticDirectoryWithEdits(t *testing.T) {
	t.Cleanup(func() {
		err := os.RemoveAll(filepath.Join("static", "test", ".omgdtmp"))

		if err != nil {
			LogDebug(fmt.Sprint(err))
			t.Fail()
		}
	})

	err := os.Mkdir(filepath.Join("static", "test", ".omgdtmp"), 0755)
	if err != nil && !os.IsExist(err) {
		LogDebug(fmt.Sprint(err))
		t.Fail()
	}

	scpp := StaticCodeCopyPlan{
		filePathAlterations: []StaticCodeFilePathAlteration{{
			filePathToRead:  filepath.Join("static", "test", "test_dir_to_copy", "test_three.md"),
			filePathToWrite: filepath.Join("static", "test", ".omgdtmp", "test_dir_post_copying", "test_trifecta.md"),
		}},
	}

	scpp.CopyStaticDirectory(
		filepath.Join("static", "test", "test_dir_to_copy"),
		filepath.Join("static", "test", ".omgdtmp", "test_dir_post_copying"),
	)

	testForFileAndRegexpMatch(t, filepath.Join("static", "test", ".omgdtmp", "test_dir_post_copying", "test_trifecta.md"), `is everything`)

	file, err := os.ReadFile(filepath.Join("static", "test", ".omgdtmp", "test_dir_post_copying", "test_one.md"))
	if err != nil {
		LogDebug(fmt.Sprint(err))
		t.Fail()
	}

	fileThree, err := os.ReadFile(filepath.Join("static", "test", ".omgdtmp", "test_dir_post_copying", "folder", "test_one.md"))
	if err != nil {
		LogDebug(fmt.Sprint(err))
		t.Fail()
	}

	expected := "test_one" + getNewLine()
	received := string(file)

	if expected != received {
		LogDebug(
			fmt.Sprintf(
				"File static %s doesn't match expected contents",
				filepath.Join("test", ".omgdtmp", "test_dir_post_copying", "test_one.md"),
			),
		)
		t.Fail()

		testLogComparison(expected, received)
	}

	expected = "test_one" + getNewLine()
	received = string(fileThree)

	if expected != received {
		LogDebug(
			fmt.Sprintf(
				"File static %s doesn't match expected contents",
				filepath.Join("test", ".omgdtmp", "test_dir_post_copying", "folder", "test_one.md"),
			),
		)
		t.Fail()

		testLogComparison(expected, received)
	}
}

func TestStaticCopyStaticDirectoryCreatesHiddenFiles(t *testing.T) {
	t.Cleanup(func() {
		err := os.RemoveAll(filepath.Join("static", "test", ".omgdtmp"))

		if err != nil {
			LogDebug(fmt.Sprint(err))
			t.Fail()
		}
	})

	err := os.Mkdir(filepath.Join("static", "test", ".omgdtmp"), 0755)
	if err != nil && !os.IsExist(err) {
		LogDebug(fmt.Sprint(err))
		t.Fail()
	}

	scpp := StaticCodeCopyPlan{}

	scpp.CopyStaticDirectory(
		filepath.Join("static", "test", "test_dir_to_copy"),
		filepath.Join("static", "test", ".omgdtmp", "test_dir_post_copying"),
	)

	testForFileAndRegexpMatch(t,
		filepath.Join("static", "test", ".omgdtmp", "test_dir_post_copying", ".hiddenfile"),
		`test hidden file`,
	)
}

func TestStaticCopyStaticDirectorySkipsFiles(t *testing.T) {
	t.Cleanup(func() {
		err := os.RemoveAll(filepath.Join("static", "test", ".omgdtmp"))

		if err != nil {
			LogDebug(fmt.Sprint(err))
			t.Fail()
		}
	})

	err := os.Mkdir(filepath.Join("static", "test", ".omgdtmp"), 0755)
	if err != nil && !os.IsExist(err) {
		LogDebug(fmt.Sprint(err))
		t.Fail()
	}

	scpp := StaticCodeCopyPlan{
		skipPaths: []string{
			filepath.Join("static", "test", "test_dir_to_copy", "test_three.md"),
			filepath.Join("static", "test", "test_dir_to_copy", "folder"),
		},
	}

	scpp.CopyStaticDirectory(
		filepath.Join("static", "test", "test_dir_to_copy"),
		filepath.Join("static", "test", ".omgdtmp", "test_dir_post_copying"),
	)

	testFileShouldNotExist(t, filepath.Join("static", "test", ".omgdtmp", "test_dir_post_copying", "test_three.md"))
	testFileShouldNotExist(t, filepath.Join("static", "test", ".omgdtmp", "test_dir_post_copying", "folder"))
}

func TestStaticCopyStaticDirectoryInSameDirectory(t *testing.T) {
	t.Cleanup(func() {
		err := os.RemoveAll(filepath.Join("static", "test", "test_dir_to_copy", ".omgdtmp"))

		if err != nil {
			LogDebug(fmt.Sprint(err))
			t.Fail()
		}
	})

	scpp := StaticCodeCopyPlan{
		skipPaths: []string{
			filepath.Join("static", "test", "test_dir_to_copy", ".omgdtmp"),
		},
	}

	scpp.CopyStaticDirectory(
		filepath.Join("static", "test", "test_dir_to_copy"),
		filepath.Join("static", "test", "test_dir_to_copy", ".omgdtmp"),
	)

	testFileShouldExist(t, filepath.Join("static", "test", "test_dir_to_copy", ".omgdtmp"))
	testFileShouldExist(t, filepath.Join("static", "test", "test_dir_to_copy", ".omgdtmp", "test_three.md"))
	testFileShouldExist(t, filepath.Join("static", "test", "test_dir_to_copy", ".omgdtmp", "folder"))
}

type testCopyStaticDirectoryResponse struct {
	pathToCopy   string
	pathToCopyTo string
}

var testCopyStaticDirectoryResponses = []testCopyStaticDirectoryResponse{}
var testCopyStaticDirectoryValidResponseSet = []testCopyStaticDirectoryResponse{}

var testCopyStaticDirectory = func(pathToCopy string, pathToCopyTo string) error {
	testCopyStaticDirectoryResponses = append(testCopyStaticDirectoryResponses, testCopyStaticDirectoryResponse{
		pathToCopy:   pathToCopy,
		pathToCopyTo: pathToCopyTo,
	})

	return nil
}

func testCopyStaticDirectoryValidCmdSet(t *testing.T, method string) {
	if len(testCopyStaticDirectoryValidResponseSet) != len(testCopyStaticDirectoryResponses) {
		t.Errorf("%s didn't create enough commands", method)
		testLogComparison(strconv.Itoa(len(testCopyStaticDirectoryValidResponseSet)), strconv.Itoa(len(testCopyStaticDirectoryResponses)))
	}

	for i := range testCopyStaticDirectoryValidResponseSet {
		if !reflect.DeepEqual(testCopyStaticDirectoryValidResponseSet[i], testCopyStaticDirectoryResponses[i]) {
			t.Errorf("%s failed on step %s", method, strconv.Itoa(i))
			testLogComparison(testCopyStaticDirectoryValidResponseSet[i], testCopyStaticDirectoryResponses[i])
		}
	}

	testCopyStaticDirectoryValidResponseSet = nil
}
