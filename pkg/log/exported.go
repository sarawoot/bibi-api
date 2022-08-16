package log

import "github.com/sirupsen/logrus"

var logger *Logger

func init() {
	logger = new()
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func SetDebugLevel() {
	logger.SetDebugLevel()
}

func SetErrorLevel() {
	logger.SetErrorLevel()
}

func SetInfoLevel() {
	logger.SetInfoLevel()
}

func WithFields(kv map[string]interface{}) *Logger {
	l := new()
	l.entry = l.entry.WithFields(logrus.Fields(kv))

	return l
}
