package logger

import "log"

type Interface interface {
	Info(args ...interface{})
	Error(args ...interface{})
}

type Logger struct{}

func New() *Logger {
	return &Logger{}
}

func (l *Logger) Info(args ...interface{}) {
	log.Println(args...)
}

func (l *Logger) Error(args ...interface{}) {
	log.Println("[ERROR]", args)
}
