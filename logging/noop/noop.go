package noop

// Logger is a logger that writes nothing.
type Logger struct{}

// Trace logs a message at LevelTrace level.
func (*Logger) Trace(format string, args ...interface{}) {}

// Debug logs a message at LevelDebug level.
func (*Logger) Debug(format string, args ...interface{}) {}

// Info logs a message at LevelInfo level.
func (*Logger) Info(format string, args ...interface{}) {}

// Warn logs a message at LevelWarn level.
func (*Logger) Warn(format string, args ...interface{}) {}

// Error logs a message at LevelError level.
func (*Logger) Error(format string, args ...interface{}) {}
