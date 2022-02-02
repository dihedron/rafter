package golang

import (
	"fmt"
	golang "log"
	"os"
	"strings"

	"github.com/dihedron/rafter/logging"
)

// Logger is te type wrapping the default Golang logger.
type Logger struct {
	logger *golang.Logger
}

// NewLogger returns a new Golang Logger.
func NewLogger(prefix string) *Logger {
	return &Logger{
		logger: golang.New(os.Stderr, prefix, golang.Ltime|golang.Ldate|golang.Lmicroseconds),
	}
}

func (l *Logger) Trace(msg string, args ...interface{}) {
	if logging.GetLevel() <= logging.LevelTrace {
		message := fmt.Sprintf("[TRC] "+msg, args...)
		message = strings.TrimRight(message, "\n\r")
		golang.Print(message)
	}
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	if logging.GetLevel() <= logging.LevelDebug {
		message := fmt.Sprintf("[DBG] "+msg, args...)
		message = strings.TrimRight(message, "\n\r")
		golang.Print(message)
	}
}

func (l *Logger) Info(msg string, args ...interface{}) {
	if logging.GetLevel() <= logging.LevelInfo {
		message := fmt.Sprintf("[INF] "+msg, args...)
		message = strings.TrimRight(message, "\n\r")
		golang.Print(message)
	}
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	if logging.GetLevel() <= logging.LevelWarn {
		message := fmt.Sprintf("[WARN] "+msg, args...)
		message = strings.TrimRight(message, "\n\r")
		golang.Print(message)
	}
}

func (l *Logger) Error(msg string, args ...interface{}) {
	if logging.GetLevel() <= logging.LevelError {
		message := fmt.Sprintf("[ERROR] "+msg, args...)
		message = strings.TrimRight(message, "\n\r")
		golang.Print(message)
	}
}
