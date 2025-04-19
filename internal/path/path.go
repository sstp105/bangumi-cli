package path

import "runtime"

const (
	AppDir = "bangumi-cli"
)

var (
	runningOS = runtime.GOOS

	osPathProviders = map[string]Provider{
		"windows": WindowsPath{},
		"darwin":  MacOSPath{},
	}
)

type Provider interface {
	ConfigPath() (string, error)
}

type WindowsPath struct{}

type MacOSPath struct{}
