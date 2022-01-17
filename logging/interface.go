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

	// With returns a new logger with an added value to its context.
	With(values ...interface{}) Logger
}

/*
type LogFunc func(format string, args ...interface{})

type Logger2 struct {
	Trace   LogFunc
	Debug   LogFunc
	Info    LogFunc
	Warn    LogFunc
	Error   LogFunc
	Context []interface{}
}

type ConsoleLogger2 struct {
	Logger2

}

func NewConsoleLogger2(where Where) *ConsoleLogger {
	switch where {
	case StdOut:
		return &Logger{
			stream: os.Stdout,
		}
	case StdErr:
		return &ConsoleLogger{
			stream: os.Stderr,
		}
	}
	return nil
}
*/
