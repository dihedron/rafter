package logging

import "go.uber.org/zap"

// ZapLogger is an adapter that allows to log using Uber's Zap
// wherever a Logger interface is expected.
type ZapLogger struct {
	BaseLogger
	logger *zap.Logger
}

func NewZapLogger() *ZapLogger {
	return &ZapLogger{
		BaseLogger: BaseLogger{
			Values: []interface{}{},
		},
		logger: zap.New(zap.L().Core(), zap.AddCallerSkip(1)),
	}
}

func (l *ZapLogger) Trace(format string, args ...interface{}) {
	zap.S().Debugf(format, args...)
}

func (l *ZapLogger) Debug(format string, args ...interface{}) {
	zap.S().Debugf(format, args)
}

func (l *ZapLogger) Info(format string, args ...interface{}) {
	zap.S().Infof(format, args)
}

func (l *ZapLogger) Warn(format string, args ...interface{}) {
	zap.S().Warnf(format, args)
}

func (l *ZapLogger) Error(format string, args ...interface{}) {
	zap.S().Errorf(format, args)
}
