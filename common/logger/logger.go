package logger

import (
	"fmt"
	"log"
	"os"
)

var logger = log.New(os.Stderr, "", log.Ldate|log.Lmicroseconds|log.Llongfile)

var (
	level      = LevelDebug
	LevelDebug = 1
	LevelInfo  = 2
	LevelWarn  = 3
	LevelError = 4
)

func SetLogLevel(l int) {
	level = l
}

func Debug(format string, args ...interface{}) {
	if level <= LevelDebug {
		printf("Debug", format, args...)
	}
}

func Info(format string, args ...interface{}) {
	if level <= LevelInfo {
		printf("Info", format, args...)
	}
}

func Warn(format string, args ...interface{}) {
	if level <= LevelWarn {
		printf("Warn", format, args...)
	}
}

func Error(format string, args ...interface{}) {
	if level <= LevelError {
		printf("Error", format, args...)
	}
}

func printf(level string, format string, args ...interface{}) {
	prefix := fmt.Sprintf("[%s] ", level)
	logger.Output(3, fmt.Sprintf(prefix+format, args...))
}
