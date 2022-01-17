package logging

import "os"

type ConsoleLogger StreamLogger

type Where int8

const (
	StdOut Where = iota
	StdErr
)

func NewConsoleLogger(where Where) *ConsoleLogger {
	switch where {
	case StdOut:
		return &ConsoleLogger{
			stream: os.Stdout,
		}
	case StdErr:
		return &ConsoleLogger{
			stream: os.Stderr,
		}
	}
	return nil
}
