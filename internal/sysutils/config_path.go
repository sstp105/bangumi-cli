package sysutils

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

const (
	AppDir = "bangumi-cli"

	BangumiCredentialFile = "bangumi_creds.json"
)

var (
	ErrWindowsAppDataEnvNotFound = errors.New("APPDATA env is not found")
)

var osConfigProviders = map[string]ConfigPathProvider{
	WindowsOS: WindowsConfig{},
	LinuxOS:   LinuxConfig{},
	MacOS:     MacOSConfig{},
}

type ConfigPathProvider interface {
	ConfigPath() ([]string, error)
}

type WindowsConfig struct{}

func (w WindowsConfig) ConfigPath() ([]string, error) {
	dir := os.Getenv("APPDATA")
	if dir == "" {
		return nil, ErrWindowsAppDataEnvNotFound
	}
	return []string{dir, AppDir}, nil
}

type LinuxConfig struct{}

func (l LinuxConfig) ConfigPath() ([]string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return []string{dir, ".config", AppDir}, nil
}

type MacOSConfig struct{}

func (m MacOSConfig) ConfigPath() ([]string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return []string{dir, ".config", AppDir}, nil
}

// SaveJSONConfigFile saves the value in json format under app's config directory.
// For Windows, the file will be saved under $env:APPDATA.
// For UNIX, the file will be saved under $HOME.
// An error will be returned for unsupported OS.
// Parameters:
//   - fn: The file name (e.g., "setting.json").
//   - v: The value to be saved to the file. It will be marshaled into JSON format.
//
// Returns:
//   - error: If any error occurs during the process, an error is returned.
func SaveJSONConfigFile(fn string, v any) error {
	path, err := configPath(fn)
	if err != nil {
		return err
	}

	data, err := MarshalJSONIndented(v)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600) // owner r&w
}

func ReadJSONConfigFile(fn string, v any) error {
	path, err := configPath(fn)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		return err
	}

	return nil
}

func configPath(fn string) (string, error) {
	provider, supported := osConfigProviders[runtime.GOOS]
	if !supported {
		return "", ErrUnsupportedOS
	}

	path, err := provider.ConfigPath()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(path...)
	if err := os.MkdirAll(dir, 0700); err != nil { // create the config folder if does not exist
		return "", err
	}

	return filepath.Join(dir, fn), nil
}
