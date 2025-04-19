package libs

import (
	"errors"
	"os/exec"
	"runtime"
)

const (
	WindowsOS = "windows"
	MacOS     = "darwin"
)

var (
	osOpenCommands = map[string][]string{
		WindowsOS: {"rundll32", "url.dll,FileProtocolHandler"},
		MacOS:     {"open"},
	}

	ErrUnsupportedOS = errors.New("unsupported os")
)

func OpenBrowser(url string) error {
	v, supported := osOpenCommands[runtime.GOOS]
	if !supported {
		return ErrUnsupportedOS
	}

	command := v[0]
	args := append(v[1:], url)

	return exec.Command(command, args...).Start()
}
