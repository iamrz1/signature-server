package util

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"time"
)

var logger = zerolog.New(zerolog.ConsoleWriter{
	Out:        os.Stdout,
	TimeFormat: time.RFC3339,
	FormatTimestamp: func(i interface{}) string {
		return fmt.Sprintf("%v", i)
	},
	FormatLevel: func(i interface{}) string {
		return fmt.Sprintf("[%-5v]", i)
	},
	FormatFieldName: func(i interface{}) string {
		return fmt.Sprintf("%v:", i)
	},
	FormatFieldValue: func(i interface{}) string {
		return fmt.Sprintf("%v", i)
	},
}).With().Timestamp().Logger().Hook(ridHook{})

// Info prints a new message with info level
func Info(msg string) {
	logger.Info().Msg(msg)
}

// Infof prints a new formated message with info level
func Infof(format string, v ...interface{}) {
	logger.Info().Msgf(format, v...)
}

// Debug prints a new message with debug level
func Debug(msg string) {
	logger.Debug().Msg(msg)
}

// Debugf prints a new formatted message with debug level
func Debugf(format string, v ...interface{}) {
	logger.Debug().Msgf(format, v...)
}

// Warn prints a new message with warn level
func Warn(msg string) {
	logger.Warn().Msg(msg)
}

// Warnf prints a new formatted message with warn level
func Warnf(format string, v ...interface{}) {
	logger.Warn().Msgf(format, v...)
}

// Error prints a new message with error level
func Error(msg string) {
	logger.Error().Msg(msg)
}

// Errorf prints a new formatted message with error level
func Errorf(format string, v ...interface{}) {
	logger.Error().Msgf(format, v...)
}

// Panic prints a new message with panic level and throws a panic
func Panic(msg string) {
	logger.Panic().Msg(msg)
}

// Panicf prints a new formatted message with panic level and throws a panic
func Panicf(format string, v ...interface{}) {
	logger.Panic().Msgf(format, v...)
}

// Fatal prints a new message with fatal level and kills the system
func Fatal(msg string) {
	logger.Fatal().Msg(msg)
}

// Fatalf prints a new formatted message with fatal level and kills the system
func Fatalf(format string, v ...interface{}) {
	logger.Fatal().Msgf(format, v...)
}
