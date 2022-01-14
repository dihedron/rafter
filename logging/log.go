package logging

import (
	"fmt"
	"log"
	"strings"
)

func NewLogLogger() *LogLogger {
	return &LogLogger{}
}

type LogLogger struct{}

func (l *LogLogger) Trace(msg string, args ...interface{}) {
	message := fmt.Sprintf("[TRACE] "+msg, args...)
	message = strings.TrimRight(message, "\n\r")
	log.Print(message)
}

func (l *LogLogger) Debug(msg string, args ...interface{}) {
	message := fmt.Sprintf("[DEBUG] "+msg, args...)
	message = strings.TrimRight(message, "\n\r")
	log.Print(message)
}

func (l *LogLogger) Info(msg string, args ...interface{}) {
	message := fmt.Sprintf("[INFO] "+msg, args...)
	message = strings.TrimRight(message, "\n\r")
	log.Print(message)
}

func (l *LogLogger) Warn(msg string, args ...interface{}) {
	message := fmt.Sprintf("[WARN] "+msg, args...)
	message = strings.TrimRight(message, "\n\r")
	log.Print(message)
}

func (l *LogLogger) Error(msg string, args ...interface{}) {
	message := fmt.Sprintf("[ERROR] "+msg, args...)
	message = strings.TrimRight(message, "\n\r")
	log.Print(message)
}
