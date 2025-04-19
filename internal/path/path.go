package path

const (
	AppDir = "bangumi-cli"
)

var (
	osPathProviders = map[string]Provider{
		"windows": WindowsPath{},
		"linux":   LinuxPath{},
		"darwin":  MacOSPath{},
	}
)

type Provider interface {
	ConfigPath() (string, error)
}

type WindowsPath struct{}

type LinuxPath struct{}

type MacOSPath struct{}
