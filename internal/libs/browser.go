package libs

import (
	"errors"
	"os/exec"
	"runtime"
)

var (
	osOpenCommands = map[string][]string{
		LinuxOS:   {"xdg-open"},
		WindowsOS: {"rundll32", "url.dll,FileProtocolHandler"},
		MacOS:     {"open"},
	}
)

var ErrUnsupportedOS = errors.New("unsupported os")

func OpenBrowser(url string) error {
	v, supported := osOpenCommands[runtime.GOOS]
	if !supported {
		return ErrUnsupportedOS
	}

	command := v[0]
	args := append(v[1:], url)

	return exec.Command(command, args...).Start()
}
