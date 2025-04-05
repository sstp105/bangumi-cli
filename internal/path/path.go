package path

import (
	"errors"
	"github.com/sstp105/bangumi-cli/internal/libs"
)

const (
	AppDir = "bangumi-cli"

	SubscribedBangumiConfigFile = "subscribed_bangumi.json"
	BangumiCredentialFile       = "bangumi_creds.json"
)

var (
	ErrWindowsAppDataEnvNotFound = errors.New("APPDATA env is not found")

	osPathProviders = map[string]Provider{
		libs.WindowsOS: WindowsPath{},
		libs.LinuxOS:   LinuxPath{},
		libs.MacOS:     MacOSPath{},
	}
)

type Provider interface {
	ConfigPath() ([]string, error)
	LogPath() ([]string, error)
}

type WindowsPath struct{}

type LinuxPath struct{}

type MacOSPath struct{}
