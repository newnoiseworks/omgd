package utils

import (
	"fmt"
	"log"
	"runtime/debug"
)

type LogLevel int

const (
	FATAL_LOG LogLevel = iota
	ERROR_LOG LogLevel = iota
	WARN_LOG  LogLevel = iota
	INFO_LOG  LogLevel = iota
	DEBUG_LOG LogLevel = iota
	TRACE_LOG LogLevel = iota
)

// Logger._logFn used internally for logging, SetLogFn adjusts for tests
var _logFn = log.Println

func SetLogFn(logFn func(v ...any)) {
	_logFn = logFn
}

// Logger._printFn used internally for logging, SetPrintFn adjusts for tests
var _printFn = fmt.Println

func SetPrintFn(printFn func(v ...any) (n int, err error)) {
	_printFn = printFn
}

var envLogLevel = FATAL_LOG
var envLogOutputWithTimestamp = false

func TraceLog(message string) {
	processMessageAgainstLogLevel(message, TRACE_LOG)
}

func processMessageAgainstLogLevel(message string, logLevel LogLevel) {
	if envLogLevel <= logLevel {
		outputLogMessage(message)
	}
}

func outputLogMessage(message string) {
	if envLogOutputWithTimestamp {
		_logFn(message)
	} else {
		_, err := _printFn(message)

		if err != nil {
			debug.PrintStack()
			log.Fatalln("Fatal error with logging, could not print passed statement")
		}
	}
}
