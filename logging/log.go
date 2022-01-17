package logging

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type LogLogger struct {
	BaseLogger
	logger *log.Logger
}

func NewLogLogger(prefix string) *LogLogger {
	return &LogLogger{
		logger: log.New(os.Stderr, prefix, log.Ltime|log.Ldate|log.Lmicroseconds),
		BaseLogger: BaseLogger{
			Values: []interface{}{},
		},
	}
}

func (l *LogLogger) Trace(msg string, args ...interface{}) {
	message := fmt.Sprintf("[TRC] "+msg, args...)
	message = strings.TrimRight(message, "\n\r")
	log.Print(message)
}

func (l *LogLogger) Debug(msg string, args ...interface{}) {
	message := fmt.Sprintf("[DBG] "+msg, args...)
	message = strings.TrimRight(message, "\n\r")
	log.Print(message)
}

func (l *LogLogger) Info(msg string, args ...interface{}) {
	message := fmt.Sprintf("[INF] "+msg, args...)
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
