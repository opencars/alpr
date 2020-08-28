package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

type zeroLogLogger struct {
	logger zerolog.Logger
}

func zeroLogLevel(level LogLevel) zerolog.Level {
	switch level {
	case Info:
		return zerolog.InfoLevel
	case Debug:
		return zerolog.DebugLevel
	case Error:
		return zerolog.ErrorLevel
	case Fatal:
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}

func newZeroLogger(logLevel LogLevel, dev bool) Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	level := zeroLogLevel(logLevel)
	logger := zerolog.New(os.Stdout).Level(level).With().Timestamp().Logger()

	if dev {
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		logger = zerolog.New(os.Stdout).Level(level).Output(output).With().Timestamp().Logger()
	}

	return &zeroLogLogger{
		logger: logger,
	}
}

// Debugf logs message with DEBUG log level.
func (l *zeroLogLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debug().Msgf(format, args...)
}

// Infof logs message with INFO log level.
func (l *zeroLogLogger) Infof(format string, args ...interface{}) {
	l.logger.Info().Msgf(format, args...)
}

// Errorf logs message with ERROR log level.
func (l *zeroLogLogger) Errorf(format string, args ...interface{}) {
	l.logger.Error().Msgf(format, args...)
}

// Fatalf logs message with FATAL log level.
func (l *zeroLogLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatal().Msgf(format, args...)
}

// WithFields builds nested logger with specified fields.
func (l *zeroLogLogger) WithFields(fields Fields) Logger {
	return &zeroLogLogger{
		logger: l.logger.With().Fields(fields).Logger(),
	}
}
