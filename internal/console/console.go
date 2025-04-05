package console

import (
	"fmt"
	"github.com/fatih/color"
)

var (
	info    = color.New(color.FgCyan).SprintfFunc()
	success = color.New(color.FgGreen).SprintfFunc()
	warning = color.New(color.FgYellow).SprintfFunc()
	failure = color.New(color.FgRed).SprintfFunc()

	infoF    = color.New(color.FgCyan).SprintfFunc()
	successF = color.New(color.FgGreen).SprintfFunc()
	warningF = color.New(color.FgYellow).SprintfFunc()
	failureF = color.New(color.FgRed).SprintfFunc()
)

func Infof(format string, args ...interface{}) {
	fmt.Printf(infoF(format, args...) + "\n")
}

func Successf(format string, args ...interface{}) {
	fmt.Printf(successF(format, args...) + "\n")
}

func Warningf(format string, args ...interface{}) {
	fmt.Printf(warningF(format, args...) + "\n")
}

func Errorf(format string, args ...interface{}) {
	fmt.Printf(failureF(format, args...) + "\n")
}

func Info(format string, args ...interface{}) {
	fmt.Print(infoF(format, args...) + "\n")
}

func Success(format string, args ...interface{}) {
	fmt.Print(successF(format, args...) + "\n")
}

func Warning(format string, args ...interface{}) {
	fmt.Print(warningF(format, args...) + "\n")
}

func Error(format string, args ...interface{}) {
	fmt.Print(failureF(format, args...) + "\n")
}
