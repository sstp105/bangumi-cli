package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	path "github.com/sstp105/bangumi-cli/internal/path"
	"os"
	"time"
)

var (
	log = logrus.New()
)

func init() {
	log.SetLevel(logrus.DebugLevel)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:    true,
		TimestampFormat:  "2006-01-02 15:04:05",
		DisableTimestamp: false,
		DisableSorting:   false,
		ForceColors:      true,
	})

	date := time.Now().Format("2006-01-02")

	fn := fmt.Sprintf("%s.log", date)
	dir, err := path.LogPath(fn)
	if err != nil {
		log.Fatalf("error getting log file path: %s", err)
	}

	logFile, err := os.OpenFile(dir, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %s", err)
	}

	log.SetOutput(logFile)
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
