package logger

import (
	"fmt"
	"io"
	"log"
)

var DebugMode bool

type Level string

const (
	LevelDebug Level = "DEBUG"
	LevelWarn  Level = "WARN "
	LevelError Level = "ERROR"
)

func SetOutput(writer io.Writer) {
	log.SetOutput(writer)
}

func Debug(format string, args ...interface{}) {
	if DebugMode {
		write(LevelDebug, format, args...)
	}
}

func Warn(format string, args ...interface{}) {
	write(LevelWarn, format, args...)
}

func Error(format string, args ...interface{}) {
	write(LevelError, format, args...)
}

func write(level Level, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)

	log.Printf("%s %s", level, message)
}
