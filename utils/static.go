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
		log.Fatal(err)
	}

	err = os.Mkdir(pathToCopyTo, FILE_WRITE_PERMS)

	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(files); i++ {
		file := files[i]

		if file.IsDir() {
			CopyStaticDirectory(
				fmt.Sprintf("%s/%s", pathToCopy, file.Name()),
				fmt.Sprintf("%s/%s", pathToCopyTo, file.Name()),
			)
		} else {
			fileBytes, err := staticFiles.ReadFile(
				fmt.Sprintf("%s/%s", pathToCopy, file.Name()),
			)

			if err != nil {
				// NOTE: Not catching error on purpose to pass above
				os.RemoveAll(pathToCopyTo)
				return err
			}

			os.WriteFile(
				fmt.Sprintf("%s/%s", pathToCopyTo, file.Name()),
				fileBytes,
				FILE_WRITE_PERMS,
			)
		}
	}

	return nil
}
