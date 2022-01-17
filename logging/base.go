package logging

type BaseLogger struct {
	Values []interface{}
}

func (*BaseLogger) Trace(format string, args ...interface{}) {}

func (*BaseLogger) Debug(format string, args ...interface{}) {}

func (*BaseLogger) Info(format string, args ...interface{}) {}

func (*BaseLogger) Warn(format string, args ...interface{}) {}

func (*BaseLogger) Error(format string, args ...interface{}) {}

func (l *BaseLogger) With(values ...interface{}) Logger {
	l.Values = append(l.Values, values...)
	return l
}
