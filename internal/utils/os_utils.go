package utils

import (
	"errors"
	"os/exec"
	"runtime"
)

const (
	WindowsOS = "windows"
	LinuxOS   = "linux"
	MacOS     = "darwin"
)

var (
	supportedOS = map[string]struct{}{
		LinuxOS:   {},
		WindowsOS: {},
		MacOS:     {},
	}

	osOpenCommands = map[string][]string{
		LinuxOS:   {"xdg-open"},
		WindowsOS: {"rundll32", "url.dll,FileProtocolHandler"},
		MacOS:     {"open"},
	}
)

var ErrUnsupportedOS = errors.New("unsupported os")

func IsSupportedOS() bool {
	_, exist := supportedOS[runtime.GOOS]
	if !exist {
		return false
	}
	return true
}

func OpenBrowser(url string) error {
	goos := runtime.GOOS
	v, supported := osOpenCommands[goos]
	if !supported {
		return ErrUnsupportedOS
	}

	command := v[0]
	args := append(v[1:], url)

	return exec.Command(command, args...).Start()
}
