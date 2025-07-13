package path

import "runtime"

const (
	DefaultDownloadDir = "Downloads"
	AppDir = "bangumi-cli"
)

var (
	RunningOS = runtime.GOOS

	OSPathProviders = map[string]Provider{
		"windows": WindowsPath{},
		"darwin":  MacOSPath{},
	}
)

type Provider interface {
	ConfigPath() (string, error)
	DownloadPath() (string, error)
}

type WindowsPath struct{}

type MacOSPath struct{}
