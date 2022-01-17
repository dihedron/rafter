package logging

type Logger interface {
	// Trace emits a message at the TRACE level.
	Trace(format string, args ...interface{})

	// Debug emits a message at the DEBUG level.
	Debug(format string, args ...interface{})

	// Info emits a message at the INFO level.
	Info(format string, args ...interface{})

	// Warn emits a message at the WARN level.
	Warn(format string, args ...interface{})

	// Error emits a message at the ERROR level.
	Error(format string, args ...interface{})
}
