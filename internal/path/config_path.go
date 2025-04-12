package path

import (
	"encoding/json"
	"github.com/sstp105/bangumi-cli/internal/libs"
	"os"
	"path/filepath"
	"runtime"
)

const (
	BangumiCredentialConfigFile = "bangumi_creds.json"
	SubscriptionConfigFile      = "subscriptions.json"
)

// ConfigPath returns the Windows path to the app's config directory in %AppData%\{APP_NAME}.
func (w WindowsPath) ConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, AppDir), nil
}

// ConfigPath returns the Linux path to the app's config directory in $HOME/.config/<APP_NAME>.
func (l LinuxPath) ConfigPath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, ".config", AppDir), nil
}

// ConfigPath returns the macOS path to the app's config directory in $HOME/.config/<APP_NAME>.
func (m MacOSPath) ConfigPath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, ".config", AppDir), nil
}

// SaveJSONConfigFile saves the value in json format under app's config directory.
// For Windows, the file will be saved under $env:APPDATA.
// For UNIX, the file will be saved under $HOME.
// An error will be returned for unsupported OS.
// Parameters:
//   - fn: The file name (e.g., "setting.json").
//   - v: The value to be saved to the file. It will be marshaled into JSON format.
func SaveJSONConfigFile(fn string, v any) error {
	path, err := configPath(fn)
	if err != nil {
		return err
	}

	data, err := libs.MarshalJSONIndented(v)
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

func DeleteJSONConfigFile(fn string) error {
	path, err := configPath(fn)
	if err != nil {
		return err
	}

	return os.Remove(path)
}

func configPath(fn string) (string, error) {
	provider, supported := osPathProviders[runtime.GOOS]
	if !supported {
		return "", libs.ErrUnsupportedOS
	}

	path, err := provider.ConfigPath()
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(path, 0700); err != nil { // create the config folder if it does not exist
		return "", err
	}

	return filepath.Join(path, fn), nil
}
