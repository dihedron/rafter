package logging

import "os"

type Where int8

const (
	StdOut Where = iota
	StdErr
)

func NewConsoleLogger(where Where) *StreamLogger {
	switch where {
	case StdOut:
		return &StreamLogger{
			stream: os.Stdout,
		}
	case StdErr:
		return &StreamLogger{
			stream: os.Stderr,
		}
	}
	return nil
}
