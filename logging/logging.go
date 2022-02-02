package logging

import "sync"

// Level represents the logging level.
type Level uint8

const (
	LevelTrace Level = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelOff
)

var (
	lock    sync.RWMutex
	current Level = LevelInfo
)

// SetLevel sets the logging level globally.
func SetLevel(new Level) {
	lock.Lock()
	defer lock.Unlock()
	current = new
}

// GetLevel retrieves the current global logging level.
func GetLevel() Level {
	lock.RLock()
	defer lock.RUnlock()
	return current
}

// Logger is the common interface to all loggers.
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
