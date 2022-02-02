package stream

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dihedron/rafter/logging"
	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
)

const TimeFormat = "2006-01-02T15:04:05.999-0700"

// Logger is a logger that write sits messages to a stream.
type Logger struct {
	stream *os.File
}

// NewLogger returns an instance of a stream Logger.
func NewLogger(stream *os.File) *Logger {
	return &Logger{
		stream: stream,
	}
}

// Trace logs a message at LevelTrace level.
func (l *Logger) Trace(msg string, args ...interface{}) {
	if logging.GetLevel() <= logging.LevelTrace {
		if isatty.IsTerminal(l.stream.Fd()) {
			l.write(color.HiWhiteString("TRC"), msg, args...)
		} else {
			l.write("TRC", msg, args...)
		}
	}
}

// Debug logs a message at LevelDebug level.
func (l *Logger) Debug(msg string, args ...interface{}) {
	if logging.GetLevel() <= logging.LevelDebug {
		if isatty.IsTerminal(l.stream.Fd()) {
			l.write(color.HiBlueString("DBG"), msg, args...)
		} else {
			l.write("DBG", msg, args...)
		}
	}
}

// Info logs a message at LevelInfo level.
func (l *Logger) Info(msg string, args ...interface{}) {
	if logging.GetLevel() <= logging.LevelInfo {
		if isatty.IsTerminal(l.stream.Fd()) {
			l.write(color.HiGreenString("INF"), msg, args...)
		} else {
			l.write("INF", msg, args...)
		}
	}
}

// Warn logs a message at LevelWarn level.
func (l *Logger) Warn(msg string, args ...interface{}) {
	if logging.GetLevel() <= logging.LevelWarn {
		if isatty.IsTerminal(l.stream.Fd()) {
			l.write(color.HiYellowString("WRN"), msg, args...)
		} else {
			l.write("WRN", msg, args...)
		}
	}
}

// Error logs a message at LevelError level.
func (l *Logger) Error(msg string, args ...interface{}) {
	if logging.GetLevel() <= logging.LevelError {
		if isatty.IsTerminal(l.stream.Fd()) {
			l.write(color.HiRedString("ERR"), msg, args...)
		} else {
			l.write("ERR", msg, args...)
		}
	}
}

func (l *Logger) write(level string, msg string, args ...interface{}) {
	message := fmt.Sprintf(strings.TrimSpace(msg), args...)
	fmt.Fprintf(l.stream, "%s [%s] %s\n", time.Now().Format(TimeFormat), level, message)
}
