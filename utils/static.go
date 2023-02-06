package utils

import (
	"embed"
	"fmt"
	"log"
	"os"
	"strings"
)

//go:embed static/*
var staticFiles embed.FS

const FILE_WRITE_PERMS = 0755

type StaticCodeFilePathAlteration struct {
	filePathToRead          string
	filePathToWrite         string
	stringToReadForReplace  string
	stringToWriteForReplace string
}

type StaticCodeCopyPlan struct {
	filePathAlterations []StaticCodeFilePathAlteration
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

	err = os.Mkdir(pathToCopyTo, FILE_WRITE_PERMS)
	if err != nil && !os.IsExist(err) {
		return err
	}

	for i := 0; i < len(files); i++ {
		file := files[i]
		filePathToRead := fmt.Sprintf("%s/%s", pathToCopy, file.Name())
		filePathToWrite := fmt.Sprintf("%s/%s", pathToCopyTo, file.Name())

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

			err = sccp.CopyStaticFile(filePathToRead, filePathToWrite)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (sccp *StaticCodeCopyPlan) CopyStaticFile(
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

	adjustmentIdx := sccp.doesFilePathNeedChanging(filePathToRead)
	if adjustmentIdx > -1 {
		adjustment := sccp.filePathAlterations[adjustmentIdx]
		fileBytes = sccp.performFileAdjustments(fileBytes, adjustment)
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
			log.Printf("Attempting to overwrite file at %s, skipping\n", filePathToWrite)
		} else {
			return err
		}
	}

	return nil
}

func (sccp *StaticCodeCopyPlan) performFileAdjustments(
	fileBytes []byte,
	alteration StaticCodeFilePathAlteration,
) []byte {
	newFileString := ""

	if alteration.stringToReadForReplace != "" && alteration.stringToWriteForReplace != "" {
		newFileString = strings.ReplaceAll(
			string(fileBytes),
			alteration.stringToReadForReplace,
			alteration.stringToWriteForReplace,
		)
	}

	return []byte(newFileString)
}

func (sccp *StaticCodeCopyPlan) doesFilePathNeedChanging(filePath string) int {
	for i := 0; i < len(sccp.filePathAlterations); i++ {
		if filePath == sccp.filePathAlterations[i].filePathToRead {
			return i
		}
	}

	return -1
}
