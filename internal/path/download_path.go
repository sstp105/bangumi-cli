package path

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/sstp105/bangumi-cli/internal/libs"
)

func (w WindowsPath) DownloadPath() (string, error) {
	dir := os.Getenv("USERPROFILE")
	if dir == "" {
		return "", errors.New("USERPROFILE is not defined")
	}
	return filepath.Join(dir, DefaultDownloadDir), nil
}

func (m MacOSPath) DownloadPath() (string, error) {
	dir := os.Getenv("HOME")
	if dir == "" {
		return "", errors.New("$HOME is not defined")
	}
	return filepath.Join(dir, DefaultDownloadDir), nil
}

func DownloadPath() (string, error) {
	provider, supported := OSPathProviders[RunningOS]
	if !supported {
		return "", libs.ErrUnsupportedOS
	}

	path, err := provider.DownloadPath()
	if err != nil {
		return "", err
	}

	return path, nil
}
