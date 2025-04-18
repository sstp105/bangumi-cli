package path

import "runtime"

const (
	AppDir = "bangumi-cli"
)

var (
	runningOS = runtime.GOOS

	osPathProviders = map[string]Provider{
		"windows": WindowsPath{},
		"linux":   LinuxPath{},
		"darwin":  MacOSPath{},
	}
)

type Provider interface {
	ConfigPath() (string, error)
	LogPath() (string, error)
}

type WindowsPath struct{}

type LinuxPath struct{}

type MacOSPath struct{}
