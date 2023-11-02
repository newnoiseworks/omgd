package utils

import (
	"embed"
	"fmt"
	"os"
	"strings"
)

//go:embed static/*
var staticFiles embed.FS

const FILE_WRITE_PERMS = 0755

type StaticCodeFilePathAlteration struct {
	filePathToRead  string
	filePathToWrite string
}

type StaticCodeCopyPlan struct {
	filePathAlterations []StaticCodeFilePathAlteration
	skipPaths           []string
}

func GetStaticFile(file string) (string, error) {
	byteArr, err := staticFiles.ReadFile(file)

	retVal := string(byteArr)

	return retVal, err
}

func (sccp *StaticCodeCopyPlan) CopyStaticDirectory(pathToCopy string, pathToCopyTo string) error {
	files, err := staticFiles.ReadDir(pathToCopy)
	if err != nil {
		files, err = os.ReadDir(pathToCopy)

		if err != nil {
			return err
		}
	}

	err = os.MkdirAll(pathToCopyTo, FILE_WRITE_PERMS)
	if err != nil && !os.IsExist(err) {
		return err
	}

	for i := 0; i < len(files); i++ {
		file := files[i]
		filePathToRead := fmt.Sprintf("%s/%s", pathToCopy, file.Name())
		filePathToWrite := fmt.Sprintf("%s/%s", pathToCopyTo, file.Name())

		if sccp.shouldSkipFilePath(filePathToRead) {
			continue
		}

		// looks like golang's fs.embed thing doesn't add hidden files in subdirectories
		// which is trifling if you ask me but whatever. the below is to get around that
		// by replacing the . with OMGD_DOT_FILE, which will be rewritten to a . into
		// userland by the below.
		filePathToWrite = strings.ReplaceAll(filePathToWrite, "/OMGD_DOT_FILE", "/.")

		if file.IsDir() {
			err = sccp.CopyStaticDirectory(filePathToRead, filePathToWrite)
			if err != nil {
				return err
			}
		} else {
			adjustmentIdx := sccp.doesFilePathNeedChanging(filePathToRead)
			if adjustmentIdx > -1 {
				alteration := sccp.filePathAlterations[adjustmentIdx]

				if alteration.filePathToWrite != "" {
					filePathToWrite = alteration.filePathToWrite
				}
			}

			err = CopyStaticFile(filePathToRead, filePathToWrite)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func CopyStaticFile(
	filePathToRead string,
	filePathToWrite string,
) error {
	fileBytes, err := staticFiles.ReadFile(filePathToRead)
	if err != nil {
		fileBytes, err = os.ReadFile(filePathToRead)

		if err != nil {
			return err
		}
	}

	err = writeFile(filePathToWrite, fileBytes)

	if err != nil {
		return err
	}

	return nil
}

func writeFile(filePathToWrite string, fileBytes []byte) error {
	err := os.WriteFile(filePathToWrite, fileBytes, FILE_WRITE_PERMS)

	if err != nil {
		if os.IsExist(err) {
			LogDebug(fmt.Sprintf("Attempting to overwrite file at %s, skipping\n", filePathToWrite))
		} else {
			return err
		}
	}

	return nil
}

func (sccp *StaticCodeCopyPlan) doesFilePathNeedChanging(filePath string) int {
	for i := 0; i < len(sccp.filePathAlterations); i++ {
		if filePath == sccp.filePathAlterations[i].filePathToRead {
			return i
		}
	}

	return -1
}

func (sccp *StaticCodeCopyPlan) shouldSkipFilePath(filePath string) bool {
	for i := 0; i < len(sccp.skipPaths); i++ {
		if filePath == sccp.skipPaths[i] {
			return true
		}
	}

	return false
}
