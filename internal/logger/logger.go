// internal/logger/logger.go
package logger

import (
	"log"
	"os"
)

type Logger struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	DebugLog *log.Logger
}

func New() *Logger {
	return &Logger{
		InfoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		DebugLog: log.New(os.Stdout, "DEBUG\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *Logger) Info(v ...interface{}) {
	l.InfoLog.Println(v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.ErrorLog.Println(v...)
}

func (l *Logger) Debug(v ...interface{}) {
	if os.Getenv("DEBUG") == "true" {
		l.DebugLog.Println(v...)
	}
}
