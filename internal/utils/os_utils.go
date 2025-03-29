package utils

import "runtime"

const (
	WindowsOS = "windows"
	LinuxOS   = "linux"
	MacOS     = "darwin"
)

var supportedOS = map[string]struct{}{
	LinuxOS:   {},
	WindowsOS: {},
	MacOS:     {},
}

func IsSupportedOS() bool {
	_, exist := supportedOS[runtime.GOOS]
	if !exist {
		return false
	}
	return true
}
