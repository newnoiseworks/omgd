package utils

import (
	"embed"
)

//go:embed static/*
var StaticFiles embed.FS

func GetStaticFile(file string) (string, error) {
	byteArr, err := StaticFiles.ReadFile(file)

	retVal := string(byteArr)

	return retVal, err
}

func CopyStaticDirectory(pathToCopy string, pathToCopyTo string) error {
	return nil
}
