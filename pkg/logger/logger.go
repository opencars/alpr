package logger

// A global variable so that log functions can be directly accessed.
var log Logger

func init() {
	NewLogger(Debug, true)
}

func Log() Logger {
	return log
}

// Fields to be passed when we want to call WithFields for structured logging.
type Fields map[string]interface{}

// LogLevel represent level of logging.
type LogLevel string

const (
	// Debug has verbose message.
	Debug LogLevel = "debug"
	// Info is default log level.
	Info LogLevel = "info"
	// Error is for logging errors.
	Error LogLevel = "errors"
	// Fatal is for logging fatal messages. The system shutdown after logging the message.
	Fatal LogLevel = "fatal"
)

// Logger is our contract for the logger.
type Logger interface {
	Debugf(format string, args ...interface{})

	Infof(format string, args ...interface{})

	Errorf(format string, args ...interface{})

	Fatalf(format string, args ...interface{})

	WithFields(keyValues Fields) Logger
}

// NewLogger returns implementation of Logger interface.
func NewLogger(level LogLevel, dev bool) {
	logger := newZeroLogger(level, dev)
	log = logger
}

// Debugf logs message with DEBUG log level.
func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

// Infof logs message with INFO log level.
func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

// Errorf logs message with ERROR log level.
func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

// Fatalf logs message with FATAL log level.
func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

// WithFields builds nested logger with specified fields.
func WithFields(keyValues Fields) Logger {
	return log.WithFields(keyValues)
}
