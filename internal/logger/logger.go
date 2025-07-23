package logger

import (
	"fmt"
	"log"
	"os"
)

var (
	Info  *log.Logger
	Error *log.Logger
)

func Init() {
	Info = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
}

type LogWriter struct {
	file *os.File
}

func NewLogWriter(path string) (*LogWriter, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return &LogWriter{file: f}, nil
}

func (lw *LogWriter) Write(step string, output string) error {
	_, err := fmt.Fprintf(lw.file, "[%s] %s\n", step, output)
	return err
}

func (lw *LogWriter) Close() error {
	return lw.file.Close()
}
