package logging

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
)

const TimeFormat = "2006-01-02T15:04:05.999-0700"

type StreamLogger struct {
	stream *os.File
}

func NewStreamLogger(stream *os.File) *StreamLogger {
	return &StreamLogger{
		stream: stream,
	}
}

func (l *StreamLogger) Trace(msg string, args ...interface{}) {
	if isatty.IsTerminal(l.stream.Fd()) {
		l.write(color.HiWhiteString("TRC"), msg, args...)
	} else {
		l.write("TRC", msg, args...)
	}
}

func (l *StreamLogger) Debug(msg string, args ...interface{}) {
	if isatty.IsTerminal(l.stream.Fd()) {
		l.write(color.HiBlueString("DBG"), msg, args...)
	} else {
		l.write("DBG", msg, args...)
	}
}

func (l *StreamLogger) Info(msg string, args ...interface{}) {
	if isatty.IsTerminal(l.stream.Fd()) {
		l.write(color.HiGreenString("INF"), msg, args...)
	} else {
		l.write("INF", msg, args...)
	}
}

func (l *StreamLogger) Warn(msg string, args ...interface{}) {
	if isatty.IsTerminal(l.stream.Fd()) {
		l.write(color.HiYellowString("WRN"), msg, args...)
	} else {
		l.write("WRN", msg, args...)
	}
}

func (l *StreamLogger) Error(msg string, args ...interface{}) {
	if isatty.IsTerminal(l.stream.Fd()) {
		l.write(color.HiRedString("ERR"), msg, args...)
	} else {
		l.write("ERR", msg, args...)
	}
}

func (l *StreamLogger) write(level string, msg string, args ...interface{}) {
	message := fmt.Sprintf(strings.TrimSpace(msg), args...)
	fmt.Fprintf(l.stream, "%s [%s] %s\n", time.Now().Format(TimeFormat), level, message)
}
