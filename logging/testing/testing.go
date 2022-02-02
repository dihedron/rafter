package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/dihedron/rafter/logging"
)

// Logger wraps the Golang testing framework logger.
type Logger struct {
	t *testing.T
}

// NewLOgger returns a Logger wrapping a testing logger.
func NewLogger(t *testing.T) *Logger {
	return &Logger{
		t: t,
	}
}

// Trace logs a message at LevelTrace level.
func (l *Logger) Trace(msg string, args ...interface{}) {
	if logging.GetLevel() <= logging.LevelTrace {
		message := l.format("TRC", msg, args...)
		l.t.Log(message)
	}
}

// Debug logs a message at LevelDebug level.
func (l *Logger) Debug(msg string, args ...interface{}) {
	if logging.GetLevel() <= logging.LevelDebug {
		message := l.format("DBG", msg, args...)
		l.t.Log(message)
	}
}

// Info logs a message at LevelInfo level.
func (l *Logger) Info(msg string, args ...interface{}) {
	if logging.GetLevel() <= logging.LevelInfo {
		message := l.format("INF", msg, args...)
		l.t.Log(message)
	}
}

// Warn logs a message at LevelWarn level.
func (l *Logger) Warn(msg string, args ...interface{}) {
	if logging.GetLevel() <= logging.LevelWarn {
		message := l.format("WRN", msg, args...)
		l.t.Log(message)
	}
}

// Error logs a message at LevelError level.
func (l *Logger) Error(msg string, args ...interface{}) {
	if logging.GetLevel() <= logging.LevelError {
		message := l.format("ERR", msg, args...)
		l.t.Log(message)
	}
}

func (l *Logger) format(level string, msg string, args ...interface{}) string {
	message := fmt.Sprintf("["+level+"] "+strings.TrimSpace(msg), args...)
	return strings.TrimRight(message, "\n\r")
}
