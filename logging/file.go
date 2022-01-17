package logging

import "os"

func NewFileLogger(path string) *StreamLogger {
	file, err := os.Create(path)
	if err != nil {
		return nil
	}
	return &StreamLogger{
		stream: file,
	}
}
