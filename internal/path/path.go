package path

import (
	"github.com/sstp105/bangumi-cli/internal/libs"
)

const (
	AppDir = "bangumi-cli"
)

var (
	osPathProviders = map[string]Provider{
		libs.WindowsOS: WindowsPath{},
		libs.LinuxOS:   LinuxPath{},
		libs.MacOS:     MacOSPath{},
	}
)

type Provider interface {
	ConfigPath() (string, error)
	LogPath() (string, error)
}

type WindowsPath struct{}

type LinuxPath struct{}

type MacOSPath struct{}
