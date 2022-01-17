package logging

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-hclog"
)

type HCLLogger struct {
	BaseLogger
	logger hclog.Logger
}

func NewHCLLogger(logger hclog.Logger) *HCLLogger {
	return &HCLLogger{
		BaseLogger: BaseLogger{
			Values: []interface{}{},
		},
		logger: logger,
	}
}

func (l *HCLLogger) Trace(msg string, args ...interface{}) {
	message := l.format(msg, args...)
	l.logger.Trace(message)
}

func (l *HCLLogger) Debug(msg string, args ...interface{}) {
	message := l.format(msg, args...)
	l.logger.Debug(message)
}

func (l *HCLLogger) Info(msg string, args ...interface{}) {
	message := l.format(msg, args...)
	l.logger.Info(message)
}

func (l *HCLLogger) Warn(msg string, args ...interface{}) {
	message := l.format(msg, args...)
	l.logger.Warn(message)
}

func (l *HCLLogger) Error(msg string, args ...interface{}) {
	message := l.format(msg, args...)
	l.logger.Error(message)
}

func (l *HCLLogger) format(msg string, args ...interface{}) string {
	message := ""
	if len(l.Values) > 0 {
		args = append(args, l.Values)
		message = fmt.Sprintf(msg+" (context: %+v)", args...)
	} else {
		message = fmt.Sprintf(msg, args...)
	}
	return strings.TrimRight(message, "\n\r")
}
