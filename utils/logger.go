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

// Logger.logFn used internally for logging, SetLogFn adjusts for tests
var logFn = log.Println

// Sets Logger.logFn for tests
func SetLogFn(_logFn func(v ...any)) {
	logFn = _logFn
}

// Logger.printFn used internally for logging, SetPrintFn adjusts for tests
var printFn = fmt.Println

// Sets Logger.printFn for tests
func SetPrintFn(_printFn func(v ...any) (n int, err error)) {
	printFn = _printFn
}

// Contains log level according to environment var, adjusted for tests
var envLogLevel LogLevel

// Sets environment log level variable, used for tests
func SetEnvLogLevel(logLevel LogLevel) {
	envLogLevel = logLevel
}

// Gets log level being used
func GetEnvLogLevel() LogLevel {
	return envLogLevel
}

var envLogOutputWithTimestamp bool

func init() {
	switch os.Getenv("OMGD_LOG_LEVEL") {
	case "FATAL":
		envLogLevel = FATAL_LOG
		break
	case "ERROR":
		envLogLevel = ERROR_LOG
		break
	case "WARN":
		envLogLevel = WARN_LOG
		break
	case "INFO":
		envLogLevel = INFO_LOG
		break
	case "DEBUG":
		envLogLevel = DEBUG_LOG
		break
	case "TRACE":
		envLogLevel = TRACE_LOG
		break
	default:
		envLogLevel = WARN_LOG
		break
	}

	timestampVal, timestampValSet := os.LookupEnv("OMGD_LOG_WITH_TIMESTAMP")

	if timestampValSet && timestampVal != "false" {
		envLogOutputWithTimestamp = true
	} else {
		envLogOutputWithTimestamp = false
	}
}

func LogFatal(message string) {
	if envLogLevel >= DEBUG_LOG {
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
	if logLevel <= envLogLevel {
		outputLogMessage(message)
	}
}

func outputLogMessage(message string) {
	if envLogOutputWithTimestamp {
		logFn(message)
	} else {
		_, err := printFn(message)

		if err != nil {
			debug.PrintStack()
			log.Fatalln("Fatal error with logging, could not print passed statement")
		}
	}
}
