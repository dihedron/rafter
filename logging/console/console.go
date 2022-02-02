package console

import (
	"os"

	"github.com/dihedron/rafter/logging/stream"
)

type Where int8

const (
	StdOut Where = iota
	StdErr
)

// NewLogger returns a stream.Logger writing to either
// StdOut or StdErr.
func NewLogger(where Where) *stream.Logger {
	switch where {
	case StdOut:
		return stream.NewLogger(os.Stdout)
	case StdErr:
		return stream.NewLogger(os.Stderr)
	}
	return nil
}
