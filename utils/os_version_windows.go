package utils

import (
	"log"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func IsWindowsServer() bool {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		log.Fatal(err)
	}
	defer k.Close()

	pn, _, err := k.GetStringValue("ProductName")
	if err != nil {
		log.Fatal(err)
	}

	return strings.Contains(pn, "Server")
}
