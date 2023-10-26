package utils

import (
	"fmt"
	"os"
	"testing"
)

var testLogFnOutput = []string{}
var testPrintFnOutput = []string{}

func testLogFn(v ...any) {
	testLogFnOutput = append(
		testLogFnOutput,
		fmt.Sprintf(v[0].(string), v[1:]...),
	)
}

func testPrintFn(v ...any) (n int, err error) {
	testPrintFnOutput = append(
		testPrintFnOutput,
		fmt.Sprintf(v[0].(string), v[1:]...),
	)

	return 0, nil
}

func TestMain(m *testing.M) {
	SetLogFn(testLogFn)
	SetPrintFn(testPrintFn)

	code := m.Run()

	os.Exit(code)
}

func TestLoggerTraceLog(t *testing.T) {
	TraceLog("message")

	if testPrintFnOutput[0] != "message" {
		fmt.Println("failed test on trace level logging")
		t.Fail()
	}

	testLogFnOutput = []string{}
}
