package log

import (
	"time"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	entry *logrus.Entry
}

func new() *Logger {
	l := logrus.New()
	l.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: "ts",
			logrus.FieldKeyMsg:  "msg",
		},
	})

	return &Logger{entry: logrus.NewEntry(l)}
}

func (l *Logger) Info(args ...interface{}) {
	l.entry.Info(args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.entry.Warn(args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.entry.Error(args...)
}

func (l *Logger) Debug(args ...interface{}) {
	l.entry.Debug(args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.entry.Fatal(args...)
}

func (l *Logger) SetDebugLevel() {
	l.entry.Logger.SetLevel(logrus.DebugLevel)
}

func (l *Logger) SetErrorLevel() {
	l.entry.Logger.SetLevel(logrus.ErrorLevel)
}

func (l *Logger) SetInfoLevel() {
	l.entry.Logger.SetLevel(logrus.InfoLevel)
}
