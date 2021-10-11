package log

import (
	"os"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

// InitLog loads our logger
func InitLog(lvl string) *logrus.Logger {

	var level logrus.Level
	level, _ = logrus.ParseLevel(lvl)
	logger := &logrus.Logger{
		Out:   os.Stdout,
		Level: level,
	}

	if lvl == "info" {
		logger.SetFormatter(&nested.Formatter{
			HideKeys:    true,
			FieldsOrder: []string{"component", "category"},
		})
	} else {
		logger.SetFormatter(&nested.Formatter{
			HideKeys:    true,
			FieldsOrder: []string{"component", "category"},
		})
	}

	Log = logger
	return Log
	
}

func Info(args... interface{}) {
	Log.Info(args)
}

func Debug(args... interface{}) {
	Log.Debug(args)
}

func Error(args... interface{}) {
	Log.Error(args)
}

func Trace(args... interface{}) {
	Log.Trace(args)
}
