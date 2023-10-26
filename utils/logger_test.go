package utils

import (
	"fmt"
	"log"
	"os"
	"testing"
)

var testLogFnOutput = []string{}
var testPrintFnOutput = []string{}

func testLogFn(v ...any) {
	testLogFnOutput = append(
		testLogFnOutput,
		fmt.Sprintf(fmt.Sprint(v[0]), v[1:]...),
	)
}

func testPrintFn(v ...any) (n int, err error) {
	testPrintFnOutput = append(
		testPrintFnOutput,
		fmt.Sprintf(fmt.Sprint(v[0]), v[1:]...),
	)

	return 0, nil
}

var defaultLogLevel = envLogLevel

func TestMain(m *testing.M) {
	SetLogFn(testLogFn)
	SetPrintFn(testPrintFn)

	code := m.Run()

	SetEnvLogLevel(defaultLogLevel)
	SetLogFn(log.Println)
	SetPrintFn(fmt.Println)

	os.Exit(code)
}

func TestLoggerLogTrace(t *testing.T) {
	t.Cleanup(func() {
		testPrintFnOutput = []string{}
		SetEnvLogLevel(defaultLogLevel)
	})

	SetEnvLogLevel(TRACE_LOG)

	LogTrace("message")

	if testPrintFnOutput[0] != "message" {
		fmt.Println("failed test on trace level logging")
		t.Fail()
	}

}

func TestLoggerLogTraceDoesntFire(t *testing.T) {
	t.Cleanup(func() {
		testPrintFnOutput = []string{}
		SetEnvLogLevel(defaultLogLevel)
	})

	SetEnvLogLevel(WARN_LOG)

	LogTrace("messages")

	if len(testPrintFnOutput) > 0 {
		fmt.Println("trace log emitted when it should not have")
		t.Fail()
	}
}
