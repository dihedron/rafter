package logging

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/fatih/color"
)

func NewConsoleLogger(stream io.Writer) *ConsoleLogger {
	return &ConsoleLogger{
		stream: stream,
	}
}

const TimeFormat = "2006-01-02T15:04:05.999-0700"

type ConsoleLogger struct {
	stream io.Writer
}

func (l *ConsoleLogger) Trace(msg string, args ...interface{}) {
	message := strings.TrimSpace(fmt.Sprintf(msg, args...))
	fmt.Fprintf(l.stream, "%s [%s] %s\n", time.Now().Format(TimeFormat), color.HiWhiteString("TRACE"), message)
}

func (l *ConsoleLogger) Debug(msg string, args ...interface{}) {
	message := strings.TrimSpace(fmt.Sprintf(msg, args...))
	fmt.Fprintf(l.stream, "%s [%s] %s\n", time.Now().Format(TimeFormat), color.HiBlueString("DEBUG"), message)
}

func (l *ConsoleLogger) Info(msg string, args ...interface{}) {
	message := strings.TrimSpace(fmt.Sprintf(msg, args...))
	fmt.Fprintf(l.stream, "%s [%s] %s\n", time.Now().Format(TimeFormat), color.HiGreenString("INFO"), message)
}

func (l *ConsoleLogger) Warn(msg string, args ...interface{}) {
	message := strings.TrimSpace(fmt.Sprintf(msg, args...))
	fmt.Fprintf(l.stream, "%s [%s] %s\n", time.Now().Format(TimeFormat), color.HiYellowString("WARN"), message)
}

func (l *ConsoleLogger) Error(msg string, args ...interface{}) {
	message := strings.TrimSpace(fmt.Sprintf(msg, args...))
	fmt.Fprintf(l.stream, "%s [%s] %s\n", time.Now().Format(TimeFormat), color.HiRedString("ERROR"), message)
}
