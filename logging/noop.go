package logging

type NoOpLogger struct{}

func (*NoOpLogger) Trace(format string, args ...interface{}) {}

func (*NoOpLogger) Debug(format string, args ...interface{}) {}

func (*NoOpLogger) Info(format string, args ...interface{}) {}

func (*NoOpLogger) Warn(format string, args ...interface{}) {}

func (*NoOpLogger) Error(format string, args ...interface{}) {}
