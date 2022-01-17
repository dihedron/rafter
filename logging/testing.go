package logging

import (
	"fmt"
	"strings"
	"testing"
)

type TestingLogger struct {
	t *testing.T
}

func NewTestingLogger(t *testing.T) *TestingLogger {
	return &TestingLogger{
		t: t,
	}
}

func (l *TestingLogger) Trace(msg string, args ...interface{}) {
	message := l.format("TRC", msg, args...)
	l.t.Log(message)
}

func (l *TestingLogger) Debug(msg string, args ...interface{}) {
	message := l.format("DBG", msg, args...)
	l.t.Log(message)
}

func (l *TestingLogger) Info(msg string, args ...interface{}) {
	message := l.format("INF", msg, args...)
	l.t.Log(message)
}

func (l *TestingLogger) Warn(msg string, args ...interface{}) {
	message := l.format("WRN", msg, args...)
	l.t.Log(message)
}

func (l *TestingLogger) Error(msg string, args ...interface{}) {
	message := l.format("ERR", msg, args...)
	l.t.Log(message)
}

func (l *TestingLogger) format(level string, msg string, args ...interface{}) string {
	message := fmt.Sprintf("["+level+"] "+strings.TrimSpace(msg), args...)
	return strings.TrimRight(message, "\n\r")
}
