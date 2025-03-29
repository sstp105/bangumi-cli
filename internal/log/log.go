package log

import "github.com/sirupsen/logrus"

var log = logrus.New()

func init() {
	log.SetLevel(logrus.DebugLevel)
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func Debug(args ...interface{}) {
	log.Debug(args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}