package hcl

import (
	"fmt"
	"strings"

	"github.com/dihedron/rafter/logging"
	"github.com/hashicorp/go-hclog"
)

// Logger is the tpe warring an HCL logger.
type Logger struct {
	logger hclog.Logger
}

// NewLogger returns an instance of HCL logger wrapper
// that complies with the logging.Logger interface.
func NewLogger(logger hclog.Logger) *Logger {
	return &Logger{
		logger: logger,
	}
}

// Trace logs a message at LevelTrace level.
func (l *Logger) Trace(msg string, args ...interface{}) {
	if logging.GetLevel() <= logging.LevelTrace {
		message := l.format(msg, args...)
		l.logger.Trace(message)
	}
}

// Debug logs a message at LevelDebug level.
func (l *Logger) Debug(msg string, args ...interface{}) {
	if logging.GetLevel() <= logging.LevelDebug {
		message := l.format(msg, args...)
		l.logger.Debug(message)
	}
}

// Info logs a message at LevelInfo level.
func (l *Logger) Info(msg string, args ...interface{}) {
	if logging.GetLevel() <= logging.LevelInfo {
		message := l.format(msg, args...)
		l.logger.Info(message)
	}
}

// Warn logs a message at LevelWarn level.
func (l *Logger) Warn(msg string, args ...interface{}) {
	if logging.GetLevel() <= logging.LevelWarn {
		message := l.format(msg, args...)
		l.logger.Warn(message)
	}
}

// Error logs a message at LevelError level.
func (l *Logger) Error(msg string, args ...interface{}) {
	if logging.GetLevel() <= logging.LevelError {
		message := l.format(msg, args...)
		l.logger.Error(message)
	}
}

func (l *Logger) format(msg string, args ...interface{}) string {
	message := fmt.Sprintf(msg, args...)
	return strings.TrimRight(message, "\n\r")
}
