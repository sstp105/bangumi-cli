package log

import (
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Logger

	info    = color.New(color.FgCyan).SprintFunc()
	debug   = color.New(color.FgWhite).SprintFunc()
	success = color.New(color.FgGreen).SprintFunc()
	warn    = color.New(color.FgYellow).SprintFunc()
	failure = color.New(color.FgRed).SprintFunc()
	prompt  = color.New(color.FgBlue, color.Bold).SprintFunc()

	infoF    = color.New(color.FgCyan).SprintfFunc()
	debugF   = color.New(color.FgWhite).SprintfFunc()
	successF = color.New(color.FgGreen).SprintfFunc()
	warnF    = color.New(color.FgYellow).SprintfFunc()
	failureF = color.New(color.FgRed).SprintfFunc()
	promptF  = color.New(color.FgBlue, color.Bold).SprintfFunc()
)

func init() {
	log = logrus.New()

	log.SetLevel(logrus.DebugLevel)
	log.SetFormatter(&PlainFormatter{})
}

func Info(args ...interface{}) {
	log.Info(info(args...))
}

func Success(args ...interface{}) {
	log.Info(success(args...))
}

func Prompt(args ...interface{}) {
	log.Info(prompt(args...))
}

func Debug(args ...interface{}) {
	log.Debug(debug(args...))
}

func Warn(args ...interface{}) {
	log.Warn(warn(args...))
}

func Error(args ...interface{}) {
	log.Error(failure(args...))
}

func Fatal(args ...interface{}) {
	log.Fatal(failure(args...))
}

func Infof(format string, args ...interface{}) {
	log.Info(infoF(format, args...))
}

func Successf(format string, args ...interface{}) {
	log.Info(successF(format, args...))
}

func Promptf(format string, args ...interface{}) {
	log.Info(promptF(format, args...))
}

func Debugf(format string, args ...interface{}) {
	log.Debug(debugF(format, args...))
}

func Warnf(format string, args ...interface{}) {
	log.Warn(warnF(format, args...))
}

func Errorf(format string, args ...interface{}) {
	log.Error(failureF(format, args...))
}

func Fatalf(format string, args ...interface{}) {
	log.Fatal(failureF(format, args...))
}
