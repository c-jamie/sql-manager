package log

import (
	"os"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func New(level string) *logrus.Logger {
	logLevel, _ := logrus.ParseLevel(level)
	logger := &logrus.Logger{
		Out:   os.Stdout,
		Level: logLevel,
	}
	logger.SetFormatter(&logrus.TextFormatter{DisableColors: true, DisableTimestamp: true})
	Log = logger
	logger.Info("logging started")
	return logger
}

func Info(args ...interface{}) {
	Log.Info(args)
}

func Debug(args ...interface{}) {
	Log.Debug(args)
}

func Error(args ...interface{}) {
	Log.Error(args)
}

func Fatal(args ...interface{}) {
	Log.Fatal(args)
}

func Errorf(format string, args ...interface{}) {
	Log.Errorf(format, args)
}
func Debugf(format string, args ...interface{}) {
	Log.Debugf(format, args)
}