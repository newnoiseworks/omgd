package utils

import (
	"fmt"
	"log"
	"os"
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

// Sets Logger._logFn for tests
func SetLogFn(logFn func(v ...any)) {
	_logFn = logFn
}

// Logger._printFn used internally for logging, SetPrintFn adjusts for tests
var _printFn = fmt.Println

// Sets Logger._printFn for tests
func SetPrintFn(printFn func(v ...any) (n int, err error)) {
	_printFn = printFn
}

// Contains log level according to environment var, adjusted for tests
var _envLogLevel LogLevel

// Sets environment log level variable, used for tests
func SetEnvLogLevel(logLevel LogLevel) {
	_envLogLevel = logLevel
}

var _envLogOutputWithTimestamp bool

func init() {
	switch os.Getenv("OMGD_LOG_LEVEL") {
	case "FATAL":
		_envLogLevel = FATAL_LOG
		break
	case "ERROR":
		_envLogLevel = ERROR_LOG
		break
	case "WARN":
		_envLogLevel = WARN_LOG
		break
	case "INFO":
		_envLogLevel = INFO_LOG
		break
	case "DEBUG":
		_envLogLevel = DEBUG_LOG
		break
	case "TRACE":
		_envLogLevel = TRACE_LOG
		break
	default:
		_envLogLevel = WARN_LOG
		break
	}

	timestampVal, timestampValSet := os.LookupEnv("OMGD_LOG_WITH_TIMESTAMP")

	if timestampValSet && timestampVal != "false" {
		_envLogOutputWithTimestamp = true
	} else {
		_envLogOutputWithTimestamp = false
	}
}

func LogFatal(message string) {
	if _envLogLevel >= DEBUG_LOG {
		debug.PrintStack()
	}

	processMessageAgainstLogLevel(message, FATAL_LOG)
	os.Exit(1)
}

func LogError(message string) {
	processMessageAgainstLogLevel(message, ERROR_LOG)
}

func LogWarn(message string) {
	processMessageAgainstLogLevel(message, WARN_LOG)
}

func LogInfo(message string) {
	processMessageAgainstLogLevel(message, INFO_LOG)
}

func LogDebug(message string) {
	processMessageAgainstLogLevel(message, DEBUG_LOG)
}

func LogTrace(message string) {
	processMessageAgainstLogLevel(message, TRACE_LOG)
}

func processMessageAgainstLogLevel(message string, logLevel LogLevel) {
	if logLevel <= _envLogLevel {
		outputLogMessage(message)
	}
}

func outputLogMessage(message string) {
	if _envLogOutputWithTimestamp {
		_logFn(message)
	} else {
		_, err := _printFn(message)

		if err != nil {
			debug.PrintStack()
			log.Fatalln("Fatal error with logging, could not print passed statement")
		}
	}
}
