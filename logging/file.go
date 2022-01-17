package logging

import "os"

type FileLogger StreamLogger

func NewFileLogger(path string) *FileLogger {
	file, err := os.Create(path)
	if err != nil {
		return nil
	}
	return &FileLogger{
		stream: file,
	}
}
