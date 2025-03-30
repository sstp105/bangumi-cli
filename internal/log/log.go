package log

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type PlainFormatter struct {
}

func init() {
	log.SetLevel(logrus.DebugLevel)
	log.SetFormatter(&PlainFormatter{})
}

func (f *PlainFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("%s\n", entry.Message)), nil
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

func Fatal(args ...interface{}) {
	log.Fatal(args...)
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

func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}
