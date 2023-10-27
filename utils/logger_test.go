package utils

import (
	"fmt"
	"log"
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

func setupLoggingTests() {
	SetLogFn(testLogFn)
	SetPrintFn(testPrintFn)
}

func cleanupLoggingTests() {
	testPrintFnOutput = []string{}
	testLogFnOutput = []string{}
	SetEnvLogLevel(defaultLogLevel)
	SetLogFn(log.Println)
	SetPrintFn(fmt.Println)
}

func TestLoggerLogTrace(t *testing.T) {
	setupLoggingTests()
	t.Cleanup(cleanupLoggingTests)

	SetEnvLogLevel(TRACE_LOG)

	LogTrace("message")

	if testPrintFnOutput[0] != "message" {
		fmt.Println("failed test on trace level logging")
		t.Fail()
	}
}

func TestLoggerLogTraceDoesntFire(t *testing.T) {
	setupLoggingTests()

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
