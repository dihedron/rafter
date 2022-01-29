package logging

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-hclog"
)

type HCLLogger struct {
	logger hclog.Logger
}

func NewHCLLogger(logger hclog.Logger) *HCLLogger {
	return &HCLLogger{
		logger: logger,
	}
}

func (l *HCLLogger) Trace(msg string, args ...interface{}) {
	if GetLevel() <= LevelTrace {
		message := l.format(msg, args...)
		l.logger.Trace(message)
	}
}

func (l *HCLLogger) Debug(msg string, args ...interface{}) {
	if GetLevel() <= LevelDebug {
		message := l.format(msg, args...)
		l.logger.Debug(message)
	}
}

func (l *HCLLogger) Info(msg string, args ...interface{}) {
	if GetLevel() <= LevelInfo {
		message := l.format(msg, args...)
		l.logger.Info(message)
	}
}

func (l *HCLLogger) Warn(msg string, args ...interface{}) {
	if GetLevel() <= LevelWarn {
		message := l.format(msg, args...)
		l.logger.Warn(message)
	}
}

func (l *HCLLogger) Error(msg string, args ...interface{}) {
	if GetLevel() <= LevelError {
		message := l.format(msg, args...)
		l.logger.Error(message)
	}
}

func (l *HCLLogger) format(msg string, args ...interface{}) string {
	message := fmt.Sprintf(msg, args...)
	return strings.TrimRight(message, "\n\r")
}
