package logger

import (
	"log"
	"os"
	"path/filepath"

	"DistributeX/pkg/common"
)

type Logger struct {
	*log.Logger
}

func NewLogger(filename string) *Logger {
	dir := filepath.Dir(filename)
	common.CreateDirectory(dir)

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file %s: %v", filename, err)
	}

	return &Logger{
		Logger: log.New(file, "", log.LstdFlags),
	}
}

func (l *Logger) Info(msg string) {
	l.Println("INFO: " + msg)
}

func (l *Logger) Error(msg string) {
	l.Println("ERROR: " + msg)
}

func (l *Logger) Debug(msg string) {
	l.Println("DEBUG: " + msg)
}
