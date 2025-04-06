package logger

import (
	"fmt"
	"os"
	"time"
)

type ILogger interface {
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
	Close()
}

type FileLogger struct {
	file *os.File
}

func NewFileLogger(filename string) (ILogger, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &FileLogger{file: file}, nil
}

func (l *FileLogger) Info(format string, v ...interface{}) {
	l.log("INFO", format, v...)
}

func (l *FileLogger) Error(format string, v ...interface{}) {
	l.log("ERROR", format, v...)
}

func (l *FileLogger) log(level, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	log := fmt.Sprintf("[%s] [%s] %s\n", 
		time.Now().Format("2006-01-02 15:04:05"),
		level,
		msg,
	)
	l.file.WriteString(log)
}

func (l *FileLogger) Close() {
	l.file.Close()
} 