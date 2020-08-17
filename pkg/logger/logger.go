package logger

import "os"

// A global variable so that log functions can be directly accessed
var Log Logger

func init() {
	Log = NewLogger(os.Stdout)
}

// Fields Type to pass when we want to call WithFields for structured logging
type Fields map[string]interface{}

// Logger is our contract for the logger
type Logger interface {
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Fatal(args ...interface{})
	Panic(format string, args ...interface{})
	WithFields(keyValues Fields) Logger
}

func Debug(format string, args ...interface{}) {
	Log.Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	Log.Info(format, args...)
}

func Warn(format string, args ...interface{}) {
	Log.Warn(format, args...)
}

func Error(format string, args ...interface{}) {
	Log.Error(format, args...)
}

func Fatal(args ...interface{}) {
	Log.Fatal(args...)
}

func Panic(format string, args ...interface{}) {
	Log.Panic(format, args...)
}

func WithFields(keyValues Fields) Logger {
	return Log.WithFields(keyValues)
}
