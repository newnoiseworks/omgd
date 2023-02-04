package utils

import (
	"embed"
	"fmt"
	"log"
	"os"
)

//go:embed static/*
var staticFiles embed.FS

const FILE_WRITE_PERMS = 0755

func GetStaticFile(file string) (string, error) {
	byteArr, err := staticFiles.ReadFile(file)

	retVal := string(byteArr)

	return retVal, err
}

func CopyStaticDirectory(pathToCopy string, pathToCopyTo string) error {
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

		if file.IsDir() {
			err = CopyStaticDirectory(
				fmt.Sprintf("%s/%s", pathToCopy, file.Name()),
				fmt.Sprintf("%s/%s", pathToCopyTo, file.Name()),
			)
			if err != nil {
				return err
			}
		} else {
			filePathToRead := fmt.Sprintf("%s/%s", pathToCopy, file.Name())
			fileBytes, err := staticFiles.ReadFile(filePathToRead)

			if err != nil {
				fileBytes, err = os.ReadFile(filePathToRead)

				if err != nil {
					return err
				}
			}

			filePath := fmt.Sprintf("%s/%s", pathToCopyTo, file.Name())

			err = os.WriteFile(
				filePath,
				fileBytes,
				FILE_WRITE_PERMS,
			)

			if err != nil {
				if os.IsExist(err) {
					log.Printf("Attempting to overwrite file at %s, skipping\n", filePath)
				} else {
					return err
				}
			}
		}
	}

	return nil
}
