package logging

type Logger interface {
	// Emit a message at the TRACE level
	Trace(format string, args ...interface{})

	// Emit a message at the DEBUG level
	Debug(format string, args ...interface{})

	// Emit a message at the INFO level
	Info(format string, args ...interface{})

	// Emit a message at the WARN level
	Warn(format string, args ...interface{})

	// Emit a message at the ERROR level
	Error(format string, args ...interface{})
}
