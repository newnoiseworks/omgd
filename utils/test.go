package utils

import (
	"log"
)

func testLogComparison(expected interface{}, received interface{}) {
	log.Printf("received %s", received)
	log.Printf("expected %s", expected)
}
