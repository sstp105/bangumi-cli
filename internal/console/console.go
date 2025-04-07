package console

import (
	"fmt"
	"github.com/fatih/color"
)

var (
	plain   = color.New(color.FgWhite).SprintFunc()
	info    = color.New(color.FgCyan).SprintFunc()
	success = color.New(color.FgGreen).SprintFunc()
	warning = color.New(color.FgYellow).SprintFunc()
	failure = color.New(color.FgRed).SprintFunc()

	plainF   = color.New(color.FgWhite).SprintfFunc()
	infoF    = color.New(color.FgCyan).SprintfFunc()
	successF = color.New(color.FgGreen).SprintfFunc()
	warningF = color.New(color.FgYellow).SprintfFunc()
	failureF = color.New(color.FgRed).SprintfFunc()
)

func Plainf(format string, args ...interface{}) {
	fmt.Println(plainF(format, args...))
}

func Infof(format string, args ...interface{}) {
	fmt.Println(infoF(format, args...))
}

func Successf(format string, args ...interface{}) {
	fmt.Println(successF(format, args...))
}

func Warningf(format string, args ...interface{}) {
	fmt.Println(warningF(format, args...))
}

func Errorf(format string, args ...interface{}) {
	fmt.Println(failureF(format, args...))
}

func Plain(args ...interface{}) {
	fmt.Print(plain(args...) + "\n")
}

func Info(args ...interface{}) {
	fmt.Print(info(args...) + "\n")
}

func Success(args ...interface{}) {
	fmt.Print(success(args...) + "\n")
}

func Warning(args ...interface{}) {
	fmt.Print(warning(args...) + "\n")
}

func Error(args ...interface{}) {
	fmt.Print(failure(args...) + "\n")
}
