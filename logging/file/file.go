package file

import (
	"os"

	"github.com/dihedron/rafter/logging/stream"
)

// NewLogger returns a stream.Logger writing to a file at
// the given path.
func NewLogger(path string) *stream.Logger {
	file, err := os.Create(path)
	if err != nil {
		return nil
	}
	return stream.NewLogger(file)
}
